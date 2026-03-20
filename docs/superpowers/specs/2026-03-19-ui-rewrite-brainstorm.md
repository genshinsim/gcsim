# UI Rewrite - Brainstorming Document

## Overview

Complete rewrite of the gcsim web UI: a monorepo containing multiple independently deployable web apps that share components, types, data, and localization. Docusaurus docs site stays as-is but consumes shared generated data.

---

## Apps to Recreate

### 1. Simulator App (main site - gcsim.app)
- Config text editor with custom gcsim syntax highlighting
- Visual team builder (up to 4 characters, drag-and-drop)
- Action list editor
- Real-time config validation (debounced)
- Dual execution: WASM (web worker pool, 1-30 workers) or remote server
- Run controls with worker readiness indicator

### 2. Results Viewer (part of main site)
- Loaded from: local state, local dev server, shared link (`/sh/:id`), or DB link (`/db/:id`)
- **Results tab**: team header, metadata, rollup stats, target info, distribution chart, damage timeline, cumulative damage, per-character DPS, element/target/source DPS breakdowns, character actions, field time, energy tracking, reactions, aura uptime
- **Config tab**: view/edit config from loaded results, re-run capability
- **Sample tab**: seed selection (sample/min/max/percentiles/custom), per-action event log with filtering and search

### 3. Database Browser (db.gcsim.app)
- Landing page with tag descriptions
- Searchable simulation list with advanced filtering (characters, tags, date range)
- Infinite scroll pagination (25/page)
- Entry cards with portraits, metadata, action buttons
- Moderator actions (approve/reject)
- Query builder producing MongoDB-style filter JSON

### 4. Embed/Screenshot Generator
- Routes: `/db/:id` and `/sh/:id`
- Renders preview card, signals completion for headless capture (Puppeteer)
- Minimal app, error boundary with JSON output

### 5. Tag Helper (low priority / optional)
- Moderator tool for managing DB submissions
- Find similar submissions by team composition
- Copy approve/reject/replace commands

---

## Shared Infrastructure

### Types
- Protobuf-generated types are the single source of truth (simulation results, samples, models, enums)
- No shared type aliases - apps can define local convenience types if needed for their own design patterns

### Localization
- 7 languages: en, zh, ja, ko, es, ru, de
- UI translation strings + generated game entity names (characters, weapons, artifacts)

### Game Data
- Character metadata (element, weapon class, rarity, latest versions)
- Weapon/artifact pipeline data
- Community tags

### Executors
- `Executor` interface: ready(), running(), validate(), sample(), run(), cancel(), buildInfo()
- `WasmExecutor`: loads WASM binary, manages web worker pool
- `ServerExecutor`: HTTP-based remote execution
- `ExecutorSupplier`: lazy init with stale detection

### User Settings
- Persistent user preferences (local storage only, no auth)

---

## Proposed Tech Stack

### Core
| Concern | Choice | Rationale |
|---------|--------|-----------|
| Language | **TypeScript** (strict) | Type safety, existing protobuf types |
| Framework | **React 19** | Team familiarity, ecosystem maturity, WASM integration proven |
| Build Tool | **Vite 6** | Fast dev server, good monorepo support, already in use |
| Monorepo | **Turborepo** | Simpler than Nx, good caching, works well with pnpm |
| Package Manager | **pnpm** | Fast, disk-efficient, strict dependency resolution |

### Styling & UI
| Concern | Choice | Rationale |
|---------|--------|-----------|
| CSS | **Tailwind CSS v5** | Utility-first, already familiar, CSS-first `@theme` config |
| Component Primitives | **shadcn/ui** | Radix + Tailwind + CVA pre-wired; copied into repo, full control, no node_modules dependency |
| Icons | **Lucide React** | Tree-shakeable, consistent style |
| Charts | **Recharts** | Already proven in existing codebase, fits the data viz needs |
| Code Editor | **CodeMirror 6** | Modern, extensible, better than Ace for custom language modes |

