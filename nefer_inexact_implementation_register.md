# Nefer Inexact Implementation Register

## Rule

- One bullet below equals one unresolved frame row, one unresolved hitbox or area definition, one unresolved metadata mapping, or one unresolved semantic assumption.

## [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go)

- `AnimationYelanN0StartDelay = 10` uses workbook 2 Yelan row `10`, but workbook 2 does not identify whether `10` maps to `info.AnimationYelanN0StartDelay` rather than another Yelan N0 hook.
- No Nefer-specific Xingqiu N0 override is implemented in the package; workbook 2 summary suggests `11`, but workbook 2 does not identify which gcsim hook that `11` should drive.

## [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go): unresolved frame rows

- `Normal 1 hitmark` currently uses `10`; workbook 1 row is mixed `10-13`.
- `Normal 2 hitmark` currently uses `8`; workbook 1 row is mixed `8-11`.
- `Normal 4 hitmark` currently uses `22`; workbook 1 row is mixed `22-24`.
- `Normal 1 -> Normal CA` currently uses `40`; workbook 1 row is mixed `40-41` including CA windup.
- `Normal 1 -> Walk` currently uses `33`; workbook 1 row is mixed `32-33`.
- `Normal 1 -> Skill` has no workbook 1 row.
- `Normal 1 -> Burst` has no workbook 1 row.
- `Normal 1 -> Dash` has no workbook 1 row.
- `Normal 1 -> Jump` has no workbook 1 row.
- `Normal 1 -> Swap` has no workbook 1 row.
- `Normal 2 -> Normal 3` currently uses `20`; workbook 1 row is mixed `18-21`.
- `Normal 2 -> Phantasm CA` currently shares `ActionCharge = 38` with ordinary CA; workbook 1 row is mixed `38-39`.
- `Normal 2 -> Skill` has no workbook 1 row.
- `Normal 2 -> Burst` has no workbook 1 row.
- `Normal 2 -> Dash` has no workbook 1 row.
- `Normal 2 -> Jump` has no workbook 1 row.
- `Normal 2 -> Swap` has no workbook 1 row.
- `Normal 3 -> Normal 4` currently uses `49`; workbook 1 row is mixed `48-53`.
- `Normal 3 -> Normal CA` currently uses `62`; workbook 1 row is mixed `61-62` including CA windup.
- `Normal 3 -> Skill` has no workbook 1 row.
- `Normal 3 -> Burst` has no workbook 1 row.
- `Normal 3 -> Dash` has no workbook 1 row.
- `Normal 3 -> Jump` has no workbook 1 row.
- `Normal 3 -> Swap` has no workbook 1 row.
- `Normal 4 -> Normal 1` currently uses `51`; workbook 1 row is mixed `50-54`.
- `Normal 4 -> Walk` currently uses `63`; workbook 1 row is mixed `62-64`.
- `Normal 4 -> Skill` has no workbook 1 row.
- `Normal 4 -> Burst` has no workbook 1 row.
- `Normal 4 -> Dash` has no workbook 1 row.
- `Normal 4 -> Jump` has no workbook 1 row.
- `Normal 4 -> Swap` has no workbook 1 row.

## [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go): unresolved hitboxes and areas

- `Normal 1 hitbox` currently uses `NewBoxHit(..., Y=0, width=2, length=8)`; no source-backed per-hit melee geometry has been identified.
- `Normal 2 hitbox` currently uses `NewBoxHit(..., Y=0, width=2, length=8)`; no source-backed per-hit melee geometry has been identified.
- `Normal 3 hit 1 hitbox` currently uses `NewBoxHit(..., Y=0, width=2.5, length=9)`; no source-backed per-hit melee geometry has been identified.
- `Normal 3 hit 2 hitbox` currently reuses `NewBoxHit(..., Y=0, width=2.5, length=9)` from `Normal 3 hit 1`; no source-backed per-hit melee geometry has been identified.
- `Normal 4 hitbox` currently uses `NewBoxHit(..., Y=-0.5, width=2.8, length=8)`; no source-backed per-hit melee geometry has been identified.
- `Low Plunge hitbox` currently uses `NewCircleHitOnTarget(..., radius=4.5)`; no Nefer-specific plunge AoE geometry has been identified.
- `High Plunge hitbox` currently uses `NewCircleHitOnTarget(..., radius=4.5)`; no Nefer-specific plunge AoE geometry has been identified.

## [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go): unresolved frame rows

