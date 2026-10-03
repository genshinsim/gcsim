import { DEFAULT_HERO_ID, HERO_IMAGES } from "../features/dashboard/heroImages";
import { type KeyStorage, LEGACY_APP_DATA_KEY, readJSON } from "./storage";
import { DEFAULT_THEME_ID, THEMES } from "./themes";

const PREFS_KEY = "gcsim-prefs";
export const THEME_KEY = "gcsim-theme";
export const HERO_KEY = "gcsim-hero";

export function loadSampleOnLoad(storage: KeyStorage): boolean {
	const saved =
		readJSON(storage, PREFS_KEY) ?? readJSON(storage, LEGACY_APP_DATA_KEY);
	return saved?.sampleOnLoad ?? false;
}

export function saveSampleOnLoad(storage: KeyStorage, sampleOnLoad: boolean) {
	storage.setItem(PREFS_KEY, JSON.stringify({ sampleOnLoad }));
}

function loadId(
	storage: KeyStorage,
	key: string,
	options: { id: string }[],
	fallback: string,
): string {
	const saved = storage.getItem(key);
	return options.find((o) => o.id === saved)?.id ?? fallback;
}

export function loadTheme(storage: KeyStorage): string {
	return loadId(storage, THEME_KEY, THEMES, DEFAULT_THEME_ID);
}

export function saveTheme(storage: KeyStorage, theme: string) {
	storage.setItem(THEME_KEY, theme);
}

export function loadHero(storage: KeyStorage): string {
	return loadId(storage, HERO_KEY, HERO_IMAGES, DEFAULT_HERO_ID);
}

export function saveHero(storage: KeyStorage, hero: string) {
	storage.setItem(HERO_KEY, hero);
}
