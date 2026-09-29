#!/usr/bin/env node
// gcsim wasm benchmark runner + correctness dumper. No npm dependencies.
//
// Drives cmd/wasm through the same JS API and call sequence as the web UI
// (ui/packages/executors/src/WasmExecutor.ts + Workers/{worker,aggregator}.ts):
//   sim worker:  go.run(instance); initializeWorker(cfg); simulate() -> Uint8Array (msgpack), repeat
//   aggregator:  go.run(instance); initializeAggregator(cfg) -> JSON; aggregate(Uint8Array); flush() -> JSON
//
// Each wasm instance gets its own JS realm (a vm context, or the worker_thread's own global in
// pool mode), like each browser Worker has its own global scope. Results passed from a sim
// instance to an aggregator are copied into the aggregator's realm, as the UI's postMessage does.
//
// Seeds: cmd/wasm seeds each iteration from crypto/rand, which on js/wasm reads
// crypto.getRandomValues. While simulate() runs we serve those bytes from a splitmix64 stream
// keyed by (--seed, iteration index), so every binary simulates the same seeds. Outside
// simulate() (runtime init, sample_seed) the real CSPRNG is used, as in the UI.
//
// Subcommands:
//   run   --wasm F --config F [--iters N] [--warmup N] [--seed N] [--workers N] [--label S] [--round N] [--json]
//         --workers 0 (default): sim + aggregator instance on one thread; simulate() is called
//                                back to back like a UI worker, then the batch is aggregated.
//                                Primary metric: ms/iter of simulate().
//         --workers N (N>=1):    N sim worker_threads + 1 aggregator thread, dispatched exactly like
//                                WasmExecutor.run (incl. 100ms throttled flush). Metric: wall ms / iters.
//   ab    --wasm A --wasm B [...] --config F [--block N] [--warmup N] [--rounds R] [--proc P] [--json]
//         paired comparison in one process: every binary is loaded side by side, then each round
//         runs the same --block seeds on each binary back to back (order rotates per round).
//         Emits one record per (binary, round). This is what compare.sh uses by default.
//   dump  --wasm F --config F [--iters N] [--seed N] [--flush-every K] [--out F]
//         canonical per-iteration + aggregate hashes; --flush-every K also flushes (and hashes the
//         stats) after every K aggregated iterations, as the UI's throttled flushes do
//   diff  golden.json candidate.json                            exit 1 on mismatch
//   report results.jsonl...                                      comparison table (markdown)
//
// A binary's JS glue defaults to <name>.wasm_exec.js next to it (written by build.sh), then
// $(go env GOROOT)/lib/wasm/wasm_exec.js. Override with --glue.

import { execFileSync } from "node:child_process";
import { createHash } from "node:crypto";
import fs from "node:fs";
import path from "node:path";
import vm from "node:vm";
import { isMainThread, parentPort, Worker, workerData } from "node:worker_threads";

// ---------------------------------------------------------------------------------------------
// deterministic seeds

const MASK64 = (1n << 64n) - 1n;
const GOLDEN = 0x9e3779b97f4a7c15n;

function splitmix64(state) {
	let s = BigInt.asUintN(64, state);
	return () => {
		s = (s + GOLDEN) & MASK64;
		let z = s;
		z = ((z ^ (z >> 30n)) * 0xbf58476d1ce4e5b9n) & MASK64;
		z = ((z ^ (z >> 27n)) * 0x94d049bb133111ebn) & MASK64;
		return z ^ (z >> 31n);
	};
}

// First 8 bytes served during iteration i; this is the uint64 that ends up in Result.seed.
function iterStreamState(base, i) {
	return (BigInt(base) * GOLDEN + BigInt(i)) & MASK64;
}
function expectedSeed(base, i) {
	return splitmix64(iterStreamState(base, i))();
}

let rngOverride = null;
const realGetRandomValues = globalThis.crypto.getRandomValues.bind(globalThis.crypto);
function seededGetRandomValues(arr) {
	if (rngOverride == null) return realGetRandomValues(arr);
	const bytes = new Uint8Array(arr.buffer, arr.byteOffset, arr.byteLength);
	for (let off = 0; off < bytes.length; off += 8) {
		let v = rngOverride();
		for (let j = 0; j < 8 && off + j < bytes.length; j++) {
			bytes[off + j] = Number(v & 0xffn);
			v >>= 8n;
		}
	}
	return arr;
}
// for instances started in this realm (pool mode worker_threads)
globalThis.crypto.getRandomValues = seededGetRandomValues;

function withSeed(base, i, fn) {
	rngOverride = splitmix64(iterStreamState(base, i));
	try {
		return fn();
	} finally {
		rngOverride = null;
	}
}

// ---------------------------------------------------------------------------------------------
// loading a Go wasm instance (same as the UI: new Go(); instantiate; go.run; call globals)

const API = [
	"sample",
	"validateConfig",
	"initializeWorker",
	"simulate",
	"initializeAggregator",
	"aggregate",
	"flush",
];

