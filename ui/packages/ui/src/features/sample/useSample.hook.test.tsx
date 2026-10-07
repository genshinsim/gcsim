// @vitest-environment jsdom
import { renderHook, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { useSample } from "./useSample";

const sampleFn = vi.fn();

vi.mock("@gcsim/components", () => ({
	useExecutor: () => ({ exec: () => ({ sample: sampleFn }) }),
}));

type Props = { config?: string; autoSeed: string | null };

describe("useSample auto sampling", () => {
	beforeEach(() => {
		sampleFn.mockReset();
		sampleFn.mockImplementation((_cfg: string, seed: string) =>
			Promise.resolve({ seed }),
		);
	});

	it("samples as soon as the first partial result arrives", async () => {
		const { result, rerender } = renderHook((p: Props) => useSample(p), {
			initialProps: { config: undefined, autoSeed: null } as Props,
		});
		rerender({ config: "cfg", autoSeed: "1" });
		await waitFor(() => expect(result.current.sample).toEqual({ seed: "1" }));
		expect(sampleFn).toHaveBeenCalledTimes(1);
	});

	it("samples a rerun even if the previous run was never seen ending", async () => {
		const { result, rerender } = renderHook((p: Props) => useSample(p), {
			initialProps: { config: "cfg", autoSeed: "1" } as Props,
		});
		await waitFor(() => expect(result.current.sample).toEqual({ seed: "1" }));
		rerender({ config: "cfg", autoSeed: "2" });
		await waitFor(() => expect(result.current.sample).toEqual({ seed: "2" }));
		expect(sampleFn).toHaveBeenCalledTimes(2);
	});

	it("drops a sample that resolves after a new run started", async () => {
		let resolveOld: (v: unknown) => void = () => {};
		sampleFn.mockImplementationOnce(
			() => new Promise((resolve) => (resolveOld = resolve)),
		);
		const { result, rerender } = renderHook((p: Props) => useSample(p), {
			initialProps: { config: "cfg", autoSeed: "1" } as Props,
		});
		rerender({ config: undefined, autoSeed: null });
		rerender({ config: "cfg", autoSeed: "2" });
		await waitFor(() => expect(result.current.sample).toEqual({ seed: "2" }));
		resolveOld({ seed: "1" });
		await Promise.resolve();
		expect(result.current.sample).toEqual({ seed: "2" });
		expect(result.current.generating).toBe(false);
	});
});
