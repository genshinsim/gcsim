import { describe, expect, it } from "vitest";
import { uptimeFormat } from "./uptimeFormat";

describe("uptimeFormat", () => {
	it("adds a percent unit without rescaling the 0-100 value", () => {
		expect(uptimeFormat("en")(45)).toBe("45.00%");
		expect(uptimeFormat("en")(100)).toBe("100.00%");
	});

	it("rounds to two decimals", () => {
		expect(uptimeFormat("en")(12.3456)).toBe("12.35%");
	});

	it("uses the locale's number format", () => {
		expect(uptimeFormat("de")(1234.5)).toBe("1.234,50%");
	});

	it("leaves a missing value empty", () => {
		expect(uptimeFormat("en")(undefined)).toBeUndefined();
	});
});
