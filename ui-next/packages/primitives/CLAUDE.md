# @gcsim/primitives

## Purpose

Design system foundation for the gcsim web UI. Contains:
- Tailwind v4 `@theme` design tokens with OKLCH colors (`theme.css`)
- `cn()` class merge utility (clsx + tailwind-merge)
- shadcn/ui primitive components installed via CLI

## How to add a new primitive

Use the shadcn CLI from this package directory:
```bash
cd ui-next/packages/primitives
pnpm dlx shadcn@4.1.0 add <component-name> -y -s
```

Then add the export to `src/index.ts`.

## Canonical example

`src/components/ui/button.tsx` — shadcn Button with CVA variants, Radix Slot for `asChild`

## Public API

All public exports go through `src/index.ts`:
- `cn()` — class merge utility (from `src/lib/utils.ts`)
- 11 component families: Button, Card, Input, Tabs, Select, Badge, Dialog, DropdownMenu, Tooltip, ScrollArea, Skeleton

Apps import the theme CSS separately: `@import '@gcsim/primitives/theme.css';`

## Dependencies

- `radix-ui` — accessible component primitives (unified package)
- `class-variance-authority` — component variants
- `clsx` + `tailwind-merge` — class merging
- `lucide-react` — icons
- `shadcn` — component registry

## Theme

Uses Tailwind v4 CSS-first `@theme inline` configuration with OKLCH colors.
Custom tokens include Genshin element colors (anemo, geo, electro, hydro, pyro, cryo, dendro).

## Don'ts

- Don't import from other packages' `src/` — only from their package index
- Don't add app-specific logic here
- Don't hardcode colors — use theme tokens
- Don't use JS-based Tailwind config — theme is CSS-first via `@theme`
- Don't hand-write shadcn components — use the CLI to install them
