import { AppThunk } from "~src/store";
import { pool, simActions } from ".";
import { Result, ResultsSummary } from "~src/types";

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
    const startTime = window.performance.now();
    console.time("runSim");
    let debug: string;
    let avg = 0;
    let results: Result[] = [];
    //extract the number of iterations from the config file
    const iters = extractItersFromConfig(cfg);

    //promise for debug run
    const setConfig = () => {
      const p = new Promise((resolve, reject) => {
        pool.setCfg(cfg, (val) => {
          console.log("set config callback: " + val);
          if (val !== "ok") {
            reject(val);
          } else {
            resolve(null);
          }
        });
      });
      return p;
    };

    const debugRun = () =>
      new Promise((resolve, reject) => {
        console.time("debug");
        const debugCB = (val: any) => {
          const res = JSON.parse(val);
          console.timeEnd("debug");
          if (res.err) {
            reject(res.err);
          } else {
            debug = res;
            resolve(null);
          }
          console.log("finish debug run: ", res);
        };
        pool.queue({ cmd: "debug", cb: debugCB });
      });

    const sims = () =>
      new Promise((resolve, reject) => {
        console.time("sim");
        let queued = 0;
        let done = 0;
        let progress = 0;
        const cbFunc = (val: any) => {
          //parse the result
          let res;
          try {
            res = JSON.parse(val);
          } catch {
            console.log(val);
          }
          //stop if we hit an error
          if (res && res.err) {
            console.timeEnd("sim");
            reject(res);
            return;
          }
          results.push(res);
          avg += res.dps;

          done++;

          if (done === iters) {
            //stop call back chain if done
            console.timeEnd("sim");
            resolve(null);
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

        //queue up 20 out of iters
        let count = 20;
        if (count > iters) {
          count = iters;
        }
        for (; queued < count; queued++) {
          pool.queue({ cmd: "run", cb: cbFunc });
        }
      });

    const aggregateResults = () =>
      new Promise((resolve, reject) => {
        let s = JSON.stringify(results);
        pool.queue({
          cmd: "collect",
          payload: s,
          cb: (val) => {
            //convert it back
            const res = JSON.parse(val);

            if (res.err) {
              reject(res.err);
            } else {
              console.log(res);
              resolve(res);
            }
          },
        });
      });

    //run the sim
    dispatch(
      simActions.setRunStats({
        progress: 0,
        result: -1,
        time: -1,
      })
    );

    setConfig()
      .then(() => {
        console.log("configs done");
        return Promise.all([debugRun(), sims()]);
      })
      // .then(() => {
      //   console.log("debug done, running next");
      //   return sims();
      // })
      .then(() => {
        console.log("all iters done, collecting results");
        //aggregate the result here
        return aggregateResults();
      })
      .then((summary) => {
        console.log("do something with results here");
        const end = window.performance.now();
        dispatch(
          simActions.setRunStats({
            progress: -1,
            result: avg / iters,
            time: end - startTime,
          })
        );
        //add debug to the summary
        //@ts-ignore
        console.log(summary.dps);
        //@ts-ignore
        summary.debug = debug;

        //summary can now be passed to viewer
      })
      .catch((res) => {
        console.log(res);
        const end = window.performance.now();
        dispatch(
          simActions.setRunStats({
            progress: -1,
            result: 0,
            time: end - startTime,
          })
        );
      });

    // setRuntime(end - startTime);
  };
}
