import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { AppThunk } from "./store";
import { UserInfo, UserSettings } from "../Types/user";
import { AuthProvider, DiscordProvider, MockProvider } from "../Data/Provider";

// export const authProvider: AuthProvider = new DiscordProvider();
export const authProvider: AuthProvider = new MockProvider();

const initialState: UserInfo = {
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

//thunks
export function logout(): AppThunk {
  return function (dispatch) {
    authProvider
      .logout()
      .then(() => dispatch(userActions.setUser(initialState)))
      .catch((err) => {
        //log out the user
        console.warn("Error occured logging out: ", err);
        dispatch(userActions.setUser(initialState));
      });
  };
}

export const userActions = userSlice.actions;

export type UserSlice = {
  [userSlice.name]: ReturnType<typeof userSlice["reducer"]>;
};
