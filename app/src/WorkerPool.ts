import { pool } from "./Pages/Sim";
import { ParsedResult } from "./types";
import { ResultsSummary } from "./Types/stats";
import { Aggregator, SimWorker } from "./Workers/common";

export type Task = {
  cmd: string;
  payload?: string;
  cb: (value: any) => void;
};

export class WorkerPool {
  private aggregator: Worker;
  private aggregatorReady: boolean;
  private workers: Worker[];
  private workersReady: boolean[];

  constructor() {
    this.aggregatorReady = false;
    this.aggregator = new Worker(new URL("./Workers/aggregator.ts", import.meta.url));
    this.aggregator.onmessage = (ev) => {
      switch (ev.data.type as Aggregator.Response) {
        case Aggregator.Response.Ready:
          this.aggregatorReady = true;
          break;
      }
    };

    this.workers = [];
    this.workersReady = [];
  }

  public count(): number {
    return this.loaded().length;
  }

  public ready(): boolean {
    return this.aggregatorReady && this.count() > 0;
  }

  private loaded(): Worker[] {
    return this.workers.filter((_, i) => this.workersReady[i]);
  }

  private createWorker(): Promise<number> {
    return new Promise((resolve, reject) => {
      const worker = new Worker(new URL("./Workers/worker.ts", import.meta.url));
      const idx = this.workers.push(worker) - 1;
      this.workersReady.push(false);
      worker.onmessage = (ev) => {
        switch (ev.data.type as SimWorker.Response) {
          case SimWorker.Response.Ready:
            this.workersReady[idx] = true;
            resolve(idx);
            return;
          case SimWorker.Response.Failed:
            reject("Worker " + idx + " " + (ev.data as SimWorker.FailedResponse).reason);
            return;
          default:
            reject("Worker " + idx + " - unknown response: " + ev.data);
        }
      };
    });
  }

  public setWorkerCount(count: number, readycb: (count: number) => void) {
    console.log("loading workers", count, this);
    const diff = count - this.workers.length;

    if (diff < 0) {
      this.workersReady.splice(diff);
      this.workers.splice(diff).forEach((w) => w.terminate());
      console.log(pool);
      return readycb(count);
    }

    console.log("loading " + diff + " workers");
    for (let i = 0; i < diff; i++) {
      this.createWorker().then((w) => {
        console.log("worker " + w + " is now ready");
        readycb(this.count());
      });
    }
  }

  public validate(cfg: string): Promise<ParsedResult> {
    return new Promise((resolve, reject) => {
      this.aggregator.onmessage = (ev) => {
        switch (ev.data.type as Aggregator.Response) {
          case Aggregator.Response.Validate:
            resolve((ev.data as Aggregator.ValidateResponse).cfg);
            return;
          case Aggregator.Response.Failed:
            reject((ev.data as Aggregator.FailedResponse).reason);
            return;
          default:
            reject("unknown validate response: " + ev.data);
        }
      };
      this.aggregator.postMessage(Aggregator.ValidateRequest(cfg));
    });
  }

  public run(cfg: string, iterations: number, cb: (completed: number) => void): Promise<ResultsSummary> {
    if (!this.ready()) {
      return Promise.reject("aggregators and/or workers are not ready!");
    }

    return new Promise((resolve, reject) => {
      let completed = 0;
      this.aggregator.onmessage = (ev) => {
        switch (ev.data.type as Aggregator.Response) {
          case Aggregator.Response.Result:
            resolve((ev.data as Aggregator.ResultResponse).result);
            return;
          case Aggregator.Response.Done:
            cb(++completed);
            if (completed >= iterations) {
              this.aggregator.postMessage(Aggregator.FlushRequest());
            }
            break;
          case Aggregator.Response.Failed:
            reject((ev.data as Aggregator.FailedResponse).reason);
            return;
        }
      };
      this.aggregator.postMessage(Aggregator.InitializeRequest(cfg));

      let requested = 0;
      this.loaded().forEach((worker, index) => {
        worker.onmessage = (ev) => {
          switch (ev.data.type as SimWorker.Response) {
            case SimWorker.Response.Initialized:
              if (requested < iterations) {
                worker.postMessage(SimWorker.RunRequest(requested++));
              }
              break;
            case SimWorker.Response.Done:
              const resp: SimWorker.RunResponse = ev.data;
              this.aggregator.postMessage(Aggregator.AddRequest(resp.result, resp.itr));
              if (requested < iterations) {
                worker.postMessage(SimWorker.RunRequest(requested++));
              }
              break;
            case SimWorker.Response.Failed:
              reject((ev.data as Aggregator.FailedResponse).reason);
              return;
          }
        };
        worker.postMessage(SimWorker.InitializeRequest(cfg));
      });
    });
  }
}
