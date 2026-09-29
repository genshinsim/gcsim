import type { model, ParsedResult, Sample } from "@gcsim/types";
import type { Executor } from "./Executor";
import { Aggregator, Helper, SimWorker } from "./Workers/common";

// Every flush re-serializes and re-signs the whole result on the aggregator, the thread
// that also has to aggregate every iteration; one costs 10-35 ms for a typical team. While
// a run is in progress, flush no more often than FLUSH_COST_RATIO times the last flush's
// cost, clamped to [MIN_FLUSH_INTERVAL, MAX_FLUSH_INTERVAL] ms, so flushing takes at most
// about a tenth of the aggregator. The final flush is sent as soon as the last iteration is
// aggregated.
//
// Intermediate results stay signed: a cancelled, failed or interrupted (pagehide) run keeps
// the last one it received, and that result must still be shareable.
const MIN_FLUSH_INTERVAL = 100;
const MAX_FLUSH_INTERVAL = 1000;
const FLUSH_COST_RATIO = 10;

// Run requests queued per worker. With more than one, a worker starts its next iteration
// right away instead of waiting a round trip through the main thread, which may be busy
// rendering results.
const PIPELINE_DEPTH = 2;

// Cap on iterations requested but not yet aggregated, per worker. It only binds when the
// aggregator falls behind, and then bounds its backlog: the memory held by queued results
// (100-200 KB each) and the leftover work a cancelled run leaves in its queue.
const MAX_OUTSTANDING_PER_WORKER = 4;

// Default number of sim workers: one per logical core, minus one for the aggregator and one
// for the page, clamped to [3, 8]. 3 was the fixed default before, so no machine gets fewer
// workers than it used to. The cap bounds memory: every worker holds its own wasm instance
// (~40-60 MB of linear memory on top of the shared compiled module). Devices that report under
// 4 GB of memory (navigator.deviceMemory, Chromium only) keep 3. Users can still set 1-30.
export function defaultWorkerCount(): number {
	const cores = navigator.hardwareConcurrency;
	const memory = (navigator as { deviceMemory?: number }).deviceMemory;
	if (!cores || (memory != null && memory < 4)) {
		return 3;
	}
	return Math.min(Math.max(cores - 2, 3), 8);
}

export class WasmExecutor implements Executor {
	private helper: HelperExecutor;
	private aggregator: PoolWorker | null;
	private workers: PoolWorker[];
	private workerCount: number;
	private isRunning: boolean;
	private runId: number;
	private runStarted: number;

	constructor(wasm: string) {
		this.helper = new HelperExecutor(wasm);

		this.aggregator = null;
		this.workers = [];
		this.workerCount = defaultWorkerCount();
		this.isRunning = false;
		this.runId = 0;
		this.runStarted = 0;
	}

	public ready(): Promise<boolean> {
		return new Promise((resolve) => resolve(!this.isRunning));
	}

	public running(): boolean {
		return this.isRunning;
	}

	public setWorkerCount(count: number) {
		this.workerCount = count;
	}

	// The aggregator and the sim workers for a run: those kept from earlier runs, minus any that
	// failed to load, plus new ones up to the worker count. Some may still be loading.
	private pool(module: WebAssembly.Module) {
		if (this.aggregator == null || this.aggregator.failed) {
			this.aggregator = new PoolWorker(
				new Worker(new URL("./Workers/aggregator.ts", import.meta.url)),
				Aggregator.ReadyRequest(module),
			);
		}
		this.workers = this.workers.filter((w) => !w.failed);
		for (const worker of this.workers.splice(this.workerCount)) {
			worker.terminate();
		}
		while (this.workers.length < this.workerCount) {
			this.workers.push(
				new PoolWorker(
					new Worker(new URL("./Workers/worker.ts", import.meta.url)),
					SimWorker.ReadyRequest(module),
				),
			);
		}
		return { aggregator: this.aggregator, workers: [...this.workers] };
	}

