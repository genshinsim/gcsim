import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface ResultState {
  text: string;
  haveResult: boolean;
  data: SingleModeSummary | AvgModeSummary | null;
}

const initialState: ResultState = {
  text: "",
  haveResult: false,
  data: null,
};

export const resultSlice = createSlice({
  name: "result",
  initialState,
  reducers: {
    setResultData: (
      state,
      action: PayloadAction<{
        text: string;
        data: SingleModeSummary | AvgModeSummary | null;
      }>
    ) => {
      state.text = action.payload.text;
      state.data = action.payload.data;
      state.haveResult = state.data !== null;
    },
  },
});

export const { setResultData } = resultSlice.actions;
export default resultSlice.reducer;

export interface AvgModeSummary {
  iter: number;
  avg_duration: number;
  dps: ResultSummary;
  damage_by_char: { [key: string]: ResultSummary }[];
  char_active_time: ResultSummary[];
  abil_usage_count_by_char: { [key: string]: ResultSummary }[];
  reactions_triggered: { [key: string]: ResultSummary };
  char_names: string[];
}

export interface SingleModeSummary {
  sim_duration: number;
  dps: number;
  damage_by_char: { [key: string]: number }[];
  char_active_time: number[];
  abil_usage_count_by_char: { [key: string]: number }[];
  reactions_triggered: { [key: string]: number };
  char_names: string[];
  detailed: {
    damage_hist: number[];
    char_active_frame: number[][];
    element_active_frame: { [key: string]: number[] };
    abil_usage_frame: {
      actor: string;
      action: string;
      param: string;
    }[];
  };
}

export interface ResultSummary {
  mean: number;
  min: number;
  max: number;
  sd?: number;
}
