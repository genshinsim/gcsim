import { cn } from "@gcsim/primitives";
import type React from "react";
import { useEffect, useEffectEvent, useRef, useState } from "react";
import { isEvent } from "../events/types";
import { laneColorVar } from "./lane";
import { type Chip, damageOf, type TimelineModel } from "./model";

const DMG_H = 22;
const LANES_TOP = DMG_H + 4;
const LANE_H = 7;
const LANE_GAP = 2;
const AXIS_H = 13;
const V_AXIS_W = 26;
const BUCKET_PX = 2;

type Props = {
	model: TimelineModel;
	enabled: Set<string>;
	/** frames in view on the strip */
	range: [number, number];
	onPan: (start: number) => void;
	onWheel: (delta: number, mode: number) => void;
	matchFrames: number[];
	/** a highlighted status duration */
	selection: [number, number] | null;
	/** sits beside the overview, or under it on narrow screens; above it when vertical */
	aside: React.ReactNode;
	/** time runs down instead of across */
	vertical?: boolean;
};

type Palette = {
	lanes: string[];
	ink: string;
	mute: string;
	line: string;
	warn: string;
};

type Scale = {
	ctx: CanvasRenderingContext2D;
	/** px along the time axis */
	len: number;
	frames: number;
	buckets: number;
	bucketOf: (frame: number) => number;
	x: (frame: number) => number;
};

function cssVar(el: Element, name: string, fallback: string) {
	const v = getComputedStyle(el).getPropertyValue(name).trim();
	return v === "" ? fallback : v;
}

const laneY = (lane: number) => LANES_TOP + lane * (LANE_H + LANE_GAP);

function drawDamage(s: Scale, model: TimelineModel, p: Palette) {
	const dmg = model.lanes.map(() => new Float64Array(s.buckets));
	const total = new Float64Array(s.buckets);
	for (const c of model.chips) {
		const d = damageOf(c);
		if (d > 0) {
			const b = s.bucketOf(c.frame);
			dmg[c.lane][b] += d;
			total[b] += d;
		}
	}
	const maxTotal = Math.max(1, ...total);
	for (let b = 0; b < s.buckets; b++) {
		let y = DMG_H;
		for (let l = 0; l < model.lanes.length; l++) {
			const h = (dmg[l][b] / maxTotal) * (DMG_H - 2);
			if (h > 0) {
				s.ctx.fillStyle = p.lanes[l];
				s.ctx.fillRect(b * BUCKET_PX, y - h, BUCKET_PX - 0.5, h);
				y -= h;
			}
		}
	}
}

/** per lane: a band where it's on field and a shade for event density */
function drawLanes(
	s: Scale,
	model: TimelineModel,
	enabled: Set<string>,
	p: Palette,
) {
	const { ctx } = s;
	const density = model.lanes.map(() => new Float64Array(s.buckets));
	for (const c of model.chips) {
		if (enabled.has(c.event.type)) {
			density[c.lane][s.bucketOf(c.frame)]++;
		}
	}
	model.lanes.forEach((_, l) => {
		const y = laneY(l);
		ctx.fillStyle = p.line;
		ctx.globalAlpha = 0.35;
		ctx.fillRect(0, y, s.len, LANE_H);
		ctx.globalAlpha = 0.45;
		ctx.fillStyle = p.lanes[l];
		for (const stint of model.onField) {
			if (stint.lane === l) {
				const x0 = s.x(stint.start);
				ctx.fillRect(x0, y, s.x(stint.end) - x0, LANE_H);
			}
		}
		ctx.fillStyle = p.ink;
		for (let b = 0; b < s.buckets; b++) {
			const n = density[l][b];
			if (n > 0) {
				ctx.globalAlpha = Math.min(0.9, 0.12 + n / 14);
				ctx.fillRect(b * BUCKET_PX, y + 2, BUCKET_PX - 0.5, LANE_H - 4);
			}
		}
		ctx.globalAlpha = 1;
	});
}

const isBurst = (c: Chip) =>
	!c.expired && isEvent(c.event, "action") && c.event.action === "burst";

