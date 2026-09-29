import { type KeyStorage, LEGACY_APP_DATA_KEY, readJSON } from "./storage";

const PREFS_KEY = "gcsim-prefs";

export function loadSampleOnLoad(storage: KeyStorage): boolean {
	const saved =
		readJSON(storage, PREFS_KEY) ?? readJSON(storage, LEGACY_APP_DATA_KEY);
	return saved?.sampleOnLoad ?? false;
}

export function saveSampleOnLoad(storage: KeyStorage, sampleOnLoad: boolean) {
	storage.setItem(PREFS_KEY, JSON.stringify({ sampleOnLoad }));
}
