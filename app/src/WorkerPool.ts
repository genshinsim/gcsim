import { pool } from "./Pages/Sim";

export type Task = {
  cmd: string;
  payload?: string;
  cb: (value: any) => void;
};

export class WorkerPool {
  _workers: Worker[];
  _avail: boolean[];
  _queue: Task[];
  _loaded: number;
  _active: number;

  constructor() {
    this._workers = [];
    this._avail = [];
    this._queue = [];
    this._active = 0; //no active worker
    this._loaded = 0;
  }

  count(): number {
    return this._workers.length;
  }

  setWorkerCount(count: number, readycb: (count: number) => void) {
    console.log("loading workers", count, this);
    const extras = count - this._workers.length;

    if (extras === 0) {
      //do nothing
      return;
    }
    if (extras < 0) {
      //truncate extra workers

      this._workers.splice(extras);
      this._avail.splice(extras);
      this._loaded = count;

      console.log(pool);
      readycb(count);
      return;
    }

    //load one first, once first is ready, then load the rest
    console.log("start loading workers");
    let start = this._workers.length;
    let i = start;

    const loading = () =>
      new Promise((resolve) => {
        const w = new Worker(new URL("worker.ts", import.meta.url));
        this._workers.push(w);
        this._avail.push(false);

        const x = i;
        i++;
        w.onmessage = (ev) => {
          if (ev.data === "ready") {
            this._avail[x] = true;
            this._loaded++;
            console.log("worker " + x + " is now ready");
            //we're technically ready to work as long as just one worker is ready
            readycb(this._loaded);
            //start the chain
            resolve("ok");
          }
        };
      });
    loading().then(() => {
      for (; i < count; i++) {
        const w = new Worker(new URL("worker.ts", import.meta.url));
        this._workers.push(w);
        this._avail.push(false);

        let x = i;
        w.onmessage = (ev) => {
          if (ev.data === "ready") {
            this._avail[x] = true;
            this._loaded++;
            readycb(this._loaded);
            console.log("worker " + x + " is now ready");
          }
        };
      }
    });
  }

  setCfg(cfg: string, cb: (val: string) => void) {
    let count = 0;
    const max = this._workers.length;
    for (let i = 0; i < max; i++) {
      const t = i;
      this._workers[t].onmessage = (ev) => {
        count++;
        // we need to count how many are ok...
        if (count === max) {
          console.log("all configs loaded");
          console.log(pool);
          cb(ev.data);
        }
      };
      this._workers[t].postMessage({
        cmd: "cfg",
        payload: cfg,
      });
    }
  }

  queue(t: Task) {
    // console.log("got a task: ", t);
    //add it to queue
    this._queue.push(t);

    //try popping
    this.pop();
  }

  private pop() {
    // console.log("looking for worker to do work: ", this);
    if (this._queue.length == 0) {
      return;
    }
    //find free worker
    let ind = -1;
    for (let i = 0; i < this._avail.length; i++) {
      if (this._avail[i]) {
        ind = i;
        break;
      }
    }

    if (ind === -1) {
      return;
    }

    //pop from slice
    const task = this._queue[0];
    this._queue = this._queue.slice(1, this._queue.length);
    this.run(ind, task);
  }

  //ask current worker to run a task
  private run(workerIndex: number, task: Task) {
    this._avail[workerIndex] = false;
    let w = this._workers[workerIndex];
    w.onmessage = (ev) => {
      task.cb(ev.data);
      this._avail[workerIndex] = true;
      this._active--; //done work so subtract one from active
      if (this._active < 0) {
        console.warn("unexpected active worker count < 0");
        this._active = 0; //sanity check
      }
      // console.log("worker done: ", this);
      //try popping maybe there's more
      this.pop();
    };
    w.postMessage({
      cmd: task.cmd,
      payload: task.payload ? task.payload : "",
    });
    this._active++;
  }
}
