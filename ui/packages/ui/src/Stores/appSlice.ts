import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export interface AppState {
	isSettingsOpen: boolean;
	sampleOnLoad: boolean;

	cfg: string;
}

export const initialState: AppState = {
	isSettingsOpen: false,
	sampleOnLoad: false,
	cfg: "",
};

export const defaultStats = [
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];
export const maxStatLength = defaultStats.length;

export const charLinesRegEx =
	/^(\w+) (?:char|add) (?:lvl|weapon|set|stats).+$(?:\r\n|\r|\n)?/gm;

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
		setCfg: (
			state,
			action: PayloadAction<{ cfg: string; keepTeam: boolean }>,
		) => {
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
	},
});
export const appActions = appSlice.actions;
