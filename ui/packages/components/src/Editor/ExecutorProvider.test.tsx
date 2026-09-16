import { act, render, renderHook, waitFor } from "@testing-library/react";
import type React from "react";
import { describe, expect, it, vi } from "vitest";
import { ExecutorProvider, useExecutor } from "./ExecutorProvider";
import { makeExecutor } from "./testExecutor";

describe("ExecutorProvider", () => {
	it("throws when useExecutor is used outside a provider", () => {
		expect(() => renderHook(() => useExecutor())).toThrow(
			/must be used within an ExecutorProvider/,
		);
	});

	it("exposes the supplier and polls ready/running", async () => {
		const fake = makeExecutor({ running: true });
		const { result } = renderHook(() => useExecutor(), {
			wrapper: ({ children }: { children: React.ReactNode }) => (
				<ExecutorProvider exec={fake.supplier}>{children}</ExecutorProvider>
			),
		});

		expect(result.current.exec).toBe(fake.supplier);
		await waitFor(() => expect(result.current.isReady).toBe(true));
		await waitFor(() => expect(result.current.running).toBe(true));
	});

	it("stops polling after unmount", async () => {
		const fake = makeExecutor();
		const runningSpy = vi.spyOn(fake.executor, "running");
		const { unmount } = render(
			<ExecutorProvider exec={fake.supplier}>
				<div />
			</ExecutorProvider>,
		);
		await waitFor(() => expect(runningSpy).toHaveBeenCalled());

		unmount();
		runningSpy.mockClear();
		await act(() => new Promise((r) => setTimeout(r, 200)));
		expect(runningSpy).not.toHaveBeenCalled();
	});
});
