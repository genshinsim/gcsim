import type { model } from "@gcsim/types";
import { Group } from "@visx/group";
import { BoxPlot } from "@visx/stats";
import type { ScaleLinear } from "d3-scale";
import type { ChartTooltipState } from "../ChartTooltip/useChartTooltip";

export type HoverData<Key> = {
	key: Key;
	index: number;
	x: number;
	y: number;
	height: number;
	width: number;
};

type Props<Datum, Key> = {
	data: Datum[];
	tooltip: ChartTooltipState<HoverData<Key>>;
	scale: ScaleLinear<number, number>;
	color: (k: Key) => string;
	stat: (d: Datum, k: Key) => model.DescriptiveStats;
};

export const HoverBoxPlot = <Datum, Key>({
	data,
	tooltip,
	scale,
	color,
	stat,
}: Props<Datum, Key>) => {
	const hover = tooltip.data;
	if (!tooltip.open || hover === undefined) {
		return null;
	}

	const box = meanSdBox(stat(data[hover.index], hover.key));

	return (
		<Group
			left={hover.x}
			top={hover.y}
			onMouseLeave={tooltip.scheduleHide}
			onMouseMove={(e) => tooltip.show(e, hover)}
		>
			<BoxPlot
				horizontal
				top={10}
				boxWidth={hover.height - 20}
				valueScale={scale}
				min={box.min}
				max={box.max}
				firstQuartile={box.firstQuartile}
				median={box.median}
				thirdQuartile={box.thirdQuartile}
				stroke={"var(--g-text)"}
				fill={color(hover.key)}
			/>
		</Group>
	);
};

const meanSdBox = (s: model.DescriptiveStats) => ({
	min: s.min,
	firstQuartile: (s.mean ?? 0) - (s.sd ?? 0),
	median: s.mean,
	thirdQuartile: (s.mean ?? 0) + (s.sd ?? 0),
	max: s.max,
});
