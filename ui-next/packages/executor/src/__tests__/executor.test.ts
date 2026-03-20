import { describe, expect, it } from "vitest";
import type { Executor, ExecutorSupplier } from "../executor.js";

describe("Executor interface", () => {
  it("should be assignable from an object implementing all methods", () => {
    // Type-level test: verify the interface compiles and is assignable
    const executor: Executor = {
      ready: () => Promise.resolve(true),
      running: () => false,
      validate: () =>
        Promise.resolve({ characters: [], errors: [], player_initial_pos: { x: 0, y: 0, r: 0 } }),
      sample: () => Promise.resolve({}),
      run: () => Promise.resolve(true),
      cancel: () => {},
      buildInfo: () => ({ hash: "", date: "" }),
    };

    expect(executor).toBeDefined();
    expect(typeof executor.ready).toBe("function");
    expect(typeof executor.running).toBe("function");
    expect(typeof executor.validate).toBe("function");
    expect(typeof executor.sample).toBe("function");
    expect(typeof executor.run).toBe("function");
    expect(typeof executor.cancel).toBe("function");
    expect(typeof executor.buildInfo).toBe("function");
  });

  it("should support ExecutorSupplier type", () => {
    const supplier: ExecutorSupplier<Executor> = () => ({
      ready: () => Promise.resolve(true),
      running: () => false,
      validate: () =>
        Promise.resolve({ characters: [], errors: [], player_initial_pos: { x: 0, y: 0, r: 0 } }),
      sample: () => Promise.resolve({}),
      run: () => Promise.resolve(true),
      cancel: () => {},
      buildInfo: () => ({ hash: "", date: "" }),
    });

    const executor = supplier();
    expect(executor).toBeDefined();
    expect(executor.running()).toBe(false);
  });
});
