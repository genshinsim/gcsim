#!/usr/bin/env node
// Corpus tools: run the wasmbench correctness oracle over a directory of configs (such as the
// gcsim DB dump that corpus-extract.sh unpacks), compare two runs, pick stratified samples, and
// build config sets for compare.sh and corpus-pgo.sh. No npm dependencies.
//
// Subcommands:
//   dump     --wasm F --configs DIR --out DIR [--ids FILE] [--iters 20] [--seed 1] [--flush-every 1]
//            [--sample] [--jobs 8] [--timeout 300] [--share-module]
//            Runs bench.mjs's dump on every DIR/<id>.txt (or on the ids listed in FILE) in JOBS
//            worker threads, each with its own compiled module (--share-module: one for all);
//            every config gets fresh Go instances. A crashed pool process is restarted.
//            Writes OUT/<id>.json, or OUT/<id>.err (JSON with the error, Go's output tail and
//            whether it timed out). Configs that already have an output are skipped, so an
//            interrupted run resumes. The defaults (20 iterations, a flush after each) compare
//            every iteration's contribution even across payload formats. sample() is skipped
//            unless --sample, which also writes its per-frame log to OUT/<id>.sample.json.
//   diff     GOLDEN_DIR CAND_DIR [--ids FILE] [--rtol X]
//            Compares the dumps of each id with bench.mjs's oracle. Prints one line per id,
//            "<id>\t<ok|FAIL|golden-error|cand-error|missing>\t<detail>", then counts.
//   errors   DIR [--ids-out FILE]
//            Groups the failures of a dump run by reason. --ids-out writes the ids that ran.
//   select   --entries entries.tsv --ids FILE --n N [--cover] [--exclude FILE] [--seed S]
//            Picks N ids from FILE (one per line). --cover first runs a greedy set cover so every
//            character in the pool appears at least once. The remaining picks are random,
//            stratified by target count, team size and sim duration in proportion to the pool,
//            so frequent characters stay frequent.
//   manifest --timings DUMP_DIR --configs DIR --ids FILE --out DIR [--block-ms 150]
//            Writes OUT/manifest.tsv plus OUT/<id>.txt, for compare.sh and check.sh
//            (WASMBENCH_MANIFEST=OUT/manifest.tsv). Block and warmup sizes come from the cold
//            ms/iter the dumps in DUMP_DIR recorded.
//   speed    results.jsonl... [--tsv FILE]
//            Per-config paired ratios vs the first label (median of per-round ratios of
//            cpu_ms_per_iter) and peak wasm memory, from compare.sh's raw output, then their
//            distribution.

import { spawn } from "node:child_process";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { isMainThread, parentPort, Worker, workerData } from "node:worker_threads";
import { diff, dump, quantile, resolveGlue, sortNum } from "./bench.mjs";

function parseArgs(argv) {
	const out = { _: [] };
	for (let i = 0; i < argv.length; i++) {
		const a = argv[i];
		if (a.startsWith("--")) {
			let [k, v] = a.slice(2).split("=", 2);
			if (v === undefined) v = i + 1 < argv.length && !argv[i + 1].startsWith("--") ? argv[++i] : true;
			out[k] = v;
		} else out._.push(a);
	}
	return out;
}

const readLines = (f) =>
	fs
		.readFileSync(f, "utf8")
		.split("\n")
		.map((l) => l.trim())
		.filter((l) => l && !l.startsWith("#"));

function listIds(configs, idsFile) {
	if (idsFile) return readLines(idsFile);
	return fs
		.readdirSync(configs)
		.filter((f) => f.endsWith(".txt"))
		.map((f) => f.slice(0, -4))
		.sort();
}

function mulberry32(a) {
	return () => {
		a |= 0;
		a = (a + 0x6d2b79f5) | 0;
		let t = Math.imul(a ^ (a >>> 15), 1 | a);
		t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
		return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
	};
}

