import { act, renderHook } from "@testing-library/react";
import type React from "react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { useChartTooltip } from "./useChartTooltip";

const mouseEvent = {} as React.MouseEvent;

describe("useChartTooltip", () => {
	beforeEach(() => {
		vi.useFakeTimers();
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	it("opens with the datum passed to show", () => {
		const { result } = renderHook(() => useChartTooltip<number>());

		act(() => result.current.show(mouseEvent, 7));

		expect(result.current.open).toBe(true);
		expect(result.current.data).toBe(7);
	});

	it("hides after the configured delay when not cancelled", () => {
		const { result } = renderHook(() =>
			useChartTooltip<number>({ hideDelay: 300 }),
		);
		act(() => result.current.show(mouseEvent, 1));

		act(() => result.current.scheduleHide());
		act(() => vi.advanceTimersByTime(299));
		expect(result.current.open).toBe(true);

		act(() => vi.advanceTimersByTime(1));
		expect(result.current.open).toBe(false);
		expect(result.current.data).toBeUndefined();
	});

	it("cancels a scheduled hide even after a re-render in between", () => {
		const { result, rerender } = renderHook(() => useChartTooltip<number>());
		act(() => result.current.show(mouseEvent, 1));

		act(() => result.current.scheduleHide());
		rerender();
		act(() => result.current.cancelHide());
		act(() => vi.advanceTimersByTime(1000));

		expect(result.current.open).toBe(true);
	});

	it("anchors above the cursor by the configured offset", () => {
		const { result } = renderHook(() =>
			useChartTooltip<number>({ offsetY: -20 }),
		);

		act(() => result.current.show(mouseEvent, 1));

		expect(result.current.top).toBe(-20);
	});
});
