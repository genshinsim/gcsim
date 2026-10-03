import { describe, expect, it } from "vitest";
import type { SimEvent } from "../events/types";
import { sampleFixture } from "../testdata";
import {
	chipLabel,
	chipText,
	compactDamage,
	damageOf,
	lowerBound,
	modelFromEvents,
	modelFromSample,
	onFieldLaneAt,
} from "./model";

const team = [
	{ key: "bennett", element: "pyro" },
	{ key: "xiangling", element: "pyro" },
];

const ev = (over: Partial<SimEvent> & Record<string, unknown>): SimEvent => ({
	type: "sim",
	frame: 0,
	characterIndex: -1,
	message: "",
	raw: null,
	...over,
});

describe("modelFromEvents", () => {
	it("puts sim events in lane 0 and characters after it", () => {
		const m = modelFromEvents(
			[ev({ characterIndex: -1 }), ev({ characterIndex: 1 })],
			team,
		);
		expect(m.chips.map((c) => c.lane)).toEqual([0, 2]);
	});

	it("turns field stints into on-field segments, not chips", () => {
		const m = modelFromEvents(
			[
				ev({ type: "stint", characterIndex: 0, frame: 0, end: 60 }),
				ev({ type: "stint", characterIndex: 1, frame: 60, end: 90 }),
				ev({ frame: 90 }),
			],
			team,
		);
		expect(m.onField).toEqual([
			{ lane: 1, start: 0, end: 60 },
			{ lane: 2, start: 60, end: 90 },
		]);
		expect(m.chips.map((c) => c.event.type)).toEqual(["sim"]);
	});

	it("adds an expiry chip where a status ends inside the sample", () => {
		const m = modelFromEvents(
			[
				ev({
					type: "status",
					key: "bennettburst",
					message: "bennettburst status added",
					frame: 10,
					end: 50,
				}),
				ev({ frame: 100 }),
			],
			team,
		);
		const chips = m.chips.map((c) => [c.frame, chipText(c)]);
		expect(chips).toEqual([
			[10, "bennettburst status added"],
			[50, "bennettburst expired"],
			[100, ""],
		]);
	});

	it("adds no expiry chip for a status that outlasts the sample", () => {
		const m = modelFromEvents(
			[
				ev({ type: "status", key: "a", frame: 10, end: 500 }),
				ev({
					type: "status",
					key: "b",
					frame: 10,
					end: Number.POSITIVE_INFINITY,
				}),
				ev({ frame: 100 }),
			],
			team,
		);
		expect(m.chips.filter((c) => c.expired)).toEqual([]);
	});

	it("adds no expiry chip for a status with no key", () => {
		const m = modelFromEvents(
			[ev({ type: "status", key: "", frame: 10, end: 50 }), ev({ frame: 100 })],
			team,
		);
		expect(m.chips.filter((c) => c.expired)).toEqual([]);
	});

	it("counts damage only on the hit itself", () => {
		const m = modelFromEvents(
			[ev({ type: "damage", characterIndex: 0, damage: 1000 })],
			team,
		);
		expect(m.chips.map(damageOf)).toEqual([1000]);
	});
});

describe("compactDamage", () => {
	it.each([
		[8, "8"],
		[999, "999"],
		[1234, "1.23K"],
		[12345, "12.3K"],
		[999_400, "999K"],
		[999_600, "1M"],
		[1_234_567, "1.23M"],
		[8815.17, "8.82K"],
		[2_345_678_901, "2.35B"],
	])("shows %d as %s", (n, want) => {
		expect(compactDamage(n)).toBe(want);
	});
});

describe("chipLabel", () => {
	const chipOf = (over: Partial<SimEvent> & Record<string, unknown>) =>
		modelFromEvents([ev(over)], team).chips[0];

	it("leads a damage chip with the compact amount", () => {
		const c = chipOf({
			type: "damage",
			damage: 12345.6,
			attack: "Pyronado",
			mods: "vaporize crit",
			message: "Pyronado [12,346] (vaporize crit)",
		});
		expect(chipLabel(c)).toBe("12.3K Pyronado (vaporize crit)");
		expect(chipText(c)).toBe("Pyronado [12,346] (vaporize crit)");
	});

	it("leaves out the brackets when a hit has no modifiers", () => {
		const c = chipOf({
			type: "damage",
			damage: 999,
			attack: "Normal 1",
			mods: "",
		});
		expect(chipLabel(c)).toBe("999 Normal 1");
	});

	it("labels any other chip with its text", () => {
		const c = chipOf({ type: "action", message: "burst" });
		expect(chipLabel(c)).toBe("burst");
	});
});

describe("modelFromSample", () => {
	it("makes the sim lane then one lane per character", () => {
		const m = modelFromSample(sampleFixture);
		expect(m.lanes.map((l) => l.key)).toEqual([
			"sim",
			...(sampleFixture.character_details ?? []).map((c) => c.name),
		]);
	});
});

describe("lowerBound", () => {
	it("finds the first index whose key is at least the value", () => {
		const xs = [1, 3, 3, 7];
		expect(lowerBound(xs, 0, (x) => x)).toBe(0);
		expect(lowerBound(xs, 3, (x) => x)).toBe(1);
		expect(lowerBound(xs, 4, (x) => x)).toBe(3);
		expect(lowerBound(xs, 8, (x) => x)).toBe(4);
	});
});

describe("onFieldLaneAt", () => {
	const onField = [
		{ lane: 1, start: 0, end: 60 },
		{ lane: 2, start: 60, end: 90 },
	];

	it("gives the lane on field at a frame, the newcomer on a swap frame", () => {
		expect(onFieldLaneAt(onField, 0)).toBe(1);
		expect(onFieldLaneAt(onField, 59)).toBe(1);
		expect(onFieldLaneAt(onField, 60)).toBe(2);
	});

	it("gives -1 before anyone is on field", () => {
		expect(onFieldLaneAt([{ lane: 1, start: 10, end: 20 }], 5)).toBe(-1);
	});
});
