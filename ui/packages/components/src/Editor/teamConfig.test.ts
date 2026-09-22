import type { model } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { StatToIndexMap } from "../Cards";
import { cfgFromTeam, charToCfg } from "./teamConfig";

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

describe("charToCfg", () => {
	it("serializes a character's lvl/weapon/sets/nonzero-stats", () => {
		expect(charToCfg(amber())).toBe(
			"amber char lvl=80/90 cons=2 talent=6,6,6;\n" +
				'amber add weapon="dullblade" refine=1 lvl=1/20;\n' +
				'amber add set="gladiatorsfinale" count=4;\n' +
				"amber add stats atk=311 cr=0.5;\n",
		);
	});
});

describe("cfgFromTeam", () => {
	it("purges existing character lines and prepends the regenerated team", () => {
		const cfg = [
			"options iteration=1000;",
			"oldchar char lvl=1/1 cons=0 talent=1,1,1;",
			'oldchar add weapon="dullblade" refine=1 lvl=1/20;',
			"target lvl=100 hp=1000;",
		].join("\n");

		const out = cfgFromTeam([amber()], cfg);

		expect(out).toContain("amber char lvl=80/90");
		expect(out).not.toContain("oldchar");
		expect(out).toContain("options iteration=1000;");
		expect(out).toContain("target lvl=100 hp=1000;");
		expect(out.startsWith("amber char")).toBe(true);
	});
});