function resolveGlue(wasm, glue) {
	if (glue) return glue;
	const sibling = wasm.replace(/\.wasm$/, "") + ".wasm_exec.js";
	if (fs.existsSync(sibling)) return sibling;
	const root = execFileSync(process.env.GO || "go", ["env", "GOROOT"], { encoding: "utf8" }).trim();
	for (const c of [path.join(root, "lib/wasm/wasm_exec.js"), path.join(root, "misc/wasm/wasm_exec.js")]) {
		if (fs.existsSync(c)) return c;
	}
	throw new Error(`no wasm_exec.js found for ${wasm}; pass --glue`);
}

const fileCache = new Map();
function readCached(f, enc) {
	const k = `${f}\0${enc}`;
	if (!fileCache.has(k)) fileCache.set(k, fs.readFileSync(f, enc));
	return fileCache.get(k);
}

// Starts one Go instance. realm "new": a fresh vm context with only what a browser Worker
// offers the glue (console, performance, crypto, TextEncoder/Decoder, timers); realm "this":
// the current global (used inside worker_threads, which are separate realms already).
// Returns { api, mem(), toRealm(u8) } where api holds this instance's exported functions.
async function startGo(wasm, glue, realm = "new") {
	let g;
	if (realm === "new") {
		const ctx = vm.createContext({
			console,
			performance,
			TextEncoder,
			TextDecoder,
			setTimeout,
			clearTimeout,
			crypto: { getRandomValues: seededGetRandomValues },
		});
		// importScripts("/wasm_exec.js") equivalent: run the glue as a classic script.
		vm.runInContext(readCached(glue, "utf8"), ctx, { filename: glue });
		g = vm.runInContext("globalThis", ctx);
	} else {
		g = globalThis;
		if (typeof g.Go !== "function") vm.runInThisContext(readCached(glue, "utf8"), { filename: glue });
	}
	const go = new g.Go();
	const { instance } = await g.WebAssembly.instantiate(readCached(wasm), go.importObject);
	go.run(instance);
	const api = {};
	for (const name of API) {
		const fn = g[name];
		if (typeof fn !== "function") continue;
		api[name] = (...args) => {
			const r = fn(...args);
			if (go.exited) throw new Error(`wasm instance exited during ${name}() (Go panic or os.Exit)`);
			return r;
		};
	}
	const U8 = g.Uint8Array;
	return {
		api,
		mem: () => instance.exports.mem.buffer.byteLength,
		toRealm: (u8) => new U8(u8), // what structured clone does between workers
	};
}

function unwrapErr(resp, what) {
	if (typeof resp === "string" || resp instanceof String) {
		let msg = resp;
		try {
			msg = JSON.parse(resp).error ?? resp;
		} catch {}
		throw new Error(`${what} failed: ${msg}`);
	}
	return resp;
}

// ---------------------------------------------------------------------------------------------
// arg parsing

function parseArgs(argv) {
	const out = { _: [] };
	for (let i = 0; i < argv.length; i++) {
		const a = argv[i];
		if (a.startsWith("--")) {
			let [k, v] = a.slice(2).split("=", 2);
			if (v === undefined) v = i + 1 < argv.length && !argv[i + 1].startsWith("--") ? argv[++i] : true;
			if (k in out) out[k] = [].concat(out[k], v);
			else out[k] = v;
		} else out._.push(a);
	}
	return out;
}

const num = (v, d) => (v === undefined ? d : Number(v));

// ---------------------------------------------------------------------------------------------
// stats helpers

function quantile(sorted, q) {
	if (sorted.length === 0) return Number.NaN;
	const pos = (sorted.length - 1) * q;
	const lo = Math.floor(pos);
	const hi = Math.ceil(pos);
	return sorted[lo] + (sorted[hi] - sorted[lo]) * (pos - lo);
}
const sortNum = (a) => [...a].sort((x, y) => x - y);
const median = (a) => quantile(sortNum(a), 0.5);
const sum = (a) => a.reduce((s, x) => s + x, 0);

// ---------------------------------------------------------------------------------------------
// run: single-thread mode

async function loadPair(wasm, glue, cfg) {
	const sim = await startGo(wasm, glue);
	const agg = await startGo(wasm, glue);
	unwrapErr(sim.api.initializeWorker(cfg) ?? {}, "initializeWorker");
	const meta = JSON.parse(agg.api.initializeAggregator(cfg));
	if (meta.error) throw new Error(`initializeAggregator failed: ${meta.error}`);
	return { sim, agg, meta };
}

// Runs seeds [from, from+n) back to back on the sim instance, then aggregates the batch.
const threadCpuMs =
	typeof process.threadCpuUsage === "function"
		? () => {
				const u = process.threadCpuUsage();
				return (u.user + u.system) / 1000;
			}
		: () => Number.NaN;

function runBlock(p, seed, from, n) {
	const simMs = new Array(n);
	const results = new Array(n);
	const cpu0 = threadCpuMs();
	for (let k = 0; k < n; k++) {
		const i = from + k;
		const a = performance.now();
		results[k] = withSeed(seed, i, () => p.sim.api.simulate());
		simMs[k] = performance.now() - a;
		unwrapErr(results[k], "simulate");
	}
	const cpuMs = threadCpuMs() - cpu0;
	let bytes = 0;
	let aggMs = 0;
	for (const r of results) {
		const copy = p.agg.toRealm(r);
		bytes += r.length;
		const a = performance.now();
		const err = p.agg.api.aggregate(copy);
		aggMs += performance.now() - a;
		if (err != null) unwrapErr(err, "aggregate");
	}
	return { simMs, aggMs, bytes, cpuMs };
}

