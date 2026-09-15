import { localPoint } from "@visx/event";
import { Group } from "@visx/group";
import { Line } from "@visx/shape";
import { TooltipWithBounds } from "@visx/tooltip";
import type { ScaleLinear } from "d3-scale";
import { useTranslation } from "react-i18next";
import {
	Colors,
	DataColorsConst,
	FloatStatTooltipContent,
	useDataColors,
} from "../../../../common/gcsim";
import {
	Popover,
	PopoverAnchor,
	PopoverContent,
} from "../../../../common/ui/popover";
import type { CumulativePoint } from "./CumulativeData";

export interface TooltipData {
	index: number;
}

export interface TooltipHandles {
	mouseLeave: () => void;
	mouseHover: (e: React.MouseEvent) => void;
	clearTimeout: () => void;
}

type ShowTooltipArgs<Datum> = {
	tooltipData?: Datum;
	tooltipLeft?: number;
	tooltipTop?: number;
};

export function useTooltipHandles(
	showTooltip: (args: ShowTooltipArgs<TooltipData>) => void,
	hideTooltip: () => void,
	xScale: ScaleLinear<number, number>,
	yMax: number,
	margin: { left: number; right: number; top: number; bottom: number },
	bucketSize: number,
): TooltipHandles {
	let tooltipTimeout: number;
	const mouseLeave = () => {
		tooltipTimeout = window.setTimeout(() => {
			hideTooltip();
		}, 150);
	};

	const clearTimeout = () => {
		if (tooltipTimeout) {
			window.clearTimeout(tooltipTimeout);
		}
	};

	const mouseHover = (e: React.MouseEvent) => {
		const { x } = localPoint(e) || { x: 0 };
		const index = Math.round(
			(60 * xScale.invert(x - margin.left)) / bucketSize,
		);

		clearTimeout();
		showTooltip({
			tooltipData: { index: index },
			tooltipLeft: x,
			tooltipTop: yMax,
		});
	};

	return {
		mouseLeave: mouseLeave,
		mouseHover: mouseHover,
		clearTimeout: clearTimeout,
	};
}

type HoverLineProps = {
	data: CumulativePoint[];
	names: string[];
	yScale: ScaleLinear<number, number>;
	yMax: number;
	tooltipData?: TooltipData;
	tooltipOpen?: boolean;
	tooltipLeft?: number;
	margin: { left: number; right: number; top: number; bottom: number };
};

export const HoverLine = (props: HoverLineProps) => {
	if (!props.tooltipOpen || !props.tooltipLeft || !props.tooltipData) {
		return null;
	}

	const x = props.tooltipLeft;
	const point = props.data[props.tooltipData.index];

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
					fill="#000"
					fillOpacity={0.1}
					stroke="#000"
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
					stroke="#FFF"
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
				stroke="#FFF"
				opacity={0.5}
				strokeWidth={2}
				pointerEvents="none"
				strokeDasharray="5 2"
			/>
			{circles}
		</Group>
	);
};

type TooltipProps = {
	data: CumulativePoint[];
	names?: string[];
	tooltipOpen: boolean;
	tooltipData?: TooltipData;
	tooltipTop?: number;
	tooltipLeft?: number;
	handles: TooltipHandles;
	showTooltip: (args: ShowTooltipArgs<TooltipData>) => void;
	margin: { left: number; right: number; top: number; bottom: number };
};

export const RenderTooltip = (props: TooltipProps) => {
	const { DataColors } = useDataColors();
	const { i18n, t } = useTranslation();

	if (
		!props.tooltipOpen ||
		!props.tooltipData ||
		!props.tooltipLeft ||
		!props.names
	) {
		return null;
	}

	const point = props.data[props.tooltipData.index];

	const content = (
		// biome-ignore lint/a11y/noStaticElementInteractions: mouse-only chart tooltip hover region, no interactive semantics
		<div
			onMouseMove={() => {
				props.handles.clearTimeout();
				props.showTooltip({
					tooltipData: props.tooltipData,
					tooltipLeft: props.tooltipLeft,
					tooltipTop: props.tooltipTop,
				});
			}}
			onMouseLeave={() => props.handles.mouseLeave()}
		>
			<div className="flex flex-row px-2 py-1 font-mono text-xs gap-2 whitespace-nowrap">
				<span style={{ color: Colors.SEPIA4 }}>{t("result.time")}: </span>
				<span>{point.x + t("result.seconds_short")}</span>
			</div>
			{point.y
				.slice(0)
				.reverse()
				.map((val, char) => {
					const i = (props.names?.length ?? 0) - char - 1;
					return (
						<FloatStatTooltipContent
							key={"tooltip-" + i}
							title={props.names?.[i] + " " + t("result.contribution")}
							data={val}
							color={DataColors.characterLabel(i)}
							format={(s) =>
								s?.toLocaleString(i18n.language, {
									style: "percent",
									maximumFractionDigits: 2,
								})
							}
						/>
					);
				})}
		</div>
	);

	return (
		<TooltipWithBounds
			style={{ position: "absolute" }}
			offsetLeft={props.margin.left + 50}
			left={props.tooltipLeft}
			top={props.tooltipTop}
		>
			<Popover open>
				<PopoverAnchor />
				<PopoverContent
					side="top"
					onOpenAutoFocus={(e) => e.preventDefault()}
					className="w-auto max-w-none p-0"
				>
					{content}
				</PopoverContent>
			</Popover>
		</TooltipWithBounds>
	);
};
