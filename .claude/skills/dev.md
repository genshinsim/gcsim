---
name: dev
description: >
  Start Vite dev server(s) for a ui-next application. Use this skill whenever
  the user wants to run a dev server, start local development, preview an app,
  or launch any ui-next app for development. Also use it when the user says
  /dev.
---

# /dev

Starts the Vite dev server for a ui-next application.

## Input

Parse the user's message for:
- **app** (required): app name (e.g., `web`) or `"all"` to start all apps

## Steps

1. **Build dependencies first** — packages that the app depends on need to be compiled:
```bash
cd ui-next && turbo run build --filter=@gcsim/<app>^...
```
For "all": `cd ui-next && turbo run build --filter=./packages/*`

2. **Start the dev server**:
```bash
cd ui-next && pnpm --filter @gcsim/<app> dev
```
For "all": `cd ui-next && turbo run dev`

Run the dev server command in the background using the Bash tool's `run_in_background` parameter so the user can continue working.

## Known App Ports

| App | Port |
|-----|------|
| web | 5173 |

## Output

Report which app is starting and on which port. The dev server runs in the background — inform the user they can continue working.
