import { type Chip, lowerBound, type TimelineModel } from "./model";

export const COL_W = 168;
export const GAP_SMALL = 3;
export const GAP_BIG = 22;
const CHIP_H = 18;
const CHIP_GAP = 2;
const ROW_H = CHIP_H + CHIP_GAP;
export const MIN_CAP = 4;
const MAX_PORTRAIT = 120;
export const HEAD_H = 34;
export const CELL_PAD = 8;
const PORTRAIT_PAD = 8;
const LABEL_PAD_X = 12;
const STRIP_BORDERS_AND_SCROLLBAR = 14;
const MORE_BUTTON_ROWS = 1;
const RENDER_STEP_PX = 400;

/** width of the sticky lane-label column for a portrait of this size */
export const gutterWidth = (portrait: number) => portrait + LABEL_PAD_X;

export const secondsLabel = (frame: number, digits = 2) =>
	`${(frame / 60).toFixed(digits)}s`;

export const rowsH = (n: number) => 2 * CELL_PAD + n * ROW_H - CHIP_GAP;

export type Column = {
	frame: number;
	x: number;
	framesSincePrev: number;
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
			const framesSincePrev = col == null ? 0 : chip.frame - col.frame;
			const x =
				col == null
					? 0
					: col.x + COL_W + (framesSincePrev > 1 ? GAP_BIG : GAP_SMALL);
			col = {
				frame: chip.frame,
				x,
				framesSincePrev,
				cells: model.lanes.map(() => []),
			};
			out.push(col);
		}
		col.cells[chip.lane].push(chip);
	}
	return out;
}

/** strip x of a frame, interpolated between the columns around it */
export function xOfFrame(cols: Column[], frame: number): number {
	const i = lowerBound(cols, frame, (c) => c.frame);
	if (i === 0 || i === cols.length) {
		return cols[Math.min(i, cols.length - 1)]?.x ?? 0;
	}
	const [a, b] = [cols[i - 1], cols[i]];
	return a.x + ((frame - a.frame) / (b.frame - a.frame)) * (b.x - a.x);
}

/** the inverse of {@link xOfFrame} */
export function frameOfX(cols: Column[], x: number): number {
	const i = lowerBound(cols, x, (c) => c.x);
	if (i === 0 || i === cols.length) {
		return cols[Math.min(i, cols.length - 1)]?.frame ?? 0;
	}
	const [a, b] = [cols[i - 1], cols[i]];
	return a.frame + ((x - a.x) / (b.x - a.x)) * (b.frame - a.frame);
}

export type OnFieldRun = { lane: number; x0: number; x1: number };

/** `onField[i]` is the lane on field at `cols[i]`, or -1 for none */
export function onFieldRuns(cols: Column[], onField: number[]): OnFieldRun[] {
	const out: OnFieldRun[] = [];
	let prev = -1;
	for (const [i, c] of cols.entries()) {
		const lane = onField[i];
		if (lane >= 0 && lane === prev) {
			out[out.length - 1].x1 = c.x + COL_W;
		} else if (lane >= 0) {
			out.push({ lane, x0: c.x, x1: c.x + COL_W });
		}
		prev = lane;
	}
	return out;
}

/**
 * Splits the strip's height evenly into whole chip rows per lane, at least
 * MIN_CAP each. The rows left over go one each to the lanes with the busiest cells.
 */
export function splitRows(height: number, busiest: number[]): number[] {
	const n = busiest.length;
	const total = Math.floor(
		(height - HEAD_H - STRIP_BORDERS_AND_SCROLLBAR - n * rowsH(0)) / ROW_H,
	);
	const base = Math.max(MIN_CAP, Math.floor(total / n));
	const rows = busiest.map(() => base);
	const byBusiest = busiest
		.map((b, i) => [b, i])
		.sort((a, b) => b[0] - a[0] || a[1] - b[1]);
	for (const [, i] of byBusiest.slice(0, Math.max(0, total - base * n))) {
		rows[i]++;
	}
	return rows;
}

export const portraitSize = (cap: number) =>
	Math.min(MAX_PORTRAIT, rowsH(cap) - PORTRAIT_PAD);

export function chipsShown(
	count: number,
	rows: number,
	expanded: boolean,
): number {
	return !expanded && count > rows ? rows - MORE_BUTTON_ROWS : count;
}

export function visibleWindow(left: number, width: number) {
	const start =
		Math.floor(left / RENDER_STEP_PX) * RENDER_STEP_PX - RENDER_STEP_PX;
	return [start, start + width + 3 * RENDER_STEP_PX] as const;
}