function drawBursts(s: Scale, model: TimelineModel, p: Palette) {
	const { ctx } = s;
	for (const c of model.chips) {
		if (!isBurst(c)) {
			continue;
		}
		const y = laneY(c.lane);
		const x = s.x(c.frame);
		ctx.beginPath();
		ctx.moveTo(x, y - 1);
		ctx.lineTo(x + 3.5, y + LANE_H / 2);
		ctx.lineTo(x, y + LANE_H + 1);
		ctx.lineTo(x - 3.5, y + LANE_H / 2);
		ctx.closePath();
		ctx.fillStyle = p.lanes[c.lane];
		ctx.fill();
		ctx.strokeStyle = p.ink;
		ctx.lineWidth = 0.75;
		ctx.stroke();
	}
}

function drawMatches(s: Scale, matchFrames: number[], p: Palette) {
	s.ctx.fillStyle = p.warn;
	for (const f of matchFrames) {
		s.ctx.fillRect(s.x(f) - 0.5, 0, 1.5, 6);
	}
}

/** `at` is the axis's distance from the overview's top, or left edge when vertical */
function drawAxis(s: Scale, at: number, vertical: boolean, p: Palette) {
	const { ctx } = s;
	const secStep = s.len / (s.frames / 60) > 22 ? 5 : 10;
	ctx.font = vertical
		? "9px ui-monospace, monospace"
		: "10px ui-monospace, monospace";
	ctx.textBaseline = vertical ? "middle" : "top";
	for (let sec = 0; sec * 60 <= s.frames; sec += secStep) {
		const t = s.x(sec * 60);
		ctx.fillStyle = p.line;
		if (vertical) {
			ctx.fillRect(at, t, 4, 1);
		} else {
			ctx.fillRect(t, at, 1, 4);
		}
		ctx.fillStyle = p.mute;
		if (sec > 0) {
			if (vertical) {
				ctx.fillText(`${sec}s`, at + 5, t);
			} else {
				ctx.fillText(`${sec}s`, t + 2, at + 2);
			}
		}
	}
}

