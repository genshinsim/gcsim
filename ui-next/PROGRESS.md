# Implementation Progress

## Spec Corrections

Wherever the spec says "Tailwind v5", read "Tailwind v4". Wherever it says "Vite 6", read "Vite 8". Wherever it says "Storybook 8", read "Storybook 10".

## Phase 0: Monorepo Scaffolding + Agent Tooling + CI

| Step | Status | Description |
|------|--------|-------------|
| 0.0 | DONE | Dependency version verification |
| 0.1 | DONE | Initialize monorepo root |
| 0.2 | NOT STARTED | Root CLAUDE.md + auxiliary docs |
| 0.3 | DONE | Basic CI pipeline |
| 0.4 | DONE | Scaffolding skills |
| 0.5 | DONE | Subagent definitions |

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
