import type { UserInfo, UserSettings } from "@gcsim/types";
import { merge } from "lodash-es";
import { type KeyStorage, readJSON } from "./storage";

const SETTINGS_KEY = "gcsim-user-settings";
const LEGACY_SETTINGS_KEY = "redux-user-local-settings";

export const defaultUser: UserInfo = {
	uid: "",
	name: "",
	role: 0,
	permalinks: [],
	data: {
		settings: { showTips: true, showBuilder: true, showNameSearch: true },
	},
};

export function mergeUser(prev: UserInfo, incoming: UserInfo): UserInfo {
	return merge(structuredClone(prev), incoming);
}

export function loadUser(storage: KeyStorage): UserInfo {
	const settings: UserSettings | null =
		readJSON(storage, SETTINGS_KEY) ?? readJSON(storage, LEGACY_SETTINGS_KEY);
	return settings ? { ...defaultUser, data: { settings } } : defaultUser;
}

export function saveUserSettings(storage: KeyStorage, user: UserInfo) {
	storage.setItem(SETTINGS_KEY, JSON.stringify(user.data.settings));
}
