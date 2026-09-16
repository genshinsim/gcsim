import type { ParsedResult } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { toParsedTeam } from "./parsedTeam";
import { parsedResult } from "./testExecutor";

describe("toParsedTeam", () => {
	it("maps parsed character profiles into team characters", () => {
		const team = toParsedTeam(parsedResult(["amber", "xiangling"]));
		expect(team).toHaveLength(2);
		expect(team[0]).toMatchObject({
			name: "amber",
			level: 90,
			element: "pyro",
			cons: 0,
		});
		expect(team[0].stats).toEqual([1, 2, 3]);
		expect(team[0].snapshot).toEqual([0, 0, 0]);
	});

	it("returns an empty team when there are no characters", () => {
		const noChars = { ...parsedResult([]), characters: undefined };
		expect(toParsedTeam(noChars as unknown as ParsedResult)).toEqual([]);
	});
});
