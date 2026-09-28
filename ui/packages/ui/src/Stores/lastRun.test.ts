import type { model } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { lastRunStore } from "./lastRun";

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

const result = { sim_version: "abc" } as model.SimulationResult;

describe("lastRunStore", () => {
	it("loads nothing when nothing is stored", () => {
		expect(lastRunStore(fakeStorage()).load()).toBeNull();
	});

	it("loads a result saved under the pre-existing redux keys", () => {
		const storage = fakeStorage({
			"redux-local-results": JSON.stringify(result),
			"redux-local-results-hash": "hash-1",
		});
		expect(lastRunStore(storage).load()).toEqual({ result, hash: "hash-1" });
	});

	it("round-trips a saved run", () => {
		const store = lastRunStore(fakeStorage());
		store.save({ result, hash: "hash-2" });
		expect(store.load()).toEqual({ result, hash: "hash-2" });
	});

	it("ignores an unparseable result", () => {
		const storage = fakeStorage({ "redux-local-results": "{nope" });
		expect(lastRunStore(storage).load()).toBeNull();
	});
});
