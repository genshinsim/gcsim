// import { Intent, Position, Toaster } from "@blueprintjs/core";
import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { sendMessage } from "app/appSlice";
import { AppThunk } from "app/store";
import { setActiveName, setLogs, setNames } from "features/debug/debugSlice";
import { setResultData } from "features/results/resultsSlice";

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

export function runSim(config: simConfig): AppThunk {
  return function (dispatch, getState) {
    const cb = (resp: any) => {
      dispatch(setLoading(false));
      //check resp code
      if (resp.status !== 200) {
        //do something here
        console.log("Error from server: ", resp.payload);
        // Toaster.create({ position: Position.BOTTOM }).show({
        //   message: "Error running sim: " + resp.payload,
        //   intent: Intent.DANGER,
        // });
        dispatch(setMessage("Sim encountered error: " + resp.payload));
        dispatch(setHasErr(true));

        return;
      }
      //update
      console.log("sim/run received response");
      var data = JSON.parse(resp.payload);
      console.log(data);
      // dispatch(setResultData(data.summary));

      dispatch(setResultData({ text: data.summary, data: data.details }));

      if (data.log) {
        dispatch(setLogs(data.log));
        dispatch(setNames(data.names));
      }

      dispatch(setMessage("Simulation finished. check results"));

      // Toaster.create({ position: Position.BOTTOM }).show({
      //   message: "Simulation finished. check results",
      //   intent: Intent.SUCCESS,
      // });
    };
    dispatch(setLoading(true));
    dispatch(setResultData({ text: "", data: null }));
    dispatch(setLogs(""));
    dispatch(setNames([]));
    dispatch(setMessage(""));
    dispatch(setHasErr(false));

    //find out who the active is
    const found = config.config.match(/active\+=(\w+);/);
    if (found) {
      dispatch(setActiveName(found[1]));
    }

    dispatch(sendMessage("run", "", JSON.stringify(config), cb));
  };
}

export interface simConfig {
  log: string;
  seconds: number;
  config: string;
  hp: number;
  avg_mode: boolean;
  iter: number;
  noseed: boolean;
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
