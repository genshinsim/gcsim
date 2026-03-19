---
name: new-component
description: >
  Scaffold a new React component in a @gcsim/ package inside the ui-next
  monorepo. Use this skill whenever the user wants to create a new component,
  add a UI element, build a new widget, or set up a new React piece in any
  ui-next package. Also triggers for requests like "add a Button component to
  primitives" or "create a chart card in the ui package" or "I need a new
  component for X".
---

# /new-component

Scaffold a new React component inside an existing `ui-next/packages/` package.
Each component lives in its own directory with a co-located test and barrel
export, so the codebase stays navigable and every component is testable in
isolation.

## Input

Parse the user's message for:
- **package** (required): target package name (e.g., `primitives`, `ui`)
- **component** (required): component name in PascalCase (e.g., `Button`, `StatIcon`)
- **subdirectory** (optional): subdirectory within `src/` (e.g., `charts`, `result-cards`)

## Steps

1. Read the target package's `CLAUDE.md` at
   `ui-next/packages/<package>/CLAUDE.md` for the canonical example pattern.
   Adapt the component skeleton to match whatever conventions the package
   already follows -- this keeps the codebase consistent even as different
   packages evolve their own patterns.

2. Determine the component path:
   - Without subdirectory: `ui-next/packages/<package>/src/<Component>/`
   - With subdirectory: `ui-next/packages/<package>/src/<subdir>/<Component>/`

3. Create `<Component>.tsx`. Use a named export (not default) because the
   monorepo uses barrel re-exports and named exports make tree-shaking and
   refactoring straightforward:
```tsx
export interface <Component>Props {
  className?: string;
}

export function <Component>({ className }: <Component>Props) {
  return <div className={className}><Component></div>;
}
```

4. Create `<Component>.test.tsx`. Co-locating the test next to the source makes
   it easy to find and keeps the test running via the package's vitest config:
```tsx
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { <Component> } from "./<Component>";

describe("<Component>", () => {
  it("renders without crashing", () => {
    render(<<Component> />);
    expect(screen.getByText("<Component>")).toBeInTheDocument();
  });
});
```

5. Create `index.ts` in the component directory. This barrel file lets the
   package's top-level `src/index.ts` re-export without reaching into
   implementation files:
```typescript
export { <Component> } from "./<Component>";
export type { <Component>Props } from "./<Component>";
```

6. Update `ui-next/packages/<package>/src/index.ts` to add a re-export so
   consumers can import the new component from the package root:
```typescript
export * from "./<subdir?>/<Component>";
```
Use the Edit tool to append the export line. If a subdirectory is involved,
include it in the path (e.g., `export * from "./charts/SparkLine"`).

7. Verify the package still builds cleanly:
```bash
cd ui-next && pnpm --filter @gcsim/<package> typecheck && pnpm --filter @gcsim/<package> build
```
Both commands should pass. If either fails, fix the issue before finishing.
