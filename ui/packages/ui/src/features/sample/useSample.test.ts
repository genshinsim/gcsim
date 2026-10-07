import { describe, expect, it } from "vitest";
import { autoSampleKey, autoSampleSeed } from "./useSample";

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

describe("autoSampleKey", () => {
	it("is null until both the seed and config are known", () => {
		expect(autoSampleKey(undefined, "1")).toBeNull();
		expect(autoSampleKey("cfg", null)).toBeNull();
	});

	it("differs between runs of the same config, which get a new seed", () => {
		expect(autoSampleKey("cfg", "1")).not.toBe(autoSampleKey("cfg", "2"));
		expect(autoSampleKey("cfg", "1")).toBe(autoSampleKey("cfg", "1"));
	});
});
