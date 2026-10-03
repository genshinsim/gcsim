import { describe, expect, it } from "vitest";
import { DEFAULT_HERO_ID } from "../features/dashboard/heroImages";
import { fakeStorage } from "./fakeStorage";
import {
	HERO_KEY,
	loadHero,
	loadSampleOnLoad,
	loadTheme,
	saveHero,
	saveSampleOnLoad,
	saveTheme,
	THEME_KEY,
} from "./prefs";
import { DEFAULT_THEME_ID } from "./themes";

describe("sampleOnLoad", () => {
	it("defaults to false", () => {
		expect(loadSampleOnLoad(fakeStorage())).toBe(false);
	});

	it("round-trips a saved value", () => {
		const storage = fakeStorage();
		saveSampleOnLoad(storage, true);
		expect(loadSampleOnLoad(storage)).toBe(true);
	});

	it("loads a value saved by the redux store", () => {
		const storage = fakeStorage({
			"redux-app-data": JSON.stringify({ sampleOnLoad: true, cfg: "" }),
		});
		expect(loadSampleOnLoad(storage)).toBe(true);
	});
});

describe("theme", () => {
	it("defaults to cryo", () => {
		expect(loadTheme(fakeStorage())).toBe(DEFAULT_THEME_ID);
	});

	it("round-trips a saved value", () => {
		const storage = fakeStorage();
		saveTheme(storage, "ember-l");
		expect(loadTheme(storage)).toBe("ember-l");
	});

	it("falls back to cryo for an unknown value", () => {
		expect(loadTheme(fakeStorage({ [THEME_KEY]: "neon" }))).toBe(
			DEFAULT_THEME_ID,
		);
	});
});

describe("hero", () => {
	it("defaults to the default hero", () => {
		expect(loadHero(fakeStorage())).toBe(DEFAULT_HERO_ID);
	});

	it("round-trips a saved value", () => {
		const storage = fakeStorage();
		saveHero(storage, "ganyu");
		expect(loadHero(storage)).toBe("ganyu");
	});

	it("falls back to the default hero for an unknown value", () => {
		expect(loadHero(fakeStorage({ [HERO_KEY]: "nope" }))).toBe(DEFAULT_HERO_ID);
	});
});
