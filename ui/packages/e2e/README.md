# @gcsim/e2e

Local, Playwright-driven end-to-end smoke tests for gcsim, catching regressions
where a site *builds* but breaks at *runtime*. Four suites, each its own
Playwright config:

- **web** (`playwright.config.ts`, default) — boots the web app, waits for
  wasm + workers, validates and runs a config, and asserts the viewer renders.
  Beyond that smoke path it also covers the dash home, the viewer's Config and
  Sample tabs, and the GOOD / Enka toolbox imports (each stubbing the network it
  needs so it stays offline).
- **docs** (`playwright.docs.config.ts`) — builds and serves the Docusaurus docs
  site, then asserts one page per top-level sidebar section renders (route,
  title, `<h1>`, non-empty body) with its content images loaded, failing only on
  an uncaught exception — not on ambient console noise.
- **db** ("Simpact", `playwright.db.config.ts`) — boots the db app and asserts
  the home and browse views render. The db app is pure front-end (no wasm), and
  the spec stubs every `/api` call, so it needs neither Go nor network.
- **taghelper** (`playwright.taghelper.config.ts`) — boots the taghelper app (the
  Discord moderation helper) and asserts the `/id/:id` moderation view renders.
  Like db, it is pure front-end and stubs every `/api` call, so it needs neither
  Go nor network.

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

The **db** and **taghelper** suites also need none of the wasm toolchain. Each
runs its vite dev server and stubs every `/api` call, so it needs no Go, no
backend, and no network — only the pnpm deps and the Chromium browser above.

## Run it

From the `ui/` workspace root:

```sh
pnpm test:e2e            # web app suite  (builds wasm; needs Go + task)
pnpm test:e2e:docs       # docs site suite
pnpm test:e2e:db         # db app suite         (no wasm, no network)
pnpm test:e2e:taghelper  # taghelper app suite  (no wasm, no network)
```

`test:e2e` builds the wasm binary, boots the dev server on a fixed strict port
(`5173`), runs the suite headless in Chromium, and tears the server down after.
A running dev server on `5173` is reused (locally) instead of restarted.

`test:e2e:docs` builds the docs site, serves it on port `4173`, runs the suite,
and tears the server down after — a running server on `4173` is reused locally.

`test:e2e:db` boots the db dev server on `5273` and runs the db suite; a running
server on `5273` is reused locally.

`test:e2e:taghelper` boots the taghelper dev server on `5174` and runs the
taghelper suite; a running server on `5174` is reused locally.

Other entry points (from `ui/packages/e2e`):

- `pnpm --filter @gcsim/e2e test:headed` / `test:db:headed` /
  `test:taghelper:headed` — watch it drive a real browser.
- `pnpm --filter @gcsim/e2e report` — open the HTML report from the last run.

## On failure

An HTML report is written to `playwright-report/`, and a Playwright **trace**
(plus screenshot and video) is retained under `test-results/` for any failing
test. Open the report with `pnpm --filter @gcsim/e2e report`, or a single trace
with `pnpm --filter @gcsim/e2e exec playwright show-trace <trace.zip>`.

## What the harness helpers do

None of the apps ship **data-testids**, so every locator rides an observable
contract (a DOM id, an ARIA role, or an accessible name).

### `AppHarness` (`src/app-harness.ts`) — web app

Owns the full boot → run → complete flow and bundles the page objects and the
console monitor.

- `boot()` — open `/simulator`, then wait for wasm + workers ready.
- `run(cfg)` — set the config, wait for it to validate, click Run, and wait for
  the run to complete (`run time: N ms` in the console). The signal is matched
  against the console monitor's **buffer** rather than a fresh `waitForEvent`,
  because a one-iteration sim can log it before a post-click listener would
  attach.

### `DashPage` (`src/pages/dash-page.ts`) — the dash home (`/`)

- `goto()` — navigate and wait for React to mount into `#root`.
- `waitForLoaded()` — assert the title, the nav bar's Simulator entry, and the
  "Get started" CTA. **Structural only.**
- `waitForFeatured()` — assert the featured-submission card's "Show Detail" link
  and the "Visit the Teams DB" CTA. The card fetches `/api/db`, so a spec stubs
  that route (reusing `installDbRoutes`) before navigating.

### `SimulatorPage` (`src/pages/simulator-page.ts`) — the `/simulator` route

- `goto()` — navigate and wait for React to mount into `#root`.
- `waitForReady()` — wait for the **Run** button (accessible name "Run") to drop
  its loading spinner (`role="status"`). That transition is wasm + workers
  becoming ready (console: `aggregator loaded okay`, `loading N workers`).
- `setConfig(cfg)` — replace the Ace editor (`#config_editor`) contents by
  dispatching a native **paste** on its proxy textarea. Typing key-by-key would
  trip Ace's auto-indent/bracket-matching and corrupt the config.
- `waitForConfigValid()` — wait for Run to become enabled with no "Invalid
  Config" callout (console: `all is good`).
- `run()` — click Run and wait for the app to navigate to `/web`.
- `openImportDialog("GO" | "Enka")` — open the Toolbox "Tools" popover, click the
  matching import entry, and return the resulting dialog (`role="dialog"`). Wasm
  need not be ready — the toolbox renders on mount.

### `ViewerPage` (`src/pages/viewer-page.ts`) — the `/web` route

