import { describe, expect, it } from "vitest";
import { actionColorScale } from "./DataColors";

describe("actionColorScale", () => {
	it("gives every present category a distinct series slot when <= 7", () => {
		const present = ["attack", "skill", "burst", "dash", "swap"];
		const scale = actionColorScale(present);

		const colors = present.map((k) => scale(k));
		expect(new Set(colors).size).toBe(present.length);
	});

	it("colors present categories in order (first -> s1, second -> s2, ...)", () => {
		const present = ["attack", "skill", "burst", "dash", "swap"];
		const scale = actionColorScale(present);

		expect(present.map((k) => scale(k))).toEqual([
			"var(--g-s1)",
			"var(--g-s2)",
			"var(--g-s3)",
			"var(--g-s4)",
			"var(--g-s5)",
		]);
	});

	it("wraps past 7 present categories without error (8th -> s1)", () => {
		const present = ["a", "b", "c", "d", "e", "f", "g", "h"];
		const scale = actionColorScale(present);

		expect(scale("h")).toBe(scale("a"));
		expect(scale("h")).toBe("var(--g-s1)");
	});
});