function shuffle(a, rand) {
	for (let i = a.length - 1; i > 0; i--) {
		const j = Math.floor(rand() * (i + 1));
		[a[i], a[j]] = [a[j], a[i]];
	}
	return a;
}

// ---------------------------------------------------------------------------------------------
// dump

// Each worker thread runs bench.mjs's dump() on one config at a time, with fresh Go instances
// (own vm realms) per config. The worker compiles the wasm once (or gets the pool's shared
// module), which saves a process start and a compile per config. Go's stdout/stderr (a panic's
// message) is kept for the error record.
function dumpWorker() {
	const { wasm, glue } = workerData;
	const module = workerData.module ?? new WebAssembly.Module(fs.readFileSync(wasm));
	parentPort.on("message", async (t) => {
		const out = [];
		const keep = (...a) => {
			out.push(a.map(String).join(" "));
			if (out.length > 400) out.splice(0, out.length - 200);
		};
		const logger = { log: keep, error: keep, warn: keep, info: keep, debug: keep };
		try {
			const cfg = fs.readFileSync(t.config, "utf8");
			const d = await dump({ ...t.opts, wasm, glue, module, logger, cfg, configName: t.id });
			fs.writeFileSync(t.tmp, `${JSON.stringify(d)}\n`);
			parentPort.postMessage({ ok: true });
		} catch (e) {
			parentPort.postMessage({ ok: false, error: String(e?.message ?? e), output: out.join("\n").slice(-4096) });
		}
	});
	parentPort.postMessage({ ready: true });
}

// A worker thread slot. A config that hangs past the timeout gets its worker terminated, and
// workers are replaced every RECYCLE configs so anything an instance leaks stays bounded.
const RECYCLE = 200;
class Slot {
	constructor(o) {
		this.o = o;
		this.w = null;
	}
	async start() {
		this.tasks = 0;
		this.w = new Worker(new URL(import.meta.url), {
			workerData: { role: "corpus-dump", module: this.o.module, wasm: this.o.wasm, glue: this.o.glue },
			stdout: true,
			stderr: true,
		});
		this.w.stdout.resume();
		this.w.stderr.resume();
		await new Promise((resolve, reject) => {
			this.w.once("message", resolve);
			this.w.once("error", reject);
		});
	}
	async stop() {
		if (this.w) await this.w.terminate();
		this.w = null;
	}
	async run(task) {
		if (!this.w || this.tasks >= RECYCLE) {
			await this.stop();
			await this.start();
		}
		this.tasks++;
		const w = this.w;
		const r = await new Promise((resolve) => {
			const done = (x) => {
				clearTimeout(timer);
				w.off("message", onMsg);
				w.off("error", onErr);
				w.off("exit", onExit);
				resolve(x);
			};
			const onMsg = (m) => done(m);
			const onErr = (e) => done({ ok: false, error: `worker error: ${e?.message ?? e}`, dead: true });
			const onExit = (c) => done({ ok: false, error: `worker exited (${c})`, dead: true });
			const timer = setTimeout(() => done({ ok: false, error: `timeout after ${this.o.timeout}s`, timedOut: true, dead: true }), this.o.timeout * 1000);
			w.on("message", onMsg);
			w.on("error", onErr);
			w.on("exit", onExit);
			w.postMessage(task);
		});
		if (r.dead) await this.stop();
		return r;
	}
}

// node itself has died with SIGSEGV now and then early in a pool run (twice in about 15 runs,
// both times within the first 20 configs). The pool therefore runs in a child process that is
// restarted to resume after a crash. A config that was in flight in two crashes gets an .err.
const MAX_RESTARTS = 20;
async function cmdDump(args) {
	if (args.child) return dumpPool(args);
	for (let restarts = 0; ; restarts++) {
		const argv = [...process.execArgv, fileURLToPath(import.meta.url), ...process.argv.slice(2), "--child"];
		const { code, signal } = await new Promise((resolve) => {
			spawn(process.execPath, argv, { stdio: "inherit" }).on("close", (c, s) => resolve({ code: c, signal: s }));
		});
		if (code === 0) return;
		if (code !== null || restarts >= MAX_RESTARTS) throw new Error(`dump pool exited with ${signal ?? code}`);
		console.error(`dump pool died (${signal}); resuming`);
	}
}

