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
const BUCKET_PX = 2;

type Props = {
	model: TimelineModel;
	enabled: Set<string>;
	/** frames in view on the strip */
	range: [number, number];
	onPan: (start: number) => void;
	matchFrames: number[];
	/** a highlighted status duration */
	selection: [number, number] | null;
	/** sits beside the overview, or under it on narrow screens */
	aside: React.ReactNode;
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
	width: number;
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
		ctx.fillRect(0, y, s.width, LANE_H);
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

function drawAxis(s: Scale, top: number, p: Palette) {
	const { ctx } = s;
	const secStep = s.width / (s.frames / 60) > 22 ? 5 : 10;
	ctx.font = "10px ui-monospace, monospace";
	ctx.textBaseline = "top";
	for (let sec = 0; sec * 60 <= s.frames; sec += secStep) {
		const x = s.x(sec * 60);
		ctx.fillStyle = p.line;
		ctx.fillRect(x, top, 1, 4);
		ctx.fillStyle = p.mute;
		if (sec > 0) {
			ctx.fillText(`${sec}s`, x + 2, top + 2);
		}
	}
}

/** a zoomed-out overview of the whole sample with a draggable window onto the strip */
export function Minimap({
	model,
	enabled,
	range,
	onPan,
	matchFrames,
	selection,
	aside,
}: Props) {
	const wrap = useRef<HTMLDivElement>(null);
	const canvas = useRef<HTMLCanvasElement>(null);
	const [width, setWidth] = useState(0);
	const drag = useRef<{ anchor: number } | null>(null);
	const lanesH = model.lanes.length * (LANE_H + LANE_GAP);
	const height = LANES_TOP + lanesH + AXIS_H;
	const frames = model.maxFrame + 1;

	useEffect(() => {
		const el = wrap.current;
		if (el == null) {
			return;
		}
		const ro = new ResizeObserver(() => setWidth(el.clientWidth));
		ro.observe(el);
		setWidth(el.clientWidth);
		return () => ro.disconnect();
	}, []);

	useEffect(() => {
		const c = canvas.current;
		if (c == null || width === 0) {
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
		const palette: Palette = {
			lanes: model.lanes.map((l) => cssVar(c, laneColorVar(l), "#888")),
			ink: cssVar(c, "--g-text", "#eee"),
			mute: cssVar(c, "--g-text-mute", "#888"),
			line: cssVar(c, "--g-border", "#333"),
			warn: cssVar(c, "--g-warning", "#dda857"),
		};
		const buckets = Math.ceil(width / BUCKET_PX);
		const scale: Scale = {
			ctx,
			width,
			frames,
			buckets,
			bucketOf: (f) =>
				Math.min(buckets - 1, Math.floor(((f / frames) * width) / BUCKET_PX)),
			x: (f) => (f / frames) * width,
		};
		drawDamage(scale, model, palette);
		drawLanes(scale, model, enabled, palette);
		drawBursts(scale, model, palette);
		drawMatches(scale, matchFrames, palette);
		drawAxis(scale, LANES_TOP + lanesH, palette);
	}, [model, enabled, width, height, lanesH, frames, matchFrames]);

	const fx = (f: number) => (f / frames) * width;
	const xa = fx(Math.max(0, range[0]));
	const xb = Math.max(xa + 6, fx(Math.min(frames, range[1])));
	const viewW = range[1] - range[0];

	const frameAt = (clientX: number) => {
		const rect = wrap.current?.getBoundingClientRect();
		if (rect == null || width === 0) {
			return 0;
		}
		return Math.max(
			0,
			Math.min(frames, ((clientX - rect.left) / width) * frames),
		);
	};

	const onPointerDown = (e: React.PointerEvent) => {
		const rect = wrap.current?.getBoundingClientRect();
		if (rect == null) {
			return;
		}
		const x = e.clientX - rect.left;
		const f = frameAt(e.clientX);
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
		onPan(frameAt(e.clientX) - d.anchor);
	};

	const endDrag = () => {
		drag.current = null;
	};

	const wheelPan = useEffectEvent((e: WheelEvent) => {
		const delta = Math.abs(e.deltaX) > Math.abs(e.deltaY) ? e.deltaX : e.deltaY;
		onPan(range[0] + (delta / Math.max(1, width)) * frames * 0.5);
	});
	useEffect(() => {
		const el = wrap.current;
		if (el == null) {
			return;
		}
		const onWheel = (e: WheelEvent) => {
			e.preventDefault();
			wheelPan(e);
		};
		el.addEventListener("wheel", onWheel, { passive: false });
		return () => el.removeEventListener("wheel", onWheel);
	}, []);

	return (
		<div className="flex shrink-0 flex-col rounded-g-md border border-g-line bg-g-surface sm:flex-row">
			<div className="order-last shrink-0 border-t border-g-line py-1.5 sm:w-[232px] sm:border-t-0 sm:border-l sm:py-0">
				{aside}
			</div>
			<div
				ref={wrap}
				className="relative min-w-0 cursor-grab sm:flex-1 touch-none select-none active:cursor-grabbing"
				style={{ height }}
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
					className="pointer-events-none absolute inset-y-0 left-0 bg-g-surface/60"
					style={{ width: xa }}
				/>
				<div
					className="pointer-events-none absolute inset-y-0 right-0 bg-g-surface/60"
					style={{ left: xb }}
				/>
				<div
					className="pointer-events-none absolute inset-y-0 rounded-[3px] border-2 border-g-accent bg-g-accent/10"
					style={{ left: xa, width: xb - xa }}
				/>
				{selection != null && (
					<div
						className="pointer-events-none absolute inset-y-0 border-x-2 border-g-warning bg-g-warning/20"
						style={{
							left: fx(Math.min(selection[0], frames)),
							width: Math.max(
								2,
								fx(Math.min(selection[1], frames)) -
									fx(Math.min(selection[0], frames)),
							),
						}}
					/>
				)}
			</div>
		</div>
	);
}
