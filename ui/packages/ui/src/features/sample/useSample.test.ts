import { describe, expect, it } from "vitest";
import { autoSampleSeed } from "./useSample";

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