async function runSingle(o) {
	const t0 = performance.now();
	const p = await loadPair(o.wasm, o.glue, o.cfg);
	const tInit = performance.now();
	if (o.warmup > 0) runBlock(p, o.seed, 0, o.warmup);
	const tWarm = performance.now();
	const { simMs, aggMs, bytes, cpuMs } = runBlock(p, o.seed, o.warmup, o.iters);
	const tRun = performance.now();

	const f0 = performance.now();
	const flushed = JSON.parse(p.agg.api.flush());
	const flushMs = performance.now() - f0;
	if (flushed.error) throw new Error(`flush failed: ${flushed.error}`);

	const s = sortNum(simMs);
	return {
		load_init_ms: tInit - t0,
		warmup_ms: tWarm - tInit,
		run_ms: tRun - tWarm,
		ms_per_iter: sum(simMs) / o.iters,
		cpu_ms_per_iter: cpuMs / o.iters,
		sim_p50_ms: quantile(s, 0.5),
		sim_p10_ms: quantile(s, 0.1),
		sim_p90_ms: quantile(s, 0.9),
		sim_max_ms: s[s.length - 1],
		agg_ms_per_iter: aggMs / o.iters,
		flush_ms: flushMs,
		result_bytes_mean: bytes / o.iters,
		sim_mem_mb: p.sim.mem() / 2 ** 20,
		dps_mean: flushed.stats?.dps?.mean,
	};
}

// ---------------------------------------------------------------------------------------------
// ab: paired, interleaved comparison of several binaries inside one process

async function runAB(o) {
	const pairs = [];
	for (const b of o.bins) pairs.push({ ...b, ...(await loadPair(b.wasm, b.glue, o.cfg)) });
	const nb = pairs.length;
	// warm up in small interleaved chunks so every binary tiers up under similar conditions
	const chunk = Math.max(1, Math.ceil(o.warmup / 4));
	for (let done = 0; done < o.warmup; done += chunk) {
		for (const p of pairs) runBlock(p, o.seed, done, Math.min(chunk, o.warmup - done));
	}
	const out = [];
	for (let r = 0; r < o.rounds; r++) {
		for (let j = 0; j < nb; j++) {
			const p = pairs[(j + r) % nb];
			const { simMs, aggMs, bytes, cpuMs } = runBlock(p, o.seed, o.warmup, o.block);
			out.push({
				label: p.label,
				wasm: p.wasm,
				round: o.proc * 1000 + r,
				proc: o.proc,
				order: j,
				ms_per_iter: sum(simMs) / o.block,
				cpu_ms_per_iter: cpuMs / o.block,
				agg_ms_per_iter: aggMs / o.block,
				result_bytes_mean: bytes / o.block,
			});
		}
	}
	return out;
}

// ---------------------------------------------------------------------------------------------
// run: worker pool mode (mirrors WasmExecutor.run)

const VIEWER_THROTTLE = 100;

function throttle(fn, wait) {
	// lodash.throttle(fn, wait, { leading: true, trailing: true })
	let last = Number.NEGATIVE_INFINITY;
	let timer = null;
	return () => {
		const now = performance.now();
		const remaining = wait - (now - last);
		if (remaining <= 0) {
			if (timer) {
				clearTimeout(timer);
				timer = null;
			}
			last = now;
			fn();
		} else if (!timer) {
			timer = setTimeout(() => {
				last = performance.now();
				timer = null;
				fn();
			}, remaining);
		}
	};
}

function spawn(role, o) {
	return new Worker(new URL(import.meta.url), {
		workerData: { role, wasm: o.wasm, glue: o.glue, seed: o.seed },
	});
}

function once(w, want) {
	return new Promise((resolve, reject) => {
		const h = (m) => {
			w.off("message", h);
			if (m.type === "failed") reject(new Error(m.reason));
			else if (m.type !== want) reject(new Error(`expected ${want}, got ${m.type}`));
			else resolve(m);
		};
		w.on("message", h);
		w.once("error", reject);
	});
}

