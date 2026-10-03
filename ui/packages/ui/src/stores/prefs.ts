import { DEFAULT_HERO_ID, HERO_IMAGES } from "../features/dashboard/heroImages";
import { type KeyStorage, LEGACY_APP_DATA_KEY, readJSON } from "./storage";
import { DEFAULT_THEME_ID, THEMES } from "./themes";

const PREFS_KEY = "gcsim-prefs";
// web/index.html reads this key before first paint
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

export function loadTheme(storage: KeyStorage): string {
	const saved = storage.getItem(THEME_KEY);
	return THEMES.some((t) => t.id === saved)
		? (saved as string)
		: DEFAULT_THEME_ID;
}

export function saveTheme(storage: KeyStorage, theme: string) {
	storage.setItem(THEME_KEY, theme);
}

export function loadHero(storage: KeyStorage): string {
	const saved = storage.getItem(HERO_KEY);
	return HERO_IMAGES.some((h) => h.id === saved)
		? (saved as string)
		: DEFAULT_HERO_ID;
}

export function saveHero(storage: KeyStorage, hero: string) {
	storage.setItem(HERO_KEY, hero);
}
