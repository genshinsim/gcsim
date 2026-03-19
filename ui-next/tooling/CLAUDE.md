# ui-next/tooling — Shared Configuration

Shared base configs that packages and apps extend. This keeps configuration DRY across the monorepo.

## TypeScript Base Config (`typescript/base.json`)

Shared strict tsconfig. Packages and apps extend it:

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

Key settings: `strict: true`, `composite: true`, `moduleResolution: "bundler"`, `jsx: "react-jsx"`.

## Vitest Base Config (`vitest/base.ts`)

Shared Vitest config. Packages extend it:

```ts
import baseConfig from "../../tooling/vitest/base.ts";
import { defineConfig, mergeConfig } from "vitest/config";

export default mergeConfig(baseConfig, defineConfig({
  // package-specific overrides
}));
```

Key settings: `environment: "jsdom"`, `globals: true`, coverage via `v8`.

## Test Fixtures (`test-fixtures/`)

Canonical mock data used by ALL packages for testing. Prevents fixture drift across packages.

- Import from `../../tooling/test-fixtures` (relative path)
- Contains realistic mock objects: `SimResult`, character data, etc.
- When adding a new fixture: create the file, add a re-export in `index.ts`
- Do NOT create package-local fixtures for shared data shapes — use these instead

## Don'ts

- Don't duplicate base config values in package configs — extend and override only what differs
- Don't add package-specific test utilities here — those belong in the package
- Don't import from tooling at runtime — these are dev/build-time only