- `Charge Attack -> Normal 1` currently uses `ActionAttack = 29` under the accepted branch rule that ordinary non-Phantasm CA is the full no-hold Slither-release route and workbook 1 ordinary-CA rows are measured button-to-button across that whole route; no source-backed sub-breakdown for the embedded Slither-entry segment has been identified.
- `Charge Attack -> Charge Attack` currently uses `ActionCharge = 48`; workbook 1 row is mixed `48-49`.
- `Charge Attack -> Skill` currently uses generic `48`; workbook 1 has no row.
- `Charge Attack -> Burst` currently uses generic `48`; workbook 1 has no row.
- `Charge Attack -> Dash` currently uses generic `48`; workbook 1 has no row.
- `Charge Attack -> Jump` currently uses generic `48`; workbook 1 has no row.
- `Charge Attack -> Swap` currently uses `28` under the accepted branch rule that ordinary non-Phantasm CA is the full no-hold Slither-release route and workbook 1 ordinary-CA rows are measured button-to-button across that whole route; no source-backed sub-breakdown for the embedded Slither-entry segment has been identified.
- `Phantasm Performance hit 4` currently uses frame `44`; workbook 1 row is mixed `44-45` after adding the 20f CA windup to the mixed `24-25` sub-row.
- `Phantasm Performance hit 5` currently uses frame `45`; workbook 1 row allows only `44-45` after adding the 20f CA windup to the mixed `0-1` sub-row.
- `Phantasm Performance -> Skill` currently falls back to the generic post-Phantasm branch because workbook 1 provides no source-backed row.
- `Phantasm Performance -> Burst` currently falls back to the generic post-Phantasm branch because workbook 1 provides no source-backed row.
- `Phantasm Performance -> Dash` currently falls back to the generic post-Phantasm branch because workbook 1 provides no source-backed row.
- `Phantasm Performance -> Jump` currently falls back to the generic post-Phantasm branch because workbook 1 provides no source-backed row.
- `Phantasm Performance -> Walk` currently follows `phantasmAnimationLength = 106`; workbook 1 row is mixed `104-105`.
- `Phantasm Performance -> Swap` is not encoded as a dedicated `0f` row; workbook 1 confirms `0`.
- `C6 post-Phantasm extra hit timing` currently uses `phantasmAnimationLength = 106`; no source-backed row for the extra hit timing has been identified.
- `Slither entry stamina threshold` currently uses `slitherActivationFrames = 60` times the sourced `18.15/s` drain, interpreted as roughly one second of available Slither stamina; the branch now enforces that as a readiness requirement without consuming it up front, but no source-backed activation threshold row has been identified.
- `Slither minimum cancel frame` currently uses `24`; no source-backed row for this exact cancel gate has been identified.
- `Slither movement cadence` currently uses `slitherMoveInterval = 1`; no source-backed row for this exact cadence has been identified.
- `Slither movement distance per tick` currently uses `slitherMoveDistance = 0.1`; no source-backed distance mapping has been identified.
- `Slither stamina drain cadence` currently uses `slitherStamInterval = 1`; no source-backed per-tick cadence has been identified.

## [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go): unresolved hitboxes and areas

- `Charge Attack hitbox` currently uses `NewBoxHit(..., Y=-2, width=3, length=9)`; no source-backed release hitbox geometry has been identified.
- `Phantasm Performance (Nefer 1) hitbox` currently uses `NewCircleHitOnTarget(..., radius=5)`; no source-backed Phantasm AoE radius has been identified.
- `Phantasm Performance (Shade 1) hitbox` currently uses `NewCircleHitOnTarget(..., radius=5)`; no source-backed Phantasm AoE radius has been identified.
- `Phantasm Performance (Shade 2) hitbox` currently uses `NewCircleHitOnTarget(..., radius=5)`; no source-backed Phantasm AoE radius has been identified.
- `Phantasm Performance (Nefer 2) hitbox` currently uses `NewCircleHitOnTarget(..., radius=5)`; no source-backed Phantasm AoE radius has been identified.
- `Phantasm Performance (Shade 3) hitbox` currently uses `NewCircleHitOnTarget(..., radius=5)`; no source-backed Phantasm AoE radius has been identified.
- `Phantasm Performance (C6 End Hit) hitbox` currently uses `NewCircleHitOnTarget(..., radius=5)`; no source-backed C6 extra-hit AoE radius has been identified.

## [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go): unresolved frame rows

- `Skill -> Normal CA` currently shares `ActionCharge = 29`; workbook 1 row is mixed `47-49` including CA windup.
- `Skill -> Phantasm CA` currently shares `ActionCharge = 29`; workbook 1 row is mixed `28-29` without CA windup.
- `Skill -> Skill` currently falls back to the generic slice value from `InitAbilSlice(52)`; workbook 1 row is mixed `69-74`.
- `Skill -> Burst` currently uses `38`; workbook 1 row is mixed `26-28`.
- `Skill -> Dash` currently uses `38`; workbook 1 row is mixed `37-39`.
- `Skill -> Swap` currently uses `25`; workbook 1 row is mixed `24-25`.

