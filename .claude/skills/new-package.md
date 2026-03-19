---
name: new-package
description: >
  Scaffold a new @gcsim/ shared package in the ui-next monorepo. Use this skill
  whenever the user wants to create a new package, add a new library, set up a
  new workspace module, or start a new shared utility in ui-next/packages/. Also
  use it when the user says things like "I need a new package for X" or "let's
  add a types package" or "create a shared library for Y".
---

# /new-package

Scaffold a new shared package in `ui-next/packages/`. Every package in this
monorepo follows the same structure so that tooling (turbo, biome, vitest) works
uniformly -- this skill ensures new packages are consistent from the start.

## Input

Parse the user's message for:
- **name** (required): package name (e.g., `types`, `primitives`, `i18n`)
- **dependencies** (optional): comma-separated list of other `@gcsim/` packages this depends on
- **--with-tailwind** (optional flag): include Tailwind CSS setup

## Steps

1. Create directory `ui-next/packages/<name>/`.

2. Create `package.json`. The monorepo uses pnpm workspaces and ESM throughout,
   so `"type": "module"` and the `exports` map are required for packages to
   resolve correctly across the workspace:
```json
{
  "name": "@gcsim/<name>",
  "version": "0.0.0",
  "private": true,
  "type": "module",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "exports": {
    ".": {
      "import": "./dist/index.js",
      "types": "./dist/index.d.ts"
    }
  },
  "scripts": {
    "build": "tsc --build",
    "test": "vitest run",
    "typecheck": "tsc --noEmit",
    "lint": "biome check"
  },
  "devDependencies": {
    "typescript": "5.9.3",
    "vitest": "4.1.0",
    "@testing-library/react": "16.3.0",
    "@testing-library/jest-dom": "6.6.3"
  }
}
```
If dependencies are provided, add them under `"dependencies"` as
`"@gcsim/<dep>": "workspace:*"` so pnpm links to the local workspace copy.

3. Create `tsconfig.json`. Extend the shared base config so all packages share
   the same strict TypeScript settings:
```json
{
  "extends": "../../tooling/typescript/base.json",
  "compilerOptions": {
    "outDir": "dist",
    "rootDir": "src"
  },
  "include": ["src"]
}
```
If dependencies are provided, add a `"references"` array with entries like
`{ "path": "../<dep>" }` -- this enables TypeScript project references so
`tsc --build` resolves cross-package types correctly.

4. Create `src/index.ts` with a comment. This barrel file is the single public
   API surface for the package -- all consumers import from here, never from
   internal paths:
```typescript
// @gcsim/<name> public API
```

5. Create `vitest.config.ts`. Merge with the shared base config so test
   behavior (reporters, coverage thresholds) stays consistent across packages:
```typescript
import { defineConfig, mergeConfig } from "vitest/config";
import baseConfig from "../../tooling/vitest/base.ts";

export default mergeConfig(baseConfig, defineConfig({}));
```

6. Create `CLAUDE.md` so future Claude sessions understand this package:
```markdown
# @gcsim/<name>

## Purpose

TODO: describe what this package does

## How to add a new module

TODO: describe the pattern for adding new code

## Canonical example

TODO: point to a representative file

## Public API

All public exports go through `src/index.ts`.

## Dependencies

TODO: list consumed @gcsim/ packages

## Don'ts

- Don't import from other packages' `src/` — only from their package index
- Don't add app-specific logic here
```

7. If `--with-tailwind` was requested, create `src/styles.css`:
```css
@import "tailwindcss";
```

8. Run `pnpm install` from `ui-next/` to link the new workspace package into
   the pnpm lockfile.

9. Verify the package builds cleanly:
```bash
cd ui-next && pnpm --filter @gcsim/<name> typecheck && pnpm --filter @gcsim/<name> build
```
Both commands should pass. If either fails, fix the issue before finishing.
