import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { ResultsSummary } from "~src/types";

export interface Viewer {
  data: { [key in string]: ResultsSummary };
  selected: string;
}

export const viewerInitialState: Viewer = {
  data: {},
  selected: "",
};

export const viewerSlice = createSlice({
  name: "viewer",
  initialState: viewerInitialState,
  reducers: {
    addViewerData: (
      state,
      action: PayloadAction<{ key: string; data: ResultsSummary }>
    ) => {
      const now = new Date().getTime();
      state.data[action.payload.key] = action.payload.data;
      state.selected = action.payload.key;

      return state;
    },
    // addViewerDataAndSetSelected: (
    //   state,
    //   action: PayloadAction<{ key: string; data: ResultsSummary }>
    // ) => {
    //   const l = state.data.length;
    //   state.data[action.payload.key] = action.payload.data;
    //   state.selected = action.payload.key;
    //   return state;
    // },
    setSelected: (state, action: PayloadAction<string>) => {
      state.selected = action.payload;
      return state;
    },
  },
});

export const viewerActions = viewerSlice.actions;

export type ViewerSlice = {
  [viewerSlice.name]: ReturnType<typeof viewerSlice["reducer"]>;
};
