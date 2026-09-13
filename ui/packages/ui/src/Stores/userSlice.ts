import type { UserInfo, UserSettings } from "@gcsim/types";
import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import axios from "axios";
import { merge } from "lodash-es";
import type { AppThunk } from "./store";

export const initialState: UserInfo = {
	uid: "",
	name: "",
	role: 0,
	permalinks: [],
	data: {
		settings: { showTips: true, showBuilder: true, showNameSearch: true },
	},
};

export function saveUserSettings(): AppThunk {
	return (dispatch, getState) => {
		axios
			.post("/api/user/save", getState().user.data)
			.then((res) => {
				console.log("save ok");
			})
			.catch((error) => {
				console.log("save failed");
			});
	};
}

export const userSlice = createSlice({
	name: "user",
	initialState: initialState,
	reducers: {
		setUser: (state, action: PayloadAction<UserInfo>) => {
			state = action.payload;
			return state;
		},
		mergeUser: (state, action: PayloadAction<UserInfo>) => {
			merge(state, action.payload);
			return state;
		},
		setUserSettings: (state, action: PayloadAction<UserSettings>) => {
			state.data.settings = action.payload;
			return state;
		},
	},
});

export const userActions = userSlice.actions;

export type UserSlice = {
	[userSlice.name]: ReturnType<(typeof userSlice)["reducer"]>;
};
