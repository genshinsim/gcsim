import {
  createSlice,
  PayloadAction,
  createListenerMiddleware,
  isAnyOf,
  TypedStartListening,
} from "@reduxjs/toolkit";
import { AppThunk, RootState } from "./store";
import { Character } from "../Types";
import { charToCfg } from "../Pages/Simulator/helper";
import { Executor } from "@gcsim/executors";
import { pool } from "../App";

export interface AppState {
  ready: number;
  workers: number;
  cfg: string;
  cfg_err: string;
  team: Character[];
}

export const initialState: AppState = {
  ready: 0,
  workers: 3,
  cfg: "",
  cfg_err: "",
  team: [],
};

export const defaultStats = [
  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];
export const maxStatLength = defaultStats.length;

export const charLinesRegEx =
  /^(\w+) (?:char|add) (?:lvl|weapon|set|stats).+$(?:\r\n|\r|\n)?/gm;

export function cfgFromTeam(team: Character[], cfg: string): string {
  let next = "";
  //generate new
  team.forEach((c) => {
    next += charToCfg(c) + "\n";
  });

  //purge existing characters:
  cfg = cfg.replace(charLinesRegEx, "");
  cfg = next + cfg;
  //stirp extra new lines
  cfg = cfg.replace(/(\r\n|\r|\n){2,}/g, "$1\n");

  return cfg;
}

export function updateCfg(cfg: string, keepTeam?: boolean): AppThunk {
  return function (dispatch, getState) {
    // console.log(cfg);
    if (keepTeam) {
      // purge char stat from incoming
      let next = cfg;
      //purge existing characters:
      next = next.replace(charLinesRegEx, "");
      //pull out existing

      let old = "";
      let lastChar = "";
      const matches = getState().app.cfg.matchAll(charLinesRegEx);
      for (const match of matches) {
        const line = match[0];
        if (match[1] !== lastChar) {
          old += "\n";
          lastChar = match[1];
        }
        console.log(match);
        old += line;
      }

      next = old + "\n" + next;

      //strip extra new lines
      cfg = next.replace(/(\r\n|\r|\n){2,}/g, "$1\n");
    }
    dispatch(appActions.setCfg(cfg));
    pool.validate(cfg).then(
      (res) => {
        console.log("all is good");
        dispatch(appActions.setCfgErr(""));
        //if successful then we're going to update the team based on the parsed results
        let team: Character[] = [];
        if (res.characters) {
          team = res.characters.map((c) => {
            return {
              name: c.base.key,
              level: c.base.level,
              element: c.base.element,
              max_level: c.base.max_level,
              cons: c.base.cons,
              weapon: c.weapon,
              talents: c.talents,
              stats: c.stats,
              snapshot: defaultStats,
              sets: c.sets,
            };
          });
        }
        //check if there are any warning msgs
        if (res.errors) {
          let msg = "";
          res.errors.forEach((err) => {
            msg += err + "\n";
          });
          dispatch(appActions.setCfgErr(msg));
        }
        dispatch(appActions.setTeam(team));
      },
      (err) => {
        //set error state
        dispatch(appActions.setCfgErr(err));
      }
    );
  };
}

export function setTotalWorkers(pool: Executor, count: number): AppThunk {
  return function (dispatch, getState) {
    //do nothing if ready
    pool.setWorkerCount(count, (x: number) => {
      //call back for ready
      dispatch(appActions.setWorkerReady(x));
    });
    dispatch(appActions.setWorkers(count));
  };
}

export function ready(pool: Executor): boolean {
  return pool.ready();
}

export const appSlice = createSlice({
  name: "app",
  initialState: initialState,
  reducers: {
    setWorkers: (state, action: PayloadAction<number>) => {
      state.workers = action.payload;
      return state;
    },
    setWorkerReady: (state, action: PayloadAction<number>) => {
      state.ready = action.payload;
      return state;
    },
    setCfg: (state, action: PayloadAction<string>) => {
      state.cfg = action.payload;
      return state;
    },
    setCfgErr: (state, action: PayloadAction<string>) => {
      state.cfg_err = action.payload;
      return state;
    },
    addCharacter: (state, action: PayloadAction<{ character: Character }>) => {
      if (state.team.length >= 4) return state;
      state.team.push(action.payload.character);

      const cfg = cfgFromTeam(state.team, state.cfg);
      state.cfg = cfg;
      return state;
    },
    deleteCharacter: (state, action: PayloadAction<{ index: number }>) => {
      if (
        action.payload.index < 0 ||
        action.payload.index >= state.team.length
      ) {
        return state;
      }
      state.team.splice(action.payload.index, 1);
      const cfg = cfgFromTeam(state.team, state.cfg);
      state.cfg = cfg;
      return state;
    },
    editCharacter: (
      state,
      action: PayloadAction<{ char: Character; index: number }>
    ) => {
      if (
        action.payload.index < 0 ||
        action.payload.index >= state.team.length
      ) {
        return state;
      }
      state.team[action.payload.index] = action.payload.char;
      const cfg = cfgFromTeam(state.team, state.cfg);
      state.cfg = cfg;
      return state;
    },
    setTeam: (state, action: PayloadAction<Character[]>) => {
      state.team = action.payload;
      return state;
    },
  },
});

export const appActions = appSlice.actions;

export type ViewerSlice = {
  [appSlice.name]: ReturnType<typeof appSlice["reducer"]>;
};

export const listenerMiddleware = createListenerMiddleware();
const appStartListening =
  listenerMiddleware.startListening as TypedStartListening<RootState>;
appStartListening({
  matcher: isAnyOf(
    appSlice.actions.addCharacter,
    appSlice.actions.deleteCharacter,
    appSlice.actions.editCharacter
  ),
  effect: async (action, listenerApi) => {
    const cfg = listenerApi.getState().app.cfg;
    console.log("middleware triggered on: ", action.type);
    console.log("cfg updated: ", cfg);
    listenerApi.dispatch(updateCfg(cfg));
  },
});
