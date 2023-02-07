import { ParsedResult, Sample, SimResults } from "@gcsim/types";
import { throttle } from "lodash-es";
import { Executor } from "./Executor";
import { Aggregator, Helper, SimWorker } from "./Workers/common";

const VIEWER_THROTTLE = 100;

export class WasmExecutor implements Executor {
  private wasmPath: string;
  private helper: HelperExecutor;
  private aggregator: Worker | null;
  private workers: Worker[];
  private workerCount: number;
  private isRunning: boolean;

  constructor(wasm: string) {
    this.wasmPath = wasm;
    this.helper = new HelperExecutor(wasm);

    this.aggregator = null;
    this.workers = [];
    this.workerCount = 3;
    this.isRunning = false;
  }

  public ready(): boolean {
    return !this.isRunning;
  }

  public running(): boolean {
    return this.isRunning;
  }

  public setWorkerCount(count: number) {
    this.workerCount = count;
  }

  private createAggregator(): Promise<boolean> {
    return new Promise((resolve, reject) => {
      if (this.aggregator) {
        resolve(true);
        return;
      }

      this.aggregator = new Worker(new URL("./Workers/aggregator.ts", import.meta.url));
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
    console.log("loading workers", this.workerCount, this);
    const diff = this.workerCount - this.workers.length;

    if (diff < 0) {
      this.workers.splice(diff).forEach((w) => w.terminate());
      return Promise.resolve(true);
    }

    console.log("loading " + diff + " workers");
    const promises: Promise<boolean>[] = [];
    for (let i = 0; i < diff; i++) {
      promises.push(new Promise<boolean>((resolve, reject) => {
        const worker = new Worker(new URL("./Workers/worker.ts", import.meta.url));
        worker.postMessage(SimWorker.ReadyRequest(this.wasmPath));

        const idx = this.workers.push(worker) - 1;
        worker.onmessage = (ev) => {
          switch (ev.data.type as SimWorker.Response) {
            case SimWorker.Response.Ready:
              resolve(true);
              return;
            case SimWorker.Response.Failed:
              reject("Worker " + idx + " " + (ev.data as SimWorker.FailedResponse).reason);
              return;
          }
        };
      }));
    }
    return Promise.all(promises).then(() => true);
  }

  public run(
        cfg: string, updateResult: (result: SimResults, hash: string) => void
      ): Promise<boolean | void> {
    this.isRunning = true;

    // 1. Create Aggregator & Workers
    const created = Promise.all([this.createAggregator(), this.createWorkers()]);

    let result: SimResults | null = null;
    let maxIterations = 0;

    // 2. Initialize Aggregator & Workers
    const initialized = created.then(() => {
      const promises: Promise<boolean>[] = [];

      // initialize aggregator
      promises.push(new Promise<boolean>((resolve, reject) => {
        if (this.aggregator == null) {
          reject("Aggregator is null!");
          return;
        }
  
        this.aggregator.onmessage = (ev) => {
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
        this.aggregator.postMessage(Aggregator.InitializeRequest(cfg));
      }));

      // initialize workers
      this.workers.forEach((worker) => {
        promises.push(new Promise<boolean>((resolve, reject) => {
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
        }));
      });

      return Promise.all(promises);
    });

    const throttledFlush = throttle(() => {
      if (this.isRunning) {
        this.aggregator?.postMessage(Aggregator.FlushRequest());
      }
    }, VIEWER_THROTTLE, { leading: true, trailing: true });

    // 3. start execution
    return initialized.then(() => {
      if (this.aggregator == null) {
        return Promise.reject("Aggregator is null!");
      }

      let completed = 0;
      this.aggregator.onmessage = (ev) => {
        switch (ev.data.type as Aggregator.Response) {
          case Aggregator.Response.Result:
            const { hash, stats } = (ev.data as Aggregator.ResultResponse).result;

            const out = Object.assign({}, result);
            out.statistics = stats;
            updateResult(out, hash);

            if (completed >= maxIterations) {
              this.isRunning = false;
              return Promise.resolve(true);
            }
            return;
          case Aggregator.Response.Done:
            completed += 1;
            throttledFlush();
            return;
          case Aggregator.Response.Failed:
            // TODO: bug with throttled flush where a flush may happen after a cancel request.
            //    When this happens, the existing aggregator has no data and fails to flush.
            //    this doesnt cause any problems (yet) and just produces an error in console.
            if (this.isRunning) {
              return Promise.reject((ev.data as Aggregator.FailedResponse).reason);
            }
        }
      };

      let requested = 0;
      this.workers.forEach((worker) => {
        worker.onmessage = (ev) => {
          switch (ev.data.type as SimWorker.Response) {
            case SimWorker.Response.Done:
              const resp: SimWorker.RunResponse = ev.data;
              this.aggregator?.postMessage(Aggregator.AddRequest(resp.result));
              if (requested < maxIterations) {
                worker.postMessage(SimWorker.RunRequest(requested++));
              }
              return;
            case SimWorker.Response.Failed:
              return Promise.reject((ev.data as Aggregator.FailedResponse).reason);
          }
        };

        if (requested < maxIterations) {
          worker.postMessage(SimWorker.RunRequest(requested++));
        }
      });
    });
  }

  public cancel(): void {
    if (!this.isRunning || this.aggregator == null) {
      return;
    }

    this.isRunning = false;
    console.log("execution canceled");
    this.workers.forEach((worker) => {
      worker.onmessage = null;
    });

    // It is possible that there are N AddRequests in the aggregator queue that we have no control
    // over. Even if we set the onmessage here to null, the aggregator will still process through
    // all N requests. Since there is no way to clear the worker queue, recreating the worker is the
    // next best thing.
    //
    // Downside of this approach is any memory allocation/optimizations from previous runs will not
    // carry over, making executions after a cancel "less optimal".
    this.aggregator.terminate();
    this.aggregator = null;
  }

  public validate(cfg: string): Promise<ParsedResult> {
    return this.helper.validate(cfg);
  }

  public sample(cfg: string, seed: string): Promise<Sample> {
    return this.helper.sample(cfg, seed);
  }

  public buildInfo(): { hash: string; date: string; } {
    return this.helper.buildInfo();
  }
}

class HelperExecutor {
  private wasmPath: string;
  private helper: Worker | undefined;
  private responses = new Map<number, MessageEvent>;
  private id = 0;

  constructor(wasm: string) {
    this.wasmPath = wasm;
  }

  private initialize() {
    if (this.helper != null) {
      return;
    }

    this.helper = new Worker(new URL("./Workers/helper.ts", import.meta.url));
    this.helper.postMessage(Helper.ReadyRequest(this.wasmPath));
    this.helper.onmessage = (ev) => {
      this.responses.set(ev.data.id, ev);
    };
  }

  private requestId() {
    return this.id++;
  }

  private waitForResponse(id: number, cb: (event: MessageEvent) => void) {
    const event = this.responses.get(id);
    if (event != null) {
      cb(event);
      this.responses.delete(id);
      return;
    }
    setTimeout(() => this.waitForResponse(id, cb), 100);
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
      this.waitForResponse(id, handleResponse);
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
      this.waitForResponse(id, handleResponse);
      this.helper?.postMessage(Helper.SampleRequest(id, cfg, seed));
    });
  }

  public buildInfo(): { hash: string; date: string; } {
    throw new Error("Method not implemented.");
  }
}