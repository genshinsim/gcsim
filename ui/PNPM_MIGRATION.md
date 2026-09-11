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

## Notes / findings
(filled in as migration proceeds)
