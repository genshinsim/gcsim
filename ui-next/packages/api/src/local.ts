import type { Sim } from "@gcsim/types";
import { apiFetch } from "./client.js";

const DEFAULT_LOCAL_URL = "http://127.0.0.1:8381";

/**
 * Fetch simulation results from the local dev server.
 *
 * @param baseUrl - Base URL of the local dev server (defaults to http://127.0.0.1:8381)
 * @param options - Optional abort signal
 * @returns The simulation result
 */
export function fetchLocalResult(
  baseUrl?: string,
  options?: { signal?: AbortSignal },
): Promise<Sim.SimResults> {
  const url = `${baseUrl ?? DEFAULT_LOCAL_URL}/data`;
  return apiFetch<Sim.SimResults>(url, { signal: options?.signal });
}
