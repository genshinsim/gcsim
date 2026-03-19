---
name: check
description: >
  Run the full lint, typecheck, test, dependency-cruiser, and build pipeline for
  the ui-next monorepo. Use this skill whenever the user wants to verify their
  changes, run CI checks locally, validate the build, check for lint or type
  errors, or run tests across ui-next packages. Also use it when the user says
  /check.
---

# /check

Runs the full verification pipeline for the ui-next monorepo. Stops on first failure.

## Input

Parse the user's message for:
- **--filter=\<package\>** (optional): scope checks to a specific package (e.g., `--filter=@gcsim/types`, `--filter=types`)

If a bare name is given (e.g., `types`), prefix it with `@gcsim/` for turbo filter.

## Steps

Run each command sequentially from the `ui-next/` directory. **Stop on first failure** and report the error.

1. **Lint**:
```bash
cd ui-next && pnpm lint
```
If filtered, narrow the biome scope: `biome check packages/<name>` or `biome check apps/<name>`.

2. **Typecheck**:
```bash
cd ui-next && turbo run typecheck
```
If filtered: `turbo run typecheck --filter=@gcsim/<name>`

3. **Test**:
```bash
cd ui-next && turbo run test
```
If filtered: `turbo run test --filter=@gcsim/<name>`

4. **Dependency cruiser**:
```bash
cd ui-next && pnpm depcruise
```
This already handles the "no source files" case gracefully.

5. **Build**:
```bash
cd ui-next && turbo run build
```
If filtered: `turbo run build --filter=@gcsim/<name>`

## Output

Report each step's result (pass/fail). On failure, show the error output and stop — do not continue to subsequent steps.
