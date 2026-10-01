import { Tooltip, TooltipContent, TooltipTrigger } from "@gcsim/primitives";
import type { ReactNode } from "react";
import type { ChartTooltipState } from "./useChartTooltip";

type Props<TData> = {
	tooltip: ChartTooltipState<TData>;
	children: (data: TData) => ReactNode;
};

export const ChartTooltip = <TData,>({ tooltip, children }: Props<TData>) => {
	if (!tooltip.open || tooltip.data === undefined) {
		return null;
	}

	const content = children(tooltip.data);
	if (content == null) {
		return null;
	}

	return (
		<div
			className="pointer-events-none absolute"
			style={{ left: tooltip.left, top: tooltip.top }}
		>
			<Tooltip open>
				<TooltipTrigger asChild>
					<div />
				</TooltipTrigger>
				<TooltipContent
					side="top"
					onPointerDownOutside={(e) => e.preventDefault()}
					onMouseEnter={tooltip.cancelHide}
					onMouseLeave={tooltip.scheduleHide}
				>
					{content}
				</TooltipContent>
			</Tooltip>
		</div>
	);
};
