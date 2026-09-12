import type { model } from "@gcsim/types";
import { scaleLinear } from "d3-scale";
import { HISTOGRAM_ACCENT, HISTOGRAM_BAR } from "../colors";
import { ChartImg, NoData } from "./util";

type Props = {
	data?: model.OverviewStats | null;
	width: number;
	height: number;
};

const margin = { left: 8, right: 8, top: 8, bottom: 1 };
const PADDING_INNER = 0.05;

// Damage-distribution histogram. Reproduces the live Graphs/Histogram: linear
// bar heights, with the bar containing the mean accented. Pure math, emitted as
// a standalone SVG string.
export const HistogramChart = ({ data, width, height }: Props) => {
	const hist = data?.histogram;
	if (
		hist == null ||
		hist.length === 0 ||
		data?.min == null ||
		data?.max == null
	) {
		return <NoData width={width} height={height} />;
	}

	const xMax = width - margin.left - margin.right;
	const yMax = height - margin.top - margin.bottom;

	const maxCount = Math.max(...hist);
	const yScale = scaleLinear()
		.domain([0, maxCount + 0.05 * maxCount || 1])
		.range([yMax, 0])
		.clamp(true);

	const step = xMax / hist.length;
	const barWidth = step * (1 - PADDING_INNER);

	// Index of the bar the mean falls into.
	const delta = data.max > data.min ? hist.length / (data.max - data.min) : 0;
	const meanBin =
		data.mean != null ? Math.floor(delta * (data.mean - data.min)) : -1;

	const rects = hist
		.map((c, i) => {
			const barHeight = yMax - yScale(c);
			if (c <= 0 || barHeight < 0) {
				return "";
			}
			const barX = i * step;
			const barY = yMax - barHeight;
			const fill = i === meanBin ? HISTOGRAM_ACCENT : HISTOGRAM_BAR;
			return `<rect x="${barX}" y="${barY}" width="${barWidth}" height="${barHeight}" fill="${fill}"/>`;
		})
		.join("");

	const svg =
		`<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">` +
		`<g transform="translate(${margin.left},${margin.top})">${rects}</g>` +
		`</svg>`;

	return <ChartImg svg={svg} width={width} height={height} />;
};
