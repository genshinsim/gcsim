import type { model } from "@gcsim/types";
import { act, render, renderHook, waitFor } from "@testing-library/react";
import type React from "react";
import { describe, expect, it, vi } from "vitest";
import {
	ExecutorProvider,
	type ExecutorProviderProps,
	useExecutor,
	useRunResult,
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

	it("exposes the supplier and polls ready/busy", async () => {
		const fake = makeExecutor({ running: true });
		const { result } = renderHook(() => useExecutor(), {
			wrapper: ({ children }: { children: React.ReactNode }) => (
				<ExecutorProvider exec={fake.supplier}>{children}</ExecutorProvider>
			),
		});

		expect(result.current.exec).toBe(fake.supplier);
		await waitFor(() => expect(result.current.isReady).toBe(true));
		await waitFor(() => expect(result.current.busy).toBe(true));
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

	it("runs and fires navigation when validation passes", async () => {
		const fake = makeExecutor();
		const navigateOnRun = vi.fn();
		const { result } = renderHook(() => useExecutor(), {
			wrapper: wrapper({ exec: fake.supplier, navigateOnRun }),
		});

		act(() => result.current.run("config"));

		await waitFor(() =>
			expect(fake.run).toHaveBeenCalledWith("config", expect.any(Function)),
		);
		expect(navigateOnRun).toHaveBeenCalledTimes(1);
	});
});

function simResult(n: number): model.SimulationResult {
	return { sim_version: String(n) };
}

function useBoth() {
	return { executor: useExecutor(), run: useRunResult() };
}

async function startRun(
	props: Omit<ExecutorProviderProps, "children" | "exec"> = {},
) {
	const fake = makeExecutor();
	const hook = renderHook(() => useBoth(), {
		wrapper: wrapper({ exec: fake.supplier, ...props }),
	});
	act(() => hook.result.current.executor.run("config"));
	await waitFor(() => expect(fake.run).toHaveBeenCalled());
	const sink = fake.run.mock.calls[0][1];
	return { fake, hook, sink };
}

describe("useRunResult", () => {
	it("throws when used outside a provider", () => {
		expect(() => renderHook(() => useRunResult())).toThrow(
			/must be used within an ExecutorProvider/,
		);
	});

	it("records the config that was run", async () => {
		const { hook } = await startRun();
		expect(hook.result.current.run.config).toBe("config");
	});

	it("updates result and hash from streamed callbacks", async () => {
		const { hook, sink } = await startRun();

		act(() => sink(simResult(1), "hash-1"));
		await waitFor(() =>
			expect(hook.result.current.run.result).toEqual(simResult(1)),
		);
		expect(hook.result.current.run.hash).toBe("hash-1");
	});

	it("coalesces bursts and always applies the final callback", async () => {
		const onResult = vi.fn();
		const { hook, sink } = await startRun({ onResult });

		act(() => {
			for (let i = 1; i <= 20; i++) {
				sink(simResult(i), `hash-${i}`);
			}
		});

		await waitFor(() =>
			expect(hook.result.current.run.result).toEqual(simResult(20)),
		);
		expect(hook.result.current.run.hash).toBe("hash-20");
		expect(onResult.mock.calls.length).toBeLessThan(20);
		expect(onResult).toHaveBeenLastCalledWith(simResult(20), "hash-20");
	});

	it("sets error when the run rejects", async () => {
		const fake = makeExecutor();
		fake.run.mockImplementation(() => Promise.reject("boom"));
		const { result } = renderHook(() => useBoth(), {
			wrapper: wrapper({ exec: fake.supplier }),
		});

		act(() => result.current.executor.run("config"));

		await waitFor(() => expect(result.current.run.error).toBe("boom"));
	});

	it("clears the previous result and error on a new run", async () => {
		const { fake, hook, sink } = await startRun();
		act(() => sink(simResult(1), "hash-1"));
		await waitFor(() => expect(hook.result.current.run.result).not.toBeNull());

		fake.setRunning(false);
		fake.run.mockImplementationOnce(() => Promise.reject("boom"));
		act(() => hook.result.current.executor.run("config 2"));
		await waitFor(() => expect(hook.result.current.run.error).toBe("boom"));
		expect(hook.result.current.run.result).toBeNull();
		expect(hook.result.current.run.hash).toBeNull();

		act(() => hook.result.current.executor.run("config 3"));
		await waitFor(() => expect(hook.result.current.run.error).toBeNull());
		expect(hook.result.current.run.config).toBe("config 3");
	});

	it("does not re-render executor-only consumers when the result changes", async () => {
		const fake = makeExecutor();
		let executorRenders = 0;
		function ExecutorOnly() {
			useExecutor();
			executorRenders++;
			return null;
		}
		const runResult: { current: ReturnType<typeof useRunResult> | null } = {
			current: null,
		};
		function ResultReader() {
			runResult.current = useRunResult();
			return null;
		}
		const executor: { current: ReturnType<typeof useExecutor> | null } = {
			current: null,
		};
		function Runner() {
			executor.current = useExecutor();
			return null;
		}
		const tree = (
			<ExecutorProvider exec={fake.supplier}>
				<ExecutorOnly />
				<ResultReader />
				<Runner />
			</ExecutorProvider>
		);
		render(tree);

		act(() => executor.current?.run("config"));
		await waitFor(() => expect(fake.run).toHaveBeenCalled());
		await waitFor(() => expect(executor.current?.busy).toBe(true));
		const sink = fake.run.mock.calls[0][1];
		const before = executorRenders;

		for (let i = 1; i <= 3; i++) {
			act(() => sink(simResult(i), `hash-${i}`));
			await act(() => new Promise((r) => setTimeout(r, 150)));
		}

		expect(runResult.current?.result).toEqual(simResult(3));
		expect(executorRenders).toBe(before);
	});
});