/** a zoomed-out overview of the whole sample with a draggable window onto the strip */
export function Minimap({
	model,
	enabled,
	range,
	onPan,
	onWheel,
	matchFrames,
	selection,
	aside,
	vertical = false,
}: Props) {
	const wrap = useRef<HTMLDivElement>(null);
	const canvas = useRef<HTMLCanvasElement>(null);
	/** px along the time axis */
	const [len, setLen] = useState(0);
	const drag = useRef<{ anchor: number } | null>(null);
	const lanesH = model.lanes.length * (LANE_H + LANE_GAP);
	/** px across the time axis */
	const depth = LANES_TOP + lanesH + (vertical ? V_AXIS_W : AXIS_H);
	const [width, height] = vertical ? [depth, len] : [len, depth];
	const frames = model.maxFrame + 1;

	useEffect(() => {
		const el = wrap.current;
		if (el == null) {
			return;
		}
		const measure = () => setLen(vertical ? el.clientHeight : el.clientWidth);
		const ro = new ResizeObserver(measure);
		ro.observe(el);
		measure();
		return () => ro.disconnect();
	}, [vertical]);

	useEffect(() => {
		const c = canvas.current;
		if (c == null || len === 0) {
			return;
		}
		const dpr = window.devicePixelRatio || 1;
		c.width = width * dpr;
		c.height = height * dpr;
		const ctx = c.getContext("2d");
		if (ctx == null) {
			return;
		}
		ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
		ctx.clearRect(0, 0, width, height);
		if (vertical) {
			// draw as if horizontal, mirrored across the diagonal so time runs down
			ctx.setTransform(0, dpr, dpr, 0, 0, 0);
		}
		const palette: Palette = {
			lanes: model.lanes.map((l) => cssVar(c, laneColorVar(l), "#888")),
			ink: cssVar(c, "--g-text", "#eee"),
			mute: cssVar(c, "--g-text-mute", "#888"),
			line: cssVar(c, "--g-border", "#333"),
			warn: cssVar(c, "--g-warning", "#dda857"),
		};
		const buckets = Math.ceil(len / BUCKET_PX);
		const scale: Scale = {
			ctx,
			len,
			frames,
			buckets,
			bucketOf: (f) =>
				Math.min(buckets - 1, Math.floor(((f / frames) * len) / BUCKET_PX)),
			x: (f) => (f / frames) * len,
		};
		drawDamage(scale, model, palette);
		drawLanes(scale, model, enabled, palette);
		drawBursts(scale, model, palette);
		drawMatches(scale, matchFrames, palette);
		ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
		drawAxis(scale, LANES_TOP + lanesH, vertical, palette);
	}, [
		model,
		enabled,
		len,
		width,
		height,
		lanesH,
		frames,
		matchFrames,
		vertical,
	]);

	const fx = (f: number) => (f / frames) * len;
	const xa = fx(Math.max(0, range[0]));
	const xb = Math.max(xa + 6, fx(Math.min(frames, range[1])));
	const viewW = range[1] - range[0];
	const along = (e: { clientX: number; clientY: number }, rect: DOMRect) =>
		vertical ? e.clientY - rect.top : e.clientX - rect.left;
	const span = (start: number, size: number) =>
		vertical ? { top: start, height: size } : { left: start, width: size };

	const frameAt = (e: React.PointerEvent) => {
		const rect = wrap.current?.getBoundingClientRect();
		if (rect == null || len === 0) {
			return 0;
		}
		return Math.max(0, Math.min(frames, (along(e, rect) / len) * frames));
	};

	const onPointerDown = (e: React.PointerEvent) => {
		const rect = wrap.current?.getBoundingClientRect();
		if (rect == null) {
			return;
		}
		const x = along(e, rect);
		const f = frameAt(e);
		if (x >= xa && x <= xb) {
			drag.current = { anchor: f - range[0] };
		} else {
			onPan(f - viewW / 2);
			drag.current = { anchor: viewW / 2 };
		}
		e.currentTarget.setPointerCapture(e.pointerId);
	};

	const onPointerMove = (e: React.PointerEvent) => {
		const d = drag.current;
		if (d == null) {
			return;
		}
		onPan(frameAt(e) - d.anchor);
	};

	const endDrag = () => {
		drag.current = null;
	};

	const forwardWheel = useEffectEvent((e: WheelEvent) => {
		// Firefox reports lines only when deltaMode is read before the deltas
		const mode = e.deltaMode;
		const delta = Math.abs(e.deltaX) > Math.abs(e.deltaY) ? e.deltaX : e.deltaY;
		onWheel(delta, mode);
	});
	useEffect(() => {
		const el = wrap.current;
		if (el == null) {
			return;
		}
		const listener = (e: WheelEvent) => {
			e.preventDefault();
			forwardWheel(e);
		};
		el.addEventListener("wheel", listener, { passive: false });
		return () => el.removeEventListener("wheel", listener);
	}, []);

	return (
		<div
			className={cn(
				"flex shrink-0 flex-col rounded-g-md border border-g-line bg-g-surface",
				!vertical && "sm:flex-row",
			)}
			style={vertical ? { width: width + 2 } : undefined}
		>
			<div
				className={cn(
					"shrink-0 border-g-line",
					vertical
						? "flex h-8 items-center justify-center border-b"
						: "order-last border-t py-1.5 sm:w-[232px] sm:border-t-0 sm:border-l sm:py-0",
				)}
			>
				{aside}
			</div>
			<div
				ref={wrap}
				className={cn(
					"relative cursor-grab touch-none select-none active:cursor-grabbing",
					vertical ? "min-h-0 flex-1" : "min-w-0 sm:flex-1",
				)}
				style={vertical ? undefined : { height }}
				onPointerDown={onPointerDown}
				onPointerMove={onPointerMove}
				onPointerUp={endDrag}
				onPointerCancel={endDrag}
			>
				<canvas
					ref={canvas}
					className="absolute inset-0"
					style={{ width, height }}
				/>
				<div
					className={cn(
						"pointer-events-none absolute bg-g-surface/60",
						vertical ? "inset-x-0 top-0" : "inset-y-0 left-0",
					)}
					style={vertical ? { height: xa } : { width: xa }}
				/>
				<div
					className={cn(
						"pointer-events-none absolute bg-g-surface/60",
						vertical ? "inset-x-0 bottom-0" : "inset-y-0 right-0",
					)}
					style={vertical ? { top: xb } : { left: xb }}
				/>
				<div
					className={cn(
						"pointer-events-none absolute rounded-[3px] border-2 border-g-accent bg-g-accent/10",
						vertical ? "inset-x-0" : "inset-y-0",
					)}
					style={span(xa, xb - xa)}
				/>
				{selection != null && (
					<div
						className={cn(
							"pointer-events-none absolute border-g-warning bg-g-warning/20",
							vertical ? "inset-x-0 border-y-2" : "inset-y-0 border-x-2",
						)}
						style={span(
							fx(Math.min(selection[0], frames)),
							Math.max(
								2,
								fx(Math.min(selection[1], frames)) -
									fx(Math.min(selection[0], frames)),
							),
						)}
					/>
				)}
			</div>
		</div>
	);
}
