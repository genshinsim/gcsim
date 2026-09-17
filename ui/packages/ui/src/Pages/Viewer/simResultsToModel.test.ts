import type { SimResults } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { simResultsToModel } from "./simResultsToModel";

const sample = {
	schema_version: { major: "4", minor: "2" },
	sim_version: "test",
	mode: 1,
	config_file: "cfg",
	character_details: [
		{
			name: "bennett",
			element: "pyro",
			level: 90,
			max_level: 90,
			cons: 6,
			weapon: { name: "aquila", refine: 1, level: 90, max_level: 90 },
			talents: { attack: 9, skill: 9, burst: 9 },
			stats: [1, 2, 3],
			snapshot: [4, 5, 6],
			sets: { noblesse: 4 },
		},
	],
	target_details: [
		{
			level: 100,
			hp: 1000,
			resist: { pyro: 0.1 },
			position: { x: 0, y: 0, r: 1 },
			particle_drop_threshold: 0,
			particle_drop_count: 0,
			particle_element: "electro",
			name: "hilichurl",
			modified: false,
		},
	],
	statistics: {
		iterations: 1000,
		dps: {
			min: 1,
			max: 3,
			mean: 2,
			sd: 0.5,
			q1: 1.5,
			q2: 2,
			q3: 2.5,
			histogram: [1, 2],
		},
		target_dps: { "0": { min: 1, max: 3, mean: 2, sd: 0.5 } },
		dps_by_target: [{ targets: { "0": { min: 1, max: 3, mean: 2, sd: 0.5 } } }],
		cumu_damage: {
			bucket_size: 60,
			targets: {
				"0": { overall: { min: [1], max: [3], q1: [1], q2: [2], q3: [3] } },
			},
		},
		shields: {
			bennett: {
				hp: { pyro: { min: 1, max: 3, mean: 2, sd: 0.5 } },
				uptime: { min: 0, max: 1, mean: 0.5, sd: 0.1 },
			},
		},
	},
} as unknown as SimResults;

describe("simResultsToModel", () => {
	const result = simResultsToModel(sample);

	it("keeps mode as its SimMode numeric value", () => {
		expect(result.mode).toBe(1);
	});

	it("emits particle_element as a string", () => {
		const enemy = result.target_details?.[0];
		expect(typeof enemy?.particle_element).toBe("string");
		expect(enemy?.particle_element).toBe("electro");
	});

	it("coerces a numeric particle_element to a string", () => {
		const numeric = simResultsToModel({
			target_details: [{ particle_element: 1 }],
		} as unknown as SimResults);
		expect(numeric.target_details?.[0]?.particle_element).toBe("1");
	});

	it("reconciles target-keyed maps into the model shape", () => {
		expect(result.statistics?.target_dps?.[0]?.mean).toBe(2);
		expect(result.statistics?.dps_by_target?.[0]?.targets?.[0]?.mean).toBe(2);
		expect(result.statistics?.cumu_damage?.targets?.[0]?.overall?.min).toEqual([
			1,
		]);
	});

	it("emits shield hp as a plain string-keyed object", () => {
		const hp = result.statistics?.shields?.bennett?.hp;
		expect(hp).not.toBeInstanceOf(Map);
		expect(hp?.pyro?.mean).toBe(2);
	});

	it("passes summary stats and metadata through unchanged", () => {
		expect(result.statistics?.dps?.q2).toBe(2);
		expect(result.statistics?.iterations).toBe(1000);
		expect(result.character_details?.[0]?.name).toBe("bennett");
	});
});
