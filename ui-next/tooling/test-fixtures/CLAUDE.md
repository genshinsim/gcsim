# test-fixtures

Canonical mock data for testing across the ui-next monorepo.

## Available fixtures

- **`mockHutao`** / **`mockXingqiu`** — `Sim.Character` objects with realistic builds
- **`mockCharacters`** — array of both mock characters
- **`mockSimResult`** — complete `Sim.SimResults` with 2 characters, statistics, DPS data, and a config string

## How to import

```typescript
import { mockSimResult, mockCharacters } from "../../tooling/test-fixtures/index.js";
```

## Adding new fixtures

1. Create a new file in this directory (kebab-case)
2. Re-export from `index.ts`
3. Use types from `@gcsim/types` (imported via relative path)

## Don'ts

- Don't create package-local fixtures for shared data shapes — add them here instead
- Don't use these fixtures at runtime — they are for tests only