async function runWorkers(o) {
	const t0 = performance.now();
	const aggregator = spawn("agg", o);
	const workers = Array.from({ length: o.workers }, () => spawn("sim", o));
	await Promise.all([aggregator, ...workers].map((w) => once(w, "ready")));
	const tReady = performance.now();

	const initAgg = once(aggregator, "initialized");
	aggregator.postMessage({ type: "initialize", cfg: o.cfg });
	const initW = workers.map((w) => {
		const p = once(w, "initialized");
		w.postMessage({ type: "initialize", cfg: o.cfg });
		return p;
	});
	await Promise.all([initAgg, ...initW]);
	const tInit = performance.now();

	const maxIterations = o.iters;
	let tDone = 0;
	let flushes = 0;
	let lastStats = null;
	const tResult = await new Promise((resolve, reject) => {
		let completed = 0;
		let requested = 0;
		let running = true;
		const throttledFlush = throttle(() => {
			if (running) aggregator.postMessage({ type: "flush" });
		}, VIEWER_THROTTLE);

		aggregator.on("message", (m) => {
			switch (m.type) {
				case "result": {
					flushes++;
					lastStats = m.result.stats;
					if (completed >= maxIterations) {
						running = false;
						resolve(performance.now());
					}
					return;
				}
				case "done":
					completed += 1;
					if (completed === maxIterations) tDone = performance.now();
					throttledFlush();
					return;
				case "failed":
					if (running) reject(new Error(m.reason));
			}
		});
		for (const w of workers) {
			w.on("message", (m) => {
				switch (m.type) {
					case "done":
						aggregator.postMessage({ type: "add", result: m.result });
						if (requested < maxIterations) w.postMessage({ type: "run", itr: requested++ });
						return;
					case "failed":
						reject(new Error(m.reason));
				}
			});
			if (requested < maxIterations) w.postMessage({ type: "run", itr: requested++ });
		}
	});

	await Promise.all([aggregator, ...workers].map((w) => w.terminate()));
	return {
		ready_ms: tReady - t0,
		init_ms: tInit - tReady,
		wall_ms: tDone - tInit,
		ui_wall_ms: tResult - tInit,
		total_ms: tResult - t0,
		ms_per_iter: (tDone - tInit) / maxIterations,
		iters_per_sec: maxIterations / ((tDone - tInit) / 1000),
		flushes,
		dps_mean: lastStats?.dps?.mean,
	};
}

// worker_threads side: the UI's Workers/worker.ts and Workers/aggregator.ts
async function workerMain() {
	const { role, wasm, glue, seed } = workerData;
	let go;
	try {
		go = await startGo(wasm, glue, "this");
	} catch (e) {
		parentPort.postMessage({ type: "failed", reason: String(e?.message ?? e) });
		return;
	}
	let meta = null;
	parentPort.postMessage({ type: "ready" });
	parentPort.on("message", (req) => {
		if (role === "sim") {
			switch (req.type) {
				case "initialize": {
					const resp = go.api.initializeWorker(req.cfg);
					if (resp != null) return parentPort.postMessage({ type: "failed", reason: JSON.parse(resp).error });
					return parentPort.postMessage({ type: "initialized" });
				}
				case "run": {
					const resp = withSeed(seed, req.itr, () => go.api.simulate());
					if (typeof resp === "string") return parentPort.postMessage({ type: "failed", reason: JSON.parse(resp).error });
					return parentPort.postMessage({ type: "done", result: resp, itr: req.itr });
				}
			}
		} else {
			switch (req.type) {
				case "initialize": {
					meta = JSON.parse(go.api.initializeAggregator(req.cfg));
					if (meta.error) return parentPort.postMessage({ type: "failed", reason: meta.error });
					return parentPort.postMessage({ type: "initialized" });
				}
				case "add": {
					const resp = go.api.aggregate(req.result);
					if (resp != null) return parentPort.postMessage({ type: "failed", reason: JSON.parse(resp).error });
					return parentPort.postMessage({ type: "done" });
				}
				case "flush": {
					const resp = JSON.parse(go.api.flush());
					if (resp.error) return parentPort.postMessage({ type: "failed", reason: resp.error });
					return parentPort.postMessage({ type: "result", result: resp });
				}
			}
		}
	});
}

// ---------------------------------------------------------------------------------------------
// canonical forms for the correctness oracle

