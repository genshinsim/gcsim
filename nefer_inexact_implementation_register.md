# Nefer Inexact Implementation Register

## Purpose

This document tracks only the currently live Nefer implementation points that are approximate, assumed, incomplete, or intentionally scaffolded.

It is not a changelog. Implemented work that is no longer meaningfully inexact should be recorded in [nefer_implementation_progress.md](nefer_implementation_progress.md) instead of repeated here.

## Current Overall State

- Canonical generated data exists and is in use.
- Core registration exists and is in use.
- Combat logic now covers the main gameplay loop, but several timing, geometry, passive, and constellation details are still not final.
- Anything listed below should not be treated as final behavior.

## Data and Pipeline

### [pipeline/pkg/data/dm/model.go](pipeline/pkg/data/dm/model.go)

- The loader supports multiple obfuscated `AvatarSkillDepot` passive-open field names.
- This is a compatibility fix based on observed dump structure, not a stable schema guarantee.
- If the datamine format changes again, this code may need another compatibility pass.

### [pipeline/pkg/data/avatar/avatar.go](pipeline/pkg/data/avatar/avatar.go)

- Passive scaling extraction now skips `proudSkillGroupId == 0`.
- This is a defensive compatibility behavior and assumes `0` means no passive scaling group should be parsed.

## Character State and Registration

### [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go)

- `AnimationYelanN0StartDelay = 10` is still a best-fit integration based on available sheet notes, not a fully validated engine-semantic confirmation.
- `maxVeilStacks` now supports the base 3-stack cap and the current C2 5-stack cap.
- Veil stacks now use sourced independent per-stack timers with a 9s base duration and the current C2 extension to 14s.
- `ascendantGleam` is currently inferred from `Moonsign >= 2`; broader moonsign-specific behavior is not fully implemented.
- Veil threshold buffs now retrigger when the capped third or fifth stack is refreshed.
- When new stacks are added at cap, the current implementation refreshes the oldest active stack first; this is a reasonable independent-timer model, but the exact refresh-target semantics are still not source-confirmed.

### [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go)

- [nefer_ingame_observations.md](nefer_ingame_observations.md) now records stronger in-game evidence that seeds should have a 12s lifetime and that conversion into a seed appears to reset lifetime.
- The current implementation models that 12s lifetime reset behavior by replacing a core with a newly created seed gadget on conversion, but the exact gameplay semantics of that reset are still inferred from observation rather than confirmed source data.
- Seed absorption currently consumes all nearby seeds found inside an assumed absorb radius on Charged Attack or Phantasm Performance.
- Seed absorption currently assumes a simplified effective radius instead of verified geometry.
- In-game observation also indicates a shared field limit of 5 across cores and seeds, with the oldest object disappearing or exploding first; the current implementation relies on the shared `GadgetTypDendroCore` limit path to approximate that behavior rather than a Nefer-specific queue.

### [internal/characters/nefer/seed_gadget.go](internal/characters/nefer/seed_gadget.go)

- Seed gadgets currently reuse `GadgetTypDendroCore` intentionally so they share the existing field cap with ordinary Dendro Cores.
- This matches the current in-game observation better than the previous counter model, but still depends on the engine's generic gadget-limit behavior rather than a Nefer-specific queue implementation.
- Seed gadgets currently do not model any separate collision, visibility, or target interaction beyond existing as absorbable world objects.

## Normal Attacks and Plunge

### [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go)

- `attackHitmarks`, `attackFrames`, `attackHitboxes`, and `attackOffsets` are not fully validated against final frame data.
- N3 attack chaining uses provisional timing derived from mixed measurements.
- The corrected sheet label is `N3 -> N4`, but the timing is still mixed and not final.
- Datamine now confirms lock-shape metadata for Nefer combat entries, but it still does not provide a full per-hit melee geometry specification for NA.
- Hitbox sizes and offsets are still current approximations rather than verified geometry.
- N3 two-hit scheduling exists, but exact intra-string timing and geometry still need validation.
- Plunge uses a generic circle hit and does not yet have verified Nefer-specific geometry, hitlag, or strike typing nuances.

## Charged Attack and Phantasm

### [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go)

