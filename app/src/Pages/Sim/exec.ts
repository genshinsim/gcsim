import { throttle } from 'lodash';
import { AppThunk } from '~src/store';
import { pool } from '.';
import { VIEWER_THROTTLE } from '../Viewer';
import { SimResults } from '../Viewer/SimResults';
import { viewerActions } from '../Viewer/viewerSlice';

// TODO: move to viewer?
export function runSim(cfg: string): AppThunk {
  return function (dispatch) {
    console.log('starting run');
    dispatch(viewerActions.start());

    const updateResult = throttle((res: SimResults) => {
      dispatch(viewerActions.setResult({ data: res }));
    } , VIEWER_THROTTLE, { leading: true, trailing: true });
    
    pool.run(cfg, (result) => {
      updateResult(result);
    }).catch((err) => {
      dispatch(viewerActions.setError({ error: err }));
    });
  };
}