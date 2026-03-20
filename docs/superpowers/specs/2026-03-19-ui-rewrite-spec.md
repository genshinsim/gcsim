# UI Rewrite — Formal Specification

## 1. Overview

Complete rewrite of the gcsim web UI as a pnpm + Turborepo monorepo. Multiple independently deployable web apps share feature packages, types, data, and localization. Docusaurus docs site is kept as-is but consumes shared generated data. No user authentication.

## 2. Tech Stack

**IMPORTANT: Before starting Phase 0, verify all dependency versions against context7 and npm to confirm latest stable releases. The versions below reflect intent, not pinned versions.**

| Concern | Choice | Notes |
|---------|--------|-------|
| Language | TypeScript (strict) | |
| Framework | React 19 | |
| Build | Vite 6 | |
| Monorepo | Turborepo | |
| Package Manager | pnpm | |
| CSS | **Tailwind CSS v5** | Uses CSS-first `@theme` configuration, NOT JS config files |
| Primitives | shadcn/ui (Radix + CVA + Tailwind, copied into repo) | Verify shadcn compatibility with Tailwind v5 |
| Icons | Lucide React | |
| Charts | Recharts | For charts that need finer control than Recharts allows, fall back to raw SVG |
| Code Editor | CodeMirror 6 | Spike the gcsim language mode early (see Phase 3 notes) |
| Server State | TanStack Query | |
| Client State | Zustand (with localStorage middleware) | |
| Routing | TanStack Router | |
| i18n | i18next + react-i18next | |
| Testing | Vitest + Testing Library | |
| E2E | Playwright | |
| Component Dev | Storybook 8 | |
| Linting/Format | Biome | |
| Protobuf | buf / protobuf-es | Migration from `protobufjs-cli` — verify generated output matches existing types |
| API Client | Native fetch (via TanStack Query) | |
| Compression | pako (zlib) | Existing API may return compressed data — must handle decompression |

## 3. Folder Structure

```
ui-next/
├── apps/
│   ├── web/                    # Main simulator + viewer (gcsim.app) — port 5173
│   │   ├── src/
│   │   │   ├── pages/
│   │   │   │   ├── simulator/
│   │   │   │   ├── viewer/
│   │   │   │   └── sample/
│   │   │   ├── stores/         # Zustand stores (simulator, viewer, settings)
│   │   │   ├── routes.ts
│   │   │   ├── app.tsx
│   │   │   ├── app.css         # Tailwind v5 @theme + @import directives
│   │   │   └── main.tsx
│   │   ├── .env                # VITE_API_BASE_URL, VITE_WASM_BASE_URL
│   │   ├── index.html
│   │   ├── vite.config.ts
│   │   └── package.json
│   │
│   ├── db/                     # Database browser (db.gcsim.app) — port 5174
│   │   ├── src/
│   │   │   ├── pages/
│   │   │   │   ├── home/
│   │   │   │   └── database/
│   │   │   ├── stores/
│   │   │   ├── routes.ts
│   │   │   ├── app.tsx
│   │   │   ├── app.css
│   │   │   └── main.tsx
│   │   ├── .env
│   │   ├── index.html
│   │   ├── vite.config.ts
│   │   └── package.json
│   │
│   ├── embed/                  # Screenshot generator — port 5175
│   │   ├── src/
│   │   ├── vite.config.ts
│   │   └── package.json
│   │
│   ├── storybook/              # Storybook 8 for component development
│   │   ├── .storybook/
│   │   ├── vite.config.ts
│   │   └── package.json
│   │
│   ├── taghelper/              # Moderator helper (low priority)
│   │   └── ...
│   │
│   └── docs/                   # Docusaurus (kept as-is, stays in ui/ until cutover)
│       └── ...
│
├── packages/
│   ├── primitives/             # shadcn/ui primitives + design tokens
│   │   ├── src/
│   │   │   ├── theme.css       # Tailwind v5 @theme tokens (colors, spacing, etc.)
│   │   │   └── ...
│   │   └── package.json
│   ├── viewer/                 # Shared viewer feature components
│   ├── editor/                 # CodeMirror wrapper + gcsim language
│   ├── avatar/                 # Character/team display components
│   ├── preview/                # Preview/embed card
│   ├── executor/               # WASM + Server executor abstraction
│   ├── types/                  # Protobuf-generated types (sole source of truth)
│   ├── i18n/                   # Localization (7 languages)
│   ├── data/                   # Game data (characters, weapons, artifacts, tags)
│   └── api/                    # API client layer (includes pako for decompression)
│
├── assets/                     # Shared static assets (character portraits, element
│   ├── characters/             # icons, weapon images, artifact images)
│   ├── elements/               # Referenced by apps via import or public/ symlink
│   ├── weapons/
│   └── artifacts/
│
├── tooling/
│   ├── typescript/             # Shared tsconfig bases
│   ├── vitest/                 # Shared Vitest config
│   └── test-fixtures/          # Canonical mock data (SimResult, etc.)
│
├── turbo.json
├── pnpm-workspace.yaml
├── package.json
├── biome.json
└── .env.example                # Template for app-level .env files
```

## 4. Architecture Rules

1. **Protobuf types are the single source of truth.** No shared type aliases across packages. Apps may define local convenience types.
2. **Import from package index only.** Never import from internal paths like `@gcsim/viewer/src/charts/...`. Enforced by `dependency-cruiser` in CI (see Phase 0).
3. **Components live in apps until a second consumer appears.** Only then extract to a shared feature package.
4. **TanStack Query for all server data.** No raw fetch + useEffect patterns.
5. **Zustand for client state.** With localStorage middleware for persistence.
6. **All apps are CSR.** No SSR/SSG.
7. **No authentication.** Discord OAuth is removed.
8. **Error boundaries required** in all apps — at minimum around the viewer (large data parsing/rendering), chart components, and any component that fetches external data.
9. **Route-level code splitting** — all page components loaded via lazy imports. Required from Phase 4 onward.
10. **Tailwind v5 CSS-first theming** — design tokens defined via `@theme` in CSS, NOT JavaScript config files.

## 4.1. Environment Configuration

Each app has a `.env` file (not committed) and a `.env.example` (committed) defining:

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_BASE_URL` | Backend API base URL | `""` (same origin) |
| `VITE_WASM_BASE_URL` | WASM binary location | `/api/wasm` |
| `VITE_LOCAL_DEV_URL` | Local dev server for LocalViewer | `http://127.0.0.1:8381` |

Apps access via `import.meta.env.VITE_*`.

## 4.2. Static Asset Management

Game assets (character portraits, element icons, weapon/artifact images) live in `ui-next/assets/` as a shared directory. Apps reference them via:
- Import in components: `import portrait from '@gcsim/assets/characters/...'`
- Or symlink `assets/` into each app's `public/` directory at build time

The existing assets from `ui/` are copied to `ui-next/assets/` in Phase 0.

