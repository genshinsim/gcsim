# Yarn 3 → pnpm Migration Checklist (Proof of Concept)

Branch: `explore/pnpm-migration`. This is an **exploratory PoC** to see what a pnpm
monorepo migration for the `ui/` workspace actually requires. Not intended to merge as-is.

Current state (pre-migration):
- Yarn 3.2.4 (Berry), `nodeLinker: node-modules` (not PnP)
- Native Yarn workspaces (`packages/*`), Lerna 6 present but **vestigial** (nothing calls `lerna run`)
- Cross-package deps use `workspace:^` (pnpm-compatible as-is)
- pnpm 11.5.1 available locally; Node v24

## 1. pnpm config plumbing
- [ ] Add `ui/pnpm-workspace.yaml` with `packages: ["packages/*"]`
- [ ] `ui/package.json`: `packageManager` → `pnpm@...`, remove `workspaces` array
- [ ] Add `ui/.npmrc` (hoisting strategy — start strict, add escape hatch only if needed)
- [ ] Delete `.yarnrc.yml`, `.yarn/` (releases/plugins/sdks), `yarn.lock`, `lerna.json`
- [ ] Rewrite `ui/.gitignore` (drop `.yarn/*` allowlist block)
- [ ] `pnpm install` → generates `pnpm-lock.yaml`

## 2. Rewrite internal `yarn run` / `yarn workspace` calls
Package scripts (`yarn run X` → `pnpm run X`):
- [ ] `executors` `watch:web`
- [ ] `types` `watch`
- [ ] `web` / `db` / `embed` `preview`

Root scripts (`yarn workspace @gcsim/X Y` → `pnpm --filter @gcsim/X Y`):
- [ ] All `start:*`, `build:*`, `preview:*`, `deploy:*`, `storybook` scripts in `ui/package.json`

## 3. CI + shell scripts (repo-wide, outside ui/)
- [ ] `.github/actions/yarn-setup-and-test/action.yml` (install, cache, + `pnpm/action-setup`)
- [ ] Deploy composite actions: `deploy-ui`, `deploy-db`, `deploy-docs`, `deploy-taghelper`, `deploy-workers`
- [ ] `.github/actions/containers/embedgenerator/action.yml`
- [ ] `Taskfile.yml` (`cd ui && yarn install`)
- [ ] `scripts/build_preview.sh`
- [ ] `containers/embedgenerator/ci/build-linux-amd64.sh`
- [ ] `ui/packages/docs/README.md` (doc-only)

## 4. Phantom dependencies (the real risk)
pnpm strict linking only exposes **declared** deps. Fix undeclared-but-used imports:
- [ ] React consumed transitively (esp. `ui`, `db`) — add explicit `react`/`react-dom` where needed
- [ ] Mixed TS toolchains (TS 4.8/ESLint 5 vs `components`/`embed` on TS 5.2/ESLint 7) may install multiple versions side-by-side → surface type conflicts
- [ ] Bridge option if blocked: `shamefully-hoist=true` (temporary) or targeted `public-hoist-pattern`

## 5. Verify
- [ ] `pnpm install` clean
- [ ] `pnpm build` (wasm + web) green
- [ ] `pnpm --filter @gcsim/embed build`, `db`, `docs`, `taghelper`
- [ ] `pnpm lint`, dep-hygiene scripts (`syncpack`, `madge`)

## Status (PoC result on `explore/pnpm-migration`)

**Working under pnpm 11.5.1 / Node 24:**
- [x] `pnpm install` clean; `pnpm install --frozen-lockfile` from scratch exits 0 (CI path)
- [x] `pnpm build` (wasm via `task` + web) — exit 0, real `dist/` output
- [x] `pnpm --filter @gcsim/{embed,db,taghelper} build` — all exit 0
- [x] Root scripts + internal `yarn run` calls rewritten to pnpm
- [x] CI actions, Taskfile, shell scripts, docs README migrated

**Key mechanics learned:**
- pnpm v11 gates lifecycle scripts via an **`allowBuilds` map** (name → true/false) in
  `pnpm-workspace.yaml`, NOT the older `onlyBuiltDependencies` list. With `strictDepBuilds`
  on (default), an unresolved build **fails install (exit 1)** — this propagates through
  `pnpm --filter run`'s pre-run deps check and breaks builds until every native/codegen
  dep (swc, esbuild, sharp, workerd, protobufjs, core-js) is explicitly allowed.
- web/db bundle workspace **source** directly (vite-tsconfig-paths), so every transitive
  package's phantom deps surface, not just the leaf app's.

**Phantom deps fixed** (undeclared, previously satisfied by Yarn hoisting):
web, ui, db, embed, taghelper, localization, components — see per-commit messages.
Included undeclared **workspace** deps (e.g. ui→components/localization) and one
package self-import (`components` importing `@gcsim/components`), converted to relative.

## Known remaining issues (NOT migration regressions)
- **docs** (`@gcsim/docs`) build fails: Docusaurus 2.4.1 + webpack 5.110.3 `ProgressPlugin`
  schema validation error. Only one webpack copy is installed, so this is a Docusaurus/
  Node-24 compat problem, not pnpm resolution — would likely fail under Yarn on this Node.
  Fix path: bump Docusaurus, or pin a compatible webpack/Node.
- **lint**: `pnpm lint` reports 275 pre-existing eslint problems (210 errors) in untouched
  source (e.g. `workers/src`). CI has lint steps commented out. Pre-existing debt.
- **syncpack list-mismatches** (exit 1): pre-existing cross-package version divergence
  (e.g. `vite` ^4 vs ^5), plus `wouter` ^2.9.0 (db/taghelper) vs ^3.1.0 (embed) — the
  versions the PoC declared were copied from existing workspace usage. Run `fix-mismatches`
  to unify.
- **madge circular** (exit 1): 1 source-level cycle `ui/Stores/store.ts ↔ userSlice.ts`.
  Pre-existing, unrelated to package manager.

## Follow-ups before this could merge (beyond PoC scope)
- Rename `.github/actions/yarn-setup-and-test` → `pnpm-setup-and-test` and update the 3
  workflow files referencing it (tests.yml, deploy.yml, build-images.yaml).
- Confirm `corepack`/pnpm version provisioning in every CI runner.
- Resolve version mismatches (syncpack) and decide on the TS 4.8 vs 5.2 toolchain split
  now that pnpm makes divergent versions explicit.
- Verify `pnpm dev`/`start` (watch-wasm + vite) locally; delete stray repo-root `.turbo/`.
