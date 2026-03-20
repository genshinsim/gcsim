# @gcsim/executor

## Purpose

Provides the `Executor` interface and implementations for running gcsim simulations:

1. **Executor interface** — contract for running, validating, and sampling sim configs
2. **ServerExecutor** — HTTP-based executor that communicates with a backend server
3. **WasmExecutor** — Web Worker pool executor that runs WASM simulations in-browser

## Public API

All public exports go through `src/index.ts`:

- **`Executor`** (type) — interface for all executor implementations
- **`ExecutorSupplier<T>`** (type) — factory type `() => T`
- **`ServerExecutor`** — HTTP executor using native `fetch`
- **`WasmExecutor`** — WASM worker pool executor

Usage:
```typescript
import { ServerExecutor, WasmExecutor } from "@gcsim/executor";
import type { Executor } from "@gcsim/executor";

const server = new ServerExecutor("http://localhost:8080");
const wasm = new WasmExecutor("/main.wasm");
```

## Architecture

### ServerExecutor
- Uses native `fetch` (no axios)
- Instance ID for request isolation (each instance gets a unique ID)
- Polls `/results/{id}` every 100ms during runs
- Caches `ready()` result after first successful call

### WasmExecutor
- Creates an aggregator worker + N simulation workers (default 3, range 1-30)
- Workers load WASM binary and communicate via postMessage
- Uses a HelperExecutor (single worker) for validate/sample operations
- Throttled flush (100ms) to batch result updates during runs

### Worker Types
- `workers/common.ts` — message type definitions (Aggregator, SimWorker, Helper namespaces)
- `workers/aggregator.ts` — aggregator worker that merges simulation results
- `workers/worker.ts` — simulation worker that runs individual iterations
- `workers/helper.ts` — helper worker for validate/sample operations
- `workers/wasm-types.d.ts` — WASM global function declarations

## Testing

Tests mock `fetch` (for ServerExecutor) and `Worker` (for WasmExecutor) since real WASM cannot run in test environment.

Run tests: `pnpm test` or `turbo run test --filter=@gcsim/executor`

## Don'ts

- Don't use axios — use native `fetch`
- Don't use lodash — use built-in throttle utility
- Don't import from other packages' `src/` — only from their package index
