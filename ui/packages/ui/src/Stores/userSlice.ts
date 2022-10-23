import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { UserInfo, UserSettings } from "@gcsim/types";

export const initialState: UserInfo = {
  user_id: 0,
  user_name: "Guest",
  token: "",
};

export const userSlice = createSlice({
  name: "user",
  initialState: initialState,
  reducers: {
    setUser: (state, action: PayloadAction<UserInfo>) => {
      state = action.payload;
      return state;
    },
    setUserSettings: (state, action: PayloadAction<UserSettings>) => {
      state.settings = action.payload;
      return state;
    },
  },
});

export const userActions = userSlice.actions;

export type UserSlice = {
  [userSlice.name]: ReturnType<typeof userSlice["reducer"]>;
};
