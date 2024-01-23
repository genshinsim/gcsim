import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { SimResults } from "@gcsim/types";

export interface Viewer {
  data: SimResults | null;
  hash: string | null;
  recoveryConfig: string | null;
  error: string | null;
}

export const viewerInitialState: Viewer = {
  data: null,
  hash: null,
  recoveryConfig: null,
  error: null,
};

export const viewerSlice = createSlice({
  name: "viewer",
  initialState: viewerInitialState,
  reducers: {
    setResult: (state, action: PayloadAction<{ data: SimResults, hash: string | null }>) => {
      state.data = action.payload.data;
      state.hash = action.payload.hash;
      return state;
    },
    start: (state) => {
      state.data = null;
      state.error = null;
      return state;
    },
    setError: (state, action: PayloadAction<{ recoveryConfig: string | null; error: string }>) => {
      state.recoveryConfig = action.payload.recoveryConfig;
      state.error = action.payload.error;
      return state;
    },
  },
});

export const viewerActions = viewerSlice.actions;

export type ViewerSlice = {
  [viewerSlice.name]: ReturnType<typeof viewerSlice["reducer"]>;
};
