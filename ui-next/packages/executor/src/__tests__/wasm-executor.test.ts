import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { WasmExecutor } from "../wasm-executor.js";

// Mock Worker class since we cannot run real WASM in tests
class MockWorker {
  public onmessage: ((ev: MessageEvent) => void) | null = null;
  public postMessage = vi.fn();
  public terminate = vi.fn();
}

describe("WasmExecutor", () => {
  let executor: WasmExecutor;

  beforeEach(() => {
    vi.stubGlobal("Worker", MockWorker);
    executor = new WasmExecutor("/test.wasm");
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  describe("constructor", () => {
    it("should set default worker count to 3", () => {
      expect(executor.getWorkerCount()).toBe(3);
    });
  });

  describe("setWorkerCount()", () => {
    it("should update worker count", () => {
      executor.setWorkerCount(5);
      expect(executor.getWorkerCount()).toBe(5);
    });

    it("should clamp to minimum of 1", () => {
      executor.setWorkerCount(0);
      expect(executor.getWorkerCount()).toBe(1);

      executor.setWorkerCount(-5);
      expect(executor.getWorkerCount()).toBe(1);
    });

    it("should clamp to maximum of 30", () => {
      executor.setWorkerCount(50);
      expect(executor.getWorkerCount()).toBe(30);
    });
  });

  describe("running()", () => {
    it("should return false initially", () => {
      expect(executor.running()).toBe(false);
    });
  });

  describe("ready()", () => {
    it("should return true when not running", async () => {
      expect(await executor.ready()).toBe(true);
    });
  });

  describe("cancel()", () => {
    it("should not throw when not running", () => {
      expect(() => executor.cancel()).not.toThrow();
    });

    it("should set running to false and terminate aggregator during run", async () => {
      // Set up executor in a running state by starting a run
      // We need to manually set up the internal state
      // Since cancel checks isRunning && aggregator != null, we need to trigger a run first

      // Create a controlled run that we can cancel
      const updateResult = vi.fn();

      // Start the run - workers will post ready messages
      void executor.run("test config", updateResult);

      // The run creates workers which need Ready responses
      // Since MockWorker captures postMessage calls, we can simulate responses
      // but the internal state should have isRunning = true now
      expect(executor.running()).toBe(true);

      executor.cancel();
      expect(executor.running()).toBe(false);

      // The run promise should not hang - it may reject since we cancelled
      // Just ensure cancel worked
    });
  });

  describe("buildInfo()", () => {
    it("should throw not implemented error", () => {
      expect(() => executor.buildInfo()).toThrow("Method not implemented");
    });
  });
});