## [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go): unresolved hitboxes, areas, and particles

- `Senet Strategy: Dance of a Thousand Nights hitbox` currently uses `NewCircleHit(..., Y=1, radius=3)`; datamine points to skill-family lock shapes `CircleLockEnemyR8H6HC` and `CircleLockEnemyR15H10HC`, but the exact mapping from those shapes to this hit has not been closed.
- `Skill particle proc timing` currently uses the skill hit callback on frame `24`; no in-game validation has been recorded for whether particle generation is tied to that exact contact frame in every case.
- `Skill particle proc gating` currently uses `particleICDKey` with `0.2s`; no in-game validation has been recorded for multi-target or repeated-contact edge cases.
- `Skill particle count` currently uses `2` with `34%` and `3` with `66%`; no in-game validation has been recorded for that exact distribution.

## [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go): unresolved frame rows

- `Sacred Vow: True Eye's Phantasm (Hit 1) hitmark` currently uses `26`; workbook 1 row is mixed `99-101`.
- `Sacred Vow: True Eye's Phantasm (Hit 2) hitmark` currently uses `46`; workbook 1 row is mixed `40-44`.
- `Burst energy drain timing` currently consumes energy immediately through `ConsumeEnergy(60)`; workbook 1 row is mixed around `6-7` and includes an outlier.
- `Burst -> Normal 1` currently uses `119`; workbook 1 row is mixed `118-119`.
- `Burst -> Normal CA` currently uses `131` when the ordinary CA route is selected; workbook 1 row is mixed `129-131` including CA windup.
- `Burst -> Phantasm CA` currently uses `124` when the Phantasm route is selected; workbook 1 row is mixed `123-124` including CA windup.
- `Burst -> Dash` currently uses `120`; workbook 1 row is mixed `119-120`.
- `Burst -> Jump` currently uses `120`; workbook 1 row is mixed `119-120`.
- `Burst -> Walk` currently uses `119`; workbook 1 row is mixed `118-119`.
- `Burst -> Swap` currently uses `118`; workbook 1 row is mixed `117-118`.

## [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go): unresolved hitboxes and areas

- `Sacred Vow: True Eye's Phantasm (Hit 1) hitbox` currently uses `NewBoxHit(..., width=6, length=10)`; datamine points to `CircleLockEnemyR15H10HC`, but the exact mapping from that lock shape to this hitbox has not been closed.
- `Sacred Vow: True Eye's Phantasm (Hit 2) hitbox` currently reuses `NewBoxHit(..., width=6, length=10)` from Hit 1; datamine points to `CircleLockEnemyR15H10HC`, but the exact mapping from that lock shape to this hitbox has not been closed.

## [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go) and [internal/characters/nefer/seed_gadget.go](internal/characters/nefer/seed_gadget.go)

- `Seed absorption area` currently uses `seedAbsorbRadius = 6` centered on the player through `NewCircleHitOnTarget`; no source-backed absorb radius has been identified.
- `Seed absorption target-selection rule` currently absorbs every `seedGadget` inside that player-centered circle; no source-backed rule has been identified for whether vertical separation, line-of-travel, or another target-selection filter should apply.
- `Core/seed shared field-cap removal order` currently relies on the generic `GadgetTypDendroCore` engine path after converting cores into `seedGadget`; no source-backed Nefer-specific queue implementation has been identified for mixed core/seed eviction order.
- `Seed gadget collision and target interaction` currently have no Nefer-specific implementation beyond existing as absorbable gadgets; no source-backed collision or interactability rule has been identified.

## [internal/characters/nefer/cons.go](internal/characters/nefer/cons.go)

- `C4 nearby-opponent area` currently uses `c4NearbyRadius = 10` plus enemy circle radius; no source-backed area definition for “nearby opponents” has been identified.

## [internal/characters/nefer/*.go]: explicit poise and hitlag mappings still missing

- `Normal 1` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Normal 2` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Normal 3 hit 1` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Normal 3 hit 2` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Normal 4` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Low Plunge` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `High Plunge` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Charge Attack` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Senet Strategy: Dance of a Thousand Nights` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Phantasm Performance (Nefer 1)` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Phantasm Performance (Shade 1)` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Phantasm Performance (Shade 2)` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Phantasm Performance (Nefer 2)` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Phantasm Performance (Shade 3)` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Phantasm Performance (C6 End Hit)` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Sacred Vow: True Eye's Phantasm (Hit 1)` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
- `Sacred Vow: True Eye's Phantasm (Hit 2)` has no explicit `PoiseDMG`, `HitlagFactor`, or `HitlagHaltFrames`.
