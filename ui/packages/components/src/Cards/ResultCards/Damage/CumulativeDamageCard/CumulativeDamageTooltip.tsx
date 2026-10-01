import { Group } from "@visx/group";
import { Line } from "@visx/shape";
import type { ScaleLinear } from "d3-scale";
import type { MutableRefObject } from "react";
import { useTranslation } from "react-i18next";
import {
	type ChartTooltipState,
	DataColorsConst,
	PathDataPoint,
	TooltipList,
	TooltipRow,
} from "../../../../common/gcsim";
import type { Point } from "./CumulativeDamageData";

type HoverLineProps = {
	data: Point[];
	xScale: ScaleLinear<number, number>;
	yScale: ScaleLinear<number, number>;
	yMax: number;
	minRef: MutableRefObject<SVGPathElement | null>;
	maxRef: MutableRefObject<SVGPathElement | null>;
	q1Ref: MutableRefObject<SVGPathElement | null>;
	q2Ref: MutableRefObject<SVGPathElement | null>;
	q3Ref: MutableRefObject<SVGPathElement | null>;
	tooltip: ChartTooltipState<number>;
	margin: { left: number; right: number; top: number; bottom: number };
};

export const HoverLine = (props: HoverLineProps) => {
	const point =
		props.tooltip.data === undefined
			? undefined
			: props.data[props.tooltip.data];
	if (
		!props.tooltip.open ||
		point === undefined ||
		!props.minRef.current ||
		!props.maxRef.current ||
		!props.q1Ref.current ||
		!props.q2Ref.current ||
		!props.q3Ref.current
	) {
		return null;
	}

	const x = props.tooltip.left;

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
				x={props.xScale(point.x)}
				fill={DataColorsConst.qualitative2(3)}
				path={props.minRef}
			/>
			<PathDataPoint
				cx={x}
				x={props.xScale(point.x)}
				fill={DataColorsConst.qualitative2(1)}
				path={props.maxRef}
			/>
			<PathDataPoint
				cx={x}
				x={props.xScale(point.x)}
				fill={DataColorsConst.qualitative2(4)}
				path={props.q1Ref}
			/>
			<PathDataPoint
				cx={x}
				x={props.xScale(point.x)}
				fill={DataColorsConst.qualitative3(8)}
				path={props.q2Ref}
			/>
			<PathDataPoint
				cx={x}
				x={props.xScale(point.x)}
				fill={DataColorsConst.qualitative2(5)}
				path={props.q3Ref}
			/>
		</Group>
	);
};

export const CumulativeDamageTooltipContent = ({ point }: { point: Point }) => {
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
					color={DataColorsConst.qualitative2(3)}
					name={t("result.stat_min")}
					value={point.y.min}
				/>
				<TooltipRow
					color={DataColorsConst.qualitative2(1)}
					name={t("result.stat_max")}
					value={point.y.max}
				/>
				<TooltipRow
					color={DataColorsConst.qualitative2(4)}
					name={t("result.stat_p25")}
					value={point.y.q1}
				/>
				<TooltipRow
					color={DataColorsConst.qualitative3(8)}
					name={t("result.stat_p50")}
					value={point.y.q2}
				/>
				<TooltipRow
					color={DataColorsConst.qualitative2(5)}
					name={t("result.stat_p75")}
					value={point.y.q3}
				/>
			</TooltipList>
		</div>
	);
};
