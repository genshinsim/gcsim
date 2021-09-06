import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface ResultState {
  text: string;
  haveResult: boolean;
  data: AvgModeSummary | null;
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
    setResultData: (state, action: PayloadAction<AvgModeSummary | null>) => {
      state.data = action.payload;
      state.haveResult = state.data !== null;
    },
  },
});

export const { setResultData } = resultSlice.actions;
export default resultSlice.reducer;

export interface AvgModeSummary {
  is_damage_mode: boolean;
  char_names: string[];
  damage_by_char: { [key: string]: ResultSummary }[];
  char_active_time: ResultSummary[];
  abil_usage_count_by_char: { [key: string]: ResultSummary }[];
  reactions_triggered: { [key: string]: ResultSummary };
  sim_duration: ResultSummary;
  damage: ResultSummary;
  dps: ResultSummary;
  iter: number;
  text: string;
  debug: string;
}

export interface ResultSummary {
  mean: number;
  min: number;
  max: number;
  sd?: number;
}