	public run(
		cfg: string,
		updateResult: (result: model.SimulationResult, hash: string) => void,
	): Promise<boolean | void> {
		this.isRunning = true;
		this.runStarted = performance.now();

		// The workers and the aggregator are reused across runs, so their queues can still hold
		// requests from a cancelled run. Every request carries the run id and every response
		// echoes it; responses from any other run are dropped. A cancelled run sends nothing
		// more, and its promise never settles.
		const run = ++this.runId;
		const current = () => this.runId === run;
		const isCurrent = (ev: MessageEvent) => current() && ev.data.run === run;

		return new Promise<boolean>((resolve, reject) => {
			const stop = () => {
				this.runId++;
				this.isRunning = false;
			};
			// a failed run frees the executor for the next one
			const fail = (reason: unknown) => {
				if (current()) {
					stop();
					reject(reason);
				}
			};

			this.helper.load().then((module) => {
				if (!current()) {
					return;
				}
				const { aggregator, workers } = this.pool(module);

				let result: model.SimulationResult | null = null;
				let maxIterations = 0;
				// the aggregator and the workers yet to answer the initialize request
				let initializing = workers.length + 1;
				let completed = 0;
				let flushing = false; // an intermediate flush is outstanding
				let nextFlush = 0;

				const maxOutstanding = workers.length * MAX_OUTSTANDING_PER_WORKER;
				const queued = workers.map(() => 0);
				let requested = 0;
				const dispatch = () => {
					// Near the end, queue one at a time so no worker sits on the last iterations
					// while the others go idle.
					const depth =
						maxIterations - requested > workers.length * PIPELINE_DEPTH
							? PIPELINE_DEPTH
							: 1;
					workers.forEach((worker, i) => {
						while (
							queued[i] < depth &&
							requested < maxIterations &&
							requested - completed < maxOutstanding
						) {
							worker.post(SimWorker.RunRequest(run, requested++));
							queued[i]++;
						}
					});
				};
				const initialized = () => {
					initializing--;
					if (initializing === 0) {
						dispatch();
					}
				};

				aggregator.listen((ev) => {
					if (!isCurrent(ev)) {
						return;
					}
					switch (ev.data.type as Aggregator.Response) {
						case Aggregator.Response.Initialized:
							result = (ev.data as Aggregator.InitializeResponse).result;
							maxIterations = result?.simulator_settings?.iterations ?? 1000;
							initialized();
							return;
						case Aggregator.Response.Result: {
							const resp = ev.data as Aggregator.ResultResponse;
							const { hash, stats } = resp.result;

							const out = Object.assign({}, result);
							out.statistics = stats;
							updateResult(out, hash);

							if (resp.final) {
								stop();
								resolve(true);
								if (this.runStarted > 0) {
									const end = performance.now();
									console.log(`run time: ${end - this.runStarted} ms`);
									this.runStarted = 0;
								}
								return;
							}
							flushing = false;
							nextFlush =
								performance.now() +
								Math.min(
									Math.max(FLUSH_COST_RATIO * resp.ms, MIN_FLUSH_INTERVAL),
									MAX_FLUSH_INTERVAL,
								);
							return;
						}
						case Aggregator.Response.Done:
							completed += 1;
							if (completed === maxIterations) {
								aggregator.post(Aggregator.FlushRequest(run, true));
							} else if (!flushing && performance.now() >= nextFlush) {
								flushing = true;
								aggregator.post(Aggregator.FlushRequest(run, false));
							}
							dispatch();
							return;
						case Aggregator.Response.Failed:
							fail((ev.data as Aggregator.FailedResponse).reason);
					}
				});

				workers.forEach((worker, i) => {
					worker.listen((ev) => {
						if (!isCurrent(ev)) {
							return;
						}
						switch (ev.data.type as SimWorker.Response) {
							case SimWorker.Response.Initialized:
								initialized();
								return;
							case SimWorker.Response.Done: {
								const resp: SimWorker.RunResponse = ev.data;
								queued[i]--;
								// transfer rather than copy; the result is only forwarded
								aggregator.post(Aggregator.AddRequest(run, resp.result), [
									resp.result.buffer,
								]);
								dispatch();
								return;
							}
							case SimWorker.Response.Failed:
								fail((ev.data as SimWorker.FailedResponse).reason);
						}
					});
				});

				// Instances kept from a cancelled run may still be loading; nothing but the ready
				// request may reach them before they have.
				Promise.all([aggregator.ready, ...workers.map((w) => w.ready)]).then(
					() => {
						if (!current()) {
							return;
						}
						aggregator.post(Aggregator.InitializeRequest(run, cfg));
						for (const worker of workers) {
							worker.post(SimWorker.InitializeRequest(run, cfg));
						}
					},
					fail,
				);
			}, fail);
		});
	}

	public cancel(): void {
		if (!this.isRunning) {
			return;
		}

		// The workers and the aggregator are kept. Requests of the cancelled run still in their
		// queues are processed, but their responses are dropped by run id, and the next run's
		// initialize request resets the aggregator after them.
		this.runId++;
		this.isRunning = false;
		console.log("execution canceled");

		if (this.runStarted > 0) {
			const end = performance.now();
			console.log(`cancelled with run time: ${end - this.runStarted} ms`);
			this.runStarted = 0;
		}
	}

	public validate(cfg: string): Promise<ParsedResult> {
		return this.helper.validate(cfg);
	}

	public sample(cfg: string, seed: string): Promise<Sample> {
		return this.helper.sample(cfg, seed);
	}

	public buildInfo(): { hash: string; date: string } {
		return this.helper.buildInfo();
	}
}

// A worker running one wasm instance, the aggregator or a sim worker. The executor keeps it
// across runs.
class PoolWorker {
	// Resolves once the instance has loaded. The worker scripts only define the Go functions
	// once go.run has returned, so no request but the ready request may reach it before then.
	readonly ready: Promise<void>;
	// it failed to load; the next run replaces it
	failed = false;
	private worker: Worker;
	private handler: (ev: MessageEvent) => void = () => {};

