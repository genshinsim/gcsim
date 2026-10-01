import { describe, expect, it } from "vitest";
import { characterChartData } from "./ByCharacter";

describe("characterChartData", () => {
	it("gives each character its share of the team's mean DPS", () => {
		const { data, total } = characterChartData(
			[{ mean: 300 }, { mean: 100 }, { mean: 600 }],
			["a", "b", "c"],
		);

		const shares = data.map((d) => (d.data.mean ?? 0) / total);
		expect(shares).toEqual([0.3, 0.1, 0.6]);
		expect(shares.reduce((a, b) => a + b, 0)).toBeCloseTo(1);
	});

	it("treats a missing mean as zero", () => {
		const { total } = characterChartData([{ mean: 50 }, {}], ["a", "b"]);
		expect(total).toBe(50);
	});

	it("is empty without stats", () => {
		expect(characterChartData(undefined, ["a"])).toEqual({
			data: [],
			keys: [],
			xMax: 0,
			total: 0,
		});
	});
});
