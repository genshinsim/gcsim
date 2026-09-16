import { act, renderHook, waitFor } from "@testing-library/react";
import type React from "react";
import { describe, expect, it, vi } from "vitest";
import { ExecutorProvider } from "./ExecutorProvider";
import { makeExecutor, parsedResult } from "./testExecutor";
import { type UseValidationOptions, useValidation } from "./useValidation";

function wrapper(supplier: ReturnType<typeof makeExecutor>["supplier"]) {
	return ({ children }: { children: React.ReactNode }) => (
		<ExecutorProvider exec={supplier}>{children}</ExecutorProvider>
	);
}

describe("useValidation", () => {
	it("debounces then validates the initial config once ready", async () => {
		const fake = makeExecutor();
		const { result } = renderHook(() => useValidation("my config"), {
			wrapper: wrapper(fake.supplier),
		});

		await waitFor(() =>
			expect(fake.validate).toHaveBeenCalledWith("my config"),
		);
		await waitFor(() => expect(result.current.parsedTeam).toHaveLength(1));
		expect(result.current.isValid).toBe(true);
		expect(result.current.error).toBeNull();
	});

	it("does not validate an empty config", async () => {
		const fake = makeExecutor();
		const { result } = renderHook(() => useValidation(""), {
			wrapper: wrapper(fake.supplier),
		});

		await act(() => new Promise((r) => setTimeout(r, 500)));
		expect(fake.validate).not.toHaveBeenCalled();
		expect(result.current.parsedTeam).toEqual([]);
	});

	it("does not validate while the executor is not ready", async () => {
		const fake = makeExecutor({ ready: false });
		renderHook(() => useValidation("config"), {
			wrapper: wrapper(fake.supplier),
		});

		await act(() => new Promise((r) => setTimeout(r, 500)));
		expect(fake.validate).not.toHaveBeenCalled();
	});

	it("surfaces validation errors as an error string", async () => {
		const fake = makeExecutor({
			validate: () => Promise.resolve(parsedResult(["amber"], ["bad action"])),
		});
		const { result } = renderHook(() => useValidation("config"), {
			wrapper: wrapper(fake.supplier),
		});

		await waitFor(() => expect(result.current.error).toBe("bad action"));
		expect(result.current.isValid).toBe(false);
	});

	it("produces independent results for two configs against the shared pool", async () => {
		const fake = makeExecutor({
			validate: (cfg) =>
				Promise.resolve(
					parsedResult(cfg === "a" ? ["amber"] : ["bennett", "xingqiu"]),
				),
		});
		const w = wrapper(fake.supplier);
		const first = renderHook(() => useValidation("a"), { wrapper: w });
		const second = renderHook(() => useValidation("b"), { wrapper: w });

		await waitFor(() =>
			expect(first.result.current.parsedTeam).toHaveLength(1),
		);
		await waitFor(() =>
			expect(second.result.current.parsedTeam).toHaveLength(2),
		);
		expect(first.result.current.parsedTeam[0].name).toBe("amber");
		expect(second.result.current.parsedTeam[0].name).toBe("bennett");
	});
});

describe("useValidation run()", () => {
	it("freshly validates at invocation instead of trusting stale isValid", async () => {
		const fake = makeExecutor({
			validate: () => Promise.resolve(parsedResult(["amber"], ["stale error"])),
		});
		const { result } = renderHook(() => useValidation("config"), {
			wrapper: wrapper(fake.supplier),
		});
		await waitFor(() => expect(result.current.isValid).toBe(false));

		fake.validate.mockImplementation(() =>
			Promise.resolve(parsedResult(["amber"])),
		);
		result.current.run();

		await waitFor(() =>
			expect(fake.run).toHaveBeenCalledWith("config", expect.any(Function)),
		);
	});

	it("refuses to run while the shared pool is already running", async () => {
		const fake = makeExecutor({ running: true });
		const onRun = vi.fn();
		const options: UseValidationOptions = { onRun };
		const { result } = renderHook(() => useValidation("config", options), {
			wrapper: wrapper(fake.supplier),
		});
		await waitFor(() => expect(fake.validate).toHaveBeenCalled());
		fake.validate.mockClear();

		result.current.run();

		await waitFor(() => expect(fake.validate).toHaveBeenCalled());
		expect(fake.run).not.toHaveBeenCalled();
		expect(onRun).not.toHaveBeenCalled();
	});

	it("runs and fires the navigation callback when fresh validation passes", async () => {
		const fake = makeExecutor();
		const onResult = vi.fn();
		const onRun = vi.fn();
		const { result } = renderHook(
			() => useValidation("config", { onResult, onRun }),
			{ wrapper: wrapper(fake.supplier) },
		);
		await waitFor(() => expect(result.current.isValid).toBe(true));

		result.current.run();

		await waitFor(() =>
			expect(fake.run).toHaveBeenCalledWith("config", onResult),
		);
		expect(onRun).toHaveBeenCalledTimes(1);
	});
});
