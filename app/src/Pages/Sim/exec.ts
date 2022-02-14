import { AppThunk } from "~src/store";
import { pool, simActions } from ".";
import { Result } from "~src/types";

const iterRegex = /iteration=(\d+)/;

function aggregateResults(results: Result[]) {
  console.log(results[0]);
  let s = JSON.stringify(results);
  pool.queue({
    cmd: "collect",
    payload: s,
    cb: (val) => {
      //convert it back
      const res = JSON.parse(val);
      console.log(res);
    },
  });
}

function extractItersFromConfig(cfg: string): number {
  let iters = 1;
  let m = iterRegex.exec(cfg);
  console.log(m);

  if (m) {
    iters = parseInt(m[1]);

    if (isNaN(iters)) {
      console.warn("no iteration found in settings: ", m);
      iters = 1000;
    }
  } else {
    console.log(cfg);
    console.warn("iter regex failed");
  }
  //   console.log("parsed iters: " + iters);

  return iters;
}

export function runSim(): AppThunk {
  return function (dispatch, getState) {
    let cfg = getState().sim.cfg;

    //extract the number of iterations from the config file
    const iters = extractItersFromConfig(cfg);

    //run the sim
    dispatch(
      simActions.setRunStats({
        progress: 0,
        result: -1,
        time: -1,
      })
    );
    let queued = 0;
    let done = 0;
    let avg = 0;
    let progress = 0;
    let results: Result[] = [];

    const startTime = window.performance.now();
    console.time("sim");
    pool.setCfg(cfg);

    const cbFunc = (val: any) => {
      //parse the result
      const res = JSON.parse(val);
      results.push(res);
      avg += res.dps;

      done++;

      if (done === iters) {
        console.timeEnd("sim");
        // setRuntime(end - startTime);
        avg = avg / iters;
        aggregateResults(results);
        const end = window.performance.now();

        dispatch(
          simActions.setRunStats({
            progress: -1,
            result: avg,
            time: end - startTime,
          })
        );
        //stop call back chain if done
        return;
      }

      //update progress in increments of 5%
      const per = Math.floor((20 * done) / iters);
      if (per > progress) {
        dispatch(
          simActions.setRunStats({
            progress: per,
            result: -1,
            time: -1,
          })
        );
        progress = per;
      }

      //queue more if we haven't queued everything yet
      if (queued < iters) {
        //queue another worker
        queued++;
        pool.queue({ cmd: "run", cb: cbFunc });
      }
    };
    const debugCB = (val: any) => {
      const res = JSON.parse(val);
      if (res.err) {
        console.error(res.err);
        return;
      }
      console.log("finish debug run: ", res);
    };

    pool.queue({ cmd: "debug", cb: debugCB });

    //queue up 20 out of iters
    let count = 20;
    if (count > iters) {
      count = iters;
    }
    for (; queued < count; queued++) {
      pool.queue({ cmd: "run", cb: cbFunc });
    }
  };
}
