import {
	createStore,
	type KeyStorage,
	LEGACY_APP_DATA_KEY,
	readJSON,
} from "./externalStore";

const PREFS_KEY = "gcsim-prefs";

export interface Prefs {
	sampleOnLoad: boolean;
	settingsOpen: boolean;
}

export function prefsStore(storage: KeyStorage) {
	const saved =
		readJSON(storage, PREFS_KEY) ?? readJSON(storage, LEGACY_APP_DATA_KEY);
	return createStore<Prefs>(
		{ sampleOnLoad: saved?.sampleOnLoad ?? false, settingsOpen: false },
		({ sampleOnLoad }) =>
			storage.setItem(PREFS_KEY, JSON.stringify({ sampleOnLoad })),
	);
}