// Minimal MessagePack decoder (tinylib/msgp output). 64-bit ints outside the safe range stay
// BigInt; ext types become {ext, hex}.
function mpDecode(buf) {
	const dv = new DataView(buf.buffer, buf.byteOffset, buf.byteLength);
	const td = new TextDecoder();
	let p = 0;
	const big = (v) => (v >= BigInt(Number.MIN_SAFE_INTEGER) && v <= BigInt(Number.MAX_SAFE_INTEGER) ? Number(v) : v);
	const str = (n) => {
		const s = td.decode(buf.subarray(p, p + n));
		p += n;
		return s;
	};
	const bin = (n) => {
		const s = Buffer.from(buf.subarray(p, p + n)).toString("hex");
		p += n;
		return { bin: s };
	};
	const ext = (n) => {
		const type = dv.getInt8(p++);
		return { ext: type, ...bin(n) };
	};
	const arr = (n) => {
		const a = new Array(n);
		for (let i = 0; i < n; i++) a[i] = read();
		return a;
	};
	const map = (n) => {
		const o = {};
		for (let i = 0; i < n; i++) {
			const k = read();
			o[String(k)] = read();
		}
		return o;
	};
	const u8 = () => dv.getUint8(p++);
	const u16 = () => ((p += 2), dv.getUint16(p - 2));
	const u32 = () => ((p += 4), dv.getUint32(p - 4));
	function read() {
		const b = buf[p++];
		if (b <= 0x7f) return b;
		if (b >= 0xe0) return b - 0x100;
		if ((b & 0xf0) === 0x80) return map(b & 0x0f);
		if ((b & 0xf0) === 0x90) return arr(b & 0x0f);
		if ((b & 0xe0) === 0xa0) return str(b & 0x1f);
		switch (b) {
			case 0xc0:
				return null;
			case 0xc2:
				return false;
			case 0xc3:
				return true;
			case 0xc4:
				return bin(u8());
			case 0xc5:
				return bin(u16());
			case 0xc6:
				return bin(u32());
			case 0xc7:
				return ext(u8());
			case 0xc8:
				return ext(u16());
			case 0xc9:
				return ext(u32());
			case 0xca:
				return (p += 4), dv.getFloat32(p - 4);
			case 0xcb:
				return (p += 8), dv.getFloat64(p - 8);
			case 0xcc:
				return u8();
			case 0xcd:
				return u16();
			case 0xce:
				return u32();
			case 0xcf:
				return (p += 8), big(dv.getBigUint64(p - 8));
			case 0xd0:
				return (p += 1), dv.getInt8(p - 1);
			case 0xd1:
				return (p += 2), dv.getInt16(p - 2);
			case 0xd2:
				return (p += 4), dv.getInt32(p - 4);
			case 0xd3:
				return (p += 8), big(dv.getBigInt64(p - 8));
			case 0xd4:
				return ext(1);
			case 0xd5:
				return ext(2);
			case 0xd6:
				return ext(4);
			case 0xd7:
				return ext(8);
			case 0xd8:
				return ext(16);
			case 0xd9:
				return str(u8());
			case 0xda:
				return str(u16());
			case 0xdb:
				return str(u32());
			case 0xdc:
				return arr(u16());
			case 0xdd:
				return arr(u32());
			case 0xde:
				return map(u16());
			case 0xdf:
				return map(u32());
		}
		throw new Error(`msgpack: bad byte 0x${b.toString(16)} at ${p - 1}`);
	}
	const v = read();
	if (p !== buf.length) throw new Error(`msgpack: ${buf.length - p} trailing bytes`);
	return v;
}

// Canonical JSON: sorted object keys, BigInt as digits, -0 as 0, and empty arrays/objects
// folded into null (Go nil vs empty slice/map is not a behavioural difference).
function canon(v) {
	if (v === null || v === undefined) return "null";
	if (Array.isArray(v)) return v.length === 0 ? "null" : `[${v.map(canon).join(",")}]`;
	switch (typeof v) {
		case "object": {
			const ks = Object.keys(v).sort();
			if (ks.length === 0) return "null";
			return `{${ks.map((k) => `${JSON.stringify(k)}:${canon(v[k])}`).join(",")}}`;
		}
		case "number":
			if (Object.is(v, -0)) return "0";
			return Number.isFinite(v) ? JSON.stringify(v) : JSON.stringify(String(v));
		case "bigint":
			return v.toString();
		default:
			return JSON.stringify(v);
	}
}
const sha = (s) => createHash("sha256").update(s).digest("hex");

// Metadata fields that legitimately differ between builds or runs.
const VOLATILE_META = ["sim_version", "modified", "build_date", "sample_seed"];

// simulate() returns either a stats.Result (msgpack map with field names) or, since the
// per-iteration reduction moved into the sim worker, an agg.Summary (msgpack tuple, seed first).
function payloadKind(decoded) {
	return Array.isArray(decoded) ? "summary" : "result";
}
function payloadSeed(decoded) {
	return Array.isArray(decoded) ? decoded[0] : decoded.seed;
}

// sample() returns the debug log in emission order. Within a frame that order can follow Go map
// iteration (SetupResonance and artifact set setup range over maps), so it differs between
// processes running the same binary. This form sorts each frame's events: any change in what is
// logged, or in which frame, still shows, but the order within a frame does not.
function sampleByFrame(sample) {
	const logs = sample.logs ?? [];
	const frames = [];
	for (let i = 0; i < logs.length; ) {
		let j = i;
		while (j < logs.length && logs[j].frame === logs[i].frame) j++;
		frames.push(logs.slice(i, j).map(canon).sort());
		i = j;
	}
	return { ...sample, logs: frames };
}

function flushStats(agg) {
	const flushed = JSON.parse(agg.api.flush());
	if (flushed.error) throw new Error(`flush failed: ${flushed.error}`);
	return JSON.parse(canon(flushed.stats));
}

