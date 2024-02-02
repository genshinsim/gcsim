import {
  createSlice,
  PayloadAction,
} from "@reduxjs/toolkit";
import { charToCfg } from "../Pages/Simulator/helper";
import { Character } from "@gcsim/types";

export interface AppState {
  isSettingsOpen: boolean;
  sampleOnLoad: boolean;

  cfg: string;
  team: Character[];
}

export const initialState: AppState = {
  isSettingsOpen: false,
  sampleOnLoad: false,
  cfg: "",
  team: [],
};

export const defaultStats = [
  0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];
export const maxStatLength = defaultStats.length;

export const charLinesRegEx =
  /^(\w+) (?:char|add) (?:lvl|weapon|set|stats).+$(?:\r\n|\r|\n)?/gm;

function cfgFromTeam(team: Character[], cfg: string): string {
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

export const appSlice = createSlice({
  name: "app",
  initialState: initialState,
  reducers: {
    setSettingsOpen: (state, action: PayloadAction<boolean>) => {
      state.isSettingsOpen = action.payload;
      return state;
    },
    setSampleOnLoad: (state, action: PayloadAction<boolean>) => {
      state.sampleOnLoad = action.payload;
      return state;
    },
    setCfg: (state, action: PayloadAction<{ cfg: string, keepTeam: boolean}>) => {
      if (!action.payload.keepTeam) {
        state.cfg = action.payload.cfg;
        return state;
      }

      //purge existing characters:
      let next = action.payload.cfg.replace(charLinesRegEx, "");

      let old = "";
      let lastChar = "";
      const matches = state.cfg.matchAll(charLinesRegEx);
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
      state.cfg = next.replace(/(\r\n|\r|\n){2,}/g, "$1\n");
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