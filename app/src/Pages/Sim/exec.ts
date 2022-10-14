import { AppThunk } from '~src/store';
import { pool, simActions } from '.';
import { Result } from '~src/types';
import { viewerActions } from '../ViewerDashboard/viewerSlice';
import { ResultsSummary } from '~src/Types/stats';

const iterRegex = /iteration=(\d+)/;

function extractItersFromConfig(cfg: string): number {
  let iters = 1;
  let m = iterRegex.exec(cfg);
  // console.log(m);

  if (m) {
    iters = parseInt(m[1]);

    if (isNaN(iters)) {
      console.warn('no iteration found in settings: ', m);
      iters = 1000;
    }
  } else {
    console.log(cfg);
    console.warn('iter regex failed');
    iters = 1000;
  }
  //   console.log("parsed iters: " + iters);

  return iters;
}

export function runSim(cfg: string): AppThunk {
  return function (dispatch) {
    console.log('starting run');
    //extract the number of iterations from the config file
    // TODO: should be a wasm function for config metadata (response of initialize?)
    const iters = extractItersFromConfig(cfg);
    const start = window.performance.now();

    //run the sim
    dispatch(simActions.setRunStats({
      progress: 0,
      result: -1,
      time: -1,
      err: '',
    }));

    const refresh = Math.max(5, (pool.count() * 2));
    pool.run(cfg, iters, (completed) => {
      if (completed % refresh == 0) {
        dispatch(simActions.setRunStats({
          progress: completed / iters,
          result: -1,
          time: -1,
          err: '',
        }));
      }
    }).then((summary) => {
      const end = window.performance.now();
      dispatch(simActions.setRunStats({
        progress: -1,
        result: summary.dps.mean,
        time: end - start,
        err: '',
      }));
      // need to set runtime ourselves due to how it's calc'd in cli
      summary.runtime = 1_000_000 * (end - start); // runtime is in ns

      dispatch(
        viewerActions.addViewerData({
          key: 'Simulation run on: ' + new Date().toLocaleString(),
          data: summary,
        })
      );
    }).catch((err) => {
      console.warn(err);
      dispatch(simActions.setRunStats({
        progress: -1,
        result: 0,
        time: -1,
        err: err,
      }));
    });
  };
}