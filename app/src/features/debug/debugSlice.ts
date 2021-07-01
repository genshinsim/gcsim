import { Intent } from "@blueprintjs/core";
import { IconNames, IconName } from "@blueprintjs/icons";
import { createSelector, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RootState } from "app/store";

interface DebugState {
  logs: string;
  names: string[];
  activeName: string;
  haveDebug: boolean;
}

const initialState: DebugState = {
  logs: "",
  names: [],
  activeName: "",
  haveDebug: false,
};

export const debugSlice = createSlice({
  name: "debug",
  initialState,
  reducers: {
    setLogs: (state, action: PayloadAction<string>) => {
      state.logs = action.payload;
      if (state.logs !== "") {
        state.haveDebug = true;
      }
    },
    setNames: (state, action: PayloadAction<string[]>) => {
      state.names = action.payload;
    },
    setActiveName: (state, action: PayloadAction<string>) => {
      state.activeName = action.payload;
    },
  },
});

export const { setLogs, setNames, setActiveName } = debugSlice.actions;
export default debugSlice.reducer;

export interface DataRow {
  F: number;
  Cols: DataCol[];
  Active: number;
}

export interface DataCol {
  Parts: DataPoint[];
}

export interface DataPoint {
  F: number;
  Event: string;
  Char: number;
  M: string;
  Raw: string;
  Amount: number;
  Intent: Intent;
  Icon: IconName;
  Right: string;
}

const selectDebugLog = (state: RootState) => state.debug.logs;
const selectDebugNames = (state: RootState) => state.debug.names;
const selectDebugActive = (state: RootState) => state.debug.activeName;

export const selectLogs = createSelector(
  [selectDebugLog, selectDebugNames, selectDebugActive],
  (logs, names, active) => {
    //find first active
    let current = names.findIndex((e) => e === active);
    //0 if -1, else 1/2/3/4
    current++;

    const charCount = names.length;

    const str = logs.split(/\r?\n/);
    let data: DataRow[] = [];
    let cols: DataCol[] = [
      {
        Parts: [],
      },
      {
        Parts: [],
      },
      {
        Parts: [],
      },
      {
        Parts: [],
      },
      {
        Parts: [],
      },
    ];

    let lastFrame = -1;
    // let isFirst = true;

    str.forEach((v) => {
      if (v === "") {
        return;
      }
      const d = JSON.parse(v);

      if (!("frame" in d)) {
        console.log("error, no frames: ", d);
        return;
      }

      if (!("event" in d)) {
        console.log("error, no event: ", d);
        return;
      }

      const event = d.event;

      let char = 0;

      if ("char" in d) {
        char = d.char + 1;
      }

      //if frame changed
      if (d.frame !== lastFrame) {
        //add it to the labels
        data.push({
          F: lastFrame,
          Cols: cols,
          Active: current,
        });
        lastFrame = d.frame;
        cols = [];
        for (var i = 0; i <= charCount; i++) {
          cols.push({ Parts: [] });
        }
      }

      let x: DataPoint = {
        F: d.frame,
        M: d.M,
        Raw: JSON.stringify(JSON.parse(v), null, 2),
        Event: event,
        Char: char,
        Intent: Intent.NONE,
        Icon: IconNames.CIRCLE,
        Amount: 0,
        Right: "",
      };

      switch (x.Event) {
        case "damage":
          x.M +=
            " [" +
            Math.round(d.damage)
              .toString()
              .replace(/\B(?=(\d{3})+(?!\d))/g, ",") +
            "]";
          let extra = "";
          if (d.amp && d.amp !== "") {
            extra += d.amp;
          }
          if (d.crit) {
            extra += " crit";
          }
          if (extra !== "") {
            x.M += " (" + extra.trim() + ")";
          }

          x.Icon = IconNames.FLAME;
          x.Intent = Intent.PRIMARY;
          x.Amount = d.damage;
          x.Right = d.target;
          break;
        case "action":
          if (d.M.includes("executed") && d.action === "swap") {
            current = names.findIndex((e) => e === d.target);
            current++;
            x.M += " to " + d.target;
          }
          x.Icon = IconNames.PLAY;
          break;
        case "element":
          switch (d.M) {
            case "expired":
              x.M = d.old_ele + " expired";
              break;
            case "application":
              x.M =
                d.applied_ele +
                " applied" +
                (d.existing_ele === "" ? "" : " to " + d.existing_ele);
              break;
            case "refreshed":
              x.M = d.ele + " refreshed";
              break;
            default:
              x.M = d.M;
          }

          x.Icon = IconNames.FLASH;
          x.Intent = Intent.WARNING;
          x.Right = d.target;
          break;
        case "energy":
          if (x.M.includes("particle")) {
            x.M =
              d.M +
              " from " +
              d.source +
              ", next: " +
              Math.round(d["post_recovery"]);
          }
          x.Icon = IconNames.IMPORT;
          x.Intent = Intent.SUCCESS;
          break;
        case "calc":
          x.Icon = IconNames.CALCULATOR;
          x.Right = d.target;
          break;
        case "character":
          x.Icon = IconNames.PERSON;
          break;
        case "snapshot":
          x.Icon = IconNames.CAMERA;
          break;
        case "heal":
          x.Icon = IconNames.HEART;
          break;
        case "hurt":
          x.Icon = IconNames.VIRUS;
          break;
        case "queue":
          x.Icon = IconNames.ADD_TO_ARTIFACT;
          break;
        case "shield":
          x.Icon = IconNames.SHIELD;
          break;
        case "hook":
          x.Icon = IconNames.PAPERCLIP;
          break;
        case "icd":
          x.Icon = IconNames.STOPWATCH;
          break;
        case "construct":
          x.Icon = IconNames.WRENCH;
          break;
        default:
          x.M = event + ": " + d.M;
      }

      if (!cols[char]) {
        console.log(cols);
        console.log(char);
        console.log(d);
      }

      cols[char].Parts.push(x);
    });

    return data;
  }
);
