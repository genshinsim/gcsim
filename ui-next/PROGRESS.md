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
