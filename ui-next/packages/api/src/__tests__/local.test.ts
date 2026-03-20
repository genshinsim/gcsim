import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { fetchLocalResult } from "../local.js";

function createMockResponse(json: unknown): Response {
  return {
    ok: true,
    status: 200,
    statusText: "OK",
    headers: new Headers(),
    json: vi.fn().mockResolvedValue(json),
  } as unknown as Response;
}

describe("fetchLocalResult", () => {
  let mockFetch: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockFetch = vi.fn();
    vi.stubGlobal("fetch", mockFetch);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("uses default URL when none provided", async () => {
    const data = { sim_version: "local" };
    mockFetch.mockResolvedValue(createMockResponse(data));

    await fetchLocalResult();

    expect(mockFetch).toHaveBeenCalledWith("http://127.0.0.1:8381/data", {
      signal: undefined,
    });
  });

  it("uses custom URL when provided", async () => {
    const data = { sim_version: "local" };
    mockFetch.mockResolvedValue(createMockResponse(data));

    await fetchLocalResult("http://localhost:9000");

    expect(mockFetch).toHaveBeenCalledWith("http://localhost:9000/data", {
      signal: undefined,
    });
  });

  it("returns typed SimResults", async () => {
    const data = { sim_version: "local", config_file: "test" };
    mockFetch.mockResolvedValue(createMockResponse(data));

    const result = await fetchLocalResult();

    expect(result).toEqual(data);
  });
});
