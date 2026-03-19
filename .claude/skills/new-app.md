---
name: new-app
description: Scaffold a new Vite + React application in the ui-next monorepo. Use this skill whenever the user wants to create a new app, add a new frontend application, start a new SPA, or bootstrap a new project inside ui-next/apps/. Also applies when they mention spinning up a new dashboard, tool, or web interface in the monorepo.
---

# /new-app

Scaffold a complete Vite + React application in `ui-next/apps/`. This gives every new app the same foundation (TypeScript strict, Tailwind, React Query, routing, Zustand, i18n, Vitest) so that tooling, CI, and developer experience stay uniform across the monorepo.

## Input

Parse the user's message for:
- **name** (required) -- app name in lowercase, e.g. `web`, `db`
- **port** (required) -- dev server port, e.g. `5173`, `5174`

If either value is ambiguous, ask before proceeding. Check existing apps under `ui-next/apps/` to avoid port collisions.

## Steps

1. **Create the directory structure** at `ui-next/apps/<name>/src/pages/` and `ui-next/apps/<name>/src/stores/`.

2. **Create `package.json`** -- pin dependency versions to match the rest of the monorepo. Check `ui-next/DEPENDENCIES.md` for the latest pinned versions before writing; the versions below are examples only:
```json
{
  "name": "@gcsim/<name>",
  "version": "0.0.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc --build && vite build",
    "test": "vitest run",
    "typecheck": "tsc --noEmit",
    "lint": "biome check",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "19.2.4",
    "react-dom": "19.2.4",
    "@tanstack/react-query": "5.91.2",
    "@tanstack/react-router": "1.167.5",
    "zustand": "5.0.12",
    "i18next": "25.8.20",
    "react-i18next": "16.5.8"
  },
  "devDependencies": {
    "vite": "8.0.1",
    "@vitejs/plugin-react": "4.7.0",
    "tailwindcss": "4.2.2",
    "@tailwindcss/vite": "4.2.2",
    "typescript": "5.9.3",
    "vitest": "4.1.0",
    "@testing-library/react": "16.3.0",
    "@testing-library/jest-dom": "6.6.3",
    "@types/react": "19.2.4",
    "@types/react-dom": "19.2.4"
  }
}
```

3. **Create `tsconfig.json`** -- extend the shared base config so all apps share the same strict TypeScript settings:
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

4. **Create `vite.config.ts`** -- Tailwind is loaded as a Vite plugin rather than PostCSS for faster HMR:
```typescript
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: <port>,
  },
});
```

5. **Create `index.html`** in the app root (Vite uses this as the entry point):
```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>gcsim — <name></title>
    <link rel="icon" type="image/png" href="/favicon.png" />
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

6. **Create `src/main.tsx`** -- sets up React Query at the root so every page can use data fetching hooks:
```tsx
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { App } from "./app";
import "./app.css";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </StrictMode>,
);
```

7. **Create `src/app.tsx`**:
```tsx
export function App() {
  return (
    <div>
      <h1>@gcsim/<name></h1>
    </div>
  );
}
```

8. **Create `src/app.css`** -- the single Tailwind v4 import that enables utility classes everywhere:
```css
@import "tailwindcss";
```

9. **Create `src/routes.ts`** -- placeholder for lazy-loaded page routes (the `/new-page` skill populates this):
```typescript
// Route definitions for @gcsim/<name>
// Use /new-page to add lazy-loaded page imports and route config here.
```

10. **Create `.env`** by copying from `ui-next/.env.example` so environment variables are available during local dev.

11. **Create `CLAUDE.md`** -- this helps Claude understand the app's structure in future sessions:
```markdown
# @gcsim/<name>

## Routes

TODO: document route structure

## Stores

Zustand stores live in `src/stores/`. Use `/new-store` to scaffold.

## Consumed Packages

TODO: list @gcsim/ packages this app depends on

## Pages

Pages live in `src/pages/`. Use `/new-page` to scaffold.
```

12. **Create `vitest.config.ts`** -- merge with the shared base config so test settings stay consistent across apps:
```typescript
import { defineConfig, mergeConfig } from "vitest/config";
import baseConfig from "../../tooling/vitest/base.ts";

export default mergeConfig(baseConfig, defineConfig({}));
```

13. **Install dependencies** by running `pnpm install` from the `ui-next/` root. pnpm resolves workspace packages automatically.

14. **Verify** by running `cd ui-next && pnpm --filter @gcsim/<name> typecheck && pnpm --filter @gcsim/<name> build` -- both should pass with no errors.
