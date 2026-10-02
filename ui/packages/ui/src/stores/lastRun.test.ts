import type { model } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { fakeStorage } from "./fakeStorage";
import { lastRunStore } from "./lastRun";

const result = { sim_version: "abc" } as model.SimulationResult;
const raw = JSON.stringify(result);

describe("lastRunStore", () => {
	it("loads nothing when nothing is stored", () => {
		expect(lastRunStore(fakeStorage()).load()).toBeNull();
	});

	it("loads a result saved under the pre-existing redux keys without its bytes or hash", () => {
		const storage = fakeStorage({
			"redux-local-results": raw,
			"redux-local-results-hash": "hash-1",
		});
		expect(lastRunStore(storage).load()).toEqual({
			result,
			raw: null,
			hash: null,
		});
	});

	it("round-trips a saved run with its exact bytes", () => {
		const store = lastRunStore(fakeStorage());
		const spaced = `{ "sim_version": "abc" }`;
		store.save({ result, raw: spaced, hash: "hash-2" });
		expect(store.load()).toEqual({ result, raw: spaced, hash: "hash-2" });
	});

	it("prefers the saved bytes over a pre-existing redux result", () => {
		const storage = fakeStorage({ "redux-local-results": "{}" });
		const store = lastRunStore(storage);
		store.save({ result, raw, hash: "hash-3" });
		expect(store.load()).toEqual({ result, raw, hash: "hash-3" });
	});

	it("ignores an unparseable result", () => {
		const storage = fakeStorage({ "local-results-raw": "{nope" });
		expect(lastRunStore(storage).load()).toBeNull();
	});
});
