import { describe, expect, it } from "vitest";
import type { SimEvent } from "../events/types";
import {
	COL_W,
	chipCap,
	chipsShown,
	columns,
	GAP_BIG,
	GAP_SMALL,
	laneHeights,
	MIN_CAP,
	onFieldRuns,
	rowsH,
} from "./layout";
import { modelFromEvents } from "./model";

const team = [
	{ key: "bennett", element: "pyro" },
	{ key: "xiangling", element: "pyro" },
];

const ev = (over: Partial<SimEvent> & Record<string, unknown>): SimEvent => ({
	type: "action",
	frame: 0,
	characterIndex: -1,
	message: "",
	raw: null,
	...over,
});

describe("columns", () => {
	const model = modelFromEvents(
		[
			ev({ frame: 0, characterIndex: 0 }),
			ev({ frame: 0, characterIndex: 1 }),
			ev({ frame: 1, type: "damage", damage: 1 }),
			ev({ frame: 10 }),
		],
		team,
	);

	it("makes one column per frame with events, a cell per lane", () => {
		const cols = columns(model, new Set(["action", "damage"]));
		expect(cols.map((c) => c.frame)).toEqual([0, 1, 10]);
		expect(cols[0].cells.map((c) => c.length)).toEqual([0, 1, 1]);
	});

	it("spaces neighbouring frames narrowly and skipped frames widely", () => {
		const cols = columns(model, new Set(["action", "damage"]));
		expect(cols.map((c) => c.x)).toEqual([
			0,
			COL_W + GAP_SMALL,
			2 * COL_W + GAP_SMALL + GAP_BIG,
		]);
		expect(cols.map((c) => c.gap)).toEqual([0, 1, 9]);
	});

	it("leaves out filtered types, and frames left empty by them", () => {
		const cols = columns(model, new Set(["action"]));
		expect(cols.map((c) => c.frame)).toEqual([0, 10]);
	});
});

describe("onFieldRuns", () => {
	const model = modelFromEvents(
		[
			ev({ type: "field", characterIndex: 0, frame: 0, end: 20 }),
			ev({ type: "field", characterIndex: 1, frame: 20, end: 40 }),
			ev({ frame: 0 }),
			ev({ frame: 10 }),
			ev({ frame: 20 }),
			ev({ frame: 40 }),
		],
		team,
	);

	it("joins neighbouring columns where the same lane is on field", () => {
		const cols = columns(model, new Set(["action"]));
		expect(onFieldRuns(cols, model.onField)).toEqual([
			{ lane: 1, x0: cols[0].x, x1: cols[1].x + COL_W },
			{ lane: 2, x0: cols[2].x, x1: cols[3].x + COL_W },
		]);
	});
});

describe("chipCap", () => {
	it("fits every lane in the available height", () => {
		const cap = chipCap(2000, 5);
		expect(cap).toBeGreaterThan(MIN_CAP);
		expect(5 * rowsH(cap)).toBeLessThanOrEqual(2000);
	});

	it("never drops below the minimum", () => {
		expect(chipCap(0, 5)).toBe(MIN_CAP);
	});
});

describe("chipsShown", () => {
	it("shows everything when it fits or the lane is expanded", () => {
		expect(chipsShown(3, 4, false)).toBe(3);
		expect(chipsShown(9, 4, true)).toBe(9);
	});

	it("leaves the last row for the +N more button", () => {
		expect(chipsShown(9, 4, false)).toBe(3);
	});
});

describe("laneHeights", () => {
	it("fits chip rows, keeping character lanes tall enough for the portrait", () => {
		const portrait = 60;
		const [sim, char, busy] = laneHeights([0, 0, 10], portrait);
		expect(sim).toBeLessThan(portrait);
		expect(char).toBe(portrait + 8);
		expect(busy).toBe(rowsH(10));
	});
});
