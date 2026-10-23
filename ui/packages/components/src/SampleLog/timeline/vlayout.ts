import { type Column, chipsShown } from "./layout";
import { type Chip, lowerBound } from "./model";

export const CHIP_H = 18;
const CHIP_GAP = 2;
export const V_CELL_PAD_Y = 4;
export const MIN_ROW_H = 36;
const RULE = 1;
export const V_GAP_BIG = 16;
export const V_GUTTER_W = 78;
export const V_HEAD_H = 40;
export const V_MIN_COL_W = 128;
export const V_CELL_PAD_X = 8;
const WRAP_LINE_H = 14;
const WRAP_PAD_Y = 2;
// chip padding, the type icon (wider on duration chips) and its gap, plus rounding slack
const CHIP_CHROME_X = 8 + 14 + 4 + 3;

/** a frame with events, laid out as a row; `y` is its top, gap band included */
export type VRow = {
	i: number;
	col: Column;
	y: number;
	h: number;
	/** height of the "+Nf" band at the top of the row */
	gap: number;
};

export type ChipHeight = (c: Chip) => number;

export const slimChip: ChipHeight = () => CHIP_H;

/** lane columns share the width beyond the gutter, but never get narrower than V_MIN_COL_W */
export const vColumnWidth = (scrollerWidth: number, lanes: number) =>
	Math.max(V_MIN_COL_W, (scrollerWidth - V_GUTTER_W) / lanes);

/** chips a cell shows before "+N more"; `cap` null shows them all */
export const vShown = (count: number, cap: number | null, expanded: boolean) =>
	cap == null ? count : chipsShown(count, cap, expanded);

const stackH = (chips: Chip[], chipH: ChipHeight, more: boolean) => {
	const n = chips.length + (more ? 1 : 0);
	const h = chips.reduce((a, c) => a + chipH(c), 0) + (more ? CHIP_H : 0);
	return n === 0 ? 0 : h + (n - 1) * CHIP_GAP;
};

export function vRows(
	cols: Column[],
	lanes: number[],
	{
		cap,
		expanded,
		chipH,
	}: { cap: number | null; expanded: Set<number>; chipH: ChipHeight },
): VRow[] {
	let y = 0;
	return cols.map((col, i) => {
		const gap = i > 0 && col.framesSincePrev > 1 ? V_GAP_BIG : 0;
		let content = 0;
		for (const l of lanes) {
			const items = col.cells[l];
			const fit = vShown(items.length, cap, expanded.has(i));
			content = Math.max(
				content,
				stackH(items.slice(0, fit), chipH, fit < items.length),
			);
		}
		const h = gap + Math.max(MIN_ROW_H, 2 * V_CELL_PAD_Y + content) + RULE;
		const row = { i, col, y, h, gap };
		y += h;
		return row;
	});
}

export const rowsHeight = (rows: VRow[]) => {
	const last = rows[rows.length - 1];
	return last == null ? 0 : last.y + last.h;
};

const topOf = (r: VRow) => r.y + r.gap;

/** y of a frame, interpolated between the rows around it */
export function yOfFrame(rows: VRow[], frame: number): number {
	const i = lowerBound(rows, frame, (r) => r.col.frame);
	if (i === 0 || i === rows.length) {
		const r = rows[Math.min(i, rows.length - 1)];
		return r == null ? 0 : topOf(r);
	}
	const [a, b] = [rows[i - 1], rows[i]];
	return (
		topOf(a) +
		((frame - a.col.frame) / (b.col.frame - a.col.frame)) *
			(topOf(b) - topOf(a))
	);
}

/** the inverse of {@link yOfFrame} */
export function frameOfY(rows: VRow[], y: number): number {
	const i = lowerBound(rows, y, topOf);
	if (i === 0 || i === rows.length) {
		return rows[Math.min(i, rows.length - 1)]?.col.frame ?? 0;
	}
	const [a, b] = [rows[i - 1], rows[i]];
	return (
		a.col.frame +
		((y - topOf(a)) / (topOf(b) - topOf(a))) * (b.col.frame - a.col.frame)
	);
}

/**
 * Lines the browser wraps `text` onto in `avail` px: greedy, breaking at spaces
 * and after hyphens, and mid-word when a piece can't fit a line by itself.
 */
export function wrapLines(
	text: string,
	avail: number,
	width: (s: string) => number,
): number {
	if (avail <= 0) {
		return 1;
	}
	const space = width(" ");
	let lines = 1;
	let cur = 0;
	for (const word of text.split(/\s+/)) {
		for (const [k, piece] of word.split(/(?<=-)/).entries()) {
			const sep = k === 0 && cur > 0 ? space : 0;
			const w = width(piece);
			if (cur > 0 && cur + sep + w > avail) {
				lines++;
				cur = 0;
			}
			if (w > avail) {
				lines += Math.ceil(w / avail) - 1;
				cur = w % avail;
			} else {
				cur += (cur > 0 ? sep : 0) + w;
			}
		}
	}
	return lines;
}

/** heights of wrapping chips laid out `chipW` px wide */
export function wrappedChip(
	chipW: number,
	text: (c: Chip) => string,
	width: (s: string) => number,
): ChipHeight {
	const avail = chipW - CHIP_CHROME_X;
	const memo = new Map<number, number>();
	return (c) => {
		let h = memo.get(c.id);
		if (h == null) {
			h = wrapLines(text(c), avail, width) * WRAP_LINE_H + 2 * WRAP_PAD_Y;
			memo.set(c.id, h);
		}
		return h;
	};
}

let measure: CanvasRenderingContext2D | null = null;
const measured = new Map<string, number>();

/** text width in the chip font, measured on a canvas */
export function chipTextWidth(font: string): (s: string) => number {
	return (s) => {
		const key = `${font}|${s}`;
		let w = measured.get(key);
		if (w == null) {
			measure ??= document.createElement("canvas").getContext("2d");
			if (measure == null) {
				return s.length * 6;
			}
			measure.font = font;
			w = measure.measureText(s).width;
			measured.set(key, w);
		}
		return w;
	};
}
