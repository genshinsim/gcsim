import type { model } from "@gcsim/types";
import { Group } from "@visx/group";
import { BoxPlot } from "@visx/stats";
import type { ScaleLinear } from "d3-scale";
import type { TooltipData, TooltipHandles } from "./Tooltip";

type Props<Datum, Key> = {
	data: Datum[];
	tooltip?: TooltipData<Key>;
	open: boolean;
	scale: ScaleLinear<number, number>;
	color: (k: Key) => string;
	handles: TooltipHandles<Key>;
	stat: (d: Datum, k: Key) => model.DescriptiveStats;
};

export const HoverBoxPlot = <Datum, Key>({
	data,
	tooltip,
	scale,
	open,
	color,
	handles,
	stat,
}: Props<Datum, Key>) => {
	if (!tooltip || !open) {
		return null;
	}

	const box = meanSdBox(stat(data[tooltip.index], tooltip.key));

	return (
		<Group
			left={tooltip.x}
			top={tooltip.y}
			onMouseLeave={() => handles.mouseLeave()}
			onMouseMove={(e) => handles.mouseHover(e, tooltip)}
		>
			<BoxPlot
				horizontal
				top={10}
				boxWidth={tooltip.height - 20}
				valueScale={scale}
				min={box.min}
				max={box.max}
				firstQuartile={box.firstQuartile}
				median={box.median}
				thirdQuartile={box.thirdQuartile}
				stroke={"#FFF"}
				fill={color(tooltip.key)}
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
