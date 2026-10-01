import { localPoint } from "@visx/event";
import type React from "react";
import { useCallback, useEffect, useRef, useState } from "react";

export type ChartTooltipOptions = {
	hideDelay?: number;
	offsetY?: number;
};

export type ChartTooltipState<TData> = {
	open: boolean;
	data?: TData;
	left: number;
	top: number;
	show: (e: React.MouseEvent, data: TData) => void;
	scheduleHide: () => void;
	cancelHide: () => void;
};

type VisibleState<TData> = Pick<
	ChartTooltipState<TData>,
	"open" | "data" | "left" | "top"
>;

export function useChartTooltip<TData>({
	hideDelay = 150,
	offsetY = -50,
}: ChartTooltipOptions = {}): ChartTooltipState<TData> {
	const [visible, setVisible] = useState<VisibleState<TData>>({
		open: false,
		left: 0,
		top: 0,
	});
	const hideTimer = useRef<number | undefined>(undefined);

	const cancelHide = useCallback(() => {
		window.clearTimeout(hideTimer.current);
		hideTimer.current = undefined;
	}, []);

	const scheduleHide = useCallback(() => {
		cancelHide();
		hideTimer.current = window.setTimeout(() => {
			hideTimer.current = undefined;
			setVisible((p) => ({ ...p, open: false, data: undefined }));
		}, hideDelay);
	}, [cancelHide, hideDelay]);

	const show = useCallback(
		(e: React.MouseEvent, data: TData) => {
			cancelHide();
			const { x, y } = localPoint(e) ?? { x: 0, y: 0 };
			setVisible({ open: true, data: data, left: x, top: y + offsetY });
		},
		[cancelHide, offsetY],
	);

	useEffect(() => cancelHide, [cancelHide]);

	return { ...visible, show, scheduleHide, cancelHide };
}