### State & Data
| Concern | Choice | Rationale |
|---------|--------|-----------|
| Server State | **TanStack Query** | Caching, pagination, background refetch - perfect for DB browser & viewer data loading |
| Client State | **Zustand** | Lightweight, simpler than Redux for app state (simulator config, user prefs) |
| Routing | **TanStack Router** | Type-safe routing, good code splitting, single router across all apps |
| Forms | **React Hook Form + Zod** | Only if needed for filter/settings UIs |
| Persistence | **localStorage** via Zustand middleware | User preferences, mode selection |

### i18n
| Concern | Choice | Rationale |
|---------|--------|-----------|
| Framework | **i18next + react-i18next** | Already in use, mature, supports 7 languages, namespace splitting |

### Testing & Quality
| Concern | Choice | Rationale |
|---------|--------|-----------|
| Unit/Component | **Vitest + Testing Library** | Fast, Vite-native, good DX |
| E2E | **Playwright** | Cross-browser, good for testing viewer/simulator flows |
| Storybook | **Storybook 8** | Keep for component development, upgrade to latest |
| Linting | **Biome** | Fast, replaces ESLint + Prettier in one tool |

### Infrastructure
| Concern | Choice | Rationale |
|---------|--------|-----------|
| Types | **buf / protobuf-es** | Modern protobuf-to-TS generation, tree-shakeable |
| API Client | **ky** or **native fetch + TanStack Query** | Lighter than axios, modern API |
| WASM | Keep existing integration pattern | Web Workers + WASM binary loading |

---

## Proposed Folder Structure

```
ui/
├── apps/
│   ├── web/                    # Main simulator + viewer app (gcsim.app)
│   │   ├── src/
│   │   │   ├── pages/
│   │   │   │   ├── simulator/  # Config editor, team builder, action list
│   │   │   │   ├── viewer/     # Results viewer (results/config/sample tabs)
│   │   │   │   │   └── sample/     # Sample upload/local
│   │   │   ├── routes.ts       # Route definitions
│   │   │   ├── app.tsx
│   │   │   └── main.tsx
│   │   ├── index.html
│   │   ├── vite.config.ts
│   │   ├── tailwind.config.ts
│   │   └── package.json
│   │
│   ├── db/                     # Database browser (db.gcsim.app)
│   │   ├── src/
│   │   │   ├── pages/
│   │   │   │   ├── home/
│   │   │   │   └── database/   # List, filters, entry cards
│   │   │   ├── routes.ts
│   │   │   ├── app.tsx
│   │   │   └── main.tsx
│   │   ├── index.html
│   │   ├── vite.config.ts
│   │   ├── tailwind.config.ts
│   │   └── package.json
│   │
│   ├── embed/                  # Screenshot generator
│   │   ├── src/
│   │   ├── vite.config.ts
│   │   └── package.json
│   │
│   ├── taghelper/              # Moderator helper (low priority)
│   │   └── ...
│   │
│   └── docs/                   # Docusaurus (kept as-is, consumes shared data)
│       └── ...
│
├── packages/
│   ├── primitives/             # Low-level UI primitives (Radix + CVA + Tailwind)
│   │   ├── src/
│   │   │   ├── button/         # Button with variants (primary, destructive, ghost, etc.)
│   │   │   ├── dialog/
│   │   │   ├── dropdown-menu/
│   │   │   ├── select/
│   │   │   ├── tabs/
│   │   │   ├── tooltip/
│   │   │   ├── card/
│   │   │   ├── input/
│   │   │   ├── scroll-area/
│   │   │   └── ...             # Other Radix-based primitives as needed
│   │   ├── theme.css           # Tailwind v5 @theme tokens (colors, spacing, dark theme)
│   │   └── package.json
│   │
│   ├── viewer/                 # Shared viewer feature components
│   │   ├── src/
│   │   │   ├── result-cards/   # Rollup stats, DPS cards, distribution
│   │   │   ├── charts/         # Damage timeline, cumulative, energy, etc.
│   │   │   ├── team-header/    # Character portraits + stats display
│   │   │   ├── sample/         # Seed selector, event log, filtering
│   │   │   ├── metadata/       # Iterations, mode, commit, warnings
│   │   │   └── index.ts
│   │   └── package.json
│   │
│   ├── editor/                 # CodeMirror wrapper + gcsim language mode
│   │   ├── src/
│   │   │   ├── codemirror/     # CM6 setup, extensions, keybindings
│   │   │   ├── language/       # gcsim syntax highlighting, autocomplete
│   │   │   └── index.ts
│   │   └── package.json
│   │
│   ├── avatar/                 # Character/team display components
│   │   ├── src/
│   │   │   ├── avatar-card/    # Character card with stats
│   │   │   ├── portrait/       # Character portrait/icon
│   │   │   ├── team-display/   # Team composition display
│   │   │   └── index.ts
│   │   └── package.json
│   │
│   ├── preview/                # Preview/embed card (used by embed app + DB)
│   │   ├── src/
│   │   │   ├── preview-card/
│   │   │   └── index.ts
│   │   └── package.json
│   │
│   ├── executor/               # Simulation execution abstraction
│   │   ├── src/
│   │   │   ├── types.ts        # Executor interface
│   │   │   ├── wasm/           # WasmExecutor + worker management
│   │   │   ├── server/         # ServerExecutor (HTTP)
│   │   │   └── supplier.ts     # ExecutorSupplier (lazy init)
│   │   └── package.json
│   │
│   ├── types/                  # Protobuf-generated types (single source of truth)
│   │   ├── src/
│   │   │   ├── generated/      # buf/protobuf-es output
│   │   │   └── index.ts        # Re-exports only, no aliases
│   │   └── package.json
│   │
│   ├── i18n/                   # Localization
│   │   ├── src/
│   │   │   ├── locales/        # en/, zh/, ja/, ko/, es/, ru/, de/
│   │   │   ├── generated/      # Game entity name translations
│   │   │   └── index.ts        # i18next setup + resource merging
│   │   └── package.json
│   │
│   ├── data/                   # Game data (characters, weapons, artifacts, tags)
│   │   ├── src/
│   │   │   ├── generated/      # Pipeline output
│   │   │   └── index.ts
│   │   └── package.json
│   │
│   └── api/                    # API client layer
│       ├── src/
│       │   ├── share.ts        # /api/share endpoints
│       │   ├── db.ts           # /api/db endpoints
│       │   └── client.ts       # Base fetch/ky setup
│       └── package.json
│
├── tooling/                    # Shared configs
│   ├── typescript/             # Shared tsconfig bases
│   └── vitest/                 # Shared Vitest config
│
├── turbo.json
├── pnpm-workspace.yaml
├── package.json
└── biome.json
```

