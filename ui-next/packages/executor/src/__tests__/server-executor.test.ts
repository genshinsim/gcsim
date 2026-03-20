import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ExecutorError, ServerExecutor } from "../server-executor.js";

describe("ServerExecutor", () => {
  let executor: ServerExecutor;
  const mockFetch = vi.fn();

  beforeEach(() => {
    executor = new ServerExecutor("http://localhost:8080");
    vi.stubGlobal("fetch", mockFetch);
    mockFetch.mockReset();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  describe("ready()", () => {
    it("should call correct URL and return true on 200", async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, status: 200 });

      const result = await executor.ready();

      expect(result).toBe(true);
      expect(mockFetch).toHaveBeenCalledTimes(1);
      const url = mockFetch.mock.calls[0][0] as string;
      expect(url).toMatch(/^http:\/\/localhost:8080\/ready\/id\d+$/);
    });

    it("should cache result after first call", async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, status: 200 });

      await executor.ready();
      const result = await executor.ready();

      expect(result).toBe(true);
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it("should return false on network error", async () => {
      mockFetch.mockRejectedValueOnce(new Error("Network error"));

      const result = await executor.ready();

      expect(result).toBe(false);
    });

    it("should return false on non-200 status", async () => {
      mockFetch.mockResolvedValueOnce({ ok: false, status: 500 });

      const result = await executor.ready();

      expect(result).toBe(false);
    });

    it("should reset cache when setUrl is called", async () => {
      mockFetch.mockResolvedValueOnce({ ok: true, status: 200 });
      await executor.ready();

      executor.setUrl("http://localhost:9090");

      mockFetch.mockResolvedValueOnce({ ok: true, status: 200 });
      await executor.ready();

      expect(mockFetch).toHaveBeenCalledTimes(2);
      const url = mockFetch.mock.calls[1][0] as string;
      expect(url).toMatch(/^http:\/\/localhost:9090\/ready\/id\d+$/);
    });
  });

  describe("running()", () => {
    it("should return false initially", () => {
      expect(executor.running()).toBe(false);
    });
  });

  describe("validate()", () => {
    it("should POST to correct URL with config body", async () => {
      const mockData = {
        characters: [{ base: { key: "test" } }],
        error_msgs: [],
        initial_player_pos: { x: 0, y: 0, r: 0 },
      };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockData),
      });

      const result = await executor.validate("test config");

      expect(mockFetch).toHaveBeenCalledTimes(1);
      const [url, options] = mockFetch.mock.calls[0];
      expect(url).toMatch(/^http:\/\/localhost:8080\/validate\/id\d+$/);
      expect(options.method).toBe("POST");
      expect(JSON.parse(options.body)).toEqual({ config: "test config" });

      expect(result).toEqual({
        characters: mockData.characters,
        errors: mockData.error_msgs,
        player_initial_pos: mockData.initial_player_pos,
      });
    });

    it("should throw ExecutorError on network failure", async () => {
      mockFetch.mockRejectedValueOnce(new Error("Network error"));

      await expect(executor.validate("test")).rejects.toThrow(/Network error/);
    });

    it("should throw ExecutorError when server returns string data", async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve("Some error from server"),
      });

      await expect(executor.validate("test")).rejects.toThrow(ExecutorError);
    });
  });

  describe("sample()", () => {
    it("should POST to correct URL with config and parsed seed", async () => {
      const mockSample = { config: "test", seed: "42" };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockSample),
      });

      const result = await executor.sample("test config", "42");

      expect(mockFetch).toHaveBeenCalledTimes(1);
      const [url, options] = mockFetch.mock.calls[0];
      expect(url).toMatch(/^http:\/\/localhost:8080\/sample\/id\d+$/);
      expect(options.method).toBe("POST");
      expect(JSON.parse(options.body)).toEqual({
        config: "test config",
        seed: 42,
      });

      expect(result).toEqual(mockSample);
    });

    it("should throw ExecutorError on network failure", async () => {
      mockFetch.mockRejectedValueOnce(new Error("Network error"));

      await expect(executor.sample("test", "42")).rejects.toThrow(ExecutorError);
    });
  });

  describe("cancel()", () => {
    it("should POST to cancel URL and set running to false", async () => {
      mockFetch.mockResolvedValueOnce({ ok: true });

      executor.cancel();

      expect(executor.running()).toBe(false);
      // cancel fires asynchronously, verify it was called
      expect(mockFetch).toHaveBeenCalledTimes(1);
      const url = mockFetch.mock.calls[0][0] as string;
      expect(url).toMatch(/^http:\/\/localhost:8080\/cancel\/id\d+$/);
    });
  });

  describe("run()", () => {
    it("should POST to run URL then poll results", async () => {
      const updateResult = vi.fn();

      // POST /run
      mockFetch.mockResolvedValueOnce({ ok: true });

      // GET /results — done on first poll
      const simResult = { statistics: {} };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () =>
          Promise.resolve({
            error: "",
            result: JSON.stringify(simResult),
            hash: "abc123",
            done: true,
          }),
      });

      const result = await executor.run("test config", updateResult);

      expect(result).toBe(true);
      expect(executor.running()).toBe(false);
      expect(updateResult).toHaveBeenCalledWith(simResult, "abc123");

      // First call is POST /run, second is GET /results
      expect(mockFetch).toHaveBeenCalledTimes(2);
      const [runUrl, runOpts] = mockFetch.mock.calls[0];
      expect(runUrl).toMatch(/\/run\/id\d+$/);
      expect(runOpts.method).toBe("POST");
    });

    it("should throw ExecutorError on run network failure", async () => {
      mockFetch.mockRejectedValueOnce(new Error("Network error"));

      await expect(executor.run("test", vi.fn())).rejects.toThrow(ExecutorError);
      expect(executor.running()).toBe(false);
    });

    it("should throw ExecutorError when results contain error", async () => {
      // POST /run succeeds
      mockFetch.mockResolvedValueOnce({ ok: true });
      // GET /results returns error
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () =>
          Promise.resolve({
            error: "simulation failed",
            result: "",
            done: false,
          }),
      });

      await expect(executor.run("test", vi.fn())).rejects.toThrow("simulation failed");
    });

    it("should throw ExecutorError on blank result", async () => {
      mockFetch.mockResolvedValueOnce({ ok: true });
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () =>
          Promise.resolve({
            error: "",
            result: "",
            done: false,
          }),
      });

      await expect(executor.run("test", vi.fn())).rejects.toThrow("blank result");
    });
  });

  describe("buildInfo()", () => {
    it("should return empty hash and date", () => {
      const info = executor.buildInfo();
      expect(info).toEqual({ hash: "", date: "" });
    });
  });
});
