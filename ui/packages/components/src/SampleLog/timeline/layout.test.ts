import { describe, expect, it } from "vitest";
import type { SimEvent } from "../events/types";
import {
	COL_W,
	chipsShown,
	columns,
	frameOfX,
	GAP_BIG,
	GAP_SMALL,
	HEAD_H,
	MIN_CAP,
	onFieldRuns,
	rowsH,
	splitRows,
	xOfFrame,
} from "./layout";
import { modelFromEvents, onFieldLaneAt } from "./model";

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
		expect(cols.map((c) => c.framesSincePrev)).toEqual([0, 1, 9]);
	});

	it("leaves out filtered types, and frames left empty by them", () => {
		const cols = columns(model, new Set(["action"]));
		expect(cols.map((c) => c.frame)).toEqual([0, 10]);
	});
});

describe("xOfFrame / frameOfX", () => {
	const cols = columns(
		modelFromEvents(
			[ev({ frame: 0 }), ev({ frame: 10 }), ev({ frame: 11 })],
			team,
		),
		new Set(["action"]),
	);

	it("lands on a column's x at its frame", () => {
		expect(cols.map((c) => xOfFrame(cols, c.frame))).toEqual(
			cols.map((c) => c.x),
		);
	});

	it("moves smoothly between columns, so a small pan always moves", () => {
		const x4 = xOfFrame(cols, 4);
		const x5 = xOfFrame(cols, 5);
		expect(x4).toBeGreaterThan(cols[0].x);
		expect(x5).toBeGreaterThan(x4);
		expect(x5).toBeLessThan(cols[1].x);
	});

	it("inverts", () => {
		for (const f of [0, 2.5, 10, 10.5, 11]) {
			expect(frameOfX(cols, xOfFrame(cols, f))).toBeCloseTo(f);
		}
	});

	it("clamps outside the columns", () => {
		expect(xOfFrame(cols, -5)).toBe(0);
		expect(xOfFrame(cols, 99)).toBe(cols[2].x);
		expect(frameOfX(cols, -5)).toBe(0);
		expect(frameOfX(cols, 1e6)).toBe(11);
	});
});

describe("onFieldRuns", () => {
	const model = modelFromEvents(
		[
			ev({ type: "stint", characterIndex: 0, frame: 0, end: 20 }),
			ev({ type: "stint", characterIndex: 1, frame: 20, end: 40 }),
			ev({ frame: 0 }),
			ev({ frame: 10 }),
			ev({ frame: 20 }),
			ev({ frame: 40 }),
		],
		team,
	);

	it("joins neighbouring columns where the same lane is on field", () => {
		const cols = columns(model, new Set(["action"]));
		const onField = cols.map((c) => onFieldLaneAt(model.onField, c.frame));
		expect(onFieldRuns(cols, onField)).toEqual([
			{ lane: 1, x0: cols[0].x, x1: cols[1].x + COL_W },
			{ lane: 2, x0: cols[2].x, x1: cols[3].x + COL_W },
		]);
	});
});

describe("splitRows", () => {
	const used = (rows: number[]) =>
		HEAD_H + 14 + rows.reduce((a, r) => a + rowsH(r), 0);

	it("fills the height evenly, to within a chip row", () => {
		const rows = splitRows(1000, [3, 3, 3, 3, 3]);
		expect(Math.max(...rows) - Math.min(...rows)).toBeLessThanOrEqual(1);
		expect(used(rows)).toBeLessThanOrEqual(1000);
		expect(used(rows)).toBeGreaterThan(1000 - 20);
	});

	it("gives the rows left over to the busiest lanes", () => {
		const even = splitRows(1000, [0, 0, 0, 0, 0]);
		const base = Math.min(...even);
		const spare = even.reduce((a, r) => a + r - base, 0);
		expect(spare).toBeGreaterThan(0);
		const rows = splitRows(1000, [0, 2, 9, 1, 12]);
		expect(rows[4]).toBe(base + 1);
		expect(rows[0]).toBe(base);
	});

	it("never drops below the minimum", () => {
		expect(splitRows(0, [9, 9, 9])).toEqual([MIN_CAP, MIN_CAP, MIN_CAP]);
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
