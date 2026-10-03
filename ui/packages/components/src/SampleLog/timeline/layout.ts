import {
	type Chip,
	onFieldLaneAt,
	type Stint,
	type TimelineModel,
} from "./model";

export const COL_W = 168;
export const GAP_SMALL = 3;
export const GAP_BIG = 22;
export const CHIP_H = 18;
export const CHIP_GAP = 2;
export const MIN_CAP = 4;
export const MAX_PORTRAIT = 120;
export const HEAD_H = 34;
export const CELL_PAD = 8;
const SIM_LANE_MIN_H = 30;
const PORTRAIT_PAD = 8;
const STRIP_BORDERS_AND_SCROLLBAR = 14;

export const rowsH = (n: number) =>
	2 * CELL_PAD + n * (CHIP_H + CHIP_GAP) - CHIP_GAP;

export type Column = {
	frame: number;
	x: number;
	/** frames since the previous column */
	gap: number;
	/** chips per lane */
	cells: Chip[][];
};

export function columns(model: TimelineModel, enabled: Set<string>): Column[] {
	const out: Column[] = [];
	for (const chip of model.chips) {
		if (!enabled.has(chip.event.type)) {
			continue;
		}
		let col = out[out.length - 1];
		if (col == null || col.frame !== chip.frame) {
			const gap = col == null ? 0 : chip.frame - col.frame;
			const x =
				col == null ? 0 : col.x + COL_W + (gap > 1 ? GAP_BIG : GAP_SMALL);
			col = { frame: chip.frame, x, gap, cells: model.lanes.map(() => []) };
			out.push(col);
		}
		col.cells[chip.lane].push(chip);
	}
	return out;
}

export type OnFieldRun = { lane: number; x0: number; x1: number };

export function onFieldRuns(cols: Column[], onField: Stint[]): OnFieldRun[] {
	const out: OnFieldRun[] = [];
	let prev = -1;
	for (const c of cols) {
		const lane = onFieldLaneAt(onField, c.frame);
		if (lane >= 0 && lane === prev) {
			out[out.length - 1].x1 = c.x + COL_W;
		} else if (lane >= 0) {
			out.push({ lane, x0: c.x, x1: c.x + COL_W });
		}
		prev = lane;
	}
	return out;
}

/** chips a collapsed cell shows so that every lane fits in `height` */
export function chipCap(height: number, lanes: number): number {
	const perLane = (height - HEAD_H - STRIP_BORDERS_AND_SCROLLBAR) / lanes;
	return Math.max(
		MIN_CAP,
		Math.floor((perLane - rowsH(0)) / (CHIP_H + CHIP_GAP)),
	);
}

export const portraitSize = (cap: number) =>
	Math.min(MAX_PORTRAIT, rowsH(cap) - PORTRAIT_PAD);

/** chips drawn in a cell; a collapsed cell that overflows keeps its last row for "+N more" */
export function chipsShown(
	count: number,
	rows: number,
	expanded: boolean,
): number {
	return !expanded && count > rows ? rows - 1 : count;
}

export function laneHeights(rows: number[], portrait: number): number[] {
	return rows.map((r, i) =>
		Math.max(i === 0 ? SIM_LANE_MIN_H : portrait + PORTRAIT_PAD, rowsH(r)),
	);
}

/** pixel range to render, padded and quantized so scrolling doesn't re-render per pixel */
export function visibleWindow(left: number, width: number) {
	const q = 400;
	const start = Math.floor(left / q) * q - q;
	return [start, start + width + 3 * q] as const;
}