async function dumpPool(args) {
	const wasm = path.resolve(args.wasm);
	const o = {
		wasm,
		glue: resolveGlue(wasm, args.glue && path.resolve(args.glue)),
		configs: path.resolve(args.configs),
		out: path.resolve(args.out),
		jobs: Number(args.jobs ?? 8),
		timeout: Number(args.timeout ?? 300),
		// Sharing one compiled module across the worker threads saves compile time and memory;
		// it is off by default because the crashes above happened with it on.
		module: args["share-module"] ? new WebAssembly.Module(fs.readFileSync(wasm)) : undefined,
	};
	const opts = {
		iters: Number(args.iters ?? 20),
		seed: Number(args.seed ?? 1),
		flushEvery: Number(args["flush-every"] ?? 1),
		sample: Boolean(args.sample),
	};
	fs.mkdirSync(o.out, { recursive: true });
	const all = listIds(o.configs, args.ids);
	const running = (id) => path.join(o.out, `${id}.running`);
	for (const id of all) {
		if (fs.existsSync(running(id)) && Number(fs.readFileSync(running(id), "utf8")) >= 2) {
			const rec = { id, reason: "node crashed twice while running this config", timedOut: false, ms: 0, output: "" };
			fs.writeFileSync(path.join(o.out, `${id}.err`), `${JSON.stringify(rec)}\n`);
			fs.rmSync(running(id));
		}
	}
	const todo = all.filter((id) => !fs.existsSync(path.join(o.out, `${id}.json`)) && !fs.existsSync(path.join(o.out, `${id}.err`)));
	console.error(`${all.length} configs, ${all.length - todo.length} done already, ${todo.length} to run, ${o.jobs} at a time`);
	let next = 0;
	let done = 0;
	let failed = 0;
	const t0 = performance.now();
	const loop = async () => {
		const slot = new Slot(o);
		while (next < todo.length) {
			const id = todo[next++];
			const final = path.join(o.out, `${id}.json`);
			const task = {
				id,
				config: path.join(o.configs, `${id}.txt`),
				tmp: `${final}.tmp`,
				opts: { ...opts, sampleOut: opts.sample ? path.join(o.out, `${id}.sample.json`) : undefined },
			};
			const attempts = fs.existsSync(running(id)) ? Number(fs.readFileSync(running(id), "utf8")) : 0;
			fs.writeFileSync(running(id), String(attempts + 1));
			const t = performance.now();
			const r = await slot.run(task);
			const ms = performance.now() - t;
			fs.rmSync(running(id), { force: true });
			if (r.ok) {
				fs.renameSync(task.tmp, final);
			} else {
				failed++;
				fs.rmSync(task.tmp, { force: true });
				const rec = { id, reason: r.error, timedOut: Boolean(r.timedOut), ms, output: r.output ?? "" };
				fs.writeFileSync(path.join(o.out, `${id}.err`), `${JSON.stringify(rec)}\n`);
			}
			done++;
			if (done % 200 === 0 || done === todo.length) {
				console.error(`${done}/${todo.length} (${failed} failed) ${((performance.now() - t0) / 1000).toFixed(0)}s`);
			}
		}
		await slot.stop();
	};
	await Promise.all(Array.from({ length: o.jobs }, loop));
}

// ---------------------------------------------------------------------------------------------
// diff

function readDump(dir, id) {
	const j = path.join(dir, `${id}.json`);
	if (fs.existsSync(j)) return { dump: JSON.parse(fs.readFileSync(j, "utf8")) };
	const e = path.join(dir, `${id}.err`);
	if (fs.existsSync(e)) return { err: JSON.parse(fs.readFileSync(e, "utf8")) };
	return {};
}