---

## Key Architecture Decisions

### 1. Feature packages over a monolithic UI library
Instead of a single shared `packages/ui/` that grows into a grab-bag, shared components are organized into **focused feature packages** (`primitives`, `viewer`, `editor`, `avatar`, `preview`). Each package has a clear domain boundary. Apps import only the feature packages they need, keeping bundles lean. Components stay in the app that owns them until a second consumer appears - only then extract to a shared package.

### 2. Separate API client package
Extract all HTTP calls into `packages/api/`. Apps import typed functions instead of scattering axios calls. TanStack Query hooks in each app wrap these functions for caching/loading states.

### 3. CodeMirror 6 over Ace Editor
Ace is largely unmaintained. CodeMirror 6 is the modern standard - modular, performant, and has excellent extension support for building a custom gcsim language mode (syntax highlighting, autocomplete, error markers).

### 4. Zustand over Redux
The current Redux setup is standard but heavyweight for the actual state needs. Zustand is simpler, less boilerplate, and handles the simulator/viewer state patterns well. Zustand's middleware supports localStorage persistence out of the box.

### 5. TanStack Query for all server data
DB queries, shared result loading, sample fetching - all benefit from TanStack Query's caching, background refetch, pagination, and abort handling. Replaces manual axios + useEffect patterns.

### 6. shadcn/ui over Blueprint
Blueprint is heavy and opinionated. shadcn/ui copies components into your repo (Radix + Tailwind + CVA pre-wired), giving full ownership and control. No node_modules styling to fight. The `packages/primitives/` package holds these components and the design tokens.

