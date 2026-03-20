import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { queryDB } from "../db.js";

function createMockResponse(json: unknown): Response {
  return {
    ok: true,
    status: 200,
    statusText: "OK",
    headers: new Headers(),
    json: vi.fn().mockResolvedValue(json),
  } as unknown as Response;
}

describe("queryDB", () => {
  let mockFetch: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockFetch = vi.fn();
    vi.stubGlobal("fetch", mockFetch);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("sends query as URL parameter", async () => {
    const data = { data: [], total: 0 };
    mockFetch.mockResolvedValue(createMockResponse(data));

    await queryDB({ query: '{"char":"xiangling"}' });

    const calledUrl = mockFetch.mock.calls[0][0] as string;
    expect(calledUrl).toContain("/api/db?");
    expect(calledUrl).toContain("q=%7B%22char%22%3A%22xiangling%22%7D");
  });

  it("supports pagination params", async () => {
    const data = { data: [{ id: "1" }], total: 50 };
    mockFetch.mockResolvedValue(createMockResponse(data));

    await queryDB({ query: "{}", page: 2, limit: 10 });

    const calledUrl = mockFetch.mock.calls[0][0] as string;
    expect(calledUrl).toContain("page=2");
    expect(calledUrl).toContain("limit=10");
  });

  it("supports abort signal", async () => {
    const controller = new AbortController();
    const data = { data: [], total: 0 };
    mockFetch.mockResolvedValue(createMockResponse(data));

    await queryDB({ query: "{}", signal: controller.signal });

    const calledOptions = mockFetch.mock.calls[0][1] as RequestInit;
    expect(calledOptions.signal).toBe(controller.signal);
  });

  it("abort signal triggers AbortError", async () => {
    const controller = new AbortController();
    mockFetch.mockImplementation(() => {
      return Promise.reject(new DOMException("The operation was aborted.", "AbortError"));
    });

    controller.abort();

    await expect(queryDB({ query: "{}", signal: controller.signal })).rejects.toThrow(
      "The operation was aborted.",
    );
  });
});
