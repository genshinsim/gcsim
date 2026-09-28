import type { UserInfo, UserSettings } from "@gcsim/types";
import { merge } from "lodash-es";
import { createStore, type KeyStorage, readJSON } from "./externalStore";

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

export function userStore(storage: KeyStorage) {
	const settings: UserSettings | null =
		readJSON(storage, SETTINGS_KEY) ?? readJSON(storage, LEGACY_SETTINGS_KEY);
	const store = createStore<UserInfo>(
		settings ? { ...defaultUser, data: { settings } } : defaultUser,
		(user) => storage.setItem(SETTINGS_KEY, JSON.stringify(user.data.settings)),
	);
	return {
		...store,
		merge: (user: UserInfo) =>
			store.set((prev) => merge(structuredClone(prev), user)),
		setSettings: (settings: UserSettings) =>
			store.set((prev) => ({ ...prev, data: { ...prev.data, settings } })),
		reset: () => store.set(defaultUser),
	};
}