## 4.3. WASM Binary Strategy

The Go simulation engine compiles to WebAssembly via `task wasm` (existing build step). For the rewrite:
- **Development:** Pre-built WASM binary placed in `ui-next/assets/wasm/` (or fetched from `VITE_WASM_BASE_URL`)
- **Production:** WASM binary served from CDN/API at `/api/wasm/{branch}/{commit}/main.wasm` (existing pattern)
- **Cache busting:** Binary URL includes commit hash (existing pattern preserved)
- **Build step:** `task wasm` output directory configured to write to `ui-next/assets/wasm/` for local dev

## 4.4. Legacy Package Disposition

| Old Package | Disposition |
|-------------|-------------|
| `@gcsim/ui` | Merged into `apps/web` + feature packages |
| `@gcsim/components` | Split into `packages/avatar`, `packages/preview`, `packages/viewer` |
| `@gcsim/executors` | Ported to `packages/executor` |
| `@gcsim/utils` | `useLocalStorage` → Zustand middleware (eliminated). Any remaining utils absorbed into the package that uses them. |
| `@gcsim/workers` | **Clarify:** If Cloudflare Workers deployment is still needed, port as `apps/workers/`. If not, drop. |
| `@gcsim/types` | Ported to `packages/types` with buf/protobuf-es (replacing protobufjs-cli) |
| `@gcsim/localization` | Ported to `packages/i18n` |
| `@gcsim/data` | Ported to `packages/data` |

## 4.5. Deployment & Cutover

1. Build `ui-next/` apps alongside `ui/` during development
2. When all apps pass E2E tests and are feature-complete:
   - Rename `ui/` → `ui-legacy/` (preserved temporarily)
   - Rename `ui-next/` → `ui/`
   - Update CI/CD build paths
   - Deploy new apps to existing domains
3. After 2 weeks with no issues, delete `ui-legacy/`

Production deployment target (verify with existing infra): static sites on Cloudflare Pages or equivalent. Each app builds to a `dist/` directory served independently.

## 5. Design System

Formalized before any component work begins. Defined using Tailwind v5's CSS-first `@theme` directive in `packages/primitives/src/theme.css`. Apps import this CSS file to inherit all tokens.

```css
/* packages/primitives/src/theme.css */
@theme {
  /* Colors — dark theme as default */
  --color-primary: ...;
  --color-secondary: ...;
  --color-destructive: ...;
  --color-muted: ...;
  --color-accent: ...;
  --color-background: ...;
  --color-foreground: ...;
  --color-card: ...;
  --color-popover: ...;
  --color-border: ...;
  --color-input: ...;
  --color-ring: ...;

  /* Typography */
  --font-sans: ...;
  --font-mono: ...;

  /* Border radii */
  --radius-sm: ...;
  --radius-md: ...;
  --radius-lg: ...;

  /* Shadows, z-index, animations, etc. */
}
```

Tokens to define:
- Color palette (primary, secondary, destructive, muted, accent, background, foreground, card, popover, border, input, ring) with dark theme as default
- Spacing scale (consistent increments)
- Typography: font families (sans, mono), size scale, weight scale, line height scale
- Border radii: sm, md, lg, xl, full
- Shadows: sm, md, lg
- Z-index scale: dropdown, sticky, modal, popover, tooltip
- Breakpoints: sm, md, lg, xl, 2xl
- Animation durations and easings

Each app imports the theme: `@import '@gcsim/primitives/theme.css';` in its `app.css`.

## 6. Apps

### 6.1 Web App (gcsim.app)

