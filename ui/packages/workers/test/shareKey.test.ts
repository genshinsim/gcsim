import { describe, expect, it } from "vitest";
import { newShareKey, shareKeyPrefix } from "../src/share/shareKey";

const SUFFIX = /^[6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz]{8}$/;

describe("shareKeyPrefix", () => {
	it.each([
		[
			"a 4-character team",
			["raidenshogun", "kaedeharakazuha", "sangonomiyakokomi", "furina"],
			"raidenshogun-kaedeharakazuha-sangonomiyakokomi-furina",
		],
		["a 1-character team", ["nahida"], "nahida"],
		[
			"names that need cleaning",
			["Hu Tao", "aether_anemo", "Mavuika!"],
			"hutao-aetheranemo-mavuika",
		],
		["more than 4 names", ["a", "b", "c", "d", "e"], "a-b-c-d"],
		["no usable names", ["", "!!", "_"], ""],
		["no names", [], ""],
	])("handles %s", (_, names, want) => {
		expect(shareKeyPrefix(names)).toBe(want);
	});
});

describe("newShareKey", () => {
	it("appends an 8-character suffix to the prefix", () => {
		const key = newShareKey("nahida");
		expect(key.startsWith("nahida-")).toBe(true);
		expect(key.slice("nahida-".length)).toMatch(SUFFIX);
	});

	it("is just the suffix when there is no prefix", () => {
		expect(newShareKey("")).toMatch(SUFFIX);
	});
});
