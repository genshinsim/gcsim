# @gcsim/e2e

Local, Playwright-driven end-to-end smoke tests for gcsim, catching regressions
where a site *builds* but breaks at *runtime*. Three suites, each its own
Playwright config:

- **web** (`playwright.config.ts`, default) — boots the web app, waits for
  wasm + workers, validates and runs a config, and asserts the viewer renders.
- **docs** (`playwright.docs.config.ts`) — builds and serves the Docusaurus docs
  site, then asserts the landing page and a nested doc render with a clean
  console.
- **db** ("Simpact", `playwright.db.config.ts`) — boots the db app and asserts
  the home and browse views render. The db app is pure front-end (no wasm), and
  the spec stubs every `/api` call, so it needs neither Go nor network.

Two layers:

- **Harness** (`src/`) — an importable, test-runner-agnostic layer that encodes
  the app's implicit protocol once (page objects + fixtures). This is the seam a
  future agent-exploration head can build on.
- **Specs** (`tests/`) — regression specs that consume the harness and assert
  the critical path.

## Prerequisites

The **web** suite runs against the **dev server**, not a production build: in a
prod build the wasm URL points at a remote origin (R2 / `/api/wasm/...`) that
isn't available locally.

The Playwright `webServer` builds the wasm binary and then starts `vite`. The
wasm build shells out to `task wasm` (`go build` for `GOOS=js`), so you need:

- **Go** and **[Task](https://taskfile.dev)** (`task`) on your `PATH` — the wasm
  build step.
- **pnpm** deps installed: from `ui/`, run `pnpm install`.
- **Playwright's Chromium** browser: from `ui/`, run
  `pnpm --filter @gcsim/e2e exec playwright install chromium`.

The **docs** suite needs none of the wasm toolchain — no Go, no `task`. It
serves a *production build* (`docusaurus build` then `docusaurus serve`): docs
has no wasm, workers, or backend, and building first makes the smoke spec catch
build-time breakage too (broken links, MDX compile errors, a bad sidebar entry).
It still needs the pnpm deps and the Chromium browser above.

The **db** suite also needs none of the wasm toolchain. It runs the db vite dev
server and stubs every `/api` call, so it needs no Go, no backend, and no
network — only the pnpm deps and the Chromium browser above.

## Run it

From the `ui/` workspace root:

```sh
pnpm test:e2e       # web app suite  (builds wasm; needs Go + task)
pnpm test:e2e:docs  # docs site suite
pnpm test:e2e:db    # db app suite    (no wasm, no network)
```

`test:e2e` builds the wasm binary, boots the dev server on a fixed strict port
(`5173`), runs the suite headless in Chromium, and tears the server down after.
A running dev server on `5173` is reused (locally) instead of restarted.

`test:e2e:docs` builds the docs site, serves it on port `4173`, runs the suite,
and tears the server down after — a running server on `4173` is reused locally.

`test:e2e:db` boots the db dev server on `5273` and runs the db suite; a running
server on `5273` is reused locally.

Other entry points (from `ui/packages/e2e`):

- `pnpm --filter @gcsim/e2e test:headed` / `test:db:headed` — watch it drive a
  real browser.
- `pnpm --filter @gcsim/e2e report` — open the HTML report from the last run.

## On failure

An HTML report is written to `playwright-report/`, and a Playwright **trace**
(plus screenshot and video) is retained under `test-results/` for any failing
test. Open the report with `pnpm --filter @gcsim/e2e report`, or a single trace
with `pnpm --filter @gcsim/e2e exec playwright show-trace <trace.zip>`.

## What the harness helpers do

Neither app ships **data-testids**, so every locator rides an observable
contract (a DOM id, a Blueprint/theme class, or an accessible name).

### `AppHarness` (`src/app-harness.ts`) — web app

Owns the full boot → run → complete flow and bundles the page objects and the
console monitor.

- `boot()` — open `/simulator`, then wait for wasm + workers ready.
- `run(cfg)` — set the config, wait for it to validate, click Run, and wait for
  the run to complete (`run time: N ms` in the console). The signal is matched
  against the console monitor's **buffer** rather than a fresh `waitForEvent`,
  because a one-iteration sim can log it before a post-click listener would
  attach.

### `SimulatorPage` (`src/pages/simulator-page.ts`) — the `/simulator` route

- `goto()` — navigate and wait for React to mount into `#root`.
- `waitForReady()` — wait for the **Run** button (Blueprint, accessible name
  "Run") to leave its loading state (`bp4-loading` spinner). That transition is
  wasm + workers becoming ready (console: `aggregator loaded okay`,
  `loading N workers`).
- `setConfig(cfg)` — replace the Ace editor (`#config_editor`) contents by
  dispatching a native **paste** on its proxy textarea. Typing key-by-key would
  trip Ace's auto-indent/bracket-matching and corrupt the config.
- `waitForConfigValid()` — wait for Run to become enabled with no "Invalid
  Config" callout (console: `all is good`).
- `run()` — click Run and wait for the app to navigate to `/web`.

### `ViewerPage` (`src/pages/viewer-page.ts`) — the `/web` route

- `waitForViewer()` — assert the page title `gcsim - viewer` and the
  **Results / Config / Sample** tab strip.
- `waitForResults()` — assert the Results tab rendered at least one card
  (`bp4-card`) and one inline `<svg>` chart. **Structural only** — never asserts
  DPS or any numeric result.

### `DocsHarness` (`src/docs-harness.ts`) — docs site

Docs is a single, uniform surface (every route is a rendered MDX page with the
same sidebar chrome), so one class owns the whole surface plus a
`ConsoleMonitor`.

- `goto(path = "/")` — navigate and wait for the doc body
  (`article .theme-doc-markdown`) to render.
- `sidebar` — the `<nav aria-label="Docs sidebar">`; `sidebarLink(name)` is a
  sidebar entry by its exact label.
- `heading(name)` — the current doc's `<h1>` by accessible name.

### `ConsoleMonitor` (`src/console-monitor.ts`)

Watches the page console for its whole lifetime. It collects uncaught
`console.error` / `pageerror` — `assertNoErrors()` fails the test if any
un-ignored error was seen across the flow — and buffers every message so
`waitForLog(re)` can match a signal even if it already fired. A short, justified
ignore-list filters known dev-mode noise: React StrictMode `Warning:`s (dev-only,
absent from prod builds — a real crash is a `pageerror` or a non-`Warning:`
error, which is **not** masked), a pre-existing visx negative-radius chart quirk,
and a Vite HMR websocket blip on teardown.

### Config fixture

`fixtures/sucrose.txt` is the known-good, minimal Sucrose config
(`iteration=1`), re-exported as `sucroseConfig` from `src/config.ts`.

### `DbHarness` (`src/db-harness.ts`) — the db app

Bundles the db page objects (`DbHomePage`, `DbDatabasePage`) and a
`ConsoleMonitor`. The `db` test fixture stubs the app's network before any
navigation via `installDbRoutes` (`src/db-fixtures.ts`):

- `/api/db` returns the two `dbEntries`, narrowed to the characters an included
  filter names — so a character filter visibly shrinks the list (2 → 1);
- `/api/assets/**` returns a 1x1 PNG (avatars, art);
- `api.github.com` (the home page's latest-release lookup) returns a fixed
  payload;
- any other `/api/**` returns an empty 200.

The db dev server proxies `/api` to production by default; these routes keep the
spec deterministic and offline.

## Out of scope

CI integration, the production/preview (R2 wasm) path, server mode, the
`/db/:id` viewer render, local/share viewer routes, Enka/GOOD import,
engine-correctness or numeric assertions, and non-Chromium browsers. See issues
#2805 and #2869. For docs: the search backend, i18n/translations, and visual
regression (issue #2868).
