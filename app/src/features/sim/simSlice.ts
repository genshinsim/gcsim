// import { Intent, Position, Toaster } from "@blueprintjs/core";
import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { sendMessage } from "app/appSlice";
import { AppThunk } from "app/store";
import { setActiveName, setLogs, setNames } from "features/debug/debugSlice";
import { setResultData } from "features/results/resultsSlice";
// eslint-disable-next-line
import Worker from "worker-loader!./SimWorker";

export function saveConfig(path: string, config: string): AppThunk {
  return function (dispatch) {
    const cb = (resp: any) => {
      //check resp code
      if (resp.status !== 200) {
        //do something here
        console.log("Error from server: ", resp.payload);
        return;
      }
      //update
      dispatch(setHasChange(false));
      console.log(resp.data);
    };

    dispatch(
      sendMessage(
        "file",
        "save/file",
        JSON.stringify({
          path: path,
          data: config,
        }),
        cb
      )
    );
  };
}

declare function sim(
  content: string,
  callback: (err: string, data: string) => void
): void;

export function runSim(config: simConfig): AppThunk {
  return function (dispatch, getState) {

    let worker: Worker = new Worker()

    worker.current.onmessage = (e) => {

    }

    worker.current.postMessage(config)


    const cb = (err: string, data: string) => {
      console.log(err)
      console.log(data)
      dispatch(setLoading(false));
      if (err == null) {
        let r = JSON.parse(data)
        dispatch(setResultData(r));
        if (r.debug) {
          dispatch(setLogs(r.debug));
          dispatch(setNames(r.char_names));
        }
        dispatch(setMessage("Simulation finished. check results"));
      } else {
        dispatch(setMessage("Sim encountered error: " + err));
        dispatch(setHasErr(true));
        console.log("err: ", err)
      }
    }

    dispatch(setLoading(true));
    dispatch(setResultData(null));
    dispatch(setLogs(""));
    dispatch(setNames([]));
    dispatch(setMessage(""));
    dispatch(setHasErr(false));

    //find out who the active is
    const found = config.config.match(/active\+=(\w+);/);
    if (found) {
      dispatch(setActiveName(found[1]));
    }

    sim(
      JSON.stringify(config),
      cb
    )

  };
}

export interface simConfig {
  config: string;
  options: {
    log_details: boolean;
    debug: boolean;
    iter: number;
    workers: number;
    duration: number;
  };
}

export interface ICharacter {
  name: string;
  ascension: number; // 0 to 6
  level: number;
  constellation: number; // 0 to 6
  talents: {
    auto: number,
    skill: number,
    burst: number,
  }
  stats: number[];
}

export interface IArtifact {
  stats: number[];
}

interface SimState {
  isLoading: boolean;
  config: string;
  hasChange: boolean;
  msg: string;
  hasErr: boolean;
}
const initialState: SimState = {
  isLoading: false,
  config: "",
  hasChange: false,
  msg: "",
  hasErr: false,
};

export const simSlice = createSlice({
  name: "sim",
  initialState,
  reducers: {
    setConfig: (state, action: PayloadAction<string>) => {
      state.config = action.payload;
      localStorage.setItem("sim-config", action.payload);
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    setHasChange: (state, action: PayloadAction<boolean>) => {
      state.hasChange = action.payload;
    },
    setMessage: (state, action: PayloadAction<string>) => {
      state.msg = action.payload;
    },
    setHasErr: (state, action: PayloadAction<boolean>) => {
      state.hasErr = action.payload;
    },
  },
});

export const { setConfig, setLoading, setHasChange, setMessage, setHasErr } =
  simSlice.actions;
export default simSlice.reducer;
