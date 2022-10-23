import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { SimResults } from "@gcsim/types";

export interface Viewer {
  data: SimResults | null;
  error: string | null;
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
