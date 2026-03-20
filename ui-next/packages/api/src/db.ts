import { apiFetch } from "./client.js";

/**
 * Options for querying the database.
 */
export interface DBQueryOptions {
  /** MongoDB-style filter JSON string */
  query: string;
  /** Page number (1-indexed) */
  page?: number;
  /** Results per page */
  limit?: number;
  /** Abort signal for cancellation */
  signal?: AbortSignal;
}

/**
 * Result shape returned by the database query endpoint.
 */
export interface DBQueryResult {
  /** Array of database entries */
  data: unknown[];
  /** Total number of matching entries */
  total: number;
}

/**
 * Query the gcsim database with a MongoDB-style filter.
 *
 * @param options - Query options including filter, pagination, and abort signal
 * @returns Query results with data array and total count
 */
export function queryDB(options: DBQueryOptions): Promise<DBQueryResult> {
  const params = new URLSearchParams();
  params.set("q", options.query);

  if (options.page !== undefined) {
    params.set("page", String(options.page));
  }
  if (options.limit !== undefined) {
    params.set("limit", String(options.limit));
  }

  return apiFetch<DBQueryResult>(`/api/db?${params.toString()}`, {
    signal: options.signal,
  });
}