async function dump(o) {
	const { sim, agg, meta } = await loadPair(o.wasm, o.glue, o.cfg);
	for (const k of VOLATILE_META) delete meta[k];

	const iterHashes = [];
	const flushHashes = [];
	const seedErrors = [];
	let payload;
	for (let i = 0; i < o.iters; i++) {
		const res = unwrapErr(withSeed(o.seed, i, () => sim.api.simulate()), "simulate");
		const decoded = mpDecode(res);
		payload ??= payloadKind(decoded);
		const want = expectedSeed(o.seed, i);
		const got = payloadSeed(decoded);
		if (got === undefined || BigInt(got) !== want) seedErrors.push({ i, want: want.toString(), got: String(got) });
		iterHashes.push(sha(canon(decoded)).slice(0, 16));
		const err = agg.api.aggregate(agg.toRealm(res));
		if (err != null) unwrapErr(err, "aggregate");
		if (o.flushEvery > 0 && (i + 1) % o.flushEvery === 0 && i + 1 < o.iters) {
			flushHashes.push(sha(canon(flushStats(agg))).slice(0, 16));
		}
	}
	const stats = flushStats(agg);

	const validated = JSON.parse(sim.api.validateConfig(o.cfg));
	const sampleSeed = expectedSeed(o.seed, 0).toString();
	const sample = JSON.parse(sim.api.sample(o.cfg, sampleSeed));
	if (sample.error) throw new Error(`sample failed: ${sample.error}`);
	for (const k of VOLATILE_META) delete sample[k];

	return {
		config: o.configName,
		iters: o.iters,
		seed: o.seed,
		payload,
		flushEvery: o.flushEvery,
		seedErrors,
		flushHashes,
		hashes: {
			iterations: sha(iterHashes.join(",")),
			stats: sha(canon(stats)),
			meta: sha(canon(meta)),
			validate: sha(canon(validated)),
			sample: sha(canon(sample)),
			sampleByFrame: sha(canon(sampleByFrame(sample))),
		},
		iterHashes,
		stats,
	};
}

function* walkDiff(a, b, p = "") {
	if (typeof a === "number" && typeof b === "number") {
		if (a !== b) yield { path: p, golden: a, candidate: b, rel: Math.abs(a - b) / Math.max(Math.abs(a), Math.abs(b)) };
		return;
	}
	if (a && b && typeof a === "object" && typeof b === "object") {
		for (const k of new Set([...Object.keys(a), ...Object.keys(b)])) yield* walkDiff(a[k], b[k], `${p}.${k}`);
		return;
	}
	if (JSON.stringify(a) !== JSON.stringify(b)) yield { path: p, golden: a, candidate: b };
}

function diff(golden, cand, rtol) {
	const problems = [];
	const notes = [];
	const gFlush = golden.flushEvery ?? 0;
	const cFlush = cand.flushEvery ?? 0;
	if (golden.iters !== cand.iters || golden.seed !== cand.seed || gFlush !== cFlush) {
		problems.push(
			`run parameters differ: golden iters=${golden.iters} seed=${golden.seed} flushEvery=${gFlush}, candidate iters=${cand.iters} seed=${cand.seed} flushEvery=${cFlush}`,
		);
	}
	if (cand.seedErrors?.length) {
		problems.push(`seed injection failed for ${cand.seedErrors.length} iterations (first: ${JSON.stringify(cand.seedErrors[0])})`);
	}
	// Goldens made before the payload field existed hold stats.Result payloads. Payloads of
	// different formats can't be compared, so the aggregated stats are the check.
	const samePayload = (golden.payload ?? "result") === (cand.payload ?? "result");
	if (!samePayload) {
		notes.push(`per-iteration payloads differ in format (${golden.payload ?? "result"} vs ${cand.payload}); compared aggregated stats only`);
	}
	// Goldens made before sampleByFrame existed compare the sample log in exact order.
	const sampleByFrameBoth = golden.hashes.sampleByFrame !== undefined && cand.hashes.sampleByFrame !== undefined;
	const fg = golden.flushHashes ?? [];
	const fc = cand.flushHashes ?? [];
	const firstFlush = fg.findIndex((h, i) => h !== fc[i]);
	if (fg.length !== fc.length || firstFlush >= 0) {
		const at = firstFlush >= 0 ? firstFlush : Math.min(fg.length, fc.length);
		problems.push(`intermediate flushes differ, first after iteration ${(at + 1) * gFlush}`);
	}
	for (const k of Object.keys(golden.hashes)) {
		if (golden.hashes[k] === cand.hashes[k]) continue;
		if (k === "iterations" && !samePayload) continue;
		if (k === "sample" && sampleByFrameBoth) {
			notes.push("sample() log events come in a different order within some frames (Go map iteration); compared per frame");
			continue;
		}
		if (k === "iterations") {
			const idx = golden.iterHashes.findIndex((h, i) => h !== cand.iterHashes[i]);
			const n = golden.iterHashes.filter((h, i) => h !== cand.iterHashes[i]).length;
			problems.push(`per-iteration results differ in ${n}/${golden.iterHashes.length} iterations (first: #${idx})`);
		} else if (k === "stats") {
			const d = [...walkDiff(golden.stats, cand.stats)];
			const bad = d.filter((x) => x.rel === undefined || x.rel > rtol);
			if (bad.length === 0) continue;
			problems.push(`aggregated stats differ at ${bad.length} paths${rtol ? ` (rtol ${rtol})` : ""}:`);
			for (const x of bad.slice(0, 15)) problems.push(`    ${x.path}: golden=${JSON.stringify(x.golden)} candidate=${JSON.stringify(x.candidate)}${x.rel !== undefined ? ` rel=${x.rel.toExponential(2)}` : ""}`);
		} else {
			problems.push(`${k} output differs`);
		}
	}
	return { problems, notes };
}

