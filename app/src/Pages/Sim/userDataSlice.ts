import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Character } from "~/src/types";

export interface UserData {
  GOODImport: { [key: string]: Character };
}

const initialState: UserData = {
  GOODImport: {},
};

export const userDataSlice = createSlice({
  name: "user_data",
  initialState: initialState,
  reducers: {
    loadFromGOOD: (state, action: PayloadAction<{ data: Character[] }>) => {
      // if there are characters, do something
      if (action.payload.data.length > 0) {
        //make it
        state.GOODImport = {};
        action.payload.data.forEach((c) => {
          state.GOODImport[c.name] = c;
        });
      }
      return state;
    },
  },
});

export const userDataActions = userDataSlice.actions;

export type UserDataSlice = {
  [userDataSlice.name]: ReturnType<typeof userDataSlice["reducer"]>;
};
