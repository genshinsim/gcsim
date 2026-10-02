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
	// The instance's Go program exits while handling a request this matches, as on a fatal
	// error such as running out of memory. It then fails every request, as wasm_exec.js does
	// once the program has exited.
	static crashOn: ((request: string, msg: Message) => boolean) | null = null;
	// "role:type": the worker script throws while handling this request, which the page sees
	// as an error event on the Worker.
	static throwOn: string | null = null;
	// ms a request takes to handle
	static delay: (request: string, msg: Message) => number = () => 0;
	// runs requested minus results aggregated, across the pool
	static ahead = 0;
	static maxAhead = 0;
	// every response the page got, from any worker
	static toPage: Message[] = [];

	role: "helper" | "aggregator" | "worker";
	onmessage: ((ev: MessageEvent) => void) | null = null;
	onerror: ((ev: ErrorEvent) => void) | null = null;
	received: Message[] = [];
	terminated = false;
	loaded = false;
	exited = false;
	queuedRuns = 0;
	maxQueuedRuns = 0;
	private inbox: Message[] = [];
	private busy = false;
	private added = 0;
	// the ends of the channels between the helper and the aggregator or a sim worker
	private ports: MessagePort[] = [];
	module: unknown = null;
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

	postMessage(msg: Message, _transfer?: Transferable[]) {
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
		for (const port of this.ports) {
			port.close();
		}
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
		setTimeout(
			() => {
				this.busy = false;
				if (this.terminated) {
					return;
				}
				const resp = this.handle(msg);
				if (resp != null) {
					this.respond(resp);
				}
				this.pump();
			},
			FakeWorker.delay(`${this.role}:${msg.type}`, msg),
		);
	}

	private respond(resp: Message) {
		FakeWorker.toPage.push(resp);
		this.onmessage?.(new MessageEvent("message", { data: resp }));
	}

	private handle(msg: Message): Message | null {
		const run = msg.run;
		const request = `${this.role}:${msg.type}`;
		if (this.role !== "helper" && msg.type !== "ready" && !this.loaded) {
			FakeWorker.early.push(request);
			return null;
		}
		if (this.exited || FakeWorker.crashOn?.(request, msg)) {
			this.exited = true;
			const reason = `${this.role} exited`;
			return { type: "failed", run, id: msg.id, reason, fatal: true };
		}
		if (request === FakeWorker.throwOn) {
			this.onerror?.({ message: `${this.role} threw` } as ErrorEvent);
			return null;
		}
		switch (request) {
			case "helper:ready":
				// the helper compiles the wasm
				return FakeWorker.failLoad
					? { type: "failed", reason: "no wasm", fatal: true }
					: { type: "ready" };
			case "helper:share": {
				// and sends the module to the aggregator or sim worker holding the other end
				// once it asks, not to the page
				const port = msg.port as MessagePort;
				this.ports.push(port);
				port.onmessage = () => {
					if (!this.terminated) {
						port.postMessage(fakeModule);
						this.respond({ type: "shared", id: msg.id });
					}
				};
				return null;
			}
			case "helper:validate":
				return { type: "validated", id: msg.id, cfg: {} };
			case "aggregator:ready":
			case "worker:ready": {
				// asks the helper for the module and loads once it has it
				const port = msg.port as MessagePort;
				this.ports.push(port);
				port.onmessage = (ev) => {
					port.close();
					this.module = ev.data;
					const loaded = () => {
						this.loaded = true;
						if (!this.terminated) {
							this.respond({ type: "ready" });
						}
					};
					if (FakeWorker.loadDelay > 0) {
						setTimeout(loaded, FakeWorker.loadDelay);
					} else {
						loaded();
					}
				};
				port.postMessage(null);
				return null;
			}
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
						signedJSON: JSON.stringify({
							simulator_settings: { iterations: this.iterations },
							statistics: { iterations: this.added },
						}),
						hash: `hash-${this.added}`,
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
// resolves to "timed out" instead of hanging the test
const within = <T>(p: Promise<T>, ms = 1000) =>
	Promise.race([p, new Promise((r) => setTimeout(r, ms, "timed out"))]);

const fakeModule = { compiled: "/main.wasm" } as unknown as WebAssembly.Module;

beforeEach(() => {
	FakeWorker.all = [];
	FakeWorker.failRun = false;
	FakeWorker.failLoad = false;
	FakeWorker.loadDelay = 0;
	FakeWorker.early = [];
	FakeWorker.crashOn = null;
	FakeWorker.throwOn = null;
	FakeWorker.delay = () => 0;
	FakeWorker.ahead = 0;
	FakeWorker.maxAhead = 0;
	FakeWorker.toPage = [];
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
	const updates: {
		result: model.SimulationResult;
		raw: string;
		hash: string;
	}[] = [];
	const done = exec.run(String(iterations), (result, raw, hash) =>
		updates.push({ result, raw, hash }),
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
		const raw = JSON.stringify({
			simulator_settings: { iterations: 25 },
			statistics: { iterations: 25 },
		});
		expect(updates.at(-1)).toEqual({
			result: JSON.parse(raw),
			raw,
			hash: "hash-25",
		});
	});

	it("shows the latest iterations within the flush interval when the next is slow", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(1);
		// the fake's flushes cost 1 ms, so they are spaced by the 100 ms minimum; the last
		// iteration takes much longer than that
		FakeWorker.delay = (request, msg) =>
			request === "worker:run" ? (msg.itr === 3 ? 600 : 10) : 0;

		const { updates, done } = start(exec, 4);
		await vi.waitFor(() => expect(updates.at(-1)?.hash).toBe("hash-3"), {
			timeout: 400,
		});
		await done;
		expect(updates.map((u) => u.hash)).toEqual(["hash-1", "hash-3", "hash-4"]);
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
		// one share per instance, each over the channel its ready request carries
		const pool = [...byRole("aggregator"), ...byRole("worker")];
		expect(pool).toHaveLength(5);
		expect(sent("helper", "share")).toHaveLength(5);
		for (const member of pool) {
			const readies = member.received.filter((m) => m.type === "ready");
			expect(readies).toHaveLength(1);
			expect(readies[0].port).toBeInstanceOf(MessagePort);
			expect(member.module).toEqual(fakeModule);
		}
	});

	// Under the page's CSP, Firefox can't deserialize a WebAssembly.Module on the page, so
	// the module must go worker to worker.
	it("never sends the module to the page", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(3);

		await start(exec, 5).done;

		expect(FakeWorker.toPage.length).toBeGreaterThan(0);
		for (const resp of FakeWorker.toPage) {
			expect(Object.values(resp)).not.toContainEqual(fakeModule);
		}
	});

	it("fails the run when the helper dies before sharing the module, then recovers", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(2);
		FakeWorker.crashOn = (request) => request === "helper:share";

		await expect(within(start(exec, 5).done)).rejects.toBe("helper exited");
		expect(exec.running()).toBe(false);
		expect(byRole("helper")[0].terminated).toBe(true);
		for (const member of [...byRole("aggregator"), ...byRole("worker")]) {
			expect(member.terminated).toBe(true);
		}

		FakeWorker.crashOn = null;
		const second = start(exec, 5);
		expect(await within(second.done)).toBe(true);
		expect(second.updates.at(-1)?.hash).toBe("hash-5");
		expect(byRole("helper")).toHaveLength(2);
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

		expect(await within(second.done)).toBe(true);
		expect(second.updates.at(-1)?.hash).toBe("hash-10");
		expect(FakeWorker.early).toEqual([]);
		// the cancelled run sent nothing once the pool had loaded
		expect(sent("aggregator", "initialize")).toHaveLength(1);
		expect(sent("worker", "initialize")).toHaveLength(2);
	});

	it("fails the run when the aggregator's Go program exits, then replaces it", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(2);
		FakeWorker.crashOn = (request) => request === "aggregator:add";

		expect(await within(start(exec, 10).done.catch((e) => e))).toBe(
			"aggregator exited",
		);
		expect(exec.running()).toBe(false);
		const [crashed] = byRole("aggregator");
		expect(crashed.terminated).toBe(true);

		FakeWorker.crashOn = null;
		const second = start(exec, 10);
		expect(await within(second.done)).toBe(true);
		expect(second.updates.at(-1)?.hash).toBe("hash-10");
		expect(byRole("aggregator")).toHaveLength(2);
		// the sim workers were fine and are kept
		expect(byRole("worker")).toHaveLength(2);
	});

	it("replaces a sim worker whose script threw and keeps the worker count", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(3);
		await start(exec, 10).done;
		FakeWorker.throwOn = "worker:initialize";

		expect(await within(start(exec, 10).done.catch((e) => e))).toContain(
			"worker threw",
		);
		expect(exec.running()).toBe(false);

		FakeWorker.throwOn = null;
		const third = start(exec, 10);
		expect(await within(third.done)).toBe(true);
		const alive = byRole("worker").filter((w) => !w.terminated);
		expect(alive).toHaveLength(3);
		expect(byRole("aggregator")).toHaveLength(1);
	});

	it("replaces sim workers that die on a cancelled run's leftover requests", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(3);
		// slow sims, so every worker has run requests queued when the run is cancelled
		FakeWorker.delay = (request) => (request === "worker:run" ? 5 : 0);
		const first = start(exec, 1000);
		await vi.waitFor(() => expect(first.updates.length).toBeGreaterThan(0));
		exec.cancel();
		const cancelled = sent("worker", "run")[0].run;
		FakeWorker.crashOn = (request, msg) =>
			request === "worker:run" && msg.run === cancelled;

		const second = start(exec, 10);
		expect(await within(second.done)).toBe(true);
		expect(second.updates.at(-1)?.hash).toBe("hash-10");
		expect(byRole("worker").filter((w) => w.exited)).toHaveLength(3);
		expect(byRole("worker").filter((w) => !w.terminated)).toHaveLength(3);
	});

	it("replaces an aggregator that dies on a cancelled run's leftover request", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(3);
		// the fake aggregator falls behind three workers, so adds are queued when the run is
		// cancelled
		const first = start(exec, 1000);
		await vi.waitFor(() => expect(first.updates.length).toBeGreaterThan(0));
		exec.cancel();
		const cancelled = sent("worker", "run")[0].run;
		FakeWorker.crashOn = (request, msg) =>
			request === "aggregator:add" && msg.run === cancelled;

		const second = start(exec, 10);
		expect(await within(second.done)).toBe(true);
		expect(second.updates.at(-1)?.hash).toBe("hash-10");
		expect(byRole("aggregator").map((w) => w.exited)).toEqual([true, false]);
	});

	it("fails the run when a sim worker's Go program exits on the run's own request", async () => {
		const exec = new WasmExecutor("/main.wasm");
		exec.setWorkerCount(2);
		FakeWorker.crashOn = (request, msg) =>
			request === "worker:run" && msg.itr === 3;

		expect(await within(start(exec, 10).done.catch((e) => e))).toBe(
			"worker exited",
		);
		FakeWorker.crashOn = null;
		expect(await within(start(exec, 10).done)).toBe(true);
		expect(byRole("worker")).toHaveLength(3);
	});

	it("replaces a helper whose Go program exited", async () => {
		const exec = new WasmExecutor("/main.wasm");
		FakeWorker.crashOn = (request) => request === "helper:validate";

		await expect(within(exec.validate("cfg"))).rejects.toBe("helper exited");
		expect(byRole("helper")[0].terminated).toBe(true);

		FakeWorker.crashOn = null;
		await expect(within(exec.validate("cfg"))).resolves.toEqual({});
		expect(byRole("helper")).toHaveLength(2);
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
	])(
		"uses %s cores -> %s workers on devices reporting 8 GB",
		(cores, expected) => {
			withNavigator(cores, 8);
			expect(defaultWorkerCount()).toBe(expected);
		},
	);

	it("keeps 3 on devices reporting under 4 GB", () => {
		withNavigator(16, 2);
		expect(defaultWorkerCount()).toBe(3);
	});

	it("uses at most 4 on devices reporting 4 GB and where the memory is unknown", () => {
		withNavigator(16, 4);
		expect(defaultWorkerCount()).toBe(4);
		// Firefox and Safari
		withNavigator(16);
		expect(defaultWorkerCount()).toBe(4);
		withNavigator(4);
		expect(defaultWorkerCount()).toBe(3);
	});
});
