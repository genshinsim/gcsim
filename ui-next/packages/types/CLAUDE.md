# @gcsim/types

## Purpose

Provides all TypeScript types for the gcsim UI:

1. **Proto-generated types** from `protos/` at repo root, generated via `buf` + `protoc-gen-es`
2. **Custom interfaces** in `src/sim.ts` for JSON result shapes returned by the simulator

## How to regenerate proto types

```bash
cd ui-next/packages/types
pnpm generate
```

This runs `buf generate ../../..` which reads proto files from `protos/` at the repo root and outputs TypeScript files into `src/generated/`.

Proto source files live in:
- `protos/model/` — SimulationResult, Character, Enemy, enums, etc.
- `protos/backend/` — ShareEntry, DBEntry, Preview types

## Public API

All public exports go through `src/index.ts`:

- **Top-level exports**: model proto types (SimulationResult, Character, etc.)
- **`share` namespace**: backend share proto types
- **`db` namespace**: backend DB proto types
- **`preview` namespace**: backend preview proto types
- **`Sim` namespace**: custom JSON interfaces (SimResults, Statistics, etc.)

Usage:
```typescript
import { SimulationResultSchema, type Sim } from "@gcsim/types";
const result: Sim.SimResults = { ... };
```

## Testing

Tests verify that:
- Generated proto schemas are defined and have correct type names
- Custom sim interfaces compile and are usable
- Backend namespace exports are accessible

Run tests: `pnpm test` or `turbo run test --filter=@gcsim/types`

## Don'ts

- Don't import from other packages' `src/` — only from their package index
- Don't manually edit files in `src/generated/` — they are overwritten by `pnpm generate`
- Don't add app-specific logic here — this package is types only
