import { describe, expect, it } from "vitest";
import type { SimEvent } from "../events/types";
import { columns } from "./layout";
import { modelFromEvents } from "./model";
import {
	CHIP_H,
	frameOfY,
	MIN_ROW_H,
	slimChip,
	V_GAP_BIG,
	vRows,
	wrapLines,
	wrappedChip,
	yOfFrame,
} from "./vlayout";

const team = [{ key: "bennett", element: "pyro" }];

const ev = (over: Partial<SimEvent> & Record<string, unknown>): SimEvent => ({
	type: "action",
	frame: 0,
	characterIndex: 0,
	message: "",
	raw: null,
	...over,
});

const burst = (frame: number, n: number) =>
	Array.from({ length: n }, () => ev({ frame }));

const model = modelFromEvents(
	[...burst(0, 8), ...burst(1, 1), ...burst(30, 2)],
	team,
);
const cols = columns(model, new Set(["action"]));
const lanes = [0, 1];
// 4px padding above and below, a 1px rule under each row
const rowH = (chips: number) => 8 + chips * CHIP_H + (chips - 1) * 2 + 1;

describe("vRows", () => {
	it("caps slim rows, keeping a slot for the +N more button", () => {
		const rows = vRows(cols, lanes, {
			cap: 4,
			expanded: new Set(),
			chipH: slimChip,
		});
		expect(rows[0].h).toBe(rowH(4));
	});

	it("grows an expanded row to fit every chip", () => {
		const rows = vRows(cols, lanes, {
			cap: 4,
			expanded: new Set([0]),
			chipH: slimChip,
		});
		expect(rows[0].h).toBe(rowH(8));
	});

	it("shows every chip when uncapped", () => {
		const rows = vRows(cols, lanes, {
			cap: null,
			expanded: new Set(),
			chipH: slimChip,
		});
		expect(rows[0].h).toBe(rowH(8));
	});

	it("keeps short rows tall enough for the time label", () => {
		const rows = vRows(cols, lanes, {
			cap: 4,
			expanded: new Set(),
			chipH: slimChip,
		});
		expect(rows[1].h).toBe(MIN_ROW_H + 1);
	});

	it("adds a gap band above a row that skips frames, and stacks rows", () => {
		const rows = vRows(cols, lanes, {
			cap: 4,
			expanded: new Set(),
			chipH: slimChip,
		});
		expect(rows[1].gap).toBe(0);
		expect(rows[2].gap).toBe(V_GAP_BIG);
		expect(rows[2].h).toBe(V_GAP_BIG + rowH(2));
		expect(rows[2].y).toBe(rows[1].y + rows[1].h);
	});

	it("sizes rows by each chip's own height", () => {
		const rows = vRows(cols, lanes, {
			cap: null,
			expanded: new Set(),
			chipH: () => 32,
		});
		expect(rows[2].h).toBe(V_GAP_BIG + 8 + 32 * 2 + 2 + 1);
	});
});

describe("yOfFrame / frameOfY", () => {
	const rows = vRows(cols, lanes, {
		cap: 4,
		expanded: new Set(),
		chipH: slimChip,
	});

	it("puts a frame with events at the top of its row, below the gap", () => {
		expect(yOfFrame(rows, 30)).toBe(rows[2].y + V_GAP_BIG);
	});

	it("interpolates frames between rows and inverts", () => {
		const y = yOfFrame(rows, 15.5);
		expect(y).toBeGreaterThan(yOfFrame(rows, 1));
		expect(y).toBeLessThan(yOfFrame(rows, 30));
		expect(frameOfY(rows, y)).toBeCloseTo(15.5);
	});

	it("clamps outside the rows", () => {
		expect(frameOfY(rows, -50)).toBe(0);
		expect(frameOfY(rows, 1e6)).toBe(30);
	});
});

describe("wrapLines", () => {
	const width = (s: string) => s.length;

	it("fits text on one line when it can", () => {
		expect(wrapLines("aaa bbb", 7, width)).toBe(1);
	});

	it("breaks at spaces", () => {
		expect(wrapLines("aaa bbb", 5, width)).toBe(2);
	});

	it("breaks after hyphens", () => {
		expect(wrapLines("aa-bbb", 4, width)).toBe(2);
	});

	it("breaks a word too long for a line by itself", () => {
		expect(wrapLines("aaaaaaa", 3, width)).toBe(3);
	});

	it("grows a wrapped chip by a line height per extra line", () => {
		const chip = model.chips[0];
		const one = wrappedChip(1000, () => "a b", width)(chip);
		const two = wrappedChip(
			30,
			() => "aaaaaa bbbbbb",
			(s) => 3 * s.length,
		)(chip);
		expect(one).toBe(CHIP_H);
		expect(two).toBeGreaterThan(CHIP_H);
	});
});
