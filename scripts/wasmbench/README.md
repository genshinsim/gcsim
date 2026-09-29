# wasmbench

Benchmark harness and correctness oracle for the web UI's wasm build (`cmd/wasm`). It drives
the binary through the same JS API and call sequence the UI's web workers use, so the numbers
include the JS↔Go boundary and msgpack/JSON serialization.

Requirements: Go, node ≥ 23.9 for thread-CPU timing (older node falls back to wall time). No npm packages.

## Quick start

```sh
export WASMBENCH_HOME=/some/shared/dir   # lock, goldens, results (default ${TMPDIR:-/tmp}/wasmbench)
B=$WASMBENCH_HOME/bin

scripts/wasmbench/build.sh $B/baseline.wasm                 # from a clean checkout of the base commit
scripts/wasmbench/check.sh --regen $B/baseline.wasm         # once: write goldens, prove they reproduce

EXTRA_BUILD_FLAGS='-gcflags=all=-B' scripts/wasmbench/build.sh $B/cand.wasm
scripts/wasmbench/check.sh $B/cand.wasm                     # same results as baseline? exit 0/1
scripts/wasmbench/compare.sh $B/baseline.wasm $B/cand.wasm  # paired single-worker ms/iter + ratios
WORKERS=3 scripts/wasmbench/compare.sh $B/baseline.wasm $B/cand.wasm   # UI-style worker pool

node scripts/wasmbench/bench.mjs run --wasm $B/cand.wasm \
  --config scripts/wasmbench/configs/national_zhongli.txt   # one ad-hoc run, all metrics (no lock)
```

## Scripts

| file | what |
|---|---|
| `build.sh <out.wasm>` | `GOOS=js GOARCH=wasm go build -trimpath -ldflags "-X main.shareKey=…" ./cmd/wasm`, the same as `cmd/wasm/build.sh` (CI deploy) and `task wasm` (`pnpm build:wasm`). Env: `GO`, `GOWASM`, `EXTRA_BUILD_FLAGS`, `BASE_BUILD_FLAGS` (default `-trimpath`), `PKG`, `WASM_OPT` (`ci` for the flags CI deploys with, or any wasm-opt arguments), `WASM_OPT_BIN`, `WASM_EXEC`. Also copies the toolchain's `wasm_exec.js` to `<out>.wasm_exec.js`. |
| `compare.sh a.wasm b.wasm [c.wasm…]` | Fixed seeds, fixed iteration counts, binaries interleaved. Prints per config: median ms/iter, IQR/median, ratio vs the first binary with a bootstrap 95% CI, and the geomean ratio. Env: `CONFIGS`, `PROCS` (3), `ROUNDS` (20 paired / 4 pool), `BLOCK`, `WARMUP`, `WORKERS` (0), `ITERS` (pool, 1000), `SEED`, `METRIC`, `OUT`, `NODE_FLAGS`. |
| `check.sh cand.wasm` / `check.sh --regen base.wasm` | Correctness oracle against goldens in `$WASMBENCH_HOME/wasmbench-golden` (`GOLDEN_DIR`). Env: `CHECK_ITERS` (100), `SEED`, `CONFIGS`, `RTOL`. |
| `bench.mjs` | The runner: `run`, `ab`, `dump`, `diff`, `report` subcommands (see its header). |
| `configs/` | Team configs plus `manifest.tsv` (per-config block size, warmup, source). |

`compare.sh` and `check.sh` take a mkdir lock at `$WASMBENCH_LOCK` (default
`$WASMBENCH_HOME/wasmbench.lock`) so only one benchmark runs on the machine at a time. They
wait while it's held, clear it if the holder died, and release it on exit or signal.
`WASMBENCH_NOLOCK=1` skips the lock. Direct `bench.mjs` calls don't take it.

## What is measured

The UI (`ui/packages/executors/src/WasmExecutor.ts`, `Workers/worker.ts`, `Workers/aggregator.ts`)
runs N sim workers plus one aggregator worker, each with its own wasm instance:

- sim worker: `initializeWorker(cfg)`, then `simulate()` → msgpack `Uint8Array`, once per iteration
- aggregator: `initializeAggregator(cfg)` → JSON, `aggregate(bytes)` per result, `flush()` → JSON
  at most every 100 ms (lodash throttle, leading + trailing)

