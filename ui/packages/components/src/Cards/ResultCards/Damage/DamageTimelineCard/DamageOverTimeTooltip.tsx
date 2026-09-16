import { Popover, PopoverAnchor, PopoverContent } from "@gcsim/primitives";
import { localPoint } from "@visx/event";
import { Group } from "@visx/group";
import { Line } from "@visx/shape";
import { TooltipWithBounds } from "@visx/tooltip";
import type { ScaleLinear } from "d3-scale";
import type { MutableRefObject } from "react";
import { useTranslation } from "react-i18next";
import { Colors, DataColorsConst } from "../../../../common/gcsim";
import type { Point } from "./DamageOverTimeData";

export interface TooltipData {
	index: number;
}

export interface TooltipHandles {
	mouseLeave: () => void;
	mouseHover: (e: React.MouseEvent) => void;
	clearTimeout: () => void;
}

export function useTooltipHandles(
	showTooltip: (args: ShowTooltipArgs<TooltipData>) => void,
	hideTooltip: () => void,
	xScale: ScaleLinear<number, number>,
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
			tooltipTop: e.nativeEvent.offsetY - 50,
		});
	};

	return {
		mouseLeave: mouseLeave,
		mouseHover: mouseHover,
		clearTimeout: clearTimeout,
	};
}

type ShowTooltipArgs<Datum> = {
	tooltipData?: Datum;
	tooltipLeft?: number;
	tooltipTop?: number;
};

type HoverLineProps = {
	yMax: number;
	minRef: MutableRefObject<SVGPathElement | null>;
	meanRef: MutableRefObject<SVGPathElement | null>;
	maxRef: MutableRefObject<SVGPathElement | null>;
	tooltipData?: TooltipData;
	tooltipOpen?: boolean;
	tooltipLeft?: number;
	margin: { left: number; right: number; top: number; bottom: number };
};

export const HoverLine = (props: HoverLineProps) => {
	if (
		!props.tooltipOpen ||
		!props.tooltipLeft ||
		!props.tooltipData ||
		!props.meanRef.current ||
		!props.minRef.current ||
		!props.maxRef.current
	) {
		return null;
	}

	const x = props.tooltipLeft;

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
		</Group>
	);
};

type TooltipProps = {
	data: Point[];
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
	const { t } = useTranslation();
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
			<div className="flex flex-col px-2 py-1 font-mono text-xs">
				<ul className="grid grid-cols-[repeat(2,_max-content)] gap-x-2 justify-start">
					<Item
						color={Colors.SEPIA4}
						name={t("result.time")}
						value={point.x}
						suffix={t("result.seconds_short")}
					/>
					<Item
						color={DataColorsConst.qualitative4(3)}
						name="min"
						value={point.y.min}
					/>
					<Item
						color={DataColorsConst.qualitative4(1)}
						name="max"
						value={point.y.max}
					/>
					<Item
						color={DataColorsConst.qualitative4(8)}
						name="mean"
						value={point.y.mean}
					/>
					<Item
						color={DataColorsConst.qualitative4(0)}
						name="std"
						value={point.y.sd}
					/>
				</ul>
			</div>
		</div>
	);

	return (
		<TooltipWithBounds
			style={{ position: "absolute" }}
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

type ItemProps = {
	name: string;
	value?: number;
	color?: string;
	suffix?: string;
};

const Item = ({ name, value, color, suffix }: ItemProps) => {
	const { i18n } = useTranslation();
	const num = value?.toLocaleString(i18n.language, {
		minimumFractionDigits: 2,
		maximumFractionDigits: 2,
	});

	return (
		<>
			<span className="text-gray-400 list-item" style={{ color: color }}>
				{name}
			</span>
			<span>
				{num}
				{suffix}
			</span>
		</>
	);
};
