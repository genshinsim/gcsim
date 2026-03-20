# Implementation Progress

## Spec Corrections

Wherever the spec says "Tailwind v5", read "Tailwind v4". Wherever it says "Vite 6", read "Vite 8". Wherever it says "Storybook 8", read "Storybook 10".

## Phase 0: Monorepo Scaffolding + Agent Tooling + CI

| Step | Status | Description |
|------|--------|-------------|
| 0.0 | DONE | Dependency version verification |
| 0.1 | DONE | Initialize monorepo root |
| 0.2 | DONE | Root CLAUDE.md + auxiliary docs |
| 0.3 | DONE | Basic CI pipeline |
| 0.4 | DONE | Scaffolding skills |
| 0.5 | DONE | Subagent definitions |
| 0.6 | DONE | Claude Code hooks |
| 0.7 | DONE | Git workflow protocol |

### Step 0.0 — Dependency Version Verification (DONE)

- Verified all dependency versions against npm (2026-03-19)
- Created `ui-next/DEPENDENCIES.md` with pinned versions
- Key corrections from spec:
  - **Tailwind CSS v4.2.2** (spec said v5 — v5 doesn't exist; CSS-first `@theme` is a v4 feature)
  - **Vite 8.0.1** (spec said v6)
  - **Storybook 10.3.1** (spec said v8)
- Confirmed shadcn v4.1.0 is fully compatible with Tailwind v4 `@theme`
- `tw-animate-css` replaces deprecated `tailwindcss-animate` for Tailwind v4
- `tailwind-merge` v3.x is the correct line for Tailwind v4

### Step 0.1 — Initialize Monorepo Root (DONE)

- Created `package.json` (private, `@gcsim/monorepo` scope) with build/test/typecheck/lint scripts
- Created `pnpm-workspace.yaml` with `apps/*` and `packages/*`
- Created `turbo.json` with `build`, `test`, `typecheck`, `lint` pipelines (build/test/typecheck depend on `^build`)
- Created `biome.json` with recommended linting, space indentation, double quotes, semicolons
- Installed root dev deps: `turbo@2.8.20`, `@biomejs/biome@2.4.8`, `vitest@4.1.0`, `dependency-cruiser@17.3.9`, `typescript@5.9.3`
- Created `tooling/typescript/base.json` — strict tsconfig with composite, bundler moduleResolution, react-jsx
- Created `tooling/vitest/base.ts` — shared Vitest config (jsdom, globals, v8 coverage)
- Created `.env.example` with `VITE_API_BASE_URL`, `VITE_WASM_BASE_URL`, `VITE_LOCAL_DEV_URL`
- Created `.gitignore` (node_modules, dist, .turbo, .env, *.wasm, coverage)
- Copied static assets from `ui/` to `assets/` (favicon, stat icons, logos, wasm_exec.js)
- Created `assets/wasm/` directory for local WASM dev builds
- Verified: `pnpm install` succeeds, `turbo run build` succeeds (empty)

### Step 0.3 — Basic CI Pipeline (DONE)

- Created `.github/workflows/ui-next.yml` — GitHub Actions workflow triggered on push/PR to `web-rewrite` branch (path-filtered to `ui-next/**`)
- Pipeline steps: pnpm install → biome check → typecheck → test → dependency-cruiser → build
- Uses `pnpm/action-setup@v4` (v10), `actions/setup-node@v4` (Node 22), pnpm cache
- Created `ui-next/.dependency-cruiser.cjs` with rules:
  - `no-deep-package-imports` — import from package index only, never `@gcsim/<pkg>/src/...`
  - `no-circular` — no circular dependencies
  - `no-app-to-app` — apps must not import from other apps
  - `no-package-to-app` — packages must not import from apps
- Fixed `biome.json` for Biome 2.x: replaced deprecated `files.ignore` with `files.includes` scoped to `apps/**`, `packages/**`, `tooling/**` (excludes vendored `assets/wasm/wasm_exec.js`)
- Verified: all CI steps pass locally against empty monorepo

### Step 0.4 — Scaffolding Skills (DONE)

Created 7 Claude Code skills in `.claude/skills/`:
- **`/new-package`** (0.4a) — scaffolds `ui-next/packages/<name>/` with package.json, tsconfig, vitest config, CLAUDE.md, optional Tailwind
- **`/new-component`** (0.4b) — scaffolds a React component with .tsx, .test.tsx, index.ts barrel, wires into package exports
- **`/new-page`** (0.4c) — scaffolds a page in an app with lazy-loaded route entry
- **`/new-app`** (0.4d) — scaffolds a full Vite + React app with TanStack Query, Router, Zustand, i18n, Tailwind
- **`/new-store`** (0.4e) — scaffolds a typed Zustand store, optional `--persist` for localStorage middleware
- **`/check`** (0.4f) — runs sequential pipeline: biome → typecheck → test → dependency-cruiser → build (stops on first failure)
- **`/dev`** (0.4g) — builds dependencies then starts Vite dev server for specified app

### Step 0.5 — Subagent Definitions (DONE)

Created 4 Claude Code agent definitions in `.claude/agents/`:
- **`package-reviewer`** (0.5a) — qualitative review of a single package for boundary violations, type alias misuse, data-fetching patterns, test quality, CLAUDE.md completeness, design token usage, error boundaries
- **`package-tester`** (0.5b) — runs typecheck + tests for a package, diagnoses failures with specific fix suggestions
- **`feature-implementer`** (0.5c) — TDD-based implementation of a single feature component (max 3 sub-components); reads canonical example, writes failing test first, implements, updates CLAUDE.md. Skeleton — to be refined after Phase 2 (primitives) and Phase 3 (feature components)
- **`cross-package-integrator`** (0.5d) — wires completed packages into an app with composition components, integration tests, route updates; requires a composition spec in the dispatch call

### Step 0.2 — Root CLAUDE.md + Auxiliary Docs (DONE)

- Expanded `ui-next/CLAUDE.md` with full monorepo documentation:
  - Monorepo structure overview (apps/ vs packages/ vs assets/ vs tooling/)
  - Turborepo commands (build, test, typecheck with --filter)
  - All 10 architecture rules
  - Naming conventions (kebab-case files, PascalCase components, camelCase functions)
  - pnpm commands for adding dependencies
  - Biome usage and config summary
  - Tailwind v4 CSS-first theming note
  - Testing protocol table (what to test per package type)
  - Environment variables reference
  - WASM binary strategy (local dev vs production)
  - Scaffolding skills and subagent references
  - Commit discipline and git workflow protocol
- Created `ui-next/apps/CLAUDE.md`:
  - Apps own pages, stores, and route definitions
  - Apps are composition layers wiring packages together
  - Error boundary and lazy import requirements
  - Dev server port assignments (web=5173, db=5174, embed=5175)
  - References to `/new-app`, `/new-page`, `/new-store` skills
- Created `ui-next/tooling/CLAUDE.md`:
  - How to extend base tsconfig (with example)
  - How to extend base vitest config (with example)
  - Test fixtures usage (import from `tooling/test-fixtures/`, don't duplicate)
  - Don'ts (no runtime imports from tooling, no config duplication)

### Step 0.6 — Claude Code Hooks (DONE)

- Created `.claude/settings.json` with hooks configuration:
  - **PreToolUse (Bash):** `.claude/hooks/pre-commit-biome.sh` — intercepts `git commit` commands, runs `biome check --write` on staged ui-next files, blocks commit (exit 2) if biome fails
  - **PostToolUse (Edit|Write):** `.claude/hooks/post-edit-test.sh` — async hook that detects which `@gcsim/` package owns the edited file and runs `turbo run test --filter=@gcsim/<pkg>`
- Both scripts are executable and handle edge cases (non-ui-next files, missing paths)

### Step 0.7 — Git Workflow Protocol (DONE)

- Documented git workflow in `ui-next/CLAUDE.md`:
  - Branch naming: `phase-{N}/step-{X.Y}-{package-name}`
  - Worktree naming: `worktree-phase{N}-{package-name}`
  - Merge order: foundation packages first at phase gates
- Created `ui-next/tooling/cleanup-worktrees.sh`:
  - Finds worktrees matching `worktree-phase*` naming convention
  - Prunes missing worktrees
  - Removes clean (no uncommitted changes) worktrees
  - Skips worktrees with uncommitted changes (safety)
  - Reports cleanup summary

### Phase 0 Review Notes

Post-review fixes applied:
- Added `@testing-library/react` and `@testing-library/jest-dom` to `/new-package` skill devDependencies and DEPENDENCIES.md
- Pinned `tw-animate-css` to `1.3.4` in DEPENDENCIES.md (was `latest`)
- Fixed unquoted `$STAGED_FILES` in pre-commit hook (now uses `xargs`)
- Spec says `.claude/hooks.json` but Claude Code uses `.claude/settings.json` — this is correct

**TODO for Phase 1:** Test dependency-cruiser rules with deliberate violations. The `no-deep-package-imports` rule's `via` usage and the `no-app-to-app` rule's `{FROM_APP}` placeholder may not work as intended. Create a test file that imports from `@gcsim/<pkg>/src/...` and verify the rule fires.

**Assets note:** Spec section 3 shows `assets/characters/`, `assets/elements/`, etc. but game assets (character portraits, weapon/artifact images) are served from `/api/assets/` at runtime, not stored in the repo. Current `assets/images/` contains only static UI assets (favicon, icons, logo). This is correct — the spec's asset directories are aspirational for when/if assets are bundled locally.

## Phase 1: Foundation Packages

| Step | Status | Description |
|------|--------|-------------|
| 1.1 | DONE | `@gcsim/types` + `tooling/test-fixtures` |
| 1.2 | DONE | `@gcsim/data` |
| 1.3 | DONE | `@gcsim/i18n` |
| 1.4 | DONE | `@gcsim/api` |
| 1.5 | DONE | `@gcsim/executor` |

### Step 1.1 — `@gcsim/types` + `tooling/test-fixtures` (DONE)

- Created `packages/types/` with buf/protobuf-es generation pipeline
- `buf.gen.yaml` points at `protos/` at repo root, generates into `src/generated/`
- Generated TypeScript types from all 10 `.proto` files (7 model + 3 backend)
- Ported all custom interfaces from `ui/packages/types/src/sim.ts` into `src/sim.ts`
  - Includes: SimResults, Statistics, SummaryStat, Character, Enemy, Weapon, Talent, etc.
  - Did NOT port `user.ts` (auth removed per spec)
- Public API via `src/index.ts`:
  - Model proto types exported at top level (SimulationResult, Character, etc.)
  - Backend proto types namespaced: `share`, `db`, `preview`
  - Custom interfaces namespaced: `Sim` (use as `Sim.SimResults`, `Sim.Character`, etc.)
- 16 tests passing, typecheck clean, build succeeds
- Created `tooling/test-fixtures/` with canonical mock data:
  - `sim-result.ts` — mock SimResults with 2 characters, stats, config
  - `characters.ts` — mock Character objects (hutao, xingqiu)
  - `index.ts` — re-exports
- Dependencies: `@bufbuild/protobuf@2.11.0`, `@bufbuild/buf@1.66.1`, `@bufbuild/protoc-gen-es@2.11.0`
- Branch `phase-1/step-1.1-types` merged into `web-rewrite`

### Fixes applied during Phase 1

- Fixed pre-commit biome hook: file paths weren't stripped of `ui-next/` prefix after `cd ui-next/` (caused biome to look for `ui-next/ui-next/...`)
- Added `.turbo/` to root `.gitignore` (turbo cache dir was showing as untracked)

### Phase 1 Learnings (for next session)

- **Worktree branch issue:** `isolation: "worktree"` agents create branches from the repo's default HEAD (usually `main`), not from the current branch (`web-rewrite`). The `ui-next/` directory only exists on `web-rewrite`. Agents need explicit `git reset --hard origin/web-rewrite` instructions, but this gets blocked by sandbox permissions.
- **Recommended approach:** Either (a) create branches and worktrees manually before dispatching agents, or (b) don't use worktree isolation and instead run agents sequentially on the main working tree, or (c) pre-create branches from `web-rewrite` and have agents use them.
- **Steps 1.2–1.5 are independent** of each other (all depend only on `@gcsim/types` which is done). They can be parallelized once the worktree issue is resolved.
- **Worktree approach (resolved):** Pre-create worktrees with `git worktree add -b <branch> .claude/worktrees/<name> web-rewrite`, then dispatch agents to work in them. Agents can write files but cannot run bash in worktrees (sandbox limitation). Main agent must verify/commit each worktree's work.

### Step 1.2 — `@gcsim/data` (DONE)

- Created `packages/data/` with typed exports for game data
- Ported `latest_chars.json` (version→character key mapping) and `tags.json` (tag ID→display info)
- Typed interfaces: `TagInfo`, `TagMap`, `LatestCharsMap`
- `tsconfig.json` requires `resolveJsonModule: true` and `include: ["src/**/*.ts", "src/**/*.json"]`
- 11 tests passing, typecheck clean, build succeeds

### Step 1.3 — `@gcsim/i18n` (DONE)

- Created `packages/i18n/` with i18next + react-i18next
- Ported 7 language JSON files + names.generated.json + names.traveler.json
- Uses spread operator instead of lodash-es merge for combining name resources
- `initI18n(lng)` function for app initialization, `resources` object, `specialLocales` array
- Two namespaces: `translation` (UI strings) and `game` (character/entity names)
- `tsconfig.json` requires `resolveJsonModule: true` and `include: ["src/**/*.ts", "src/**/*.json"]`
- 6 tests passing, typecheck clean, build succeeds
- Dependencies: `i18next@25.8.20`, `react-i18next@16.5.8`

### Step 1.4 — `@gcsim/api` (DONE)

- Created `packages/api/` with typed fetch functions (new package, not ported)
- `apiFetch<T>()` base wrapper with gzip decompression via pako
- `ApiError` class with HTTP status code
- Endpoints: `fetchShareResult()`, `fetchDBResult()`, `queryDB()` (with pagination/abort), `fetchLocalResult()`
- Uses native `fetch` (not axios)
- 16 tests passing (4 test files), typecheck clean, build succeeds
- Dependencies: `@gcsim/types`, `pako@2.1.0`

### Step 1.5 — `@gcsim/executor` (DONE)

- Created `packages/executor/` ported from `ui/packages/executors/`
- `Executor` interface using `Sim` namespace types (`Sim.SimResults`, `Sim.ParsedResult`, `Sim.Sample`)
- `ServerExecutor` — HTTP-based with native fetch (replaces axios), async/await polling
- `WasmExecutor` — web worker pool with aggregator, configurable 1-30 workers
- `ExecutorError` class with typed error codes (`NETWORK`/`SERVER`/`PARSE`/`UNKNOWN`)
- Native throttle utility (replaces lodash-es)
- Worker files excluded from tsconfig (standalone scripts with duplicate function names)
- `tsconfig.json` excludes worker files: `"exclude": ["src/workers/worker.ts", "src/workers/aggregator.ts", "src/workers/helper.ts"]`
- 28 tests passing (3 test files), typecheck clean, build succeeds

### Phase 1 Gate (PASSED)

- All 5 packages: biome clean, typecheck pass, tests pass, build succeeds
- Total tests: 16 (types) + 11 (data) + 6 (i18n) + 16 (api) + 28 (executor) = 77 tests
- Dependency-cruiser: no violations
- All branches merged into `web-rewrite`, worktrees cleaned up

## Phase 2: Design System + Primitives

| Step | Status | Description |
|------|--------|-------------|
| 2.1 | DONE | Design system tokens |
| 2.2-2.4 | DONE | All shadcn primitives (installed via CLI) |
| 2.5 | DONE | Storybook setup |

### Steps 2.1–2.4 — `@gcsim/primitives` (DONE)

- Created `packages/primitives/` using shadcn CLI (`shadcn@4.1.0 init` + `shadcn add`)
- **Used shadcn CLI** to install all components instead of hand-writing them
- `components.json` configured for: `style: "radix-nova"`, `rsc: false`, `baseColor: "neutral"`, `iconLibrary: "lucide"`
- Theme uses Tailwind v4 CSS-first `@theme inline` with OKLCH color system (light + `.dark` variants)
- Custom Genshin element colors added: anemo, geo, electro, hydro, pyro, cryo, dendro
- `cn()` utility from `src/lib/utils.ts` (clsx + tailwind-merge)
- 11 shadcn components installed: Button, Card, Input, Tabs, Select, Badge, Dialog, DropdownMenu, Tooltip, ScrollArea, Skeleton
- All components use unified `radix-ui` package with relative imports (no `@/` path aliases — tsc doesn't rewrite them)
- `tsconfig.json` has `baseUrl` + `paths` for `@/*` alias; `vitest.config.ts` mirrors with `resolve.alias`
- `vitest.config.ts` includes `test-setup.ts` for `@testing-library/jest-dom/vitest` matchers
- Barrel export in `src/index.ts` re-exports all components and `cn()`
- Apps import theme via `@import '@gcsim/primitives/theme.css';` (exported via `package.json` exports)
- 14 tests passing (5 cn utility + 9 component rendering), typecheck clean, build succeeds
- Dependencies: `radix-ui`, `class-variance-authority`, `clsx`, `tailwind-merge`, `shadcn`
- DevDependencies: `tailwindcss@4.2.2`, `@tailwindcss/vite@4.2.2`, `tw-animate-css@1.3.4`, `lucide-react@0.577.0`

### Step 2.5 — Storybook (DONE)

- Created `apps/storybook/` with Storybook 10.3.1 + `@storybook/react-vite`
- 30 stories covering all 11 primitives with autodocs
- Tailwind integration via `@tailwindcss/postcss` (PostCSS, not Vite plugin — Storybook's Vite pipeline doesn't reliably pick up the Vite plugin)
- `storybook.css` inlines theme.css content (can't `@import` across package boundaries for Tailwind processing)
- `@source "../../../packages/primitives/src"` tells Tailwind to scan component source files for class names
- Dev server: `pnpm --filter @gcsim/storybook dev` on port 6006

### Phase 2 Fixes

- Enabled `tailwindDirectives: true` in `biome.json` CSS parser (Biome 2.x doesn't parse `@theme`, `@custom-variant`, `@apply` without it)
- Replaced all `@/lib/utils` and `@/components/ui/*` imports in shadcn components with relative paths — TypeScript path aliases are NOT rewritten by `tsc --build`, causing runtime resolution failures
- Removed embedded `.git` directory created by `shadcn init`

### Phase 2 Learnings

- **shadcn CLI creates a `.git` in the package** — must delete before committing
- **`@/` path aliases don't work in library packages** — tsc emits them verbatim in `.js` output. Use relative imports instead. The `@/` alias in `tsconfig.json` + `vitest.config.ts` is kept only for test resolution.
- **Storybook + Tailwind v4**: Use `@tailwindcss/postcss` with a `postcss.config.mjs`, not `@tailwindcss/vite`. The Vite plugin doesn't reliably activate in Storybook's internal Vite pipeline.
- **`@source` paths are relative to the CSS file**, not the project root. Count `../` carefully.
- **`@import "@gcsim/primitives/theme.css"` doesn't work for Tailwind processing** across package boundaries — inline the theme CSS in each consuming app's stylesheet instead.

### Phase 2 Gate

- All primitives: typecheck pass, 14 tests pass, build succeeds
- Storybook: all 30 stories render with correct Tailwind styling
- Total tests: 77 (Phase 1) + 14 (primitives) = 91 tests