	constructor(
		worker: Worker,
		readyRequest: Aggregator.ReadyRequest | SimWorker.ReadyRequest,
	) {
		this.worker = worker;
		this.ready = new Promise((resolve, reject) => {
			worker.onmessage = (ev) => {
				// the aggregator's and the sim workers' Ready and Failed responses are the same
				switch (ev.data.type) {
					case SimWorker.Response.Ready:
						resolve();
						return;
					case SimWorker.Response.Failed:
						if (ev.data.run == null) {
							this.failed = true;
							worker.terminate();
							reject((ev.data as SimWorker.FailedResponse).reason);
							return;
						}
				}
				this.handler(ev);
			};
		});
		// a run waiting on it handles the failure; don't report it as unhandled otherwise
		this.ready.catch(() => {});
		worker.postMessage(readyRequest);
	}

	// Sends the instance's responses to handler, in place of the previous run's.
	public listen(handler: (ev: MessageEvent) => void) {
		this.handler = handler;
	}

	public post(message: unknown, transfer: Transferable[] = []) {
		this.worker.postMessage(message, transfer);
	}

	public terminate() {
		this.worker.terminate();
	}
}

class HelperExecutor {
	private wasmPath: string;
	private helper: Worker | undefined;
	private module: Promise<WebAssembly.Module> | undefined;
	private pending = new Map<number, (event: MessageEvent) => void>();
	private id = 0;

	constructor(wasm: string) {
		this.wasmPath = wasm;
	}

	// Starts the helper if needed and resolves to the compiled module. The helper fetches and
	// compiles the wasm and sends the module back, so the aggregator and the sim workers
	// instantiate it instead of each compiling their own copy. It is compiled in a worker
	// because the page's CSP (script-src without 'wasm-unsafe-eval') does not allow compiling
	// wasm on the main thread.
	public load(): Promise<WebAssembly.Module> {
		if (this.module != null) {
			return this.module;
		}

		const helper = new Worker(new URL("./Workers/helper.ts", import.meta.url));
		this.helper = helper;
		let resolveModule: (module: WebAssembly.Module) => void = () => {};
		let rejectModule: (reason: string) => void = () => {};
		const module = new Promise<WebAssembly.Module>((resolve, reject) => {
			resolveModule = resolve;
			rejectModule = reject;
		});
		// callers that only validate or sample get the failure through their own request
		module.catch(() => {});
		this.module = module;

		helper.postMessage(Helper.ReadyRequest(this.wasmPath));
		helper.onmessage = (ev) => {
			if (ev.data.type === Helper.Response.Ready) {
				resolveModule((ev.data as Helper.ReadyResponse).module);
				return;
			}
			if (ev.data.type === Helper.Response.Failed && ev.data.id == null) {
				// Loading failed. Fail whatever is still waiting on this helper and start over
				// with a new one on the next call.
				const reason = (ev.data as Helper.FailedResponse).reason;
				rejectModule(reason);
				helper.terminate();
				this.helper = undefined;
				this.module = undefined;
				for (const [id, handleResponse] of this.pending) {
					handleResponse(
						new MessageEvent("message", {
							data: Helper.FailedResponse(id, reason),
						}),
					);
				}
				this.pending.clear();
				return;
			}

			const handleResponse = this.pending.get(ev.data.id);
			if (handleResponse == null) {
				console.error("helper - response without a pending request: ", ev.data);
				return;
			}
			this.pending.delete(ev.data.id);
			handleResponse(ev);
		};
		return module;
	}

	private requestId() {
		return this.id++;
	}

	public validate(cfg: string): Promise<ParsedResult> {
		this.load();

		const id = this.requestId();
		return new Promise((resolve, reject) => {
			function handleResponse(event: MessageEvent) {
				switch (event.data.type as Helper.Response) {
					case Helper.Response.Validate:
						resolve((event.data as Helper.ValidateResponse).cfg);
						return;
					case Helper.Response.Failed:
						reject((event.data as Helper.FailedResponse).reason);
						return;
					default:
						reject("unknown validate response: " + event.data.type);
				}
			}
			this.pending.set(id, handleResponse);
			this.helper?.postMessage(Helper.ValidateRequest(id, cfg));
		});
	}

	public sample(cfg: string, seed: string): Promise<Sample> {
		this.load();
		const id = this.requestId();

		return new Promise((resolve, reject) => {
			function handleResponse(event: MessageEvent) {
				switch (event.data.type as Helper.Response) {
					case Helper.Response.Sample:
						resolve((event.data as Helper.SampleResponse).sample);
						return;
					case Helper.Response.Failed:
						reject((event.data as Helper.FailedResponse).reason);
						return;
					default:
						console.log(event.data);
						reject("unknown sample response: " + event.data.type);
				}
			}
			this.pending.set(id, handleResponse);
			this.helper?.postMessage(Helper.SampleRequest(id, cfg, seed));
		});
	}

	public buildInfo(): { hash: string; date: string } {
		throw new Error("Method not implemented.");
	}
}
