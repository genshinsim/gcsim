import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import axios from "axios";
import pako from "pako";
import { AppThunk } from "Store/store";
import { DBAvatarSimCount, DBAvatarSimDetails } from "Types/database";

type statusType = "idle" | "loading" | "done" | "error";

export interface DBState {
  characters: DBAvatarSimCount[];
  charSims: {
    [key in string]: DBAvatarSimDetails[];
  };
  all: DBAvatarSimDetails[];
  status: statusType;
  errorMsg: string;
}

export const dbInitialState: DBState = {
  characters: [],
  charSims: {},
  status: "idle",
  errorMsg: "",
  all: [],
};

export function loadAllDB(): AppThunk {
  return function (dispatch) {
    dispatch(dbActions.setStatus("loading"));
    const url = "/api/db/all";
    axios
      .get(url)
      .then((resp) => {
        let next: DBAvatarSimDetails[] = [];
        resp.data.forEach((e: any) => {
          const binaryStr = Uint8Array.from(window.atob(e.config_hash), (v) =>
            v.charCodeAt(0)
          );
          const restored = pako.inflate(binaryStr, { to: "string" });
          next.push({
            ...e,
            metadata: JSON.parse(e.metadata),
            create_time: Math.floor(new Date(e.create_time).getTime() / 1000),
            config: restored,
          });
        });
        next.sort((a,b) => {
          return (b.metadata.dps.mean / b.metadata.num_targets) - (a.metadata.dps.mean / a.metadata.num_targets)
        })
        dispatch(dbActions.setFullDB(next));
      })
      .catch(function (error) {
        // handle error
        console.log(error);
        dispatch(dbActions.setError(`Error encountered: ${error}`));
      });
  };
}


export function loadDB(): AppThunk {
  return function (dispatch) {
    dispatch(dbActions.setStatus("loading"));
    const url = "/api/db";
    axios
      .get(url)
      .then((resp) => {
        dispatch(dbActions.setCharacters(resp.data));
      })
      .catch(function (error) {
        // handle error
        console.log(error);
        dispatch(dbActions.setError(`Error encountered: ${error}`));
      });
  };
}

export function loadCharacter(char: string): AppThunk {
  return function (dispatch) {
    dispatch(dbActions.setStatus("loading"));
    axios
      .get(`/api/db/${char}`)
      .then((resp) => {
        console.log(resp.data);
        let next: DBAvatarSimDetails[] = [];
        resp.data.forEach((e: any) => {
          const binaryStr = Uint8Array.from(window.atob(e.config_hash), (v) =>
            v.charCodeAt(0)
          );
          const restored = pako.inflate(binaryStr, { to: "string" });
          next.push({
            ...e,
            metadata: JSON.parse(e.metadata),
            create_time: Math.floor(new Date(e.create_time).getTime() / 1000),
            config: restored,
          });
        });
        console.log(next);
        dispatch(dbActions.setCharSimList({ char: char, data: next }));
        dispatch(dbActions.setStatus("done"));
      })
      .catch((err) => {
        dispatch(
          dbActions.setError(
            `Error encountered loading sims for ${char}: ${err}`
          )
        );
      });
  };
}

export const dbSlice = createSlice({
  name: "db",
  initialState: dbInitialState,
  reducers: {
    setCharacters: (state, action: PayloadAction<DBAvatarSimCount[]>) => {
      state.characters = action.payload;
      state.errorMsg = "";
      state.status = "done";
      return state;
    },
    setCharSimList: (
      state,
      action: PayloadAction<{ char: string; data: DBAvatarSimDetails[] }>
    ) => {
      const { char, data } = action.payload;
      state.charSims[char] = data;
      return state;
    },
    setFullDB: (state, action: PayloadAction<DBAvatarSimDetails[]>) => {
      state.all = action.payload
      return state
    },
    setError: (state, action: PayloadAction<string>) => {
      state.errorMsg = action.payload;
      state.status = "error";
      return state;
    },
    setStatus: (state, action: PayloadAction<statusType>) => {
      state.status = action.payload;
      state.errorMsg = "";
      return state;
    },
  },
});

export const dbActions = dbSlice.actions;

export type DBSlice = {
  [dbSlice.name]: ReturnType<typeof dbSlice["reducer"]>;
};
