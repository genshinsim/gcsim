# @gcsim/api

## Purpose

Typed fetch functions for all gcsim API endpoints. Uses native `fetch` (not axios) with pako for gzip decompression.

## Public API

All exports go through `src/index.ts`:

- **`apiFetch<T>(url, options?)`** — base fetch wrapper with error handling and gzip support
- **`ApiError`** — error class with HTTP status code
- **`fetchShareResult(id, options?)`** — fetch shared sim result from `/api/share/{id}`
- **`fetchDBResult(id, options?)`** — fetch DB entry result from `/api/share/db/{id}`
- **`queryDB(options)`** — query the database with filter/pagination
- **`fetchLocalResult(baseUrl?, options?)`** — fetch from local dev server

## Endpoints

| Function | Method | URL |
|----------|--------|-----|
| `fetchShareResult` | GET | `/api/share/{id}` |
| `fetchDBResult` | GET | `/api/share/db/{id}` |
| `queryDB` | GET | `/api/db?q={query}&page={page}&limit={limit}` |
| `fetchLocalResult` | GET | `http://127.0.0.1:8381/data` (default) |

## Dependencies

- `@gcsim/types` — `Sim.SimResults` type
- `pako` — gzip decompression for compressed API responses

## Testing

Tests mock `fetch` via `vi.stubGlobal` and verify:
- Correct URL construction
- Error handling (ApiError with status)
- Gzip decompression path
- AbortSignal support

Run: `pnpm test` or `turbo run test --filter=@gcsim/api`

## Don'ts

- Don't use axios — we use native fetch
- Don't import from `@gcsim/types` internal paths
- Don't make actual HTTP calls in tests — always mock fetch
