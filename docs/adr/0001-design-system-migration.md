# 1. Design-system migration: deprecate-then-rebuild off shadcn

Date: 2026-09-18

## Status

Accepted

## Context

The UI's semantic design tokens are the stock shadcn set (`background`,
`foreground`, `primary`, `muted`, `accent`, `destructive`, `border`, `ring`,
`radius`, …). They are defined as raw CSS custom properties in a `:root` + `.dark`
block (two byte-identical copies, one per package entrypoint), mapped through the
shared Tailwind theme (`packages/tailwind.css`, via `@theme inline`) into
Tailwind's `--color-*` / `--radius-*` namespaces, and consumed as utilities
(`bg-primary`, `text-muted-foreground`, `border-border`, …) across every package.
The shadcn CLI config (`components.json`) is checked in to three packages, so a
future `shadcn add` would regenerate a component and reintroduce these tokens.

We are migrating to the **Gauge** design language. We want the new vocabulary
without a big-bang rewrite that would break every component at once, and without
leaving the old shadcn tokens as an attractive nuisance that new or regenerated
code keeps reaching for.

## Decision

Migrate **deprecate-then-rebuild**, in three steps:

1. **Quarantine** (this step). Rename every shadcn semantic token to a
   `--deprecated-<name>` form — at its definition, its `@theme` registration, and
   every consuming utility (`bg-primary` → `bg-deprecated-primary`) — moving in
   lockstep so values are unchanged and rendered output is identical. Remove the
   shadcn CLI (`components.json`) everywhere. A checked-in gate
   (`ui/scripts/check-no-shadcn-tokens.mjs`, wired into `lint`/`lint-ci`) fails if
   any bare shadcn token or utility reappears.
2. **Adopt** the Gauge handoff token vocabulary as the new, un-prefixed tokens
   (separate ticket).
3. **Rebuild** each component against the Gauge tokens, then delete the
   `--deprecated-` tokens once nothing consumes them (separate ticket).

We **keep** the parts of the shadcn stack that are just good React/Tailwind
practice and carry no design opinion: Radix primitives, CVA variants, and the
`cn()` class-merge helper. Only the shadcn CLI and (eventually) its tokens go.

Radius is a special case: shadcn overrode `rounded-lg/md/sm` to values identical
to Tailwind v4's defaults, so those overrides are dropped and `rounded-*` falls
back to Tailwind core unchanged. Only the raw `--radius` custom property is
carried forward as `--deprecated-radius`.

## Consequences

- The current components keep working unchanged through Step 1 — a pure rename,
  zero visual change, script-verifiable.
- New code cannot silently depend on shadcn tokens: the gate blocks bare names,
  forcing an explicit `--deprecated-` opt-in that reads as debt.
- The two duplicate token blocks are **not** de-duplicated here; that stays a
  known wart until the rebuild removes them.
- Removing `components.json` means `shadcn add` no longer works; new primitives
  are hand-authored (or vendored) against the Gauge tokens, which is the intent.
