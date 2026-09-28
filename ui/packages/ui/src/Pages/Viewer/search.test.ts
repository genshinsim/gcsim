import { describe, expect, it } from "vitest";
import { legacyHashSearch, validateViewerSearch } from "./search";

describe("validateViewerSearch", () => {
	it("keeps a known tab and a seed", () => {
		expect(validateViewerSearch({ tab: "sample", seed: "42" })).toEqual({
			tab: "sample",
			seed: "42",
		});
	});

	it("drops an unknown tab", () => {
		expect(validateViewerSearch({ tab: "analyze" })).toEqual({});
	});

	it("reads a numeric seed as a string", () => {
		expect(validateViewerSearch({ seed: 42 })).toEqual({ seed: "42" });
	});

	it("keeps a seed past the safe integer range exactly", () => {
		expect(validateViewerSearch({ seed: "18446744073709551615" })).toEqual({
			seed: "18446744073709551615",
		});
	});

	it("drops a seed that is not an integer", () => {
		expect(validateViewerSearch({ seed: "abc" })).toEqual({});
		expect(validateViewerSearch({ seed: 1.5 })).toEqual({});
	});
});

describe("legacyHashSearch", () => {
	it("maps #tab= and #sample= to tab and seed", () => {
		expect(legacyHashSearch("tab=sample&sample=123")).toEqual({
			tab: "sample",
			seed: "123",
		});
	});

	it("ignores a hash with neither key", () => {
		expect(legacyHashSearch("")).toBeNull();
		expect(legacyHashSearch("dps-card")).toBeNull();
	});
});
