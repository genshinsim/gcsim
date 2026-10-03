import { describe, expect, it } from "vitest";
import { DefaultSampleOptions } from "./SampleOptions";
import { loadSampleFilter, saveSampleFilter } from "./sampleFilter";

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

describe("sample filter", () => {
	it("falls back to defaults when nothing is stored", () => {
		expect(loadSampleFilter(fakeStorage())).toEqual(DefaultSampleOptions);
	});

	it("loads a filter saved under the pre-existing key", () => {
		const storage = fakeStorage({
			"gcsim-sample-settings": JSON.stringify(["damage"]),
		});
		expect(loadSampleFilter(storage)).toEqual(["damage"]);
	});

	it("ignores an unparseable filter", () => {
		const storage = fakeStorage({ "gcsim-sample-settings": "[nope" });
		expect(loadSampleFilter(storage)).toEqual(DefaultSampleOptions);
	});

	it("ignores a saved value that isn't a list of categories", () => {
		const storage = fakeStorage({ "gcsim-sample-settings": '"damage"' });
		expect(loadSampleFilter(storage)).toEqual(DefaultSampleOptions);
	});

	it("round-trips a saved filter", () => {
		const storage = fakeStorage();
		saveSampleFilter(storage, ["energy", "hitlag"]);
		expect(storage.data.get("gcsim-sample-settings")).toBe(
			JSON.stringify(["energy", "hitlag"]),
		);
		expect(loadSampleFilter(storage)).toEqual(["energy", "hitlag"]);
	});
});