// ---------------------------------------------------------------------------------------------
// report: comparison of jsonl results

function mulberry32(a) {
	return () => {
		a |= 0;
		a = (a + 0x6d2b79f5) | 0;
		let t = Math.imul(a ^ (a >>> 15), 1 | a);
		t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
		return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
	};
}

// Ratios are medians of per-round paired ratios (binary X vs the first binary in the same
// round, run back to back), which cancels load that hits both. CI: bootstrap over rounds.
function report(files, metricArg) {
	const rows = files.flatMap((f) =>
		fs
			.readFileSync(f, "utf8")
			.split("\n")
			.filter(Boolean)
			.map((l) => JSON.parse(l)),
	);
	if (rows.length === 0) throw new Error("no results");
	const metric = metricArg ?? (rows.every((r) => Number.isFinite(r.cpu_ms_per_iter)) ? "cpu_ms_per_iter" : "ms_per_iter");
	const other = metric === "ms_per_iter" ? "cpu_ms_per_iter" : "ms_per_iter";
	const hasOther = rows.every((r) => Number.isFinite(r[other]));
	const labels = [...new Set(rows.map((r) => r.label))];
	const configs = [...new Set(rows.map((r) => r.config))];
	const r0 = rows[0];
	const paired = rows.some((r) => r.proc !== undefined);
	const rand = mulberry32(12345);
	const metricName = { ms_per_iter: "wall", cpu_ms_per_iter: "thread-CPU" };

	// config -> label -> Map(round -> row)
	const byRound = (cfg, l) => {
		const m = new Map();
		for (const r of rows) if (r.config === cfg && r.label === l) m.set(`${r.proc ?? ""}:${r.round ?? m.size}`, r);
		return m;
	};
	const pairedRatios = (a, b, key) => [...b.keys()].filter((k) => a.has(k)).map((k) => b.get(k)[key] / a.get(k)[key]);
	const medRatio = (rs) => Math.exp(median(rs.map(Math.log)));
	const bootCI = (rs) => {
		if (rs.length < 4) return "";
		const logs = rs.map(Math.log);
		const bs = [];
		for (let b = 0; b < 2000; b++) {
			const xs = logs.map(() => logs[Math.floor(rand() * logs.length)]);
			bs.push(Math.exp(median(xs)));
		}
		const s = sortNum(bs);
		return `[${quantile(s, 0.025).toFixed(3)}, ${quantile(s, 0.975).toFixed(3)}]`;
	};

	const out = [];
	const mode = paired ? "paired in-process" : r0.workers ? `worker pool (${r0.workers} workers + aggregator, UI dispatch)` : "single worker, process per run";
	out.push(`mode: ${mode}; seed ${r0.seed}; ${paired ? `block ${r0.iters} iters/round` : `${r0.iters} iters/run`}, warmup ${r0.warmup}`);
	out.push(`metric: ${r0.workers ? "wall ms / iterations" : `simulate() ${metricName[metric]} ms/iter`}; ratio = median of per-round paired ratios vs ${labels[0]}`);
	out.push("");
	const extraHead = `${hasOther && !r0.workers ? ` ${metricName[other]} ratio |` : ""}${paired ? " per-process ratios |" : ""}`;
	const extraSep = `${hasOther && !r0.workers ? "---:|" : ""}${paired ? "---|" : ""}`;
	out.push(`| config | binary | samples | median ms/iter | IQR/median | ratio | 95% CI |${extraHead}`);
	out.push(`|---|---|---:|---:|---:|---:|---|${extraSep}`);
	const logRatios = Object.fromEntries(labels.slice(1).map((l) => [l, []]));
	for (const cfg of configs) {
		const base = byRound(cfg, labels[0]);
		for (const l of labels) {
			const mine = byRound(cfg, l);
			const vals = [...mine.values()].map((r) => r[metric]);
			if (vals.length === 0) continue;
			const s = sortNum(vals);
			const med = quantile(s, 0.5);
			const iqr = (quantile(s, 0.75) - quantile(s, 0.25)) / med;
			let cols = "1.000 | |";
			let extra = "";
			if (l !== labels[0]) {
				const rs = pairedRatios(base, mine, metric);
				const r = medRatio(rs);
				logRatios[l].push(Math.log(r));
				cols = `${r.toFixed(3)} | ${bootCI(rs)} |`;
				if (hasOther && !r0.workers) extra += ` ${medRatio(pairedRatios(base, mine, other)).toFixed(3)} |`;
				if (paired) {
					const procs = [...new Set([...mine.values()].map((x) => x.proc))];
					const pr = procs.map((pp) => {
						const f = (m) => new Map([...m].filter(([, x]) => x.proc === pp));
						return medRatio(pairedRatios(f(base), f(mine), metric)).toFixed(3);
					});
					extra += ` ${pr.join(" ")} |`;
				}
			} else {
				if (hasOther && !r0.workers) extra += " |";
				if (paired) extra += " |";
			}
			out.push(`| ${cfg} | ${l} | ${vals.length} | ${med.toFixed(3)} | ${(iqr * 100).toFixed(1)}% | ${cols}${extra}`);
		}
	}
	if (configs.length > 1 && labels.length > 1) {
		out.push("");
		for (const l of labels.slice(1)) {
			const g = Math.exp(sum(logRatios[l]) / logRatios[l].length);
			out.push(`geomean ratio ${l} / ${labels[0]} over ${logRatios[l].length} configs: ${g.toFixed(3)}`);
		}
	}
	return out.join("\n");
}

