import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Character } from "~/src/types";

export interface UserData {
  GOODImport: Character[];
}

const initialState: UserData = {
  GOODImport: [],
};

export const userDataSlice = createSlice({
  name: "user_data",
  initialState: initialState,
  reducers: {
    loadFromGOOD: (state, action: PayloadAction<{ data: Character[] }>) => {
      // if there are characters, do something
      if (action.payload.data.length > 0) {
        state.GOODImport = action.payload.data;
      }
      return state;
    },
  },
});

export const userDataActions = userDataSlice.actions;

export type UserDataSlice = {
  [userDataSlice.name]: ReturnType<typeof userDataSlice["reducer"]>;
};
