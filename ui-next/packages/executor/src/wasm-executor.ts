import type { Sim } from "@gcsim/types";
import type { Executor } from "./executor.js";
import { throttle } from "./throttle.js";
import { Aggregator, Helper, SimWorker } from "./workers/common.js";

const VIEWER_THROTTLE = 100;
const MIN_WORKERS = 1;
const MAX_WORKERS = 30;
const DEFAULT_WORKERS = 3;

export class WasmExecutor implements Executor {
  private wasmPath: string;
  private helper: HelperExecutor;
  private aggregator: Worker | null;
  private workers: Worker[];
  private workerCount: number;
  private isRunning: boolean;
  private runStarted: number;

  constructor(wasm: string) {
    this.wasmPath = wasm;
    this.helper = new HelperExecutor(wasm);
    this.aggregator = null;
    this.workers = [];
    this.workerCount = DEFAULT_WORKERS;
    this.isRunning = false;
    this.runStarted = 0;
  }

  public ready(): Promise<boolean> {
    return Promise.resolve(!this.isRunning);
  }

  public running(): boolean {
    return this.isRunning;
  }

  public setWorkerCount(count: number): void {
    this.workerCount = Math.max(MIN_WORKERS, Math.min(MAX_WORKERS, count));
  }

  public getWorkerCount(): number {
    return this.workerCount;
  }

  private createAggregator(): Promise<boolean> {
    return new Promise((resolve, reject) => {
      if (this.aggregator) {
        resolve(true);
        return;
      }

      this.aggregator = new Worker(new URL("./workers/aggregator.ts", import.meta.url));
      this.aggregator.postMessage(Aggregator.ReadyRequest(this.wasmPath));
      this.aggregator.onmessage = (ev) => {
        switch (ev.data.type as Aggregator.Response) {
          case Aggregator.Response.Ready:
            resolve(true);
            return;
          case Aggregator.Response.Failed:
            reject((ev.data as Aggregator.FailedResponse).reason);
            return;
        }
      };
    });
  }

  private createWorkers(): Promise<boolean> {
    const diff = this.workerCount - this.workers.length;

    if (diff < 0) {
      for (const w of this.workers.splice(diff)) {
        w.terminate();
      }
      return Promise.resolve(true);
    }

    const promises: Promise<boolean>[] = [];
    for (let i = 0; i < diff; i++) {
      promises.push(
        new Promise<boolean>((resolve, reject) => {
          const worker = new Worker(new URL("./workers/worker.ts", import.meta.url));
          worker.postMessage(SimWorker.ReadyRequest(this.wasmPath));

          const idx = this.workers.push(worker) - 1;
          worker.onmessage = (ev) => {
            switch (ev.data.type as SimWorker.Response) {
              case SimWorker.Response.Ready:
                resolve(true);
                return;
              case SimWorker.Response.Failed:
                reject(`Worker ${idx} ${(ev.data as SimWorker.FailedResponse).reason}`);
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
    updateResult: (result: Sim.SimResults, hash: string) => void,
  ): Promise<boolean | void> {
    this.isRunning = true;
    this.runStarted = performance.now();

    // 1. Create Aggregator & Workers
    const created = Promise.all([this.createAggregator(), this.createWorkers()]);

    let result: Sim.SimResults | null = null;
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
            switch (ev.data.type as Aggregator.Response) {
              case Aggregator.Response.Initialized:
                result = (ev.data as Aggregator.InitializeResponse).result as Sim.SimResults | null;
                maxIterations = (result?.simulator_settings?.iterations as number) ?? 1000;
                resolve(true);
                return;
              case Aggregator.Response.Failed:
                reject((ev.data as Aggregator.FailedResponse).reason);
                return;
            }
          };
          this.aggregator.postMessage(Aggregator.InitializeRequest(cfg));
        }),
      );

      // initialize workers
      for (const worker of this.workers) {
        promises.push(
          new Promise<boolean>((resolve, reject) => {
            worker.onmessage = (ev) => {
              switch (ev.data.type as SimWorker.Response) {
                case SimWorker.Response.Initialized:
                  resolve(true);
                  return;
                case SimWorker.Response.Failed:
                  reject((ev.data as SimWorker.FailedResponse).reason);
                  return;
              }
            };
            worker.postMessage(SimWorker.InitializeRequest(cfg));
          }),
        );
      }

      return Promise.all(promises);
    });

    const throttledFlush = throttle(() => {
      if (this.isRunning) {
        this.aggregator?.postMessage(Aggregator.FlushRequest());
      }
    }, VIEWER_THROTTLE);

    // 3. Start execution
    return initialized.then(() => {
      return new Promise((resolve, reject) => {
        if (this.aggregator == null) {
          reject("Aggregator is null!");
          return;
        }

        let completed = 0;
        this.aggregator.onmessage = (ev) => {
          switch (ev.data.type as Aggregator.Response) {
            case Aggregator.Response.Result: {
              const { hash, stats } = (ev.data as Aggregator.ResultResponse).result;

              const out = Object.assign({}, result);
              out.statistics = stats as Sim.Statistics | undefined;
              updateResult(out, hash);

              if (completed >= maxIterations) {
                this.isRunning = false;
                resolve(true);
                if (this.runStarted > 0) {
                  this.runStarted = 0;
                }
              }
              return;
            }
            case Aggregator.Response.Done:
              completed += 1;
              throttledFlush();
              return;
            case Aggregator.Response.Failed:
              // A flush may happen after a cancel request. When this happens,
              // the existing aggregator has no data and fails to flush.
              if (this.isRunning) {
                reject((ev.data as Aggregator.FailedResponse).reason);
              }
          }
        };

        let requested = 0;
        for (const worker of this.workers) {
          worker.onmessage = (ev) => {
            switch (ev.data.type as SimWorker.Response) {
              case SimWorker.Response.Done: {
                const resp = ev.data as SimWorker.RunResponse;
                this.aggregator?.postMessage(Aggregator.AddRequest(resp.result));
                if (requested < maxIterations) {
                  worker.postMessage(SimWorker.RunRequest(requested++));
                }
                return;
              }
              case SimWorker.Response.Failed:
                reject((ev.data as Aggregator.FailedResponse).reason);
            }
          };

          if (requested < maxIterations) {
            worker.postMessage(SimWorker.RunRequest(requested++));
          }
        }
      });
    });
  }