// ---------------------------------------------------------------------------------------------

function loadRunOpts(args) {
	if (!args.wasm || !args.config) throw new Error("--wasm and --config are required");
	if (Array.isArray(args.wasm)) throw new Error("only one --wasm allowed here (use ab)");
	const wasm = path.resolve(args.wasm);
	return {
		wasm,
		glue: resolveGlue(wasm, args.glue && path.resolve(args.glue)),
		cfg: fs.readFileSync(args.config, "utf8"),
		configName: path.basename(args.config).replace(/\.[^.]+$/, ""),
		iters: num(args.iters, undefined),
		warmup: num(args.warmup, 20),
		seed: num(args.seed, 1),
		workers: num(args.workers, 0),
		flushEvery: num(args["flush-every"], 0),
	};
}

async function main() {
	const args = parseArgs(process.argv.slice(2));
	const cmd = args._[0];
	switch (cmd) {
		case "run": {
			const o = loadRunOpts(args);
			o.iters ??= o.workers ? 1000 : 100;
			const res = o.workers ? await runWorkers(o) : await runSingle(o);
			const rec = {
				label: args.label ?? path.basename(o.wasm),
				round: args.round !== undefined ? Number(args.round) : undefined,
				config: o.configName,
				wasm: o.wasm,
				wasm_bytes: fs.statSync(o.wasm).size,
				iters: o.iters,
				warmup: o.workers ? 0 : o.warmup,
				seed: o.seed,
				workers: o.workers,
				node: process.version,
				ts: new Date().toISOString(),
				...res,
			};
			if (args.json) console.log(JSON.stringify(rec));
			else {
				for (const [k, v] of Object.entries(rec)) console.log(`${k.padEnd(18)} ${typeof v === "number" && !Number.isInteger(v) ? v.toFixed(3) : v}`);
			}
			return;
		}
		case "ab": {
			if (!args.wasm || !args.config) throw new Error("--wasm (repeatable) and --config are required");
			const wasms = [].concat(args.wasm).map((w) => path.resolve(w));
			const labels = args.label ? [].concat(args.label) : wasms.map((w) => path.basename(w, ".wasm"));
			const glues = args.glue ? [].concat(args.glue).map((g) => path.resolve(g)) : [];
			if (labels.length !== wasms.length) throw new Error("need one --label per --wasm");
			const o = {
				bins: wasms.map((w, i) => ({ wasm: w, label: labels[i], glue: resolveGlue(w, glues[i]) })),
				cfg: fs.readFileSync(args.config, "utf8"),
				configName: path.basename(args.config).replace(/\.[^.]+$/, ""),
				block: num(args.block, 20),
				warmup: num(args.warmup, 40),
				rounds: num(args.rounds, 10),
				seed: num(args.seed, 1),
				proc: num(args.proc, 0),
			};
			const recs = await runAB(o);
			for (const r of recs) {
				const rec = { mode: "paired", config: o.configName, iters: o.block, warmup: o.warmup, seed: o.seed, workers: 0, node: process.version, ...r };
				if (args.json) console.log(JSON.stringify(rec));
				else console.log(`${rec.config} ${rec.label.padEnd(16)} round ${String(rec.round).padStart(5)} ${rec.ms_per_iter.toFixed(3)} ms/iter`);
			}
			return;
		}
		case "dump": {
			const o = loadRunOpts(args);
			o.iters ??= 100;
			const d = await dump(o);
			const s = `${JSON.stringify(d)}\n`;
			if (args.out) fs.writeFileSync(args.out, s);
			else process.stdout.write(s);
			return;
		}
		case "diff": {
			const [g, c] = args._.slice(1);
			const golden = JSON.parse(fs.readFileSync(g, "utf8"));
			const cand = JSON.parse(fs.readFileSync(c, "utf8"));
			const { problems, notes } = diff(golden, cand, num(args.rtol, 0));
			const flushes = golden.flushEvery ? `, flush every ${golden.flushEvery}` : "";
			if (problems.length) {
				console.log(`FAIL ${golden.config}`);
				for (const p of problems) console.log(`  ${p}`);
				process.exitCode = 1;
			} else console.log(`ok   ${golden.config} (${golden.iters} iterations, seed ${golden.seed}${flushes})`);
			for (const n of notes) console.log(`  note: ${n}`);
			return;
		}
		case "report":
			console.log(report(args._.slice(1), args.metric));
			return;
		default:
			console.error("usage: bench.mjs run|dump|diff|report ... (see header comment)");
			process.exitCode = 2;
	}
}

if (isMainThread) {
	main().catch((e) => {
		console.error(e?.stack ?? String(e));
		process.exit(1);
	});
} else {
	workerMain();
}
