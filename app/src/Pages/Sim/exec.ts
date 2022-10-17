import { AppThunk } from '~src/store';
import { pool } from '.';
import { viewerActions } from '../Viewer/viewerSlice';

export function runSim(cfg: string): AppThunk {
  return function (dispatch) {
    console.log('starting run');
    dispatch(viewerActions.start());

    pool.run(cfg, (result) => {
      dispatch(viewerActions.setResult({ data: result }));
    }).catch((err) => {
      console.warn(err);
      dispatch(viewerActions.setError({ error: err }));
    });
  };
}