### 7. Single Tailwind `@theme` CSS in primitives
The `packages/primitives/` package exports a Tailwind `@theme` CSS defining design tokens (colors, spacing, dark theme). Each app extends this preset. This ensures visual consistency across apps without coupling them.

---

## Migration Strategy Notes

- **Incremental migration is not practical** - the current and new stacks are too different (Blueprint vs shadcn, Redux vs Zustand, React Router vs TanStack Router). A clean rewrite in parallel is more efficient.
- **Start with shared packages** (types, data, i18n, executor, api) since these are foundational.
- **Then build the component library** with shadcn/ui base components + gcsim-specific components.
- **Then build apps** starting with the web app (highest complexity), then db, then embed.
- **Docs stays as-is**, just point its data imports to the new `packages/data/`.

---

## Agent-Assisted Development Workflow

### Development

**Package-scoped agents:** The feature package boundaries are natural agent boundaries. When a feature spans multiple packages (e.g., adding a new chart type touches `viewer/`, `api/`, and the web app), each package can be worked on by a parallel subagent in its own worktree since they communicate through well-defined interfaces (TypeScript exports).

**TDD loop:** Vitest is fast enough for agents to run in watch mode. The cycle is: write test → implement → verify → move on. Each package has its own test suite runnable in isolation via `turbo run test --filter=@gcsim/viewer`.

**Turborepo-aware builds:** Agents should use `turbo run build` rather than manually building packages in dependency order. Turborepo handles the graph and caches results. A change to `primitives/` automatically triggers rebuilds of packages that depend on it.

**Typical agent dispatch pattern for a feature:**
```
Feature: "Add energy tracking chart to viewer"

Agent 1 (worktree): packages/api/     → Add energy data fetching endpoint
Agent 2 (worktree): packages/viewer/  → Build EnergyChart component + tests
Agent 3 (after 1+2): apps/web/        → Wire it into the viewer page
```

### Code Review

**Per-package scoped reviews:** A reviewer agent scoped to a single package can fully reason about it without drowning in monorepo context. For cross-package PRs, dispatch one reviewer per affected package.

**Automated checks (pre-review):**
- `biome check` — formatting + lint (pre-commit hook, issues never reach review)
- `tsc --noEmit` — type checking across the entire project
- `turbo run test` — all affected package tests pass

**Review focus areas:**
- Package boundary violations (importing internals instead of public API)
- No shared type aliases outside `packages/types/` (protobuf is source of truth)
- Feature components not prematurely extracted (should live in app until second consumer)
- TanStack Query usage for all server data (no raw fetch + useEffect)

### Testing Strategy

**Unit/Component tests (Vitest + Testing Library):**
- Every package has its own test suite
- Primitives: test variant rendering, accessibility, keyboard interaction
- Feature packages (viewer, editor, avatar): test component behavior with mocked data
- API package: test request construction, error handling
- Executor: test interface contract, worker lifecycle

**Integration tests (Vitest):**
- Test feature packages with real primitives (no mocking the component library)
- Test API + TanStack Query hooks together

**E2E tests (Playwright):**
- Per-app test suites: `apps/web/e2e/`, `apps/db/e2e/`
- Web app: simulator flow (edit config → validate → run → view results), viewer loading from share link
- DB app: search + filter → navigate to entry → view results
- Embed app: load preview card → verify completion signal

**Visual regression (Storybook + Chromatic or similar):**
- Primitives and feature components catalogued in Storybook
- Screenshot comparison catches unintended visual changes

### CI Pipeline

Run in this order (fail fast — cheapest checks first):

```
1. biome check          (seconds — formatting/lint)
2. tsc --noEmit         (seconds — type errors)
3. turbo run test       (fast — unit/component tests, Turborepo caches unchanged packages)
4. turbo run build      (builds only affected packages)
5. playwright test      (slow — E2E, only on PRs targeting main)
```

Turborepo remote caching means CI skips rebuilding/retesting packages with no changes. A PR touching only `packages/viewer/` won't re-test `packages/editor/`.