  public cancel(): void {
    if (!this.isRunning || this.aggregator == null) {
      return;
    }

    this.isRunning = false;
    for (const worker of this.workers) {
      worker.onmessage = null;
    }

    // Recreate the aggregator since there is no way to clear the worker message queue.
    // Any pending AddRequests will still be processed otherwise.
    this.aggregator.terminate();
    this.aggregator = null;

    this.runStarted = 0;
  }

  public validate(cfg: string): Promise<Sim.ParsedResult> {
    return this.helper.validate(cfg);
  }

  public sample(cfg: string, seed: string): Promise<Sim.Sample> {
    return this.helper.sample(cfg, seed);
  }

  public buildInfo(): { hash: string; date: string } {
    return this.helper.buildInfo();
  }
}

class HelperExecutor {
  private wasmPath: string;
  private helper: Worker | undefined;
  private responses = new Map<number, MessageEvent>();
  private nextId = 0;

  constructor(wasm: string) {
    this.wasmPath = wasm;
  }

  private initialize(): void {
    if (this.helper != null) {
      return;
    }

    this.helper = new Worker(new URL("./workers/helper.ts", import.meta.url));
    this.helper.postMessage(Helper.ReadyRequest(this.wasmPath));
    this.helper.onmessage = (ev) => {
      this.responses.set(ev.data.id, ev);
    };
  }

  private requestId(): number {
    return this.nextId++;
  }

  private waitForResponse(id: number, cb: (event: MessageEvent) => void): void {
    const event = this.responses.get(id);
    if (event != null) {
      cb(event);
      this.responses.delete(id);
      return;
    }
    setTimeout(() => this.waitForResponse(id, cb), 100);
  }

  public validate(cfg: string): Promise<Sim.ParsedResult> {
    this.initialize();

    const id = this.requestId();
    return new Promise((resolve, reject) => {
      const handleResponse = (event: MessageEvent) => {
        switch (event.data.type as Helper.Response) {
          case Helper.Response.Validate:
            resolve((event.data as Helper.ValidateResponse).cfg as Sim.ParsedResult);
            return;
          case Helper.Response.Failed:
            reject((event.data as Helper.FailedResponse).reason);
            return;
          default:
            reject(`unknown validate response: ${event.data.type}`);
        }
      };
      this.waitForResponse(id, handleResponse);
      this.helper?.postMessage(Helper.ValidateRequest(id, cfg));
    });
  }

  public sample(cfg: string, seed: string): Promise<Sim.Sample> {
    this.initialize();
    const id = this.requestId();

    return new Promise((resolve, reject) => {
      const handleResponse = (event: MessageEvent) => {
        switch (event.data.type as Helper.Response) {
          case Helper.Response.Sample:
            resolve((event.data as Helper.SampleResponse).sample as Sim.Sample);
            return;
          case Helper.Response.Failed:
            reject((event.data as Helper.FailedResponse).reason);
            return;
          default:
            reject(`unknown sample response: ${event.data.type}`);
        }
      };
      this.waitForResponse(id, handleResponse);
      this.helper?.postMessage(Helper.SampleRequest(id, cfg, seed));
    });
  }

  public buildInfo(): { hash: string; date: string } {
    throw new Error("Method not implemented.");
  }
}