Every instance gets its own JS realm (a `vm` context holding only what a Worker gives the glue,
or a worker_thread's own global), and results are copied into the aggregator's realm the way
`postMessage` does.

- **paired** (`bench.mjs ab`, `compare.sh` default): one process loads every binary side by side
  (sim + aggregator instance each). All binaries warm up (covering V8's Liftoff → TurboFan
  tier-up). Then each round runs the same fixed-seed block, about 150 ms, on each binary back
  to back, and the order rotates every round. `simulate()` is called back to back as in a
  worker, and the block is then aggregated. The ratio is the median of per-round ratios, so
  load that hits both halves of a round cancels. `PROCS` processes × `ROUNDS` rounds per
  config. The primary metric is **thread CPU time** per `simulate()`
  (`process.threadCpuUsage`). It leaves out time spent descheduled by other processes. On an
  idle machine it equals wall time, and wall ratios are shown next to it. `METRIC=ms_per_iter`
  switches the primary metric to wall.
- **single** (`bench.mjs run`): one binary, `--warmup` then `--iters` iterations, with
  per-call timings (mean, p10/p50/p90), aggregate and flush time, result size, and wasm memory.
- **pool** (`--workers N`, `WORKERS=N`): N `worker_threads` sim workers plus an aggregator
  thread, dispatched the same way as `WasmExecutor.run`, including structured clones through
  the main thread and throttled flushes. `ms_per_iter` = wall time from dispatch to the last
  `aggregate()` ÷ iterations. `ui_wall_ms` also includes the trailing flush the UI waits for.
  There is no warmup, the same as the UI, and the default is 1000 iterations, the UI's default.
  The UI defaults to 3 workers (user-adjustable, 1–30).

## Seeds and determinism

`cmd/wasm` seeds every `simulate()` from `crypto/rand`, which on js/wasm calls
`crypto.getRandomValues`. During `simulate()`, `bench.mjs` serves those bytes from splitmix64
keyed by (`--seed`, iteration index), so every binary runs the same seeds. The dump checks
that each result's `seed` field matches. Outside `simulate()` (Go runtime init, `sample_seed`)
the real CSPRNG is used, as in the browser.

With fixed seeds, wasm output is bit-for-bit reproducible across processes. That holds even
though the Go runtime's own randomness (map iteration order) differs between runs, so the
oracle compares exactly:

- every iteration's `stats.Result` (msgpack decoded, canonicalised with sorted keys; nil and empty
  slices/maps treated as equal), hashed
- the aggregated `flush()` statistics (the full JSON is kept in the golden, so a mismatch shows the paths that differ)
- `initializeAggregator()` metadata without build info and `sample_seed`, `validateConfig()`,
  and `sample(cfg, seed)`

`--regen` dumps twice in separate processes and fails if they differ. A change that scales
damage by 1+1e-12 fails the check. If a candidate legitimately changes float evaluation
order, `RTOL` relaxes the aggregated-stats comparison. Per-iteration hashes stay exact.

## Configs

| name | chars | source |
|---|---|---|
| `sucrose_solo` | 1, a few actions | `ui/packages/e2e/fixtures/sucrose.txt`. Very short sim, so it mostly measures per-iteration fixed cost (core setup, serialization, GC) |
| `raiden_national` | 4, energy-starved | docs basic tutorial, last full config |
| `overload_chevreuse` | 4, turrets | storybook `sampleConfig.ts` |
| `geo_yelan_albedo` | 4, 109 s | `pkg/simulation/example/main.go` |
| `national_zhongli` | 4, vape + geo | storybook `sampleResult.json` |

## Caveats

- CI runs binaryen version_133's `wasm-opt -Oz` (with the features Go's output needs) on the
  deployed binary. `build.sh` skips it unless `WASM_OPT` is set; `WASM_OPT=ci` uses CI's flags.
  Point `WASM_OPT_BIN` at a version_133 `wasm-opt` to match CI.
- The UI serves its own checked-in `ui/packages/web/public/wasm_exec.js`, an older copy that
  lacks a few test-only hooks. The harness uses the glue matching the binary's toolchain
  (`<name>.wasm_exec.js`, else `$(go env GOROOT)/lib/wasm/wasm_exec.js`). Override with `--glue`.
- Two binaries with identical bytes share V8's compiled module inside one process. For a
  self-comparison, append a custom section to one copy so the bytes differ.
- Node's V8 is not Chrome's V8 build, and worker_threads are not browser Workers. Ratios should
  carry over. Absolute numbers are approximate.
