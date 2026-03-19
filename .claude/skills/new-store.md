---
name: new-store
description: >
  Scaffold a typed Zustand store in a ui-next app. Use this skill whenever the
  user wants to create a new store, add state management, set up Zustand, or
  mentions needing a store for any ui-next application in apps/. Also use it
  when the user says /new-store.
---

# /new-store

Scaffolds a new Zustand store inside an existing `ui-next/apps/` application.

## Input

Parse the user's message for:
- **app** (required): app name (e.g., `web`)
- **store** (required): store name in camelCase (e.g., `simulator`, `viewer`, `settings`)
- **--persist** (optional flag): wire up localStorage persistence middleware

## Steps

1. **Convert store name to PascalCase** for type names (e.g., `simulator` → `Simulator`).

2. **Create `ui-next/apps/<app>/src/stores/<store>.ts`**:

Without `--persist`:
```typescript
import { create } from "zustand";

interface <Store>State {
  // TODO: add state fields
}

interface <Store>Actions {
  // TODO: add actions
}

type <Store>Store = <Store>State & <Store>Actions;

export const use<Store>Store = create<<Store>Store>()((set, get) => ({
  // TODO: initial state and action implementations
}));
```

With `--persist`:
```typescript
import { create } from "zustand";
import { persist } from "zustand/middleware";

interface <Store>State {
  // TODO: add state fields
}

interface <Store>Actions {
  // TODO: add actions
}

type <Store>Store = <Store>State & <Store>Actions;

export const use<Store>Store = create<<Store>Store>()(
  persist(
    (set, get) => ({
      // TODO: initial state and action implementations
    }),
    { name: "<store>-storage" },
  ),
);
```

3. **Verify**: Run `cd ui-next && pnpm --filter @gcsim/<app> typecheck` — should pass.
