import type { model } from "@gcsim/types";
import { arc, pie } from "d3-shape";
import { characterColor, elementColor } from "../colors";
import { ChartImg, NoData, wrapSvg } from "./util";

type Slice = {
	value: number;
	color: string;
};

const PIE_RADIUS = 0.8;
const OUTLINE = "#FFFFFF";
const OUTLINE_WIDTH = 0.5;

// Shared pie renderer: emits a standalone SVG string of colored arcs.
function renderPie(slices: Slice[], width: number, height: number): string {
	const radius = (Math.min(width, height) / 2) * PIE_RADIUS;

	const arcs = pie<Slice>()
		.value((d) => d.value)
		.sort(null)(slices);

	const arcGen = arc<(typeof arcs)[number]>()
		.innerRadius(0)
		.outerRadius(radius);

	const paths = arcs
		.map(
			(a) =>
				`<path d="${arcGen(a) ?? ""}" fill="${a.data.color}" stroke="${OUTLINE}" stroke-width="${OUTLINE_WIDTH}"/>`,
		)
		.join("");

	return wrapSvg(width, height, paths, `translate(${width / 2},${height / 2})`);
}

type CharacterPieProps = {
	dps?: model.DescriptiveStats[] | null;
	width: number;
	height: number;
};

// Per-character DPS pie, colored by character index (matches DataColors.character).
export const CharacterPie = ({ dps, width, height }: CharacterPieProps) => {
	if (dps == null || dps.length === 0) {
		return <NoData width={width} height={height} />;
	}

	const slices: Slice[] = dps.map((d, i) => ({
		value: d.mean ?? 0,
		color: characterColor(i),
	}));

	return (
		<ChartImg
			svg={renderPie(slices, width, height)}
			width={width}
			height={height}
		/>
	);
};

type ElementPieProps = {
	dps?: { [k: string]: model.DescriptiveStats } | null;
	width: number;
	height: number;
};

// Per-element DPS pie, colored by element name (matches DataColors.element).
export const ElementPie = ({ dps, width, height }: ElementPieProps) => {
	if (dps == null || Object.keys(dps).length === 0) {
		return <NoData width={width} height={height} />;
	}

	const slices: Slice[] = Object.entries(dps).map(([key, value]) => ({
		value: value.mean ?? 0,
		color: elementColor(key),
	}));

	return (
		<ChartImg
			svg={renderPie(slices, width, height)}
			width={width}
			height={height}
		/>
	);
};
