import type { model } from "@gcsim/types";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { defaultWorkerCount, WasmExecutor } from "./WasmExecutor";

// A stand-in for the browser Worker running Workers/{helper,aggregator,worker}.ts. Each
// fake handles its messages one at a time, in order, on later macrotasks, like a real
// worker's queue.
type Message = { type: string; [key: string]: unknown };

class FakeWorker {
	static all: FakeWorker[] = [];
	static failRun = false;
	static failLoad = false;
	// With loadDelay > 0, the aggregator and the sim workers finish loading that many ms
	// after handling "ready", like WebAssembly.instantiate resolving in a later task.
	static loadDelay = 0;
	// Requests an instance handled before it had loaded. The real worker scripts throw on
	// those: the Go functions only exist once go.run has returned.
	static early: string[] = [];
	// runs requested minus results aggregated, across the pool
	static ahead = 0;
	static maxAhead = 0;

	role: "helper" | "aggregator" | "worker";
	onmessage: ((ev: MessageEvent) => void) | null = null;
	received: Message[] = [];
	terminated = false;
	loaded = false;
	queuedRuns = 0;
	maxQueuedRuns = 0;
	private inbox: Message[] = [];
	private busy = false;
	private added = 0;
	private iterations = 0;

	constructor(url: URL | string) {
		const name = String(url);
		this.role = name.includes("helper")
			? "helper"
			: name.includes("aggregator")
				? "aggregator"
				: "worker";
		FakeWorker.all.push(this);
	}

	postMessage(msg: Message) {
		this.received.push(msg);
		if (msg.type === "run") {
			this.queuedRuns++;
			this.maxQueuedRuns = Math.max(this.maxQueuedRuns, this.queuedRuns);
			FakeWorker.ahead++;
			FakeWorker.maxAhead = Math.max(FakeWorker.maxAhead, FakeWorker.ahead);
		}
		this.inbox.push(msg);
		this.pump();
	}

	terminate() {
		this.terminated = true;
		this.inbox = [];
	}

	private pump() {
		if (this.busy || this.terminated) {
			return;
		}
		const msg = this.inbox.shift();
		if (msg == null) {
			return;
		}
		this.busy = true;
		setTimeout(() => {
			this.busy = false;
			if (this.terminated) {
				return;
			}
			const resp = this.handle(msg);
			if (resp != null) {
				this.respond(resp);
			}
			this.pump();
		}, 0);
	}

	private respond(resp: Message) {
		this.onmessage?.(new MessageEvent("message", { data: resp }));
	}

	private handle(msg: Message): Message | null {
		const run = msg.run;
		if (this.role !== "helper" && msg.type !== "ready" && !this.loaded) {
			FakeWorker.early.push(`${this.role}:${msg.type}`);
			return null;
		}
		switch (`${this.role}:${msg.type}`) {
			case "helper:ready":
				// the helper compiles the wasm and hands the module back
				return FakeWorker.failLoad
					? { type: "failed", reason: "no wasm" }
					: { type: "ready", module: fakeModule };
			case "helper:validate":
				return { type: "validated", id: msg.id, cfg: {} };
			case "aggregator:ready":
			case "worker:ready":
				if (FakeWorker.loadDelay > 0) {
					setTimeout(() => {
						this.loaded = true;
						if (!this.terminated) {
							this.respond({ type: "ready" });
						}
					}, FakeWorker.loadDelay);
					return null;
				}
				this.loaded = true;
				return { type: "ready" };
			case "aggregator:initialize":
				this.added = 0;
				this.iterations = Number(msg.cfg);
				return {
					type: "initialized",
					run,
					result: { simulator_settings: { iterations: this.iterations } },
				};
			case "aggregator:add":
				this.added++;
				FakeWorker.ahead--;
				return { type: "done", run };
			case "aggregator:flush":
				return {
					type: "result",
					run,
					final: msg.final,
					ms: 1,
					result: {
						hash: `hash-${this.added}`,
						stats: { iterations: this.added },
					},
				};
			case "worker:initialize":
				return { type: "initialized", run };
			case "worker:run":
				this.queuedRuns--;
				if (FakeWorker.failRun) {
					return { type: "failed", run, reason: "boom" };
				}
				return { type: "done", run, itr: msg.itr, result: new Uint8Array(4) };
		}
		throw new Error(`unexpected ${this.role} message ${msg.type}`);
	}
}

const byRole = (role: FakeWorker["role"]) =>
	FakeWorker.all.filter((w) => w.role === role);
const sent = (role: FakeWorker["role"], type: string) =>
	byRole(role).flatMap((w) => w.received.filter((m) => m.type === type));
const settle = () => new Promise((r) => setTimeout(r, 20));

const fakeModule = {} as WebAssembly.Module;

beforeEach(() => {
	FakeWorker.all = [];
	FakeWorker.failRun = false;
	FakeWorker.failLoad = false;
	FakeWorker.loadDelay = 0;
	FakeWorker.early = [];
	FakeWorker.ahead = 0;
	FakeWorker.maxAhead = 0;
	vi.stubGlobal("Worker", FakeWorker);
	vi.stubGlobal(
		"fetch",
		vi.fn(() => Promise.resolve(new Response())),
	);
});

afterEach(() => {
	vi.restoreAllMocks();
	vi.unstubAllGlobals();
});

// The fake aggregator reads the iteration count from the config string.
function start(exec: WasmExecutor, iterations: number) {
	const updates: { result: model.SimulationResult; hash: string }[] = [];
	const done = exec.run(String(iterations), (result, hash) =>
		updates.push({ result, hash }),
	);
	return { updates, done };
}

