import type { model } from "@gcsim/types";
import { scaleLinear } from "d3-scale";
import { area, curveBasis } from "d3-shape";
import { TIMELINE_COLOR } from "../colors";
import { ChartImg, NoData } from "./util";

type Props = {
	data?: model.BucketStats | null;
	width: number;
	height: number;
};

const margin = { left: 8, right: 8, top: 0, bottom: 1 };

// Damage-over-time area chart. Pure d3-shape/d3-scale math emitted as a
// standalone SVG string; no visx, no ParentSize, no DOM measurement.
export const TimelineChart = ({ data, width, height }: Props) => {
	if (data?.bucket_size == null || data.buckets == null) {
		return <NoData width={width} height={height} />;
	}

	const bucketSize = data.bucket_size;
	const points = data.buckets.map((b, i) => ({
		x: (i * bucketSize) / 60,
		y: b.mean ?? 0,
	}));

	if (points.length === 0) {
		return <NoData width={width} height={height} />;
	}

	const duration = ((points.length - 1) * bucketSize) / 60;
	const maxValue = points.reduce((acc, p) => Math.max(acc, p.y), 0);

	const xMax = width - margin.left - margin.right;
	const yMax = height - margin.top - margin.bottom;

	const xScale = scaleLinear()
		.domain([0, duration || 1])
		.range([0, xMax]);
	const yScale = scaleLinear()
		.domain([0, maxValue || 1])
		.range([yMax, 0]);

	const gen = area<{ x: number; y: number }>()
		.x((d) => xScale(d.x))
		.y0(yMax)
		.y1((d) => yScale(d.y))
		.curve(curveBasis);

	const path = gen(points) ?? "";

	const svg =
		`<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">` +
		`<g transform="translate(${margin.left},${margin.top})">` +
		`<path d="${path}" fill="${TIMELINE_COLOR}" stroke="${TIMELINE_COLOR}" stroke-width="2"/>` +
		`</g></svg>`;

	return <ChartImg svg={svg} width={width} height={height} />;
};
