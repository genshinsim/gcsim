import {
	type Chip,
	onFieldLaneAt,
	SIM_LANE,
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
const LABEL_PAD_X = 12;
const STRIP_BORDERS_AND_SCROLLBAR = 14;
const MORE_BUTTON_ROWS = 1;
const RENDER_STEP_PX = 400;

/** width of the sticky lane-label column for a portrait of this size */
export const gutterWidth = (portrait: number) => portrait + LABEL_PAD_X;

export const secondsLabel = (frame: number, digits = 2) =>
	`${(frame / 60).toFixed(digits)}s`;

export const rowsH = (n: number) =>
	2 * CELL_PAD + n * (CHIP_H + CHIP_GAP) - CHIP_GAP;

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

export function chipCap(height: number, lanes: number): number {
	const perLane = (height - HEAD_H - STRIP_BORDERS_AND_SCROLLBAR) / lanes;
	return Math.max(
		MIN_CAP,
		Math.floor((perLane - rowsH(0)) / (CHIP_H + CHIP_GAP)),
	);
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

export function laneHeights(rows: number[], portrait: number): number[] {
	return rows.map((r, i) =>
		Math.max(
			i === SIM_LANE ? SIM_LANE_MIN_H : portrait + PORTRAIT_PAD,
			rowsH(r),
		),
	);
}

export function visibleWindow(left: number, width: number) {
	const step = RENDER_STEP_PX;
	const start = Math.floor(left / step) * step - step;
	return [start, start + width + 3 * step] as const;
}