describe("WasmExecutor", () => {
	it("runs exactly the configured iterations and ends on the final flush", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(3);

		const { updates, done } = start(exec, 25);
		await expect(done).resolves.toBe(true);

		expect(exec.running()).toBe(false);
		expect(sent("worker", "run")).toHaveLength(25);
		expect(sent("aggregator", "add")).toHaveLength(25);
		const flushes = sent("aggregator", "flush");
		expect(flushes.filter((m) => m.final)).toHaveLength(1);
		expect(flushes.at(-1)?.final).toBe(true);
		expect(updates.at(-1)).toEqual({
			result: {
				simulator_settings: { iterations: 25 },
				statistics: { iterations: 25 },
			},
			hash: "hash-25",
		});
	});

	it("keeps at most two runs queued per worker", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(2);

		await start(exec, 40).done;

		for (const w of byRole("worker")) {
			expect(w.maxQueuedRuns).toBe(2);
		}
	});

	it("shares the helper's compiled module with the aggregator and the workers", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(4);

		await exec.validate("cfg");
		await start(exec, 5).done;
		await start(exec, 5).done;

		expect(sent("helper", "ready")).toEqual([
			{ type: "ready", wasm: "/main.wasm" },
		]);
		const readies = [
			...sent("aggregator", "ready"),
			...sent("worker", "ready"),
		];
		expect(readies).toHaveLength(5);
		for (const ready of readies) {
			expect(ready).toEqual({ type: "ready", module: fakeModule });
		}
	});

	it("fails waiting requests when the helper can't load, then retries", async () => {
		const exec = new WasmExecutor("/main.wasm");
		FakeWorker.failLoad = true;

		await expect(exec.validate("cfg")).rejects.toBe("no wasm");
		expect(byRole("helper")[0].terminated).toBe(true);

		FakeWorker.failLoad = false;
		await expect(exec.validate("cfg")).resolves.toEqual({});
		expect(byRole("helper")).toHaveLength(2);
	});

	it("resolves validate without polling timers", async () => {
		const exec = new WasmExecutor("/main.wasm");
		const timeout = vi.spyOn(globalThis, "setTimeout");

		await exec.validate("cfg");

		// only the fake worker's own queue ticks (0 ms) are scheduled
		expect(timeout.mock.calls.every(([, ms]) => !ms)).toBe(true);
	});

	it("keeps the pool across a cancel and drops the cancelled run's responses", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(3);

		const first = start(exec, 1000);
		await vi.waitFor(() => expect(first.updates.length).toBeGreaterThan(0));
		exec.cancel();
		expect(exec.running()).toBe(false);
		const seen = first.updates.length;
		const runs = sent("worker", "run").length;
		// longer than the minimum flush interval
		await new Promise((r) => setTimeout(r, 150));
		expect(first.updates).toHaveLength(seen);
		expect(sent("worker", "run")).toHaveLength(runs);

		const second = start(exec, 30);
		await expect(second.done).resolves.toBe(true);
		expect(second.updates.at(-1)?.hash).toBe("hash-30");
		expect(first.updates).toHaveLength(seen);
		expect(FakeWorker.all.some((w) => w.terminated)).toBe(false);
		expect(byRole("aggregator")).toHaveLength(1);
	});

	it("waits for a pool that is still loading, also after a cancel", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(2);
		FakeWorker.loadDelay = 100;

		start(exec, 10);
		// the helper has handed over the module and the pool exists, but hasn't loaded yet
		await vi.waitFor(() => expect(byRole("worker")).toHaveLength(2), {
			interval: 1,
		});
		exec.cancel();
		const second = start(exec, 10);

		const outcome = await Promise.race([
			second.done,
			new Promise((r) => setTimeout(r, 1000, "timed out")),
		]);
		expect(outcome).toBe(true);
		expect(second.updates.at(-1)?.hash).toBe("hash-10");
		expect(FakeWorker.early).toEqual([]);
		// the cancelled run sent nothing once the pool had loaded
		expect(sent("aggregator", "initialize")).toHaveLength(1);
		expect(sent("worker", "initialize")).toHaveLength(2);
	});

	it("bounds the work requested ahead of the aggregator", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(3);

		// three workers produce results faster than the one fake aggregator consumes them
		await start(exec, 200).done;

		expect(FakeWorker.maxAhead).toBe(3 * 4);
	});

	it("frees the executor when a run fails", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(2);
		FakeWorker.failRun = true;

		await expect(start(exec, 10).done).rejects.toBe("boom");
		expect(exec.running()).toBe(false);

		FakeWorker.failRun = false;
		await settle();
		await expect(start(exec, 10).done).resolves.toBe(true);
	});
});

describe("defaultWorkerCount", () => {
	const withNavigator = (hardwareConcurrency?: number, deviceMemory?: number) =>
		vi.stubGlobal("navigator", { hardwareConcurrency, deviceMemory });

	it.each([
		[undefined, 3],
		[2, 3],
		[4, 3],
		[6, 4],
		[8, 6],
		[10, 8],
		[32, 8],
	])("uses %s cores -> %s workers", (cores, expected) => {
		withNavigator(cores);
		expect(defaultWorkerCount()).toBe(expected);
	});

	it("keeps 3 on devices reporting under 4 GB", () => {
		withNavigator(8, 2);
		expect(defaultWorkerCount()).toBe(3);
		withNavigator(8, 4);
		expect(defaultWorkerCount()).toBe(6);
	});
});