function dumpIds(dir) {
	return [...new Set(fs.readdirSync(dir).flatMap((f) => (/^[^.]+\.(json|err)$/.test(f) ? [f.replace(/\.(json|err)$/, "")] : [])))].sort();
}

function cmdDiff(args) {
	const [gdir, cdir] = args._;
	const rtol = Number(args.rtol ?? 0);
	const ids = args.ids ? readLines(args.ids) : dumpIds(gdir);
	const counts = {};
	for (const id of ids) {
		const g = readDump(gdir, id);
		const c = readDump(cdir, id);
		let status;
		let detail = "";
		if (!g.dump) {
			status = "golden-error";
			detail = g.err?.reason ?? "no golden";
		} else if (!c.dump) {
			status = c.err ? "cand-error" : "missing";
			detail = c.err?.reason ?? "";
		} else {
			const { problems } = diff(g.dump, c.dump, rtol);
			status = problems.length ? "FAIL" : "ok";
			detail = problems.join(" | ");
		}
		counts[status] = (counts[status] ?? 0) + 1;
		console.log(`${id}\t${status}\t${detail.replace(/\s+/g, " ").slice(0, 2000)}`);
	}
	console.error(Object.entries(counts).map(([k, v]) => `${k} ${v}`).join(", "));
	if ((counts.FAIL ?? 0) + (counts["cand-error"] ?? 0) + (counts.missing ?? 0) > 0) process.exitCode = 1;
}

// ---------------------------------------------------------------------------------------------
// errors

// Buckets a failure message: identifiers, numbers and quoted text are masked so that the same
// kind of failure on different characters or lines groups together.
function bucket(reason) {
	return reason
		.replace(/^(initializeWorker|initializeAggregator|simulate|aggregate|flush|sample) failed: /, "$1: ")
		.replace(/"[^"]*"|'[^']*'|`[^`]*`/g, "<q>")
		.replace(/\bln\s*\d+\b|\bline \d+\b/gi, "line N")
		.replace(/\d+(\.\d+)?/g, "N")
		.slice(0, 160);
}

function cmdErrors(args) {
	const dir = args._[0];
	const ids = dumpIds(dir);
	const groups = new Map();
	const ok = [];
	for (const id of ids) {
		const r = readDump(dir, id);
		if (r.dump) {
			ok.push(id);
			continue;
		}
		const k = bucket(r.err.reason);
		if (!groups.has(k)) groups.set(k, []);
		groups.get(k).push({ id, reason: r.err.reason });
	}
	console.log(`${ids.length} configs: ${ok.length} ran, ${ids.length - ok.length} failed`);
	for (const [k, v] of [...groups].sort((a, b) => b[1].length - a[1].length)) {
		console.log(`${String(v.length).padStart(5)}  ${k}`);
		for (const x of v.slice(0, 2)) console.log(`         e.g. ${x.id}: ${x.reason.slice(0, 200)}`);
	}
	if (args["ids-out"]) fs.writeFileSync(args["ids-out"], `${ok.join("\n")}\n`);
}

// ---------------------------------------------------------------------------------------------
// select

function readEntries(file) {
	const [head, ...rows] = fs.readFileSync(file, "utf8").split("\n").filter(Boolean);
	const cols = head.split("\t");
	const m = new Map();
	for (const r of rows) {
		const v = r.split("\t");
		const e = Object.fromEntries(cols.map((c, i) => [c, v[i]]));
		e.charList = e.chars ? e.chars.split(",") : [];
		m.set(e.id, e);
	}
	return m;
}

function stratum(e) {
	const t = Number(e.targets);
	const d = Number(e.dur_mean);
	const tb = t <= 1 ? "1t" : t === 2 ? "2t" : "3t+";
	const nb = Number(e.nchars) >= 4 ? "4c" : "<4c";
	const db = d < 60 ? "<60s" : d < 100 ? "60-100s" : d < 140 ? "100-140s" : ">=140s";
	return `${tb}/${nb}/${db}`;
}

