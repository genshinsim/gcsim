import { describe, expect, it } from "vitest";
import {
	CAP_MAX,
	CAP_MIN,
	loadSampleView,
	type SampleView,
	saveSampleView,
} from "./sampleView";

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

const KEY = "gcsim-sample-log-view";

describe("sample view", () => {
	it("defaults to horizontal with slim, capped chips", () => {
		expect(loadSampleView(fakeStorage(), true)).toEqual({
			orientation: "horizontal",
			slim: true,
			cap: 4,
			overview: true,
		});
	});

	it("opens the overview by default only when the caller says so", () => {
		expect(loadSampleView(fakeStorage(), false).overview).toBe(false);
	});

	it("ignores an unparseable view", () => {
		const storage = fakeStorage({ [KEY]: "{nope" });
		expect(loadSampleView(storage, true).orientation).toBe("horizontal");
	});

	it("keeps the valid fields of a partly broken view", () => {
		const storage = fakeStorage({
			[KEY]: JSON.stringify({ orientation: "diagonal", slim: false, cap: "9" }),
		});
		expect(loadSampleView(storage, true)).toEqual({
			orientation: "horizontal",
			slim: false,
			cap: 4,
			overview: true,
		});
	});

	it("clamps the cap into range", () => {
		const high = fakeStorage({ [KEY]: JSON.stringify({ cap: 99 }) });
		const low = fakeStorage({ [KEY]: JSON.stringify({ cap: 1 }) });
		expect(loadSampleView(high, true).cap).toBe(CAP_MAX);
		expect(loadSampleView(low, true).cap).toBe(CAP_MIN);
	});

	it("round-trips a saved view", () => {
		const storage = fakeStorage();
		const view: SampleView = {
			orientation: "vertical",
			slim: false,
			cap: 7,
			overview: false,
		};
		saveSampleView(storage, view);
		expect(loadSampleView(storage, true)).toEqual(view);
	});
});
