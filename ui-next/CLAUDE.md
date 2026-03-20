# ui-next — gcsim Web UI Rewrite

Branch: `web-rewrite`

## Key Files

- `PROGRESS.md` — current implementation progress and what's been done
- `DEPENDENCIES.md` — pinned dependency versions (MUST reference when installing packages)

**IMPORTANT**: update `PROGRESS.md` once you have completed any implementation step. Include key details the next agent may require.

## Monorepo Structure

```
ui-next/
├── apps/          # Independently deployable web apps (web, db, embed, storybook, taghelper)
├── packages/      # Shared libraries consumed by apps (@gcsim/ scope)
├── assets/        # Shared static assets (character portraits, element icons, etc.)
└── tooling/       # Shared configs (tsconfig, vitest, test fixtures)
```

- **apps/** own pages, stores, and route definitions — they are composition layers that wire packages together
- **packages/** contain reusable logic and components — they never import from apps or other apps
- **assets/** are referenced via `@gcsim/assets/...` imports or public/ symlinks
- **tooling/** provides shared base configs that packages/apps extend

## Commands

| Command | Description |
|---------|-------------|
| `turbo run build` | Build all packages and apps |
| `turbo run build --filter=@gcsim/<pkg>` | Build a specific package |
| `turbo run test` | Run all tests |
| `turbo run test --filter=@gcsim/<pkg>` | Run tests for a specific package |
| `turbo run typecheck` | Typecheck all packages |
| `turbo run typecheck --filter=@gcsim/<pkg>` | Typecheck a specific package |
| `pnpm lint` | Run Biome linting/formatting |
| `pnpm depcruise` | Run dependency-cruiser validation |
| `pnpm add <dep> --filter=@gcsim/<pkg>` | Add a dependency to a specific package |
| `pnpm add -D <dep> --filter=@gcsim/<pkg>` | Add a dev dependency to a specific package |

## Architecture Rules

1. **Protobuf types are the single source of truth.** No shared type aliases across packages. Apps may define local convenience types.
2. **Import from package index only.** Never import from internal paths like `@gcsim/viewer/src/charts/...`. Enforced by `dependency-cruiser` in CI.
3. **Components live in apps until a second consumer appears.** Only then extract to a shared feature package.
4. **TanStack Query for all server data.** No raw `fetch` + `useEffect` patterns.
5. **Zustand for client state.** With `localStorage` middleware for persistence.
6. **All apps are CSR.** No SSR/SSG.
7. **No authentication.** Discord OAuth is removed.
8. **Error boundaries required** around viewer (large data parsing/rendering), chart components, and any component that fetches external data.
9. **Route-level code splitting** — all page components loaded via lazy imports. Required from Phase 4 onward.
10. **Tailwind v4 CSS-first theming** — design tokens defined via `@theme` in CSS, NOT JavaScript config files.

## Naming Conventions

- **Files:** kebab-case (`damage-timeline.tsx`, `sim-result.ts`)
- **Components:** PascalCase (`DamageTimeline`, `TeamHeader`)
- **Functions/variables:** camelCase (`useSimResult`, `formatDamage`)
- **Packages:** kebab-case with `@gcsim/` scope (`@gcsim/viewer`, `@gcsim/types`)

## Biome

Biome handles both linting and formatting. Config in `biome.json`:
- Space indentation (2 spaces)
- Double quotes
- Semicolons always
- Line width 100
- Run `npx biome check --write <files>` to auto-fix

## Tailwind v4

Uses CSS-first `@theme` configuration — NOT JS config files. Design tokens defined in `packages/primitives/src/theme.css`. Apps import this via `@import '@gcsim/primitives/theme.css'` in their `app.css`.

## Testing Protocol

Every implementation step follows: write failing test → implement → pass test → typecheck → lint.

| Package type | What to test | What NOT to test |
|---|---|---|
| `types` | Exports exist, types compile | Generated code internals |
| `data` | Exports exist, data shape matches types | Individual data values |
| `i18n` | Language loading, key resolution, fallback | Translation string content |
| `api` | Request construction, error handling, abort | Actual HTTP calls (mock fetch) |
| `executor` | Interface contract, state transitions | WASM internals (mock workers) |
| `primitives` | Renders, variants, accessibility, keyboard | Visual appearance |
| Feature packages | Component behavior with mocked data | Internal implementation |
| Apps (unit) | Store logic, data transforms, hooks | Component rendering |

## Environment Variables

Apps use `import.meta.env.VITE_*`. See `.env.example` for available variables:
- `VITE_API_BASE_URL` — backend API (default: same origin)
- `VITE_WASM_BASE_URL` — WASM binary location (default: `/api/wasm`)
- `VITE_LOCAL_DEV_URL` — local dev server (default: `http://127.0.0.1:8381`)

## WASM Binary

- **Local dev:** pre-built binary in `assets/wasm/` or fetched from `VITE_WASM_BASE_URL`
- **Production:** served from CDN at `/api/wasm/{branch}/{commit}/main.wasm`
- Build with `task wasm` (output to `assets/wasm/` for local dev)

## Scaffolding Skills

Use these skills to scaffold new code consistently:

- `/new-package` — create a new `@gcsim/` shared package in `packages/`
- `/new-component` — create a React component in an existing package
- `/new-page` — create a page with lazy-loaded routing in an app
- `/new-app` — create a full Vite + React app in `apps/`
- `/new-store` — create a typed Zustand store in an app
- `/check` — run full lint → typecheck → test → dependency-cruiser → build pipeline
- `/dev` — start Vite dev server for an app

## Subagents

Dispatch these via the Agent tool for specialized work:

- `package-reviewer` — qualitative review of a package (after `/check` passes)
- `package-tester` — run tests with failure diagnosis
- `feature-implementer` — TDD implementation of a single feature component
- `cross-package-integrator` — wire packages into an app with integration tests

## Commit Discipline

- **Commit often with small, standalone, tested changes** — each commit should build and pass tests on its own
- **Commit after any self-contained unit of work** — a new config file, a component + its tests, a batch of related changes
- **Don't bundle unrelated changes** into a single commit
- **Commit messages** should be concise and describe the "what"
- **When in doubt, commit sooner** — small commits are easier to review and revert
- **Each commit must be independently valid**: it should typecheck, pass tests, and not break the build

## Git Workflow

- Each feature/package works on a branch: `phase-{N}/step-{X.Y}-{package-name}`
- **Each feature must be PR'd into `web-rewrite`** before the task is considered finished
- Small commits make PRs easier to review — each commit should be standalone
- PR title should reference the spec step (e.g., "Phase 3.1: @gcsim/avatar package")
- After PR is merged, the feature branch can be deleted
