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
  _max: number;
  _active: number;

  constructor() {
    this._workers = [];
    this._avail = [];
    this._queue = [];
    this._max = 5; //default 5 workers
    this._active = 0; //no active worker
    this._loaded = 0;
  }

  load(count: number, readycb: (count: number) => void) {
    if (count === 0) {
      return;
    }

    //load one first, once first is ready, then load the rest
    console.log("start loading workers");
    let start = this._workers.length;
    let i = start;

    const loading = new Promise((resolve) => {
      const w = new Worker(new URL("worker.ts", import.meta.url));
      this._workers.push(w);
      this._avail.push(false);

      let x = i;
      w.onmessage = (ev) => {
        if (ev.data === "ready") {
          this._avail[x] = true;
          this._loaded++;
          // console.log("worker 0 is now ready");
          //we're technically ready to work as long as just one worker is ready
          readycb(this._loaded);
          //start the chain
          resolve("ok");
        }
      };
    });
    loading.then(() => {
      for (; i < count + start; i++) {
        const w = new Worker(new URL("worker.ts", import.meta.url));
        this._workers.push(w);
        this._avail.push(false);

        let x = i;
        w.onmessage = (ev) => {
          if (ev.data === "ready") {
            this._avail[x] = true;
            this._loaded++;
            readycb(this._loaded);
            // console.log("worker " + i + " is now ready");
          }
        };
      }
    });
  }

  setCfg(cfg: string) {
    for (let i = 0; i < this._workers.length; i++) {
      this._workers[i].postMessage({
        cmd: "cfg",
        payload: cfg,
      });
    }
  }

  setMaxWorker(count: number) {
    this._max = count; //if you set count to more than loaded nothing will happen
  }

  queue(t: Task) {
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
    // check to make sure active does not exceed max
    if (this._active === this._max) {
      // console.log(
      //   "max worker count reached, active: " +
      //     this._active +
      //     " max: " +
      //     this._max
      // );
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
    w.postMessage({
      cmd: task.cmd,
      payload: task.payload ? task.payload : "",
    });
    this._active++;
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
  }
}