### Claude Code Hooks

```jsonc
// .claude/hooks.json
{
  "pre-commit": [
    "pnpm biome check --write",           // Auto-fix formatting
    "pnpm turbo run typecheck --affected"  // Type check changed packages
  ],
  "post-edit": [
    // Run tests for the package containing the edited file
    "pnpm turbo run test --filter=...{changed_package}"
  ]
}
```

### Development Conventions for Agents

- **Never bypass Turborepo** — always use `turbo run <task>` instead of running package scripts directly
- **Respect package boundaries** — import from package index (`@gcsim/viewer`), never from internal paths (`@gcsim/viewer/src/charts/...`)
- **Tests before implementation** — write failing test first, then implement
- **One package at a time** — complete work in a package (including tests) before moving to the next
- **Commit per package** — keeps git history clean and makes reverts surgical

---

## Skills & Subagents Setup

### Skills (slash commands)

| Skill | Purpose | When to use |
|-------|---------|-------------|
| `/new-package` | Scaffold a new package under `packages/` with tsconfig, package.json, Turborepo config, Tailwind `@theme` CSS extension, Vitest setup | Starting a new shared package |
| `/new-component` | Add a shadcn/ui primitive to `packages/primitives/` — runs shadcn CLI, adapts to project conventions | Need a new base UI primitive |
| `/new-page` | Scaffold a route/page in an app — creates page directory, route entry, basic test file | Adding a new route to any app |
| `/dev` | Start the correct dev server(s) for the current work context (Turborepo-aware, knows which apps depend on which packages) | Local development |
| `/check` | Run the full pre-PR pipeline locally: biome → tsc → test → build | Before creating a PR |
| `/extract-component` | Move a component from an app into a shared feature package — updates imports across consumers, adds package exports, creates tests | When a second app needs a component that currently lives in one app |

### Subagents

| Agent | Scope | Purpose |
|-------|-------|---------|
| **package-reviewer** | Single package | Code review scoped to one package. Checks: package boundary violations (no deep imports), no shared type aliases outside `types/`, correct TanStack Query usage, test coverage. Dispatched per-package on cross-package PRs. |
| **package-tester** | Affected packages | Runs `turbo run test --filter=...{changed}` for affected packages. Reports failures with file context and suggests fixes. |
| **feature-implementer** | Single package | Given a feature spec + target package, implements with TDD. Writes failing test → implements → verifies. Respects package boundaries and conventions. |
| **cross-package-integrator** | App + dependencies | After parallel agents complete work in individual packages, this agent wires them together in the consuming app. Adds imports, creates page-level composition, runs E2E tests to verify integration. |
| **spec-document-reviewer** | Design docs | Reviews spec/design documents for completeness, consistency, and feasibility before implementation begins. |

### Agent Dispatch Patterns

**New feature (multi-package):**
```
1. Plan → identify affected packages
2. Parallel: feature-implementer per package (in worktrees)
3. Sequential: cross-package-integrator wires it up in the app
4. package-reviewer per affected package
5. /check to verify everything
```

**New feature (single package):**
```
1. feature-implementer in the target package
2. package-reviewer
3. /check
```

**Bug fix:**
```
1. Identify the package → write a failing test
2. Fix → verify test passes
3. package-reviewer on the changed package
4. /check
```

**New shared component (extraction):**
```
1. /extract-component to move from app → package
2. package-tester on affected packages
3. package-reviewer on new package + affected app(s)
```

---

## Reducing Agent Exploration Overhead

### CLAUDE.md Hierarchy

Every level of the monorepo gets a CLAUDE.md that agents read automatically on start. This eliminates exploration for common tasks.

