import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { fetchDBResult, fetchShareResult } from "../share.js";

function createMockResponse(json: unknown): Response {
  return {
    ok: true,
    status: 200,
    statusText: "OK",
    headers: new Headers(),
    json: vi.fn().mockResolvedValue(json),
  } as unknown as Response;
}

describe("fetchShareResult", () => {
  let mockFetch: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockFetch = vi.fn();
    vi.stubGlobal("fetch", mockFetch);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("constructs correct URL /api/share/{id}", async () => {
    const data = { schema_version: { major: 1 } };
    mockFetch.mockResolvedValue(createMockResponse(data));

    await fetchShareResult("abc123");

    expect(mockFetch).toHaveBeenCalledWith("/api/share/abc123", { signal: undefined });
  });

  it("returns typed SimResults", async () => {
    const data = { sim_version: "2.0.0", config_file: "test config" };
    mockFetch.mockResolvedValue(createMockResponse(data));

    const result = await fetchShareResult("test-id");

    expect(result).toEqual(data);
  });

  it("passes abort signal", async () => {
    const data = {};
    mockFetch.mockResolvedValue(createMockResponse(data));
    const controller = new AbortController();

    await fetchShareResult("test-id", { signal: controller.signal });

    expect(mockFetch).toHaveBeenCalledWith("/api/share/test-id", {
      signal: controller.signal,
    });
  });
});

describe("fetchDBResult", () => {
  let mockFetch: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockFetch = vi.fn();
    vi.stubGlobal("fetch", mockFetch);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("constructs correct URL /api/share/db/{id}", async () => {
    const data = {};
    mockFetch.mockResolvedValue(createMockResponse(data));

    await fetchDBResult("db-entry-1");

    expect(mockFetch).toHaveBeenCalledWith("/api/share/db/db-entry-1", { signal: undefined });
  });

  it("returns typed SimResults", async () => {
    const data = { sim_version: "2.0.0", mode: 1 };
    mockFetch.mockResolvedValue(createMockResponse(data));

    const result = await fetchDBResult("db-entry-1");

    expect(result).toEqual(data);
  });
});