**Routes:**
| Route | Page | Description |
|-------|------|-------------|
| `/` | Dash | Home/landing page |
| `/simulator` | Simulator | Config editor + team builder + action list + run controls |
| `/web` | WebViewer | Results from local state (after running simulation) |
| `/local` | LocalViewer | Results from local dev server (http://127.0.0.1:8381/data) |
| `/sh/:id` | ShareViewer | Results from `/api/share/:id` |
| `/sample/upload` | UploadSample | Sample analysis upload |
| `/sample/local` | LocalSample | Local sample analysis |

**Legacy redirects:** `/v3/viewer/share/:id`, `/viewer/share/:id`, `/s/:id`, `/viewer/web`, `/viewer/local`, `/simple`, `/advanced`, `/viewer` → appropriate new routes.

**Simulator page features:**
- Config text editor (CodeMirror 6 with gcsim syntax)
- Team builder (up to 4 characters)
- Action list editor
- Real-time validation (200ms debounce) against executor
- Mode switcher (WASM vs Server)
- WASM worker count config (1-30, default 3)
- Server URL config (default http://127.0.0.1:54321)
- Run button with worker readiness indicator

**Viewer page features (shared across WebViewer, LocalViewer, ShareViewer):**
- **Results tab:** team header, metadata (iterations, mode, commit), rollup stats, target info, distribution chart, damage timeline, cumulative damage, per-character DPS, element/target/source DPS breakdowns, character actions, field time, energy tracking, reactions, aura uptime
- **Config tab:** view/edit config from loaded results, re-run capability
- **Sample tab:** seed selection (sample/min/max/percentiles/custom), per-action event log with filtering and search

**Zustand stores:**
- `simulatorStore` — config text, team composition, validation state, execution mode, worker settings
- `viewerStore` — loaded results data, active tab, error state, recovery config
- `settingsStore` — user preferences (persisted to localStorage)

### 6.2 DB App (db.gcsim.app)

**Routes:**
| Route | Page | Description |
|-------|------|-------------|
| `/` | Home | Landing page with tag descriptions |
| `/database` | Database | Searchable list with filters |

**Features:**
- Character quick-select filter
- Advanced filter panel (characters, tags, date range)
- Query builder → MongoDB-style filter JSON
- Infinite scroll pagination (25/page) with abort controller
- Entry cards with portraits, metadata, action buttons
- Moderator actions (approve/reject)
- Links to viewer on main site

### 6.3 Embed App

**Routes:**
| Route | Page | Description |
|-------|------|-------------|
| `/db/:id` | DBPreview | Preview card for DB entry |
| `/sh/:id` | SharePreview | Preview card for shared result |

**Features:**
- Fetch result from API → render PreviewCard
- `#images_loaded` signal for headless capture (Puppeteer)
- Error boundary with JSON output

### 6.4 Tag Helper (low priority)

Moderator tool. Fetch entry → find similar → copy approve/reject/replace commands.

## 7. Shared Packages

### 7.1 `packages/primitives`
shadcn/ui components + Tailwind design token preset. Button, Card, Dialog, DropdownMenu, Input, ScrollArea, Select, Tabs, Tooltip, and others as needed.

### 7.2 `packages/viewer`
Result cards, Recharts-based charts (damage timeline, cumulative damage, DPS breakdowns, energy, distribution), team header, metadata display, sample viewer (seed selector, event log).

### 7.3 `packages/editor`
CodeMirror 6 wrapper. Custom gcsim language mode with syntax highlighting, autocomplete, error markers.

### 7.4 `packages/avatar`
Character portrait, avatar card with stats, team composition display.

### 7.5 `packages/preview`
Preview/embed card component used by embed app and DB app.

### 7.6 `packages/executor`
Executor interface, WasmExecutor (web worker pool management), ServerExecutor (HTTP), ExecutorSupplier (lazy init + stale detection).

### 7.7 `packages/types`
Protobuf-generated types via buf/protobuf-es. Re-exports only, no aliases.

### 7.8 `packages/i18n`
i18next setup, 7 language resource files (en, zh, ja, ko, es, ru, de), generated game entity name translations, resource merging.

### 7.9 `packages/data`
Generated game data: character metadata, weapon/artifact pipeline data, community tags.

### 7.10 `packages/api`
Typed fetch functions: `/api/share/:id`, `/api/share/db/:id`, `/api/db?q=...`, local dev server endpoints.

---

## 8. Testing Protocol

Every step in the implementation plan follows this protocol. Agents do not move to the next sub-step until the current one passes.

### Per-Sub-Step (agent runs continuously)
1. **Write failing test** for the component/function being built
2. **Implement** until the test passes
3. **Run package tests**: `turbo run test --filter=@gcsim/<package>`
4. **Run package typecheck**: `turbo run typecheck --filter=@gcsim/<package>`
5. **Run biome**: `biome check --write` on changed files

### Per-Step Completion (agent finishes a step)
1. All sub-step tests pass
2. Write/update CLAUDE.md for the package (while patterns are fresh)
3. `turbo run build --filter=@gcsim/<package>` succeeds

### Per-Phase Gate (before starting next phase)
1. **Full test suite**: `turbo run test` (all packages)
2. **Full typecheck**: `turbo run typecheck` (catches cross-package type breaks)
3. **Full lint**: `biome check`
4. **Dependency validation**: `dependency-cruiser` — no deep imports, no circular deps
5. **Code review**: dispatch `package-reviewer` subagent per changed package (includes CLAUDE.md completeness check)
6. All gates pass before any agent starts the next phase

### Test Expectations by Package Type
| Package type | What to test | What NOT to test |
|---|---|---|
| `types` | Exports exist, generated types compile | Generated code internals |
| `data` | Exports exist, data shape matches types | Individual data values |
| `i18n` | Language loading, key resolution, fallback | Translation string content |
| `api` | Request construction, error handling, abort | Actual HTTP calls (mock fetch) |
| `executor` | Interface contract, state transitions | WASM internals (mock workers) |
| `primitives` | Renders, variants, accessibility, keyboard | Visual appearance (Storybook covers this) |
| Feature packages | Component behavior with mocked data, user interactions | Internal implementation details |
| Apps (unit) | Store logic, data transforms, hook behavior | Component rendering (covered by feature package tests) |
| Apps (E2E) | Full user flows end-to-end | Every edge case (unit tests cover these) |

---

## 9. Implementation Plan

All steps are designed to be small, with explicit parallelism where possible. Each step must pass the testing protocol before moving to the next.

### Phase 0: Monorepo Scaffolding + Agent Tooling + CI

**Step 0.0 — Dependency version verification**
- Before any installation, verify latest stable versions of ALL dependencies via context7 and npm:
  - Tailwind CSS (v5), React, Vite, Turborepo, shadcn/ui, Radix UI, Recharts, CodeMirror 6, TanStack Query/Router, Zustand, i18next, Vitest, Playwright, Storybook, Biome, buf/protobuf-es, pako, Lucide React, dependency-cruiser
- Document pinned versions in a `DEPENDENCIES.md` for agents to reference
- Verify shadcn/ui compatibility with Tailwind v5 (CSS-first `@theme` approach)

**Step 0.1 — Initialize monorepo root**
- Create `ui-next/` directory
- `pnpm init`, create `pnpm-workspace.yaml` with `apps/*` and `packages/*`
- Install Turborepo, create `turbo.json` with `build`, `test`, `typecheck`, `lint` pipelines
- Install Biome, create `biome.json` with formatting + linting rules
- Install Vitest as root dev dependency
- Install `dependency-cruiser` for import boundary validation
- Create `tooling/typescript/base.json` — shared strict tsconfig
- Create `tooling/vitest/base.ts` — shared Vitest config
- Create `.env.example` template
- Copy static assets from `ui/` to `ui-next/assets/` (character portraits, element icons, weapon/artifact images)
- Configure WASM binary for local dev: `task wasm` output → `ui-next/assets/wasm/`
- Verify: `pnpm install` succeeds, `turbo run build` succeeds (empty)

**Step 0.2 — Root CLAUDE.md + auxiliary docs**
- Write `ui-next/CLAUDE.md` with:
  - Monorepo structure overview (apps/ vs packages/ vs assets/ vs tooling/)
  - Turborepo commands (`turbo run build`, `turbo run test --filter=...`)
  - Package boundary rules: import from package index only, never deep paths
  - Protobuf types are the single source of truth — no shared type aliases
  - Naming conventions (kebab-case files, PascalCase components, camelCase functions)
  - pnpm commands for adding dependencies (`pnpm add -D --filter=@gcsim/<pkg>`)
  - Biome usage
  - Tailwind v5: use CSS-first `@theme`, NOT JS config files
  - Testing protocol reference (see Section 8)
  - Error boundary requirement for all data-fetching/rendering components
  - Route-level code splitting required (lazy imports)
  - Environment variables: use `import.meta.env.VITE_*`, see `.env.example`
  - Asset references: import from `@gcsim/assets/...` or `ui-next/assets/`
  - WASM binary: local dev at `assets/wasm/`, prod at `VITE_WASM_BASE_URL`
- Write `ui-next/apps/CLAUDE.md`:
  - Apps own pages, stores, and route definitions
  - Apps wire packages together — they are composition layers
  - Apps use error boundaries around data-fetching and rendering sections
  - Each app has its own `.env` (see `.env.example`)
  - Dev server ports: web=5173, db=5174, embed=5175 (avoid conflicts with old `ui/`)
- Write `ui-next/tooling/CLAUDE.md`:
  - How to extend base tsconfig
  - How to extend base vitest config
  - How to use test fixtures from `tooling/test-fixtures/`

**Step 0.3 — Basic CI pipeline** *(Must exist before Phase 1)*
- GitHub Actions workflow (`.github/workflows/ui-next.yml`):
  1. `biome check`
  2. `tsc --noEmit` (via `turbo run typecheck`)
  3. `turbo run test`
  4. `dependency-cruiser --validate` (no deep imports, no circular deps)
  5. `turbo run build`
- Turborepo remote caching configuration
- Triggered on: push to `web-rewrite` branch, PRs targeting `web-rewrite`
- Playwright E2E added later in Phase 7
- Verify: CI runs on an empty push and passes

**Step 0.4 — Scaffolding skills** *(Sequential — these are used by all subsequent phases)*

**Step 0.4a — `/new-package` skill**
- Input: package name, optional dependencies, optional `--with-tailwind` flag
- Creates:
  - `packages/<name>/package.json` with `@gcsim/` scope, `exports` field (ESM), `scripts` (build, test, typecheck, lint), `main`/`types` fields
  - `tsconfig.json` extending `tooling/typescript/base.json` with `references` for declared dependencies
  - `src/index.ts`
  - `vitest.config.ts` extending `tooling/vitest/base.ts`
  - `CLAUDE.md` with 6-section template stubs (Purpose, How to add X, Canonical example, Public API, Dependencies, Don'ts)
  - If `--with-tailwind`: `src/styles.css` importing primitives theme
- Runs `pnpm install`
- Verify: skill creates a package that builds and tests (empty test suite passes)

**Step 0.4b — `/new-component` skill**
- Input: package name, component name, optional subdirectory (e.g., `charts`, `result-cards`)
- Reads the target package's CLAUDE.md for canonical example path
- Creates:
  - `packages/<pkg>/src/[subdir/]<component>/<component>.tsx` — skeleton following canonical example pattern
  - `<component>.test.tsx` — imports component, imports `render` from Testing Library, includes "renders without crashing" test
  - `index.ts` — named export
- Updates `packages/<pkg>/src/index.ts` to re-export
- Verify: skill creates files in correct locations, package still builds

**Step 0.4c — `/new-page` skill**
- Input: app name, page name
- Creates:
  - `apps/<app>/src/pages/<page>/<page>.tsx` (skeleton)
  - `<page>.test.tsx`
- Adds **lazy import** route entry to `apps/<app>/src/routes.ts` (code splitting)
- Verify: skill creates files, app still builds

**Step 0.4d — `/new-app` skill**
- Input: app name, dev server port
- Creates:
  - `apps/<name>/package.json`, `vite.config.ts` (with port), `tsconfig.json`
  - `src/main.tsx` — React root, TanStack Query provider, i18n init, error boundary
  - `src/app.tsx` — TanStack Router setup with layout shell
  - `src/app.css` — `@import '@gcsim/primitives/theme.css'` + Tailwind directives
  - `src/routes.ts` — empty route definitions
  - `.env` from `.env.example` template
  - `CLAUDE.md` with app-specific template (routes, stores, consumed packages)
- Verify: app starts on specified port, renders empty shell

**Step 0.4e — `/new-store` skill**
- Input: app name, store name, optional `--persist` flag for localStorage middleware
- Creates:
  - `apps/<app>/src/stores/<store>.ts` — typed Zustand store skeleton with actions
  - If `--persist`: wires up `persist` middleware with localStorage
- Verify: store file compiles

**Step 0.4f — `/check` skill**
- Input: optional `--filter=<package>` for per-package checks
- Runs sequentially: `biome check` → `turbo run typecheck [--filter]` → `turbo run test [--filter]` → `dependency-cruiser --validate` → `turbo run build [--filter]`
- Reports first failure and stops
- Verify: runs against empty monorepo, all pass

**Step 0.4g — `/dev` skill**
- Input: app name (or "all")
- Starts the correct Vite dev server(s) for the specified app
- If dependencies need building first, runs `turbo run build` on dependency packages
- Verify: starts web app dev server on correct port

**Step 0.5 — Subagent definitions**

**Step 0.5a — `package-reviewer` subagent**
- Reviews a single package for:
  - Boundary violations (no deep imports)
  - No type aliases outside `types/`
  - TanStack Query for server data (no raw fetch + useEffect)
  - Test quality: tests describe behavior (not implementation), test the public API, follow the testing protocol table
  - CLAUDE.md completeness: all 6 sections present and non-trivial
  - Design token usage: components use theme tokens, not hardcoded values
  - Error boundaries present where required
- Input: package path
- Output: list of issues or "approved"
- **Note:** `/check` is a pass/fail gate. `package-reviewer` is a qualitative review. Agents use `/check` first, then `package-reviewer` for deeper analysis.

**Step 0.5b — `package-tester` subagent**
- Runs `turbo run test --filter=@gcsim/<pkg>` and `turbo run typecheck --filter=@gcsim/<pkg>`
- On failure: reads failing test, identifies likely cause, suggests fix
- Input: package name or "all"
- Output: pass/fail with details

**Step 0.5c — `feature-implementer` subagent** *(skeleton only — refined after Phase 2)*
- Given: feature spec, target package, canonical example to follow
- Follows TDD: reads canonical example → writes failing test → implements → verifies
- Writes/updates package CLAUDE.md on completion
- Input: spec text, package path, canonical example path
- Output: list of files created/modified
- **Maximum scope:** one component with up to 3 sub-components and their tests. Larger work should be split into multiple dispatches.
- **This subagent's instructions are refined iteratively:**
  - After Phase 2 Step 2.2 (Button): refine with primitive-building workflow
  - After Phase 3 Steps 3.3d/3.4a (DPS card, damage timeline): refine with feature component workflow

**Step 0.5d — `cross-package-integrator` subagent**
- Wires completed packages into a consuming app
- Adds imports, creates composition components
- **Writes integration tests** verifying the wiring works (components render, data flows correctly)
- Requires a **composition spec** in the dispatch call (which components go where, how data flows)
- Input: list of packages to integrate, target app, composition spec
- Output: list of files created/modified

**Step 0.6 — Claude Code hooks**
- Create `.claude/hooks.json`:
  - `pre-commit`: `biome check --write` on staged files
  - `post-edit`: `turbo run test --filter=@gcsim/{package-of-edited-file}` (map file path to owning package via a small utility script)

**Step 0.7 — Git workflow protocol**
- Each parallel agent works on a branch: `phase-{N}/step-{X.Y}-{package-name}`
- At phase gates, branches are merged sequentially (foundation packages first, since feature packages depend on them)
- Worktree naming: `worktree-phase{N}-{package-name}`
- Add cleanup script for stale worktrees

**Deliverable:** Monorepo skeleton, CI pipeline, all skills, all subagent definitions, hooks, git workflow. Every subsequent phase uses these tools.

---

### Phase 1: Foundation Packages (no UI yet)

These packages have no React dependency (except i18n). Steps 1.1 must complete first (types is a dependency for data, api, executor). Steps 1.2-1.5 can then run in parallel.

**Step 1.1 — `packages/types`** *(Agent A — must complete first)*
- Run `/new-package types`
- Set up buf/protobuf-es generation pipeline from existing `.proto` files
- **Migration note:** Existing codebase uses `protobufjs-cli`. Verify generated protobuf-es output matches existing type shapes. Run a diff of exported type names to catch regressions.
- Generate TypeScript types
- Create `src/index.ts` with re-exports
- Test: generated types compile, exports are accessible, key types (SimResult, Character, etc.) exist
- Write CLAUDE.md: protobuf generation commands, re-export-only rule, how to regenerate after proto changes

**Step 1.1b — Shared test fixtures** *(Same agent, immediately after 1.1)*
- Create `tooling/test-fixtures/sim-result.ts` — canonical mock SimResult object with realistic data (multiple characters, damage data, energy data, sample data)
- Create `tooling/test-fixtures/characters.ts` — mock character data for 4-5 common characters
- Create `tooling/test-fixtures/index.ts` — re-exports all fixtures
- These fixtures are used by ALL subsequent packages for testing. Prevents fixture drift.
- Write `tooling/test-fixtures/CLAUDE.md`: how to add a new fixture, where they're used

**Step 1.2 — `packages/data`** *(Agent B, after 1.1)*
- Run `/new-package data` with dep: `@gcsim/types`
- Port existing generated game data (characters, weapons, artifacts, tags) from `ui/packages/`
- Create `src/index.ts` exporting all data
- Test: data shape matches protobuf types from `@gcsim/types`, exports are accessible
- Write CLAUDE.md: generated data pipeline, how data flows from Go backend → JSON → TS exports

**Step 1.3 — `packages/i18n`** *(Agent C)*
- Run `/new-package i18n` with deps: `i18next`, `react-i18next`
- Port 7 language resource files from existing `ui/packages/localization/`
- Port generated game entity name translations
- Create `src/index.ts` with i18next initialization + resource merging
- Tests (TDD):
  1. Language loading — default en loads
  2. Key resolution — known key returns correct translation
  3. Fallback — missing key in zh falls back to en
  4. Game names — character name resolves in each language
- Write CLAUDE.md: how to add a translation key, namespace conventions, generated vs manual keys

**Step 1.4 — `packages/api`** *(Agent D)*
- Run `/new-package api` with dep: `@gcsim/types`
- Implement `src/client.ts` — base fetch wrapper with error handling
- Implement `src/share.ts` — `fetchShareResult(id)`, `fetchDBResult(id)`
- Implement `src/db.ts` — `queryDB(query)` with abort signal support
- Tests (TDD, mock fetch):
  1. `fetchShareResult` constructs correct URL, returns typed result
  2. `fetchDBResult` constructs correct URL
  3. `queryDB` sends query as param, supports abort signal
  4. Error handling: 404 → typed error, network failure → typed error
  5. Abort: cancellation triggers AbortError
- Write CLAUDE.md: how to add a new endpoint, canonical example: `share.ts`, error handling pattern

**Step 1.5 — `packages/executor`** *(Agent E)*
- Run `/new-package executor` with dep: `@gcsim/types`
- **Step 1.5a** — Port `Executor` interface + `ServerExecutor` first (simpler, no browser APIs)
  - Port `Executor` interface from existing `ui/packages/executors/`
  - Port `ServerExecutor` — HTTP-based execution
  - Port `ExecutorSupplier` — lazy init, stale detection
  - Tests (TDD, mock fetch):
    1. `Executor` interface contract — all methods exist
    2. `ServerExecutor` — run sends HTTP request, cancel aborts, handles errors
    3. `ExecutorSupplier` — lazy init on first access, stale detection triggers reinit
- **Step 1.5b** — Port `WasmExecutor` (complex, browser API dependent)
  - Port `WasmExecutor` — web worker pool management, WASM binary loading
  - Tests (TDD, mock workers):
    1. Init creates correct number of workers (1-30 range)
    2. Run delegates to worker pool
    3. Cancel aborts in-flight execution
    4. **Worker creation failure** — graceful error handling
    5. **WASM binary load failure** — error surfaced to caller
    6. **Worker pool scaling** — changing worker count re-creates pool
    7. **Stale binary detection** — new commit triggers reload
- Write CLAUDE.md: executor interface contract, WASM loading strategy, how to reference old code in `ui/packages/executors/`

**Phase 1 gate:** Run `/check`. Dispatch `package-reviewer` per package (5 parallel agents). All must pass before Phase 2.

**Deliverable:** All foundation packages build, test, and export correctly. No UI yet.

---

### Phase 2: Design System + Primitives

Must be done before any feature packages or apps.

**Step 2.1 — Design system tokens** *(Single agent — design decisions needed)*
- Run `/new-package primitives --with-tailwind`
- Define `src/theme.css` using Tailwind v5 `@theme` directive (see Section 5 for full token list):
  - Color palette: primary, secondary, destructive, muted, accent, background, foreground, card, popover, border, input, ring (dark theme as default)
  - Spacing scale (consistent increments)
  - Typography: font families (sans, mono), size scale, weight scale, line height scale
  - Border radii: sm, md, lg, xl, full
  - Shadows: sm, md, lg
  - Z-index scale: dropdown, sticky, modal, popover, tooltip
  - Breakpoints: sm, md, lg, xl, 2xl
  - Animation durations and easings
- Create a `src/cn.ts` utility (classname merge helper using `clsx` + `tailwind-merge`)
- Write CLAUDE.md with token reference and usage instructions
- Verify: theme CSS is importable by other packages, Tailwind compiles with it

**Step 2.2 — Button (canonical example)** *(Single agent — do this first)*
- Run `/new-component primitives button`
- Implement: variants (default, destructive, outline, secondary, ghost, link), sizes (default, sm, lg, icon)
- Tests (TDD):
  1. Renders each variant with correct classes
  2. Click handler fires
  3. Disabled state prevents click, applies disabled styles
  4. Renders as child element when `asChild` prop is used
- Update primitives CLAUDE.md: designate `src/button/` as **the canonical example** for all future primitives
- Verify: `turbo run test --filter=@gcsim/primitives` passes

**Step 2.3 — Core primitives (batch 1)** *(Parallel agents, after Button is done — each follows Button as canonical example)*
All use `/new-component primitives <name>` to scaffold.

**Step 2.3a — Card** *(Agent A)*
- Components: Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter
- Tests: renders composition, slots content correctly

**Step 2.3b — Input** *(Agent B)*
- Components: Input with variants, forwardRef
- Tests: renders, onChange fires, disabled state, placeholder

**Step 2.3c — Tabs** *(Agent C)*
- Components: Tabs, TabsList, TabsTrigger, TabsContent (Radix Tabs)
- Tests: renders tabs, switching works, correct content shown

**Step 2.3d — Select** *(Agent D)*
- Components: Select, SelectTrigger, SelectValue, SelectContent, SelectItem (Radix Select)
- Tests: renders, selection changes value, disabled items

**Step 2.3e — Badge** *(Agent E)*
- Components: Badge with variants (default, secondary, destructive, outline)
- Tests: renders each variant with correct classes

**Step 2.4 — Core primitives (batch 2)** *(Parallel agents)*

**Step 2.4a — Dialog** *(Agent A)*
- Components: Dialog, DialogTrigger, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter
- Tests: opens/closes, renders content, accessible (role=dialog)

**Step 2.4b — DropdownMenu** *(Agent B)*
- Components: DropdownMenu, Trigger, Content, Item, Separator
- Tests: opens on trigger, items clickable, keyboard navigation

**Step 2.4c — Tooltip** *(Agent C)*
- Components: Tooltip, TooltipTrigger, TooltipContent, TooltipProvider
- Tests: shows on hover, hides on leave, accessible

**Step 2.4d — ScrollArea** *(Agent D)*
- Components: ScrollArea, ScrollBar (Radix ScrollArea)
- Tests: renders content, scrollbar appears when content overflows

**Step 2.4e — Skeleton** *(Agent E)*
- Components: Skeleton (loading placeholder)
- Tests: renders with correct dimensions, animates

**Step 2.5 — Storybook setup** *(Single agent)*
- Create Storybook 8 config in `apps/storybook/`
- Add stories for all primitives (2.2-2.4) — each primitive gets one story file with all variants
- Verify: `turbo run build --filter=@gcsim/storybook` succeeds, all primitives render

**Phase 2 gate:** Run `/check`. Dispatch `package-reviewer` on `primitives`. Must pass before Phase 3.

**Deliverable:** Full primitives package with design tokens, 11 shadcn primitives, Storybook catalog, all tested.

---

### Phase 3: Feature Packages

Each uses `/new-package` then `/new-component` for scaffolding. Each depends on `primitives` + `types` (from Phases 1-2).

**Parallelism note:** Steps 3.1 (avatar), 3.2 (editor), and 3.3 (viewer setup) run in parallel. Steps 3.4 and 3.5 (viewer charts and sample) run in parallel AFTER 3.3 completes (3.3 establishes the viewer package structure and barrel file). Step 3.6 (preview) can run anytime after 3.1.

**Step 3.1 — `packages/avatar`** *(Agent A)*
- Run `/new-package avatar` with deps: `@gcsim/primitives`, `@gcsim/types`, `@gcsim/data`
- **Step 3.1a** — Run `/new-component avatar portrait`
  - Character portrait/icon component: takes character key, renders image with element indicator
  - Tests (TDD): renders correct image, handles missing character gracefully, element indicator matches character data
- **Step 3.1b** — Run `/new-component avatar avatar-card`
  - Character card: name, element, weapon, level, constellation, stats
  - Tests (TDD): renders all fields, handles partial data, displays constellation count
- **Step 3.1c** — Run `/new-component avatar team-display`
  - Horizontal row of 1-4 character portraits
  - Tests (TDD): renders correct count, handles empty team, handles 1-4 characters
- Write CLAUDE.md: canonical example `portrait/`, how to add a new display variant
- Run `turbo run test --filter=@gcsim/avatar`

**Step 3.2 — `packages/editor`** *(Agent B)*
- Run `/new-package editor` with deps: `@gcsim/primitives`, `codemirror`, `@codemirror/lang-*`
- **Risk: CodeMirror 6 Lezer parser has a learning curve.** Start with Step 3.2a as a spike. If the language mode takes more than expected, simplify to basic keyword highlighting first and iterate.
- **Step 3.2a** — `src/language/gcsim-language.ts`
  - CodeMirror 6 language support: tokenizer for gcsim config syntax (keywords, strings, numbers, comments, character names, action names)
  - Tests (TDD): tokenizes sample configs — keywords highlighted, strings delimited, comments ignored, character names recognized
- **Step 3.2b** — `src/language/autocomplete.ts`
  - Completion source for character names, action keywords, stat names (using `@gcsim/data`)
  - Tests (TDD): provides completions for partial input, filters correctly, returns empty for unknown prefix
- **Step 3.2c** — Run `/new-component editor editor`
  - React wrapper: controlled value, onChange, readOnly mode, error gutter markers, dark theme, line numbers
  - Tests (TDD): renders editor, value changes emit onChange, readOnly prevents editing, error markers display at correct lines
- Write CLAUDE.md: canonical example `codemirror/editor.tsx`, how to add syntax rules, how to add completions
- Run `turbo run test --filter=@gcsim/editor`

**Step 3.3 — `packages/viewer` (metadata + result cards)** *(Agent C — must complete before 3.4/3.5)*
- Run `/new-package viewer --with-tailwind` with deps: `@gcsim/primitives`, `@gcsim/types`, `@gcsim/avatar`, `@gcsim/i18n`
- **This step establishes the viewer package structure, barrel file, and CLAUDE.md. Steps 3.4 and 3.5 build on this.**
- **Step 3.3a** — Run `/new-component viewer metadata`
  - Sub-components: Iterations, Mode, Commit, Warnings. Each takes relevant slice of SimResult.
  - Tests (TDD): each renders correct values, Warnings renders nothing when no warnings, handles missing optional fields
- **Step 3.3b** — Run `/new-component viewer team-header`
  - Row of avatar cards from SimResult character data
  - Tests (TDD): renders correct number of characters, stats match SimResult data
- **Step 3.3c** — Run `/new-component viewer rollup-card`
  - Generic stat rollup card (mean, min, max, std dev)
  - Tests (TDD): renders all stat fields, formats numbers, handles zero values
- **Step 3.3d** — Run `/new-component viewer dps-card`
  - Per-character DPS card with inline bar chart
  - Tests (TDD): renders character name, DPS value, bar width proportional to value
  - **This becomes the canonical example for result cards**
- **Step 3.3e** — Run `/new-component viewer target-info-card`
  - Enemy positions, auras
  - Tests (TDD): renders target list, aura indicators with correct elements
- Write CLAUDE.md: canonical example `result-cards/dps-card/`, how to add a new result card
- Run `turbo run test --filter=@gcsim/viewer`

**Step 3.4 — `packages/viewer` (charts)** *(Agent D, after 3.3 completes)*
- Same package as 3.3. Depends on: `recharts`. Uses test fixtures from `tooling/test-fixtures/`.
- **Recharts fallback:** If any chart cannot be cleanly implemented in Recharts (e.g., aura uptime, complex action timelines), fall back to raw SVG. Document the fallback pattern in CLAUDE.md.
- **Step 3.4a** — Run `/new-component viewer damage-timeline`
  - Recharts LineChart: damage over frame/time, per-character series, tooltips
  - Tests (TDD): renders with sample data, correct number of Line components, tooltip content
  - **This becomes the canonical example for charts**
- **Step 3.4b** — Run `/new-component viewer cumulative-damage`
  - Recharts AreaChart: cumulative total DPS
  - Tests (TDD): renders, values accumulate correctly
- **Step 3.4c** — Run `/new-component viewer distribution-chart`
  - Histogram of damage distribution
  - Tests (TDD): renders bars, correct bucket count from data
- **Step 3.4d** — Run `/new-component viewer element-dps-chart`
  - Stacked BarChart by element
  - Tests (TDD): renders elements with correct colors per element type
- **Step 3.4e** — Run `/new-component viewer energy-chart`
  - Energy generation over time per character
  - Tests (TDD): renders one series per character
- **Step 3.4f** — Run `/new-component viewer field-time-chart`
  - Pie or bar chart of active character field time
  - Tests (TDD): proportions sum to 100%, labels correct
- **Step 3.4g** — Run `/new-component viewer reactions-chart`
  - Elemental reaction counts
  - Tests (TDD): renders reaction types, counts match data
- Write CLAUDE.md: canonical example `charts/damage-timeline/`, how to add a new chart
- Run `turbo run test --filter=@gcsim/viewer`

**Step 3.5 — `packages/viewer` (sample viewer)** *(Agent E, after 3.3 completes)*
- Same package as 3.3. Uses test fixtures from `tooling/test-fixtures/`.
- **Step 3.5a** — Run `/new-component viewer seed-selector`
  - Dropdown: sample, min, max, p25, p50, p75, custom input
  - Tests (TDD): renders all options, selection emits value, custom input validates
- **Step 3.5b** — Run `/new-component viewer event-log`
  - Scrollable log of per-action events, filterable by action type, text search
  - Tests (TDD): renders events, filter by type reduces list, search matches text, empty state
- **Step 3.5c** — Run `/new-component viewer sample-viewer`
  - Composition: seed selector + event log, handles async sample loading
  - Tests (TDD): loading spinner while fetching, error message on failure, renders log when loaded
- Write CLAUDE.md
- Run `turbo run test --filter=@gcsim/viewer`

**Step 3.6 — `packages/preview`** *(Agent F)*
- Run `/new-package preview` with deps: `@gcsim/primitives`, `@gcsim/types`, `@gcsim/avatar`
- Run `/new-component preview preview-card`
  - Compact card: team portraits, DPS summary, metadata. Used for Discord embeds and DB entry cards.
  - Tests (TDD): renders all sections, handles missing optional data, image onLoad callback fires
- Write CLAUDE.md
- Run `turbo run test --filter=@gcsim/preview`

**Phase 3 gate:** Run `/check`. Dispatch `package-reviewer` per package (avatar, editor, viewer, preview — 4 parallel agents). All must pass before Phase 4.

**Deliverable:** All feature packages built, tested, documented. No apps yet.

---

### Phase 4: Web App (gcsim.app)

The largest app. Built incrementally page by page.

**Step 4.1 — App shell** *(Single agent)*
- Run `/new-app web 5173`
- `src/routes.ts` — all route definitions with **lazy imports** for code splitting (pages are stubs initially)
- Navigation component: links to Simulator, DB (external), Docs (external), Discord, GitHub releases
- Footer component
- Top-level error boundary wrapping the router
- Run `/new-store web simulator --persist`, `/new-store web viewer`, `/new-store web settings --persist`

**Step 4.1b — Store interface definitions** *(Same agent, before any simulator work)*
- Define **type-only interfaces** for all stores before parallel agents start:
  - `SimulatorStore`: config (string), team (Character[]), validationResult, executionMode, workerCount, serverUrl
  - `ViewerStore`: results (SimResult | null), activeTab, error, recoveryConfig
  - `SettingsStore`: user preferences
- These types are the contract that parallel agents in 4.3a-d code against
- Write CLAUDE.md (routes, stores with type definitions, consumed packages, error boundary locations)
- Verify: app starts on port 5173, navigation works, all routes render stubs

**Step 4.2 — Dash (home) page** *(Agent A)*
- `src/pages/dash/` — landing page content (welcome, quick links to simulator/docs/db)
- Verify: renders on `/`

**Step 4.3 — Simulator page** *(Sequential — complex, builds on itself)*

**Step 4.3a — Executor wiring** *(Agent B)*
- Wire `ExecutorSupplier` into app context
- `settingsStore` — execution mode (wasm/server), worker count, server URL
- Mode switcher UI component (dropdown: WASM / Server)
- WASM worker count slider
- Server URL input
- Verify: executor supplier initializes, mode switching works

**Step 4.3b — Config editor panel** *(Agent C)*
- `src/pages/simulator/config-editor.tsx` — CodeMirror editor from `@gcsim/editor`
- Wired to `simulatorStore.config`
- Debounced validation (200ms) calling `executor.validate()`
- Error display below editor (from validation result)
- Verify: typing updates store, validation runs, errors display

**Step 4.3c — Team builder panel** *(Agent D)*
- `src/pages/simulator/team-builder.tsx` — select up to 4 characters
- Character picker (searchable dropdown using `@gcsim/data` character list)
- Display selected team using `@gcsim/avatar` TeamDisplay
- Remove character button per slot
- Updates `simulatorStore.team`
- Verify: add/remove characters, store updates, team displays

**Step 4.3d — Action list editor** *(Agent E)*
- `src/pages/simulator/action-editor.tsx` — CodeMirror editor for action sequences
- Wired to `simulatorStore.actionList`
- Verify: editing updates store

**Step 4.3e — Run controls + simulator page composition** *(Single agent, after 4.3a-d)*
- `src/pages/simulator/simulator.tsx` — composes all panels (editor, team builder, action editor, run controls)
- Run button: calls `executor.run()`, disabled when not ready or already running
- Worker readiness indicator
- Progress display during execution
- On completion: populate `viewerStore` with results, navigate to `/web`
- Verify: full simulator flow works end-to-end (with mock executor)

**Step 4.4 — Viewer page** *(Sequential — tabs built incrementally)*

**Step 4.4a — Viewer shell + data loading** *(Agent A)*
- `src/pages/viewer/viewer.tsx` — tab container (Results / Config / Sample)
- `src/pages/viewer/web-viewer.tsx` — loads from `viewerStore` (local state)
- `src/pages/viewer/local-viewer.tsx` — loads from `http://127.0.0.1:8381/data` via TanStack Query
- `src/pages/viewer/share-viewer.tsx` — loads from `/api/share/:id` via TanStack Query
- All three feed into the same viewer component with loaded SimResult
- Loading/error states
- Verify: each route loads data from correct source, tab shell renders

**Step 4.4b — Results tab** *(Agent B)*
- `src/pages/viewer/results-tab.tsx` — composes all `@gcsim/viewer` components in order:
  - TeamHeader → Metadata → Rollup cards → Target info → Distribution chart → Damage timeline → Cumulative damage → Character DPS cards → Element DPS → Target DPS → Source DPS → Character actions → Field time → Energy → Reactions → Aura uptime
- Each component receives its relevant slice of SimResult
- Verify: renders all sections with mock data, scrollable

**Step 4.4c — Config tab** *(Agent C)*
- `src/pages/viewer/config-tab.tsx` — read-only CodeMirror editor showing loaded config
- "Edit" toggle to make editable
- Re-validation on edit
- "Re-run" button that runs via executor and updates viewerStore
- Verify: displays config, edit + re-run flow works

**Step 4.4d — Sample tab** *(Agent D)*
- `src/pages/viewer/sample-tab.tsx` — uses `@gcsim/viewer` SampleViewer
- Wired to executor's `sample()` method
- Seed in URL hash for shareability
- Verify: seed selection triggers sample, event log displays

**Step 4.5 — Sample upload/local pages** *(Agent E)*
- `src/pages/sample/upload.tsx` — file upload → parse → display in SampleViewer
- `src/pages/sample/local.tsx` — fetch from local server → display in SampleViewer
- Verify: upload and local flows work

**Step 4.6 — Legacy redirects** *(Agent F)*
- Implement all legacy route redirects in `routes.ts`
- Verify: each legacy URL redirects to correct new route

**Phase 4 gate:** Run `/check`. Dispatch `package-reviewer` on `apps/web`. E2E smoke test: app starts, nav works, simulator page loads. Must pass before Phase 5.

**Deliverable:** Fully functional web app with simulator, viewer, and all routes.

---

### Phase 5: DB App (db.gcsim.app)

**Step 5.1 — App shell** *(Agent A)*
- Run `/new-app db 5174`
- TanStack Router with two routes: `/` and `/database` (lazy imports)
- Navigation, layout, error boundary
- Write CLAUDE.md
- Verify: app starts on port 5174, both routes render stubs

**Step 5.2 — Home page** *(Agent B)*
- `src/pages/home/` — welcome content, tag descriptions, "Get Started" link to `/database`
- Verify: renders on `/`

**Step 5.3 — Filter system** *(Agent C)*
- `src/pages/database/filters/` — filter state management (Zustand store or local state with useReducer)
- Character quick-select dropdown (using `@gcsim/data` + `@gcsim/avatar` portraits)
- Advanced filter panel: character multi-select, tag multi-select, date range
- `craftQuery()` function: filter state → MongoDB-style JSON query string
- Verify: filter state updates, query string generates correctly (unit tests)

**Step 5.4 — Database list** *(Agent D, after 5.3. Note: Step 5.5 entry cards can run in parallel — it's a pure presentational component)*
- `src/pages/database/database.tsx` — composes filters + list
- `src/pages/database/list-view.tsx` — infinite scroll list
- TanStack Query with `useInfiniteQuery` for paginated DB fetching (25/page)
- Abort controller: new query cancels in-flight request
- Loading/error/empty states
- Verify: pagination works, query changes trigger refetch, abort works

**Step 5.5 — Entry cards** *(Agent E, parallel with 5.4 — pure presentational, takes props)*
- `src/pages/database/entry-card.tsx` — DB entry display card
- Team portraits (using `@gcsim/avatar`)
- Tags, metadata display
- Action buttons: "View" (link to main site viewer), moderator approve/reject
- Verify: renders all fields, links work

**Phase 5 gate:** Run `/check`. Dispatch `package-reviewer` on `apps/db`. Must pass before Phase 6.

**Deliverable:** Fully functional DB browser app.

---

### Phase 6: Embed App

**Step 6.1 — Embed app** *(Single agent)*
- Run `/new-app embed 5175`
- Two routes: `/db/:id` and `/sh/:id`
- Fetch result via `@gcsim/api`
- Render `@gcsim/preview` PreviewCard
- `#images_loaded` span on image load completion
- Error boundary with JSON output
- Write CLAUDE.md
- Verify: renders preview card, completion signal fires

**Phase 6 gate:** Run `/check`. Dispatch `package-reviewer` on `apps/embed`. Must pass before Phase 7.

**Deliverable:** Working embed/screenshot generator.

---

### Phase 7: E2E Tests + CI

**Step 7.1 — Playwright setup** *(Single agent. Note: 7.1 and 7.4 can run in parallel)*
- Install Playwright in monorepo root
- Configure for web app and DB app
- Install **MSW (Mock Service Worker)** for API mocking in E2E tests
- Create MSW handlers for all API endpoints (`/api/share/:id`, `/api/db`, etc.) using data from `tooling/test-fixtures/`
- Create Playwright fixture that starts MSW before tests

**Step 7.2 — Web app E2E tests** *(Agent A)*
- Simulator: open → edit config → validate → run → results display
- Viewer (share): navigate to `/sh/test-id` → results load → tabs switch
- Legacy redirects: verify all redirect correctly

**Step 7.3 — DB app E2E tests** *(Agent B)*
- Home → navigate to database
- Apply filters → results update
- Scroll → more results load
- Click entry → navigates to viewer

**Step 7.4 — CI pipeline update** *(Single agent, parallel with 7.1)*
- Add Playwright E2E step to existing CI workflow (from Phase 0):
  1. `biome check`
  2. `tsc --noEmit`
  3. `turbo run test`
  4. `dependency-cruiser --validate`
  5. `turbo run build`
  6. `playwright test` (on PRs targeting main only)
- Verify: full pipeline passes end-to-end

**Deliverable:** Full CI pipeline with E2E coverage.

---

### Phase 8: Cleanup + Polish

**Step 8.1 — Performance audit** *(Single agent)*
- Bundle size analysis per app (`vite build --report`)
- Verify route-level code splitting is working (lazy imports produce separate chunks)
- Lighthouse audit for each app
- Check for accidentally bundled large dependencies
- If issues found, fix before proceeding

**Step 8.2 — Tag helper app** *(low priority, single agent)*
- Port taghelper if still needed
- Run `/new-app taghelper 5176`

**Step 8.3 — Docs data integration** *(Single agent)*
- Docusaurus stays in `ui/packages/docs/` until cutover
- Update its imports to consume from `ui-next/packages/data/` (via symlink or path alias)
- Verify docs build still works

**Step 8.4 — `/extract-component` skill** *(Single agent)*
- Needed for ongoing development (not initial build)
- Input: source app, component path, target package (existing or new)
- Discovers all consumers via grep for import paths
- Moves component, updates imports, adds package export, creates tests
- If target package doesn't exist, runs `/new-package` first
- Adds `@gcsim/<target>` dependency to source app's package.json
- Verify: skill works on a test extraction

**Step 8.5 — Storybook catalog completion** *(Single agent)*
- Ensure all shared components have Storybook stories
- Visual regression baseline screenshots (Chromatic or similar)

**Step 8.6 — CLAUDE.md final audit** *(Parallel agents, one per package/app)*
- Review and finalize all CLAUDE.md files now that canonical examples are stable
- Ensure each has: purpose, how-to-add-X, canonical example, public API, dependencies, don'ts
- Root CLAUDE.md updated with final structure
- `apps/CLAUDE.md` and `tooling/CLAUDE.md` verified

**Step 8.7 — Cutover preparation** *(Single agent)*
- Verify all apps pass E2E tests
- Verify all feature parity with existing `ui/` (compare routes, features, API calls)
- Update CI/CD build paths for the rename
- Document cutover steps (see Section 4.5)
- Create tracking issue for the rename operation