- The basic Slither or Phantasm transition loop is implemented, including self-completing standalone special charges, but the exact release or finisher sequence and fine-grained CA timing are still not final.
- The observed three-way Charged Attack stamina split is now modeled: 50 stamina for ordinary CA, 25 stamina for ordinary CA during Shadow Dance, and 0 stamina for Phantasm Performance.
- `queuePhantasmPerformance()` still uses provisional frame timings for consume and hit sequence.
- Datamine exposes additional Nefer subSkill entries and lock shapes, which narrows the targeting picture, but Phantasm hit radii, ordering, ownership details, and spacing are still approximations.
- Seeds of Deceit absorption is now implemented against nearby world objects, but not with final geometry or validated per-hit targeting rules.
- C1 and C6-sensitive Phantasm branching is not fully implemented.

## Skill

### [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go)

- Datamine indicates Nefer skill-family entries use `CircleLockEnemyR8H6HC` and `CircleLockEnemyR15H10HC` targeting shapes for relevant skill/subSkill records.
- Current code still uses a simplified circular AoE and does not yet match datamine-informed targeting geometry closely enough to call final.
- Skill startup, cancel timings, and geometry are not fully validated.
- The current C2 implementation covers the 5-stack cap, the initial 2-stack Skill grant, the total +40% Phantasm Performance damage at 5 Veil stacks, and the +200 EM handling at the fifth stack or fifth-stack refresh.
- Skill now starts the current P1 conversion window, but the exact relationship between the 15s replacement window and any real seed lifetime after window end is still not fully verified.
- Skill particle generation is now implemented from Lunaris metadata, but still needs gameplay validation against actual proc timing and enemy-hit edge cases.

## Burst

### [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go)

- Burst is currently modeled as a two-hit scaffold with provisional hit timings.
- Datamine indicates Burst-family targeting uses `CircleLockEnemyR15H10HC`, so the current geometry is better constrained than before, but exact Burst hit timeline, geometry, and any additional nuances are still not finalized.
- Veil consumption is simplified to a direct per-stack damage bonus and does not yet cover every edge case from the final design.

## Passives and Constellations

### [internal/characters/nefer/asc.go](internal/characters/nefer/asc.go)

- Current implementation only covers the Lunar-Bloom EM bonus path.
- Seeds of Deceit conversion has moved into [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go), but the broader passive behavior is still incomplete.
- The hook applies to all direct Lunar-Bloom hits and still needs full validation against the intended source-specific behavior.

### [internal/characters/nefer/cons.go](internal/characters/nefer/cons.go)

- The current C2 logic is implemented through Veil capacity, Skill-side stack grant, duration extension, fifth-stack EM handling, and Phantasm Performance damage scaling.
- Only the C6 elevation hook is implemented beyond that.
- C1, C4, and the C6 extra-damage instances are still missing.
- Even the current C6 logic only covers elevation and not the full constellation behavior.

## Missing Mechanics

### Not yet implemented in the character package

- Fully verified Seeds of Deceit world geometry and absorb radius rules.
- Exact Veil refresh-target semantics beyond the current oldest-stack-refresh model.
- P1/P2 full behavior.
- C1 full behavior.
- C4 full Verdant Dew gain-rate behavior and linger handling.
- C6 extra damage instances.
- Final particle behavior validation.
- Final ICD, StrikeType, durability, poise, hitlag, and geometry pass.

## Validation Gaps From Source Data

### [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md)

- The workbook page inventory is now confirmed as 5 pages in workbook 1 and 3 pages in workbook 2.
- Additional `Skill`, `Burst`, `DashJump`, and `N0` page data has been folded into the summary.
- Mixed frame values are still unresolved.
- Many cancel routes remain missing.
- The corrected `N3 -> N4` row is still not fully confirmed because the measured values remain mixed.

### [nefer_readiness_assessment.md](nefer_readiness_assessment.md)

- This file is now historical source-analysis context only.
- Current live implementation gaps should be maintained in this register instead of mirrored back into the readiness snapshot.

## Exit Criteria For This Register

- Remove an item only when its code path is implemented and verified against reliable source data.
- If a placeholder becomes more precise but still not final, keep it in this document and update the wording instead of deleting it.