```
ui/
├── CLAUDE.md                        # Monorepo-wide rules
├── apps/
│   ├── web/CLAUDE.md                # Web app conventions
│   ├── db/CLAUDE.md                 # DB app conventions
│   └── embed/CLAUDE.md              # Embed app conventions
├── packages/
│   ├── primitives/CLAUDE.md         # How to add a primitive
│   ├── viewer/CLAUDE.md             # How to add a chart/card
│   ├── editor/CLAUDE.md             # CodeMirror patterns
│   ├── avatar/CLAUDE.md             # Asset/portrait conventions
│   ├── preview/CLAUDE.md            # Preview card patterns
│   ├── api/CLAUDE.md                # How to add an endpoint
│   ├── types/CLAUDE.md              # Protobuf generation, no aliases rule
│   ├── executor/CLAUDE.md           # Executor interface contract
│   ├── i18n/CLAUDE.md               # Translation key conventions
│   └── data/CLAUDE.md               # Generated data pipeline
```

**Root CLAUDE.md** (monorepo-wide):
- Monorepo structure overview (apps/ vs packages/)
- Turborepo commands (`turbo run build`, `turbo run test --filter=...`)
- Package boundary rules: import from package index only, never deep paths
- Protobuf types are the single source of truth — no shared type aliases
- Naming conventions (file names, component names, package names)
- pnpm commands for adding dependencies
- Biome usage

**Per-package CLAUDE.md** — every package CLAUDE.md must include:
1. **Purpose** — one sentence describing what this package does
2. **How to add X** — step-by-step for the most common operation (e.g., "How to add a new chart" in viewer/)
3. **Canonical example** — one well-implemented file as the reference pattern (e.g., "See `src/charts/damage-timeline/` as the reference implementation for adding a new chart")
4. **Public API** — what this package exports (so agents don't explore internals)
5. **Dependencies** — which other `@gcsim/*` packages this imports
6. **Don'ts** — common mistakes to avoid in this package

**Per-app CLAUDE.md** — same as per-package, plus:
- Route structure and how to add a new page
- State management patterns (which Zustand stores exist, what they hold)
- Which feature packages the app consumes

### Canonical Examples

The most powerful exploration-reduction tool. Each package designates ONE reference implementation that demonstrates the full pattern. When an agent needs to add something new, it reads the canonical example instead of exploring multiple files.

Examples:
- `packages/primitives/src/button/` — reference for adding a new shadcn primitive
- `packages/viewer/src/charts/damage-timeline/` — reference for adding a new chart
- `packages/viewer/src/result-cards/dps-card/` — reference for adding a new result card
- `packages/api/src/share.ts` — reference for adding a new API module
- `apps/web/src/pages/simulator/` — reference for a full page with state + features

### Scaffolding Skills Eliminate Boilerplate Exploration

The `/new-*` skills generate the full file structure so agents never need to figure out conventions:

`/new-component viewer energy-chart` should generate:
```
packages/viewer/src/charts/energy-chart/
├── energy-chart.tsx          # Component skeleton following the canonical pattern
├── energy-chart.test.tsx     # Test file with describe block and first test stub
├── index.ts                  # Named export
```
Plus update `packages/viewer/src/index.ts` to re-export it.

The agent's job becomes filling in the implementation, not figuring out where things go or what the file structure should look like.

### Summary: Agent Reads, Not Explores

| Without CLAUDE.md | With CLAUDE.md |
|---|---|
| Agent globs for files, reads 10-15 to understand patterns | Agent reads CLAUDE.md + 1 canonical example |
| Agent guesses at conventions, gets corrected in review | Agent follows documented conventions from the start |
| Agent explores package internals to understand API | Agent reads the public API section |
| Agent scaffolds files manually, may get structure wrong | `/new-*` skill generates correct structure |
| Each agent session re-discovers the same patterns | Patterns are documented once, read every time |

---

## Resolved Decisions

1. **Charts**: **Recharts** — already proven in the existing codebase, fits the needs.
2. **Virtual scrolling**: Deferred. Will use **TanStack Table** when needed for DB list and sample event logs.
3. **SSR/SSG**: Not needed. All apps are CSR only. Embed app handles OpenGraph previews.
4. **Auth**: Discord OAuth is **removed**. No user authentication in the rewrite.
5. **Design system**: **Formalize before building components.** Define tokens (colors, spacing scale, typography scale, border radii, shadows) as a prerequisite before any primitives work begins. This goes into the Tailwind `@theme` CSS in `packages/primitives/`.
