import { Group } from "@visx/group";
import { Line } from "@visx/shape";
import type { ScaleLinear } from "d3-scale";
import { useTranslation } from "react-i18next";
import {
	type ChartTooltipState,
	DataColorsConst,
	FloatStatTooltipContent,
	TooltipList,
	TooltipRow,
	useDataColors,
} from "../../../../common/gcsim";
import type { CumulativePoint } from "./CumulativeData";

type HoverLineProps = {
	data: CumulativePoint[];
	names: string[];
	yScale: ScaleLinear<number, number>;
	yMax: number;
	tooltip: ChartTooltipState<number>;
	margin: { left: number; right: number; top: number; bottom: number };
};

export const HoverLine = (props: HoverLineProps) => {
	const point =
		props.tooltip.data === undefined
			? undefined
			: props.data[props.tooltip.data];
	if (!props.tooltip.open || point === undefined) {
		return null;
	}

	const x = props.tooltip.left;

	let total = 0;
	const circles = point.y.map((val, char) => {
		total += val.mean ?? 0;
		const y = props.yScale(total);
		return (
			<g key={props.names[char]}>
				<circle
					cx={x}
					cy={y + 1}
					r={4}
					fill="var(--g-bg)"
					fillOpacity={0.1}
					stroke="var(--g-bg)"
					strokeOpacity={0.1}
					strokeWidth={2}
					pointerEvents="none"
				/>
				<circle
					cx={x}
					cy={y}
					r={4}
					fill={DataColorsConst.qualitative4(char)}
					pointerEvents="none"
					stroke="var(--g-surface)"
					strokeWidth={2}
				/>
			</g>
		);
	});

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
			{circles}
		</Group>
	);
};

type ContentProps = {
	point: CumulativePoint;
	names: string[];
};

export const CumulativeTooltipContent = ({ point, names }: ContentProps) => {
	const { DataColors } = useDataColors();
	const { i18n, t } = useTranslation();

	return (
		<div className="flex flex-col font-g-mono text-g-xs">
			<TooltipList>
				<TooltipRow
					color="var(--g-text-mute)"
					name={t("result.time")}
					value={point.x}
					suffix={t("result.seconds_short")}
				/>
			</TooltipList>
			{point.y
				.slice(0)
				.reverse()
				.map((val, char) => {
					const i = names.length - char - 1;
					return (
						<FloatStatTooltipContent
							key={names[i]}
							title={`${names[i]} ${t("result.contribution")}`}
							data={val}
							color={DataColors.characterLabel(i)}
							format={(s) =>
								s?.toLocaleString(i18n.language, {
									style: "percent",
									minimumFractionDigits: 2,
									maximumFractionDigits: 2,
								})
							}
						/>
					);
				})}
		</div>
	);
};
