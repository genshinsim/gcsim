import { describe, expect, it } from "vitest";
import { fakeStorage } from "./fakeStorage";
import { loadSampleOnLoad, saveSampleOnLoad } from "./prefs";

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
