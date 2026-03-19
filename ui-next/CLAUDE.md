# ui-next — gcsim Web UI Rewrite

Branch: `web-rewrite`

## Key Files

- `PROGRESS.md` — current implementation progress and what's been done
- `DEPENDENCIES.md` — pinned dependency versions (MUST reference when installing packages)

**IMPORTANT**: update `PROGRESS.md` once you have completed any implementation step. Include key details the next agent may require.

## Scaffolding Skills

Use these skills (via `/skill-name` or natural language) to scaffold new code consistently:

- `/new-package` — create a new `@gcsim/` shared package in `packages/`
- `/new-component` — create a React component in an existing package
- `/new-page` — create a page with lazy-loaded routing in an app
- `/new-app` — create a full Vite + React app in `apps/`
- `/new-store` — create a typed Zustand store in an app
- `/check` — run full lint → typecheck → test → dependency-cruiser → build pipeline
- `/dev` — start Vite dev server for an app

## Commit Discipline

Commit your work in sensible increments — don't let changes pile up. Guidelines:

- **Commit after completing each spec step** (e.g., after finishing step 0.1, commit before starting 0.2)
- **Commit after any self-contained unit of work** — a new config file, a package scaffold, a batch of related changes
- **Don't bundle unrelated changes** into a single commit
- **Commit messages** should be concise and describe the "what" (e.g., "scaffold monorepo root with turbo, biome, vitest")
- **When in doubt, commit sooner** — small commits are easier to review and revert
