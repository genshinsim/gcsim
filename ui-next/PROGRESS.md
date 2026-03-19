# Implementation Progress

## Spec Corrections

Wherever the spec says "Tailwind v5", read "Tailwind v4". Wherever it says "Vite 6", read "Vite 8". Wherever it says "Storybook 8", read "Storybook 10".

## Phase 0: Monorepo Scaffolding + Agent Tooling + CI

| Step | Status | Description |
|------|--------|-------------|
| 0.0 | DONE | Dependency version verification |
| 0.1 | DONE | Initialize monorepo root |
| 0.2 | NOT STARTED | Root CLAUDE.md + auxiliary docs |
| 0.3 | NOT STARTED | Basic CI pipeline |
| 0.4 | NOT STARTED | Scaffolding skills |

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
