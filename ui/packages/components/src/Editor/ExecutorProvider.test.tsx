import { act, render, renderHook, waitFor } from "@testing-library/react";
import type React from "react";
import { describe, expect, it, vi } from "vitest";
import {
	ExecutorProvider,
	type ExecutorProviderProps,
	useExecutor,
} from "./ExecutorProvider";
import { makeExecutor, parsedResult } from "./testExecutor";

function wrapper(props: Omit<ExecutorProviderProps, "children">) {
	return ({ children }: { children: React.ReactNode }) => (
		<ExecutorProvider {...props}>{children}</ExecutorProvider>
	);
}

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

describe("ExecutorProvider run()", () => {
	it("freshly validates at invocation instead of trusting stale state", async () => {
		const fake = makeExecutor({
			validate: () => Promise.resolve(parsedResult(["amber"], ["stale error"])),
		});
		const { result } = renderHook(() => useExecutor(), {
			wrapper: wrapper({ exec: fake.supplier }),
		});

		fake.validate.mockImplementation(() =>
			Promise.resolve(parsedResult(["amber"])),
		);
		act(() => result.current.run("config"));

		await waitFor(() =>
			expect(fake.run).toHaveBeenCalledWith("config", expect.any(Function)),
		);
	});

	it("refuses to run while the shared pool is already running", async () => {
		const fake = makeExecutor({ running: true });
		const navigateOnRun = vi.fn();
		const { result } = renderHook(() => useExecutor(), {
			wrapper: wrapper({ exec: fake.supplier, navigateOnRun }),
		});

		act(() => result.current.run("config"));

		await waitFor(() => expect(fake.validate).toHaveBeenCalledWith("config"));
		expect(fake.run).not.toHaveBeenCalled();
		expect(navigateOnRun).not.toHaveBeenCalled();
	});

	it("refuses to run when fresh validation reports errors", async () => {
		const fake = makeExecutor({
			validate: () => Promise.resolve(parsedResult(["amber"], ["bad action"])),
		});
		const navigateOnRun = vi.fn();
		const { result } = renderHook(() => useExecutor(), {
			wrapper: wrapper({ exec: fake.supplier, navigateOnRun }),
		});

		act(() => result.current.run("config"));

		await waitFor(() => expect(fake.validate).toHaveBeenCalledWith("config"));
		expect(fake.run).not.toHaveBeenCalled();
		expect(navigateOnRun).not.toHaveBeenCalled();
	});

	it("runs with the result sink and fires navigation when validation passes", async () => {
		const fake = makeExecutor();
		const onResult = vi.fn();
		const navigateOnRun = vi.fn();
		const { result } = renderHook(() => useExecutor(), {
			wrapper: wrapper({ exec: fake.supplier, onResult, navigateOnRun }),
		});

		act(() => result.current.run("config"));

		await waitFor(() =>
			expect(fake.run).toHaveBeenCalledWith("config", onResult),
		);
		expect(navigateOnRun).toHaveBeenCalledTimes(1);
	});
});
