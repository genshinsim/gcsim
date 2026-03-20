# @gcsim/viewer

## Purpose

Presentational components for displaying gcsim simulation results. Renders metadata, team composition, and statistical rollups from `Sim.SimResults` data. Components are pure (no data fetching) and composable -- apps wire them together with real data in Phase 4.

## How to add a new component

1. Create a directory under `src/` (kebab-case, e.g. `src/my-component/`)
2. Add `my-component.tsx` with the component, typed with `Sim.*` props from `@gcsim/types`
3. Add `my-component.test.tsx` with `@testing-library/react` tests
4. Add `index.ts` barrel export
5. Re-export from `src/index.ts`
6. Run `npx biome check --write packages/viewer/`
7. Run `turbo run typecheck test --filter=@gcsim/viewer`

## Canonical example

`src/result-cards/dps-card.tsx` -- per-character DPS card with proportional bar. Shows the pattern: typed props from `@gcsim/types`, primitives from `@gcsim/primitives`, data-testid attributes for testing, number formatting helper.

## Public API

All exports go through `src/index.ts`:

- **Metadata:** `Iterations`, `Mode`, `Commit`, `Warnings` -- small sub-components for sim metadata
- **TeamHeader** -- row of character cards showing name, level, constellation, weapon
- **RollupCard** -- generic stat rollup (mean, min, max, SD) for any `FloatStat`/`SummaryStat`
- **DPSCard** -- per-character DPS with proportional bar
- **TargetInfoCard** -- enemy info display (name, level, resistances)

Usage:
```typescript
import { DPSCard, RollupCard, TeamHeader, Iterations, Mode, Commit, Warnings } from "@gcsim/viewer";
```

## Dependencies

- `@gcsim/primitives` -- Card, Badge, cn() utility
- `@gcsim/types` -- `Sim.*` interfaces for all props
- `@gcsim/i18n` -- internationalization (for future use)

## Don'ts

- Don't fetch data in components -- they are purely presentational
- Don't import from other packages' `src/` -- only from their package index
- Don't test visual appearance -- test content rendering and behavior
- Don't use `@/` path aliases in source files -- use relative imports (aliases don't work in library builds)
- Don't add app-specific logic -- apps compose these components in Phase 4
