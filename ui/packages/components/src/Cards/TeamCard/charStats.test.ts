import type { model } from "@gcsim/types";
import type { TFunction } from "i18next";
import { describe, expect, it } from "vitest";
import { ConsolidateCharStats, StatToIndexMap } from "./charStats";

const t = ((key: string) => key) as unknown as TFunction;

function char(stats: number[], snapshot: number[] = []): model.Character {
	return {
		name: "amber",
		level: 90,
		max_level: 90,
		element: "pyro",
		cons: 0,
		weapon: { name: "w", refine: 1, level: 90, max_level: 90 },
		talents: { attack: 9, skill: 9, burst: 9 },
		sets: {},
		stats,
		snapshot,
	};
}

describe("ConsolidateCharStats", () => {
	it("keeps only rows with nonzero values, in stat order", () => {
		const stats = new Array(22).fill(0);
		stats[StatToIndexMap.ATK] = 100;
		stats[StatToIndexMap.CR] = 0.3;

		const { stats: blocks, maxRows } = ConsolidateCharStats(t, [char(stats)]);

		expect(blocks.amber.map((b) => b.key)).toEqual(["atk", "cr"]);
		expect(blocks.amber[0]).toMatchObject({ t: "both", flat: 100, percent: 0 });
		expect(blocks.amber[1]).toMatchObject({ t: "%", flat: 0, percent: 0.3 });
		expect(maxRows).toBe(2);
	});

	it("produces empty blocks for an all-zero character", () => {
		const { stats, snapshot, maxRows } = ConsolidateCharStats(t, [
			char(new Array(22).fill(0)),
		]);
		expect(stats.amber).toEqual([]);
		expect(snapshot.amber).toEqual([]);
		expect(maxRows).toBe(0);
	});

	it("reads snapshot values into the snapshot blocks", () => {
		const snapshot = new Array(22).fill(0);
		snapshot[StatToIndexMap.HP] = 31204;

		const { snapshot: blocks } = ConsolidateCharStats(t, [
			char(new Array(22).fill(0), snapshot),
		]);
		expect(blocks.amber.map((b) => b.key)).toEqual(["hp"]);
		expect(blocks.amber[0]).toMatchObject({ t: "f", flat: 31204 });
	});
});
