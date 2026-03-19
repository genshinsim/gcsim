# ui-next/apps — Application Layer

Apps are the composition layer. They wire packages together into deployable web applications.

## Responsibilities

- Own **pages**, **stores**, and **route definitions**
- Import components and logic from `@gcsim/` packages via their public API
- Define app-specific layouts, navigation, and error boundaries
- Each app has its own `.env` file (see root `.env.example` for template)

## Rules

- **Never import from another app** — apps are independent
- **Never import from internal package paths** — only from `@gcsim/<pkg>` index
- **Use error boundaries** around data-fetching and rendering sections
- **Use lazy imports** for all page components (route-level code splitting)
- **Use TanStack Query** for server data — no raw `fetch` + `useEffect`
- **Use Zustand** for client state — with `persist` middleware for localStorage

## Dev Server Ports

| App | Port | Domain |
|-----|------|--------|
| web | 5173 | gcsim.app |
| db | 5174 | db.gcsim.app |
| embed | 5175 | (screenshot generator) |
| storybook | 6006 | (component dev) |

These ports avoid conflicts with the old `ui/` dev server.

## Creating a New App

Use the `/new-app` skill: it scaffolds `package.json`, `vite.config.ts`, `tsconfig.json`, `main.tsx`, `app.tsx`, `app.css`, `routes.ts`, `.env`, and `CLAUDE.md`.

## Creating a New Page

Use the `/new-page` skill: it scaffolds the page component with a test file and adds a lazy-loaded route entry to `routes.ts`.

## Creating a New Store

Use the `/new-store` skill: it scaffolds a typed Zustand store with optional `persist` middleware.
