import type { model } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { StatToIndexMap } from "../Cards";
import { insertCharacterBlock } from "./ImportedCharacterInsert";

function amber(): model.Character {
	const stats = new Array(22).fill(0);
	stats[StatToIndexMap.ATK] = 311;
	stats[StatToIndexMap.CR] = 0.5;
	return {
		name: "amber",
		level: 80,
		max_level: 90,
		element: "pyro",
		cons: 2,
		weapon: { name: "dullblade", refine: 1, level: 1, max_level: 20 },
		talents: { attack: 6, skill: 6, burst: 6 },
		sets: { gladiatorsfinale: 4 },
		stats,
		snapshot: [],
	};
}

describe("insertCharacterBlock", () => {
	it("puts the character's lvl/weapon/sets/nonzero-stats above the config", () => {
		const config = "options iteration=1000;\n";
		expect(
			insertCharacterBlock(config, { key: "amber", character: amber() }),
		).toBe(
			"amber char lvl=80/90 cons=2 talent=6,6,6;\n" +
				'amber add weapon="dullblade" refine=1 lvl=1/20;\n' +
				'amber add set="gladiatorsfinale" count=4;\n' +
				"amber add stats atk=311 cr=0.5;\n" +
				"\n" +
				config,
		);
	});
});
