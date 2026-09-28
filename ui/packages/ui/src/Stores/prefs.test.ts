import { describe, expect, it } from "vitest";
import { fakeStorage } from "./fakeStorage";
import { prefsStore } from "./prefs";

describe("prefsStore", () => {
	it("defaults to no sample on load and settings closed", () => {
		expect(prefsStore(fakeStorage()).get()).toEqual({
			sampleOnLoad: false,
			settingsOpen: false,
		});
	});

	it("persists sampleOnLoad but not settingsOpen", () => {
		const storage = fakeStorage();
		prefsStore(storage).set({ sampleOnLoad: true, settingsOpen: true });
		expect(prefsStore(storage).get()).toEqual({
			sampleOnLoad: true,
			settingsOpen: false,
		});
	});

	it("loads sampleOnLoad saved by the redux store", () => {
		const storage = fakeStorage({
			"redux-app-data": JSON.stringify({ sampleOnLoad: true, cfg: "" }),
		});
		expect(prefsStore(storage).get().sampleOnLoad).toBe(true);
	});
});
