# Nefer Inexact Implementation Register

## Purpose

This document tracks every currently known Nefer implementation point that is approximate, assumed, incomplete, or intentionally left as a scaffold.

## Current Overall State

- Canonical generated data exists and is in use.
- Core registration exists and is in use.
- Combat logic is still only a partial vertical slice.
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
- `maxVeilStacks = 3` is only the base stack cap; full duration and refresh behavior for Veil of Falsehood is still not implemented.
- `ascendantGleam` is currently inferred from `Moonsign >= 2`; broader moonsign-specific behavior is not fully implemented.

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

- `basicChargeAttack()` is only a simple hit scaffold and does not model full Slither behavior.
- There is no real moving Slither state, no continuous stamina drain over time, and no verified exit-cost implementation.
- `ActionStam()` returns `0` during Shadow Dance with Verdant Dew, but this is only a simplified gate for the replacement CA path.
- `phantasmChargeAttack()` uses provisional frame timings for consume and hit sequence.
- Datamine exposes additional Nefer subSkill entries and lock shapes, which narrows the targeting picture, but Phantasm hit radii, ordering, ownership details, and spacing are still approximations.
- Seeds of Deceit absorption is not implemented here.
- C1, C2, and C6-sensitive Phantasm branching is not fully implemented.

## Skill

### [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go)

- Datamine indicates Nefer skill-family entries use `CircleLockEnemyR8H6HC` and `CircleLockEnemyR15H10HC` targeting shapes for relevant skill/subSkill records.
- Current code still uses a simplified circular AoE and does not yet match datamine-informed targeting geometry closely enough to call final.
- Skill startup, cancel timings, and geometry are not fully validated.
- Shadow Dance duration is implemented, but the full state interaction model with CA movement and transitions is incomplete.
- The current C2 interaction only grants 2 Veil stacks when conditions are met; the rest of C2 is not implemented.

## Burst

### [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go)

- Burst is currently modeled as a two-hit scaffold with provisional hit timings.
- Datamine indicates Burst-family targeting uses `CircleLockEnemyR15H10HC`, so the current geometry is better constrained than before, but exact Burst hit timeline, geometry, and any additional nuances are still not finalized.
- Veil consumption is simplified to a direct per-stack damage bonus and does not yet cover every edge case from the final design.

## Passives and Constellations

### [internal/characters/nefer/asc.go](internal/characters/nefer/asc.go)

- Current implementation only covers the Lunar-Bloom EM bonus path.
- Seeds of Deceit conversion and broader passive behavior are not implemented here.
- The hook applies to all direct Lunar-Bloom hits and still needs full validation against the intended source-specific behavior.

### [internal/characters/nefer/cons.go](internal/characters/nefer/cons.go)

- Only the C6 elevation hook is implemented.
- C1, C2, C4, and the C6 extra-damage instances are still missing.
- Even the current C6 logic only covers elevation and not the full constellation behavior.

## Missing Mechanics

### Not yet implemented in the character package

- Seeds of Deceit replacement logic.
- Seeds of Deceit absorption logic.
- Veil of Falsehood duration and refresh rules.
- 3-stack and 5-stack Veil trigger handling beyond the current partial grant path.
- P1/P2 full behavior.
- C1 full behavior.
- C2 full behavior beyond the initial stack grant.
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

- This file still represents the authoritative list of missing confirmations and should be kept in sync with code changes.

## Exit Criteria For This Register

- Remove an item only when its code path is implemented and verified against reliable source data.
- If a placeholder becomes more precise but still not final, keep it in this document and update the wording instead of deleting it.
