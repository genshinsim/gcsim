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

export class WasmExecutor implements Executor {
	private helper: HelperExecutor;
	private aggregator: Worker | null;
	private workers: Worker[];
	private workerCount: number;
	private isRunning: boolean;
	private runId: number;
	private runStarted: number;

	constructor(wasm: string) {
		this.helper = new HelperExecutor(wasm);

		this.aggregator = null;
		this.workers = [];
		this.workerCount = 3;
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

	private createAggregator(module: WebAssembly.Module): Promise<boolean> {
		return new Promise((resolve, reject) => {
			if (this.aggregator) {
				resolve(true);
				return;
			}

			const aggregator = new Worker(
				new URL("./Workers/aggregator.ts", import.meta.url),
			);
			this.aggregator = aggregator;
			aggregator.postMessage(Aggregator.ReadyRequest(module));
			aggregator.onmessage = (ev) => {
				switch (ev.data.type as Aggregator.Response) {
					case Aggregator.Response.Ready:
						resolve(true);
						return;
					case Aggregator.Response.Failed:
						// the aggregator is reused across runs; don't keep a broken one
						aggregator.terminate();
						this.aggregator = null;
						reject((ev.data as Aggregator.FailedResponse).reason);
						return;
				}
			};
		});
	}

	private createWorkers(module: WebAssembly.Module): Promise<boolean> {
		console.log("loading workers", this.workerCount, this);
		const diff = this.workerCount - this.workers.length;

		if (diff < 0) {
			this.workers.splice(diff).forEach((w) => {
				w.terminate();
			});
			return Promise.resolve(true);
		}

		console.log("loading " + diff + " workers");
		const promises: Promise<boolean>[] = [];
		for (let i = 0; i < diff; i++) {
			promises.push(
				new Promise<boolean>((resolve, reject) => {
					const worker = new Worker(
						new URL("./Workers/worker.ts", import.meta.url),
					);
					worker.postMessage(SimWorker.ReadyRequest(module));

					const idx = this.workers.push(worker) - 1;
					worker.onmessage = (ev) => {
						switch (ev.data.type as SimWorker.Response) {
							case SimWorker.Response.Ready:
								resolve(true);
								return;
							case SimWorker.Response.Failed:
								// drop it so the next run doesn't wait on a worker that never loaded
								this.workers = this.workers.filter((w) => w !== worker);
								worker.terminate();
								reject(
									"Worker " +
										idx +
										" " +
										(ev.data as SimWorker.FailedResponse).reason,
								);
								return;
						}
					};
				}),
			);
		}
		return Promise.all(promises).then(() => true);
	}

	public run(
		cfg: string,
		updateResult: (result: model.SimulationResult, hash: string) => void,
	): Promise<boolean | void> {
		this.isRunning = true;
		this.runStarted = performance.now();

		// The workers and the aggregator are reused across runs, so their queues can still hold
		// requests from a cancelled run. Every request carries the run id and every response
		// echoes it; responses from any other run are dropped.
		const run = ++this.runId;
		const isCurrent = (ev: MessageEvent) =>
			ev.data.run === run && this.runId === run;
		const stop = () => {
			if (this.runId === run) {
				this.runId++;
				this.isRunning = false;
			}
		};

		// 1. Create Aggregator & Workers
		const created = this.helper
			.compile()
			.then((module) =>
				Promise.all([
					this.createAggregator(module),
					this.createWorkers(module),
				]),
			);

		let result: model.SimulationResult | null = null;
		let maxIterations = 0;

		// 2. Initialize Aggregator & Workers
		const initialized = created.then(() => {
			const promises: Promise<boolean>[] = [];

			// initialize aggregator
			promises.push(
				new Promise<boolean>((resolve, reject) => {
					if (this.aggregator == null) {
						reject("Aggregator is null!");
						return;
					}

					this.aggregator.onmessage = (ev) => {
						if (!isCurrent(ev)) {
							return;
						}
						switch (ev.data.type as Aggregator.Response) {
							case Aggregator.Response.Initialized:
								result = (ev.data as Aggregator.InitializeResponse).result;
								maxIterations = result?.simulator_settings?.iterations ?? 1000;
								resolve(true);
								return;
							case Aggregator.Response.Failed:
								reject((ev.data as Aggregator.FailedResponse).reason);
								return;
						}
					};
					this.aggregator.postMessage(Aggregator.InitializeRequest(run, cfg));
				}),
			);

			// initialize workers
			this.workers.forEach((worker) => {
				promises.push(
					new Promise<boolean>((resolve, reject) => {
						worker.onmessage = (ev) => {
							if (!isCurrent(ev)) {
								return;
							}
							switch (ev.data.type as SimWorker.Response) {
								case SimWorker.Response.Initialized:
									resolve(true);
									return;
								case SimWorker.Response.Failed:
									reject((ev.data as SimWorker.FailedResponse).reason);
									return;
							}
						};
						worker.postMessage(SimWorker.InitializeRequest(run, cfg));
					}),
				);
			});

			return Promise.all(promises);
		});

		// 3. start execution
		const executed = initialized.then(() => {
			return new Promise<boolean>((resolve, reject) => {
				if (this.aggregator == null) {
					reject("Aggregator is null!");
					return;
				}
				const aggregator = this.aggregator;
				let completed = 0;
				let flushing = false; // an intermediate flush is outstanding
				let nextFlush = 0;
				aggregator.onmessage = (ev) => {
					if (!isCurrent(ev)) {
						return;
					}
					switch (ev.data.type as Aggregator.Response) {
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
								aggregator.postMessage(Aggregator.FlushRequest(run, true));
							} else if (!flushing && performance.now() >= nextFlush) {
								flushing = true;
								aggregator.postMessage(Aggregator.FlushRequest(run, false));
							}
							dispatch();
							return;
						case Aggregator.Response.Failed:
							reject((ev.data as Aggregator.FailedResponse).reason);
					}
				};

				const workers = this.workers;
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
							worker.postMessage(SimWorker.RunRequest(run, requested++));
							queued[i]++;
						}
					});
				};

				workers.forEach((worker, i) => {
					worker.onmessage = (ev) => {
						if (!isCurrent(ev)) {
							return;
						}
						switch (ev.data.type as SimWorker.Response) {
							case SimWorker.Response.Done: {
								const resp: SimWorker.RunResponse = ev.data;
								queued[i]--;
								// transfer rather than copy; the result is only forwarded
								this.aggregator?.postMessage(
									Aggregator.AddRequest(run, resp.result),
									[resp.result.buffer],
								);
								dispatch();
								return;
							}
							case SimWorker.Response.Failed:
								reject((ev.data as Aggregator.FailedResponse).reason);
						}
					};
				});
				dispatch();
			});
		});

		// a failed run frees the executor for the next one
		return executed.catch((e) => {
			stop();
			throw e;
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

class HelperExecutor {
	private wasmPath: string;
	private helper: Worker | undefined;
	private module: Promise<WebAssembly.Module> | undefined;
	private pending = new Map<number, (event: MessageEvent) => void>();
	private id = 0;

	constructor(wasm: string) {
		this.wasmPath = wasm;
	}

	// The helper fetches and compiles the wasm and sends the compiled module back, so the
	// aggregator and the sim workers instantiate it instead of each compiling their own copy.
	// It is compiled in a worker because the page's CSP (script-src without
	// 'wasm-unsafe-eval') does not allow compiling wasm on the main thread.
	public compile(): Promise<WebAssembly.Module> {
		return this.initialize();
	}

	private initialize(): Promise<WebAssembly.Module> {
		if (this.helper != null && this.module != null) {
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
		this.initialize();

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
		this.initialize();
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
