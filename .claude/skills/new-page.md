---
name: new-page
description: Scaffold a new page in a ui-next app with lazy-loaded routing. Use this skill whenever the user wants to add a page, screen, or view to an existing ui-next application, including when they say things like "add a simulator page", "create a new route", or "I need a new view for X".
---

# /new-page

Scaffold a new page inside an existing `ui-next/apps/` application. Every page follows the same structure so that routing, code-splitting, and testing stay consistent across the monorepo.

## Input

Parse the user's message for:
- **app** (required) -- app name, e.g. `web`
- **page** (required) -- page name in kebab-case, e.g. `simulator`, `viewer`, `sample-detail`

If either value is ambiguous, ask before proceeding.

## Steps

1. **Derive the PascalCase component name** from the kebab-case page name (e.g. `simulator` becomes `Simulator`, `sample-detail` becomes `SampleDetail`). This keeps file names consistent (kebab) while following React conventions (PascalCase) for components.

2. **Create the page directory** at `ui-next/apps/<app>/src/pages/<page>/`.

3. **Create `<page>.tsx`** -- the page component with a default export so `React.lazy()` can resolve it:
```tsx
export default function <PageComponent>() {
  return (
    <div>
      <h1><PageComponent></h1>
    </div>
  );
}
```

4. **Create `<page>.test.tsx`** -- a smoke test that confirms the component renders:
```tsx
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import <PageComponent> from "./<page>";

describe("<PageComponent>", () => {
  it("renders without crashing", () => {
    render(<<PageComponent> />);
    expect(screen.getByText("<PageComponent>")).toBeInTheDocument();
  });
});
```

5. **Register the route** in `ui-next/apps/<app>/src/routes.ts`. Add a lazy import and a route entry. Lazy-loading is important because it keeps the initial bundle small -- each page is only fetched when the user navigates to it.
```typescript
import { lazy } from "react";

const <PageComponent> = lazy(() => import("./pages/<page>/<page>"));
```
Use the Edit tool to insert the new import at the top and the route entry into the existing routes array/config.

6. **Verify** by running `cd ui-next && pnpm --filter @gcsim/<app> typecheck` -- it should pass with no errors.
