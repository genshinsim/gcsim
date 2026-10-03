import { describe, expect, it } from "vitest";
import { AllSampleOptions } from "./SampleOptions";
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
	it("shows every category when nothing is stored", () => {
		expect(loadSampleFilter(fakeStorage())).toEqual(AllSampleOptions);
	});

	it("ignores a filter saved by the old viewer", () => {
		const storage = fakeStorage({
			"gcsim-sample-settings": JSON.stringify(["damage"]),
		});
		expect(loadSampleFilter(storage)).toEqual(AllSampleOptions);
	});

	it("ignores an unparseable filter", () => {
		const storage = fakeStorage({ "gcsim-sample-log-filter": "[nope" });
		expect(loadSampleFilter(storage)).toEqual(AllSampleOptions);
	});

	it("ignores a saved value that isn't a list of categories", () => {
		const storage = fakeStorage({ "gcsim-sample-log-filter": '"damage"' });
		expect(loadSampleFilter(storage)).toEqual(AllSampleOptions);
	});

	it("round-trips a saved filter", () => {
		const storage = fakeStorage();
		saveSampleFilter(storage, ["energy", "hitlag"]);
		expect(storage.data.get("gcsim-sample-log-filter")).toBe(
			JSON.stringify(["energy", "hitlag"]),
		);
		expect(loadSampleFilter(storage)).toEqual(["energy", "hitlag"]);
	});
});
