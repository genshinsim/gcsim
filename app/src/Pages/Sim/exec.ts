import { AppThunk } from "~src/store";
import { pool, simActions } from ".";
import { Result, ResultsSummary } from "~src/types";
import { viewerActions } from "../ViewerDashboard/viewerSlice";

const iterRegex = /iteration=(\d+)/;

function extractItersFromConfig(cfg: string): number {
  let iters = 1;
  let m = iterRegex.exec(cfg);
  // console.log(m);

  if (m) {
    iters = parseInt(m[1]);

    if (isNaN(iters)) {
      console.warn("no iteration found in settings: ", m);
      iters = 1000;
    }
  } else {
    console.log(cfg);
    console.warn("iter regex failed");
    iters = 1000;
  }
  //   console.log("parsed iters: " + iters);

  return iters;
}

export function runSim(cfg: string): AppThunk {
  return function (dispatch) {
    console.log("starting run");
    // console.log(cfg);
    cfg = cfg + "\n";
    const startTime = window.performance.now();
    let debug: string;
    let avg = 0;
    let results: Result[] = [];
    let v: string;
    let bt: string;
    //extract the number of iterations from the config file
    const iters = extractItersFromConfig(cfg);

    //promise for debug run
    const setConfig = () => {
      const p = new Promise((resolve, reject) => {
        pool.setCfg(cfg, (val) => {
          console.log("set config callback: " + val);
          try {
            const res = JSON.parse(val);
            console.log(res);
            if (res.err) {
              reject(res.err);
            } else {
              resolve(null);
            }
          } catch {
            reject(val);
          }
        });
      });
      return p;
    };

    const debugRun = () =>
      new Promise<null>((resolve, reject) => {
        console.time("debug");
        const debugCB = (val: any) => {
          try {
            const res = JSON.parse(val);
            console.timeEnd("debug");
            if (res.err) {
              reject(res.err);
              return;
            }
            //it's a string otherwise
            // console.log(res);
            debug = val;
            resolve(null);
            // console.log("finish debug run: ", res);
          } catch {
            reject("unexpected error??");
          }
        };
        pool.queue({ cmd: "debug", cb: debugCB });
      });

    const sims = () =>
      new Promise<null>((resolve, reject) => {
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
            reject(res.err);
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
                err: "",
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
      new Promise<ResultsSummary>((resolve, reject) => {
        let s = JSON.stringify(results);
        console.log(new Blob([s]).size);
        pool.queue({
          cmd: "collect",
          payload: s,
          cb: (val) => {
            //convert it back
            const res = JSON.parse(val);

            if (res.err) {
              reject(res.err);
            } else {
              // console.log(res);
              resolve(res);
            }
          },
        });
      });

    const version = () =>
      new Promise<null>((resolve, reject) => {
        const versionCB = (val: any) => {
          const res = JSON.parse(val);
          v = res.hash;
          bt = res.date;
          resolve(null);
        };
        pool.queue({ cmd: "version", cb: versionCB });
      });

    //run the sim
    dispatch(
      simActions.setRunStats({
        progress: 0,
        result: -1,
        time: -1,
        err: "",
      })
    );

    setConfig()
      .then(() => {
        console.log("configs done");
        return Promise.all([debugRun(), sims(), version()]);
      })
      .then(() => {
        console.log("all iters done, collecting results");
        console.time("aggregate results");
        //aggregate the result here
        return aggregateResults();
      })
      .then((summary) => {
        // console.log("do something with results here");
        const end = window.performance.now();
        dispatch(
          simActions.setRunStats({
            progress: -1,
            result: avg / iters,
            time: end - startTime,
            err: "",
          })
        );
        //add debug to the summary
        summary.debug = debug;
        summary.v2 = true;
        summary.runtime = 1000000 * (end - startTime); //because cmd line outputs in nanoseconds
        summary.version = v;
        summary.build_date = bt;
        summary.iter = iters;

        console.timeEnd("aggregate results");
        //summary can now be passed to viewer
        dispatch(
          viewerActions.addViewerData({
            key: "Simulation run on: " + new Date().toLocaleString(),
            data: summary,
          })
        );
      })
      .catch((err) => {
        console.warn(err);
        const end = window.performance.now();
        dispatch(
          simActions.setRunStats({
            progress: -1,
            result: 0,
            time: end - startTime,
            err: err,
          })
        );
      });

    // setRuntime(end - startTime);
  };
}
