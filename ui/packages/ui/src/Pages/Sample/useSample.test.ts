import { DefaultSampleOptions } from "@gcsim/components";
import { describe, expect, it } from "vitest";
import {
	autoSampleSeed,
	loadSampleSettings,
	saveSampleSettings,
} from "./useSample";

function fakeStorage(init: Record<string, string> = {}) {
	const data = new Map(Object.entries(init));
	return {
		data,
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => {
			data.set(k, v);
		},
	};
}

describe("sample settings", () => {
	it("falls back to defaults when nothing is stored", () => {
		expect(loadSampleSettings(fakeStorage())).toEqual(DefaultSampleOptions);
	});

	it("loads settings saved under the pre-existing key", () => {
		const storage = fakeStorage({
			"gcsim-sample-settings": JSON.stringify(["damage"]),
		});
		expect(loadSampleSettings(storage)).toEqual(["damage"]);
	});

	it("ignores unparseable settings", () => {
		const storage = fakeStorage({ "gcsim-sample-settings": "[nope" });
		expect(loadSampleSettings(storage)).toEqual(DefaultSampleOptions);
	});

	it("round-trips saved settings", () => {
		const storage = fakeStorage();
		saveSampleSettings(storage, ["energy", "hitlag"]);
		expect(loadSampleSettings(storage)).toEqual(["energy", "hitlag"]);
	});
});

describe("autoSampleSeed", () => {
	it("prefers the linked seed", () => {
		expect(autoSampleSeed("123", false, "456")).toBe("123");
		expect(autoSampleSeed("123", true, "456")).toBe("123");
	});

	it("uses the result's sample seed when sampling on load", () => {
		expect(autoSampleSeed(null, true, "456")).toBe("456");
	});

	it("does not auto sample otherwise", () => {
		expect(autoSampleSeed(null, false, "456")).toBeNull();
		expect(autoSampleSeed(null, true, undefined)).toBeNull();
	});
});
