import pako from "pako";

/**
 * Error thrown when an API request fails with a non-ok HTTP status.
 */
export class ApiError extends Error {
  public readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

/**
 * Base fetch wrapper that handles JSON parsing and gzip-compressed responses.
 *
 * @param url - The URL to fetch
 * @param options - Standard RequestInit options (signal, headers, etc.)
 * @returns Parsed JSON response of type T
 * @throws ApiError on non-ok HTTP responses
 */
export async function apiFetch<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, options);

  if (!response.ok) {
    throw new ApiError(response.status, `API error: ${response.status} ${response.statusText}`);
  }

  const contentEncoding = response.headers.get("content-encoding");
  if (contentEncoding === "gzip") {
    const buffer = await response.arrayBuffer();
    const decompressed = pako.inflate(new Uint8Array(buffer), { to: "string" });
    return JSON.parse(decompressed) as T;
  }

  return response.json() as Promise<T>;
}
