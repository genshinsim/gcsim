import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { throttle } from "lodash";
import { AppThunk } from "./store";
import { pool } from "~/Executor";
import { VIEWER_THROTTLE } from "~/Pages/Viewer";
import { SimResults } from "~/Types";

export interface Viewer {
  data: SimResults | null;
  error: string | null;
}

export function runSim(cfg: string): AppThunk {
  return function (dispatch) {
    console.log("starting run");
    dispatch(viewerActions.start());

    const updateResult = throttle(
      (res: SimResults) => {
        dispatch(viewerActions.setResult({ data: res }));
      },
      VIEWER_THROTTLE,
      { leading: true, trailing: true }
    );

    pool
      .run(cfg, (result) => {
        updateResult(result);
      })
      .catch((err) => {
        dispatch(viewerActions.setError({ error: err }));
      });
  };
}

export const viewerInitialState: Viewer = {
  data: null,
  error: "Nothing is loaded!",
};

export const viewerSlice = createSlice({
  name: "viewer",
  initialState: viewerInitialState,
  reducers: {
    setResult: (state, action: PayloadAction<{ data: SimResults }>) => {
      state.data = action.payload.data;
      return state;
    },
    start: (state) => {
      state.data = null;
      state.error = null;
      return state;
    },
    setError: (state, action: PayloadAction<{ error: string }>) => {
      state.error = action.payload.error;
      return state;
    },
  },
});

export const viewerActions = viewerSlice.actions;

export type ViewerSlice = {
  [viewerSlice.name]: ReturnType<typeof viewerSlice["reducer"]>;
};
