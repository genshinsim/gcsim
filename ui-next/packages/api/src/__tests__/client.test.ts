import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ApiError, apiFetch } from "../client.js";

function createMockResponse(options: {
  ok?: boolean;
  status?: number;
  statusText?: string;
  json?: unknown;
  headers?: Record<string, string>;
  arrayBuffer?: ArrayBuffer;
}): Response {
  const { ok = true, status = 200, statusText = "OK", json, headers = {}, arrayBuffer } = options;

  return {
    ok,
    status,
    statusText,
    headers: new Headers(headers),
    json: vi.fn().mockResolvedValue(json),
    arrayBuffer: vi.fn().mockResolvedValue(arrayBuffer),
  } as unknown as Response;
}

describe("apiFetch", () => {
  let mockFetch: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockFetch = vi.fn();
    vi.stubGlobal("fetch", mockFetch);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("returns parsed JSON on success", async () => {
    const data = { version: "1.0", statistics: {} };
    mockFetch.mockResolvedValue(createMockResponse({ json: data }));

    const result = await apiFetch<typeof data>("/api/test");

    expect(mockFetch).toHaveBeenCalledWith("/api/test", undefined);
    expect(result).toEqual(data);
  });

  it("throws ApiError with status on non-ok response", async () => {
    mockFetch.mockResolvedValue(
      createMockResponse({ ok: false, status: 404, statusText: "Not Found" }),
    );

    await expect(apiFetch("/api/missing")).rejects.toThrow(ApiError);
    await expect(apiFetch("/api/missing")).rejects.toMatchObject({
      status: 404,
      message: "API error: 404 Not Found",
    });
  });

  it("handles gzip-compressed responses", async () => {
    const pako = await import("pako");
    const payload = { compressed: true };
    const compressed = pako.deflate(JSON.stringify(payload));

    mockFetch.mockResolvedValue(
      createMockResponse({
        headers: { "content-encoding": "gzip" },
        arrayBuffer: compressed.buffer as ArrayBuffer,
      }),
    );

    const result = await apiFetch<typeof payload>("/api/gzipped");

    expect(result).toEqual(payload);
  });
});

describe("ApiError", () => {
  it("has correct name and properties", () => {
    const error = new ApiError(500, "Internal Server Error");

    expect(error).toBeInstanceOf(Error);
    expect(error).toBeInstanceOf(ApiError);
    expect(error.name).toBe("ApiError");
    expect(error.status).toBe(500);
    expect(error.message).toBe("Internal Server Error");
  });
});
