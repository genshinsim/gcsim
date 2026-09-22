# 2. Per-run series colors over present categories; enemies keyed by roster index

Date: 2026-09-22

## Status

Accepted

## Context

The shared series palette (`--g-s1..--g-s7`) is only **7** colorblind-safe
Okabe-Ito colors (see ADR 0001-adjacent work, issue #3036). Charts hand these
colors to categories, and several category sets are larger than 7.

The Actions chart drew from a **fixed 11-name domain** (`attack, charge, aim,
skill, burst, low_plunge, high_plunge, dash, jump, walk, swap`) and colored each
by its **fixed position** in that list: `series(i) = --g-s{(i % 7)+1}`. So `dash`
(index 7) always reused `attack`'s color and `swap` (index 10) reused `skill`'s,
no matter how few actions a run actually used. A run of Normal Attack, Elemental
Skill, Elemental Burst, Sprint (dash) and Character Swap (swap) showed only **3
distinct colors for 5 categories** (#3064).

Separately, enemy identity flows to two surfaces by **parallel keying**: the
overview enemy card used `qualitative3(enemyRosterIndex)` while the target DPS
charts used `q3(targetId - 1)`. These land on the same slot only because
`targetId - 1 === enemyRosterIndex` by construction — a contract the Go side
flags as fragile (`// TODO: subject to break if target key gen changes`).

## Decision

**Actions: assign colors per run, over the categories actually present, in
canonical order.** The Actions chart already derives the present action set in
canonical order; color assignment now keys off **rank within that present set**
(first present → `s1`, second → `s2`, …) instead of fixed position in the
11-name list. With ≤ 7 distinct actions present (the common case) every category
is a distinct color. Past 7, assignment **wraps** (8th → `s1`) — accepted, not a
failure, since only 7 colorblind-safe colors exist. The present set is therefore
an **input** to the assignment (`action`/`actionLabel`/`actionHighlight` take the
present, canonically-ordered keys and return a rank-based scale).

**Enemies/targets: one shared accessor keyed by absolute roster index.** Both the
overview enemy card and the target DPS charts route through a single
`DataColorsConst.enemy(rosterIndex)` (base color = `q3(rosterIndex)`), so they
cannot silently drift. The target accessors convert `targetId → rosterIndex`
(`targetId - 1`) and delegate the color to that one accessor. Enemy colors stay
keyed by **absolute roster position** and wrap past 7 enemies (rare); the
damaged-target subset is **not** dense-ranked per chart, which would make an
enemy change color between the overview and the DPS chart.

`SERIES_COUNT` / `series()` remain the single wrap point.

## Consequences

- The #3064 repro now renders 5 distinct colors for its 5 actions.
- Global, stable-per-category coloring is **not** possible with only 7 colors;
  a category's color can differ between runs (it depends on which siblings are
  present). This is the accepted trade for distinctness within a run.
- An enemy is guaranteed the same base color across the overview card and the
  target DPS charts, independent of the Go-side target key generation.
- Only the **base** color slot is shared for enemies; the label/text shade tiers
  (enemy card `q5`, target labels `q4`) are intentionally left un-unified.
- Element colors, reactable-modifier colors, the `character` scale, and the
  hand-picked `qualitative*` indices in the timeline / cumulative-damage cards
  are unchanged.