- `waitForViewer()` — assert the page title `gcsim - viewer` and the
  **Results / Config / Sample** tab strip.
- `waitForResults()` — assert the Results tab rendered at least one card
  (`data-slot="card"`) and one inline chart (`svg[role="img"]`). **Structural
  only** — never asserts DPS or any numeric result.
- `openConfig()` — click the Config tab and assert the config editor
  (`#config_editor`) rendered and is non-empty.
- `openSample()` — click the Sample tab and assert its "Generate" control
  rendered. Does not generate a sample.

### `DocsHarness` (`src/docs-harness.ts`) — docs site

Docs is a single, uniform surface (every route is a rendered MDX page with the
same sidebar chrome), so one class owns the whole surface plus a
`ConsoleMonitor`.

- `goto(path = "/")` — navigate and wait for the doc body
  (`article .theme-doc-markdown`) to render.
- `sidebar` — the `<nav aria-label="Docs sidebar">`; `sidebarLink(name)` is a
  sidebar entry by its exact label.
- `heading(name)` — the current doc's `<h1>` by accessible name.
- `assertImagesLoaded(container)` — every `<img>` under `container` finished
  loading (`complete` and `naturalWidth > 0`); scrolls lazy images into view
  first. The deterministic replacement for the old blanket 404 console sweep.

### `ConsoleMonitor` (`src/console-monitor.ts`)

Watches the page console for its whole lifetime, and offers two contracts:

- `assertNoErrors()` fails on any un-ignored `console.error` or `pageerror` seen
  across the flow (used by the web, db, and taghelper suites);
- `assertNoCrashes()` fails only on an uncaught exception (`pageerror`), ignoring
  ambient `console.error` noise such as a transient resource 404 (used by the
  docs suite, whose production build serves noisy static assets on cold start).

It also buffers every message so `waitForLog(re)` can match a signal even if it
already fired. A short, justified ignore-list filters known dev-mode noise for
`assertNoErrors`: React StrictMode `Warning:`s (dev-only, absent from prod builds
— a real crash is a `pageerror` or a non-`Warning:` error, which is **not**
masked), a pre-existing visx negative-radius chart quirk, and a Vite HMR
websocket blip on teardown.

### Config fixture

`fixtures/sucrose.txt` is the known-good, minimal Sucrose config
(`iteration=1`), re-exported as `sucroseConfig` from `src/config.ts`.

### Import fixtures (`src/import-fixtures.ts`)

The toolbox-import specs' fixtures, both re-exported from `src/`:

- `goodImport` — a minimal, parseable GOOD payload (one character, Bennett),
  read from `fixtures/good-import.json` (committed as JSON so it reads like a
  payload a user would paste, mirroring `sucrose.txt`).
- `installEnkaRoutes(page)` — stubs `/api/enka/:uid` (which the dev server
  otherwise proxies to production) so the Enka import resolves offline. It
  returns `enkaImportPayload`, a minimal Enka response that `EnkaToGOOD` parses
  into one character (Bennett), for `ENKA_UID`.

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

#### `DbHomePage` (`src/pages/db-home-page.ts`) — the `/` route

- `goto()` — navigate and wait for React to mount into `#root`.
- `waitForLoaded()` — assert the welcome copy, the tag-list copy, and the "Get
  started" CTA rendered.

#### `DbDatabasePage` (`src/pages/db-database-page.ts`) — the `/database` route

- `goto()` / `waitForBrowse()` — open the route; assert the count, search box,
  filter funnel, and an entry card's Copy Config / Open in Viewer controls.
- `openFilterPanel()` / `expandCharacters()` — open the filter drawer; expand
  the Characters section to its portrait picker.
- `filterByCharacter(name)` — pick a character from the search box, which
  refetches `/api/db` with the narrowed query; pair with `expectShowing(n)`.

### `TaghelperHarness` (`src/taghelper-harness.ts`) — the taghelper app

Taghelper is a single surface — the `/id/:id` moderation view for one db entry —
so, like `DocsHarness`, one class owns the whole surface plus a `ConsoleMonitor`.
The `taghelper` test fixture stubs the app's network before any navigation via
`installTaghelperRoutes` (`src/taghelper-fixtures.ts`): `/api/db/id/:id` returns
the `mainEntry`, `/api/db` returns the related entries, and `/api/assets/**`
returns a 1x1 PNG. The dev server proxies `/api` to production (simimpact.app) by
default; these routes keep the spec deterministic and offline.

- `goto(id)` — navigate to `/id/:id` and wait for the main entry's heading.
- `waitForEntry(chars, sourceTag)` — assert the team portraits (by `alt`) and the
  summary stat chips (mode / target / dps / avg sim time / created / source tag).
  **Structural only** — never asserts numeric results.
- `waitForControls()` — assert the Copy Reject / Copy Approve / Result Viewer
  moderation controls.
- `waitForExistingSims()` — assert the "existing sims" section rendered either
  result rows (each with a Replace This control) or the "Nothing found" empty
  state.
- `copyCommand(name)` — click a copy button and return the resulting clipboard
  text (e.g. `/approve id:…`).

## Out of scope

CI integration, the production/preview (R2 wasm) path, server mode, the
`/db/:id` viewer render, local/share viewer routes and the share flow,
engine-correctness or numeric assertions, and non-Chromium browsers. See issues
#2805, #2869 and #2871. For docs: the search backend, i18n/translations, and
visual regression (issue #2868).
