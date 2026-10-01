import { Group } from "@visx/group";
import { Line } from "@visx/shape";
import type { MutableRefObject } from "react";
import { useTranslation } from "react-i18next";
import {
	type ChartTooltipState,
	DataColorsConst,
	PathDataPoint,
	TooltipList,
	TooltipRow,
} from "../../../../common/gcsim";
import type { Point } from "./DamageOverTimeData";

type HoverLineProps = {
	yMax: number;
	minRef: MutableRefObject<SVGPathElement | null>;
	meanRef: MutableRefObject<SVGPathElement | null>;
	maxRef: MutableRefObject<SVGPathElement | null>;
	tooltip: ChartTooltipState<number>;
	margin: { left: number; right: number; top: number; bottom: number };
};

export const HoverLine = (props: HoverLineProps) => {
	if (
		!props.tooltip.open ||
		!props.meanRef.current ||
		!props.minRef.current ||
		!props.maxRef.current
	) {
		return null;
	}

	const x = props.tooltip.left;
	const pathX = x - props.margin.left;

	return (
		<Group left={-props.margin.left}>
			<Line
				from={{ x: x, y: 0 }}
				to={{ x: x, y: props.yMax }}
				stroke="var(--g-text)"
				opacity={0.5}
				strokeWidth={2}
				pointerEvents="none"
				strokeDasharray="5 2"
			/>
			<PathDataPoint
				cx={x}
				x={pathX}
				fill={DataColorsConst.qualitative2(3)}
				path={props.minRef}
			/>
			<PathDataPoint
				cx={x}
				x={pathX}
				fill={DataColorsConst.qualitative2(1)}
				path={props.maxRef}
			/>
			<PathDataPoint
				cx={x}
				x={pathX}
				fill={DataColorsConst.qualitative3(8)}
				path={props.meanRef}
			/>
		</Group>
	);
};

export const DamageOverTimeTooltipContent = ({ point }: { point: Point }) => {
	const { t } = useTranslation();
	return (
		<div className="flex flex-col font-g-mono text-g-xs">
			<TooltipList>
				<TooltipRow
					color="var(--g-text-mute)"
					name={t("result.time")}
					value={point.x}
					suffix={t("result.seconds_short")}
				/>
				<TooltipRow
					color={DataColorsConst.qualitative4(3)}
					name={t("result.stat_min")}
					value={point.y.min}
				/>
				<TooltipRow
					color={DataColorsConst.qualitative4(1)}
					name={t("result.stat_max")}
					value={point.y.max}
				/>
				<TooltipRow
					color={DataColorsConst.qualitative4(8)}
					name={t("result.stat_mean")}
					value={point.y.mean}
				/>
				<TooltipRow
					color={DataColorsConst.qualitative4(0)}
					name={t("result.stat_std")}
					value={point.y.sd}
				/>
			</TooltipList>
		</div>
	);
};
