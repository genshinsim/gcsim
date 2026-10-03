import type React from "react";
import { useEffect, useRef, useState } from "react";
import { isEvent } from "../events/types";
import type { TimelineModel } from "./model";

const DMG_H = 22;
const LANE_H = 7;
const LANE_GAP = 2;
const AXIS_H = 13;

type Props = {
	model: TimelineModel;
	enabled: Set<string>;
	/** frames in view on the strip */
	range: [number, number];
	onPan: (start: number) => void;
	matchFrames: number[];
	/** highlighted frame range (a status duration), drawn as a band */
	selection: [number, number] | null;
	/** sits in a column on the right */
	aside: React.ReactNode;
	asideWidth: number;
};

function cssVar(el: Element, name: string, fallback: string) {
	const v = getComputedStyle(el).getPropertyValue(name).trim();
	return v === "" ? fallback : v;
}

export function Minimap({
	model,
	enabled,
	range,
	onPan,
	matchFrames,
	selection,
	aside,
	asideWidth,
}: Props) {
	const wrap = useRef<HTMLDivElement>(null);
	const canvas = useRef<HTMLCanvasElement>(null);
	const [width, setWidth] = useState(0);
	const drag = useRef<{ anchor: number } | null>(null);
	const lanesH = model.lanes.length * (LANE_H + LANE_GAP);
	const height = DMG_H + 4 + lanesH + AXIS_H;
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
		const colors = model.lanes.map((l) =>
			l.element == null
				? cssVar(c, "--g-text-mute", "#888")
				: cssVar(c, `--g-${l.element}`, "#888"),
		);
		const ink = cssVar(c, "--g-text", "#eee");
		const mute = cssVar(c, "--g-text-mute", "#888");
		const line = cssVar(c, "--g-border", "#333");
		const warn = cssVar(c, "--g-warning", "#dda857");
		const bucketPx = 2;
		const buckets = Math.ceil(width / bucketPx);
		const bucketOf = (f: number) =>
			Math.min(buckets - 1, Math.floor(((f / frames) * width) / bucketPx));

		// stacked damage per bucket, by lane
		const dmg = model.lanes.map(() => new Float64Array(buckets));
		const total = new Float64Array(buckets);
		for (const e of model.chips) {
			if (e.dmg > 0) {
				const b = bucketOf(e.frame);
				dmg[e.lane][b] += e.dmg;
				total[b] += e.dmg;
			}
		}
		const maxTotal = Math.max(1, ...total);
		for (let b = 0; b < buckets; b++) {
			let y = DMG_H;
			for (let l = 0; l < model.lanes.length; l++) {
				const h = (dmg[l][b] / maxTotal) * (DMG_H - 2);
				if (h > 0) {
					ctx.fillStyle = colors[l];
					ctx.fillRect(b * bucketPx, y - h, bucketPx - 0.5, h);
					y -= h;
				}
			}
		}

		// lanes: on-field band, event density, bursts
		const top = DMG_H + 4;
		const density = model.lanes.map(() => new Float64Array(buckets));
		for (const e of model.chips) {
			if (enabled.has(e.event.type)) {
				density[e.lane][bucketOf(e.frame)]++;
			}
		}
		model.lanes.forEach((_, l) => {
			const y = top + l * (LANE_H + LANE_GAP);
			ctx.fillStyle = line;
			ctx.globalAlpha = 0.35;
			ctx.fillRect(0, y, width, LANE_H);
			ctx.globalAlpha = 0.45;
			ctx.fillStyle = colors[l];
			for (const s of model.onField) {
				if (s.lane === l) {
					const x0 = (s.start / frames) * width;
					ctx.fillRect(x0, y, (s.end / frames) * width - x0, LANE_H);
				}
			}
			ctx.fillStyle = ink;
			for (let b = 0; b < buckets; b++) {
				const n = density[l][b];
				if (n > 0) {
					ctx.globalAlpha = Math.min(0.9, 0.12 + n / 14);
					ctx.fillRect(b * bucketPx, y + 2, bucketPx - 0.5, LANE_H - 4);
				}
			}
			ctx.globalAlpha = 1;
		});
		ctx.fillStyle = ink;
		for (const e of model.chips) {
			if (
				!e.expired &&
				isEvent(e.event, "action") &&
				e.event.action === "burst"
			) {
				const y = top + e.lane * (LANE_H + LANE_GAP);
				const x = (e.frame / frames) * width;
				ctx.beginPath();
				ctx.moveTo(x, y - 1);
				ctx.lineTo(x + 3.5, y + LANE_H / 2);
				ctx.lineTo(x, y + LANE_H + 1);
				ctx.lineTo(x - 3.5, y + LANE_H / 2);
				ctx.closePath();
				ctx.fillStyle = colors[e.lane];
				ctx.fill();
				ctx.strokeStyle = ink;
				ctx.lineWidth = 0.75;
				ctx.stroke();
			}
		}

		// search matches
		ctx.fillStyle = warn;
		for (const f of matchFrames) {
			ctx.fillRect((f / frames) * width - 0.5, 0, 1.5, 6);
		}

		// axis
		const axisY = top + lanesH;
		const secStep = width / (frames / 60) > 22 ? 5 : 10;
		ctx.font = "10px ui-monospace, monospace";
		ctx.textBaseline = "top";
		for (let s = 0; s * 60 <= frames; s += secStep) {
			const x = ((s * 60) / frames) * width;
			ctx.fillStyle = line;
			ctx.fillRect(x, axisY, 1, 4);
			ctx.fillStyle = mute;
			if (s > 0) {
				ctx.fillText(`${s}s`, x + 2, axisY + 2);
			}
		}
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

	// wheel over the minimap pans the timeline horizontally
	useEffect(() => {
		const el = wrap.current;
		if (el == null) {
			return;
		}
		const onWheel = (e: WheelEvent) => {
			e.preventDefault();
			const delta =
				Math.abs(e.deltaX) > Math.abs(e.deltaY) ? e.deltaX : e.deltaY;
			onPan(range[0] + (delta / Math.max(1, width)) * frames * 0.5);
		};
		el.addEventListener("wheel", onWheel, { passive: false });
		return () => el.removeEventListener("wheel", onWheel);
	}, [onPan, range, width, frames]);

	return (
		<div className="flex shrink-0 gap-0 rounded-g-md border border-g-line bg-g-surface">
			<div
				className="order-last shrink-0 border-l border-g-line"
				style={{ width: asideWidth, height }}
			>
				{aside}
			</div>
			<div
				ref={wrap}
				className="relative min-w-0 flex-1 cursor-grab touch-none select-none active:cursor-grabbing"
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