function cmdSelect(args) {
	const entries = readEntries(args.entries);
	const exclude = new Set(args.exclude ? readLines(args.exclude) : []);
	const pool = readLines(args.ids).filter((id) => entries.has(id) && !exclude.has(id));
	const n = Math.min(Number(args.n), pool.length);
	const rand = mulberry32(Number(args.seed ?? 1));
	const picked = [];
	const taken = new Set();
	const take = (id) => {
		picked.push(id);
		taken.add(id);
	};

	if (args.cover) {
		const uncovered = new Set(pool.flatMap((id) => entries.get(id).charList));
		const order = shuffle([...pool], rand);
		while (uncovered.size > 0 && picked.length < n) {
			let best = null;
			let bestGain = 0;
			for (const id of order) {
				if (taken.has(id)) continue;
				const gain = entries.get(id).charList.filter((c) => uncovered.has(c)).length;
				if (gain > bestGain) {
					best = id;
					bestGain = gain;
				}
			}
			if (!best) break;
			take(best);
			for (const c of entries.get(best).charList) uncovered.delete(c);
		}
		console.error(`cover: ${picked.length} configs cover every character in the pool${uncovered.size ? ` but ${uncovered.size}` : ""}`);
	}

	// Stratified fill: each stratum gets its share of the remaining picks (largest remainder).
	const rest = n - picked.length;
	const byStratum = new Map();
	for (const id of shuffle(pool.filter((id) => !taken.has(id)), rand)) {
		const s = stratum(entries.get(id));
		if (!byStratum.has(s)) byStratum.set(s, []);
		byStratum.get(s).push(id);
	}
	const total = [...byStratum.values()].reduce((s, v) => s + v.length, 0);
	const quotas = [...byStratum].map(([s, ids]) => {
		const exact = (rest * ids.length) / total;
		return { s, ids, q: Math.floor(exact), frac: exact - Math.floor(exact) };
	});
	let left = rest - quotas.reduce((s, x) => s + x.q, 0);
	for (const x of [...quotas].sort((a, b) => b.frac - a.frac)) {
		if (left <= 0) break;
		x.q++;
		left--;
	}
	for (const x of quotas) for (const id of x.ids.slice(0, x.q)) take(id);

	const chars = new Set(picked.flatMap((id) => entries.get(id).charList));
	const poolChars = new Set(pool.flatMap((id) => entries.get(id).charList));
	console.error(`selected ${picked.length} of ${pool.length}; ${chars.size}/${poolChars.size} characters`);
	console.log(picked.join("\n"));
}

// ---------------------------------------------------------------------------------------------
// manifest

// A cold 20-iteration dump runs about 1.8x slower per iteration than steady state (V8 tier-up).
const COLD_TO_WARM = 0.55;

function cmdManifest(args) {
	const out = path.resolve(args.out);
	const blockMs = Number(args["block-ms"] ?? 150);
	fs.mkdirSync(out, { recursive: true });
	const rows = [
		"# name\tblock\twarmup\tsource",
		`# written by corpus.mjs manifest; block ~${blockMs} ms and warmup ~1 s of steady-state work at the timings in ${args.timings}`,
	];
	for (const id of readLines(args.ids)) {
		const d = JSON.parse(fs.readFileSync(path.join(args.timings, `${id}.json`), "utf8"));
		const est = Math.max(0.01, (d.perf.simulate_ms / d.iters) * COLD_TO_WARM);
		const block = Math.min(1000, Math.max(2, Math.round(blockMs / est)));
		const warmup = Math.min(1000, Math.max(10, Math.round(1000 / est)));
		fs.copyFileSync(path.join(args.configs, `${id}.txt`), path.join(out, `${id}.txt`));
		rows.push(`${id}\t${block}\t${warmup}\t${path.join(args.configs, `${id}.txt`)}`);
	}
	fs.writeFileSync(path.join(out, "manifest.tsv"), `${rows.join("\n")}\n`);
	console.error(`wrote ${out}/manifest.tsv (${rows.length - 2} configs)`);
}

