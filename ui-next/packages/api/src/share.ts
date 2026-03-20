import type { Sim } from "@gcsim/types";
import { apiFetch } from "./client.js";

/**
 * Fetch a shared simulation result by ID.
 *
 * @param id - The share ID
 * @param options - Optional abort signal
 * @returns The simulation result
 */
export function fetchShareResult(
  id: string,
  options?: { signal?: AbortSignal },
): Promise<Sim.SimResults> {
  return apiFetch<Sim.SimResults>(`/api/share/${id}`, { signal: options?.signal });
}

/**
 * Fetch a database entry result by ID.
 *
 * @param id - The database entry ID
 * @param options - Optional abort signal
 * @returns The simulation result
 */
export function fetchDBResult(
  id: string,
  options?: { signal?: AbortSignal },
): Promise<Sim.SimResults> {
  return apiFetch<Sim.SimResults>(`/api/share/db/${id}`, { signal: options?.signal });
}
