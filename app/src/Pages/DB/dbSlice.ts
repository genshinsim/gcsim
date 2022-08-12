import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import axios from "axios";
import { DBItem } from "~src/types";
import { AppThunk } from "~src/store";

export interface DBState {
  data: DBItem[];
  loading: boolean;
  loaded: boolean;
}

export const dbInitialState: DBState = {
  data: [],
  loading: false,
  loaded: false,
};

export function loadDB(): AppThunk {
  return function (dispatch, getState) {
    dispatch(dbActions.setLoading(true));
    const url = "https://viewer.gcsim.workers.dev/gcsimdb";
    axios
      .get(url)
      .then((resp) => {
        console.log(resp.data);
        let data = resp.data;

        dispatch(dbActions.setData(data));
        dispatch(dbActions.setLoading(false));
        dispatch(dbActions.setLoaded(true));
      })
      .catch(function (error) {
        // handle error
        console.log(error);
        dispatch(dbActions.setData([]));
        dispatch(dbActions.setLoading(false));
      });
  };
}

export const dbSlice = createSlice({
  name: "db",
  initialState: dbInitialState,
  reducers: {
    setData: (state, action: PayloadAction<DBItem[]>) => {
      state.data = action.payload;
      return state;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
      return state;
    },
    setLoaded: (state, action: PayloadAction<boolean>) => {
      state.loaded = action.payload;
      return state;
    },
  },
});

export const dbActions = dbSlice.actions;

export type DBSlice = {
  [dbSlice.name]: ReturnType<typeof dbSlice["reducer"]>;
};