// ---------------------------------------------------------------------------------------------
// speed

function cmdSpeed(args) {
	const rows = args._.flatMap((f) =>
		fs
			.readFileSync(f, "utf8")
			.split("\n")
			.filter(Boolean)
			.map((l) => JSON.parse(l)),
	);
	const labels = [...new Set(rows.map((r) => r.label))];
	const [base, ...cands] = labels;
	const byCfg = new Map();
	for (const r of rows) {
		if (!byCfg.has(r.config)) byCfg.set(r.config, new Map());
		const m = byCfg.get(r.config);
		if (!m.has(r.label)) m.set(r.label, new Map());
		m.get(r.label).set(`${r.proc}:${r.round}`, r);
	}
	const gmed = (xs) => Math.exp(quantile(sortNum(xs.map(Math.log)), 0.5));
	const out = [];
	for (const [cfg, m] of byCfg) {
		const b = m.get(base);
		if (!b) continue;
		for (const l of cands) {
			const c = m.get(l);
			if (!c) continue;
			const rs = [...c.keys()].filter((k) => b.has(k)).map((k) => c.get(k).cpu_ms_per_iter / b.get(k).cpu_ms_per_iter);
			const s = sortNum(rs);
			const med = (x, k) => quantile(sortNum([...x.values()].map((r) => r[k])), 0.5);
			const peak = (x, k) => Math.max(...[...x.values()].map((r) => r[k] ?? Number.NaN));
			out.push({
				config: cfg,
				label: l,
				n: rs.length,
				base_ms: med(b, "cpu_ms_per_iter"),
				cand_ms: med(c, "cpu_ms_per_iter"),
				ratio: gmed(rs),
				ratio_p25: quantile(s, 0.25),
				ratio_p75: quantile(s, 0.75),
				base_sim_mb: peak(b, "sim_mem_mb"),
				cand_sim_mb: peak(c, "sim_mem_mb"),
				base_agg_mb: peak(b, "agg_mem_mb"),
				cand_agg_mb: peak(c, "agg_mem_mb"),
			});
		}
	}
	const cols = Object.keys(out[0] ?? {});
	const fmt = (v) => (typeof v === "number" && !Number.isInteger(v) ? v.toFixed(4) : v);
	const tsv = [cols.join("\t"), ...out.map((r) => cols.map((c) => fmt(r[c])).join("\t"))].join("\n");
	if (args.tsv) fs.writeFileSync(args.tsv, `${tsv}\n`);
	else console.log(tsv);
	for (const l of cands) {
		const rs = sortNum(out.filter((r) => r.label === l).map((r) => r.ratio));
		const g = Math.exp(rs.reduce((s, x) => s + Math.log(x), 0) / rs.length);
		const q = (p) => quantile(rs, p).toFixed(3);
		console.error(
			`${l} / ${base} over ${rs.length} configs: geomean ${g.toFixed(3)}, median ${q(0.5)}, p10 ${q(0.1)}, p90 ${q(0.9)}, min ${rs[0]?.toFixed(3)}, max ${rs[rs.length - 1]?.toFixed(3)}`,
		);
	}
}

// ---------------------------------------------------------------------------------------------

function main() {
	const args = parseArgs(process.argv.slice(2));
	const cmd = args._.shift();
	const cmds = { dump: cmdDump, diff: cmdDiff, errors: cmdErrors, select: cmdSelect, manifest: cmdManifest, speed: cmdSpeed };
	if (!cmds[cmd]) {
		console.error("usage: corpus.mjs dump|diff|errors|select|manifest|speed ... (see the header comment)");
		process.exit(2);
	}
	Promise.resolve(cmds[cmd](args)).catch((e) => {
		console.error(e?.stack ?? String(e));
		process.exit(1);
	});
}

if (isMainThread) main();
else if (workerData?.role === "corpus-dump") dumpWorker();
