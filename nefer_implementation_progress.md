# Nefer Implementation Progress

## Purpose

This file is a dated changelog of completed Nefer work on the branch.

It should not be used as the source of current readiness or open-gap status. For that, use [Nefer PR checklist.md](Nefer%20PR%20checklist.md), [nefer_readiness_assessment.md](nefer_readiness_assessment.md), and [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).

## Completed Work Log

### 2026-03-15

#### Phase 4. Nefer-specific mechanics

- Added a Seeds of Deceit conversion window in [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go).
- Added conversion of existing on-field Dendro Cores into world-space Seeds of Deceit on Skill use in [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go).
- Added conversion of newly spawned Dendro Cores into Seeds of Deceit during the active window by intercepting [pkg/core/event/event.go](pkg/core/event/event.go) `OnDendroCore` in [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go).
- Added a dedicated seed gadget in [internal/characters/nefer/seed_gadget.go](internal/characters/nefer/seed_gadget.go) using the shared Dendro Core gadget type so cores and seeds respect the same cap behavior.
- Added a 12s seed lifetime in [internal/characters/nefer/seed_gadget.go](internal/characters/nefer/seed_gadget.go).
- Added lifetime reset on core-to-seed conversion by replacing the original core gadget with a newly created seed gadget in [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go).
- Added nearby seed absorption on Charged Attack and Phantasm Performance in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go).
- Replaced the previous Shadow Dance no-Dew Charged Attack placeholder with a persistent Slither state in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go).
- Added Slither movement ticks, ordinary running-style continuous stamina drain, and continuous nearby seed absorption during the active CA state in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go).
- Added swap reset handling for Slither and Phantasm Performance Charges in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go), while preserving shared Verdant Dew.
- Added separate Phantasm Performance Charge tracking in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go) and refresh-on-Skill behavior in [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go).
- Reworked Shadow Dance Charged Attack into a loop that alternates between Slither and queued Phantasm Performance based on both Phantasm Performance Charges and Verdant Dew in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go).
- Adjusted Nefer charge action completion so one scripted `charge` resolves after one produced normal or special Charged Attack result, including standalone special charge completion without requiring a following scripted action, in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go).
- Added initial Veil threshold EM buff handling for 3-stack and base C2 5-stack thresholds in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go).
- Added C2 Phantasm Performance damage scaling from current Veil stacks and fifth-stack EM refresh handling in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go) and [internal/characters/nefer/cons.go](internal/characters/nefer/cons.go).
- Added sourced independent Veil stack timers with the C2 duration extension in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go).
- Added initial Skill particle generation in [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go) using the current Lunaris 66%/33% split with a 0.2s ICD.
- Corrected Slither stamina drain to the sourced 18.15/s value and split Charged Attack stamina handling into the three observed modes in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go): normal 50 stamina, Shadow Dance normal 25 stamina, and Phantasm Performance 0 stamina.
- Corrected swap handling so leaving the field clears Shadow Dance itself in addition to resetting Phantasm Performance Charges in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go).
- Added regression coverage for the three Charged Attack stamina modes and for swap-driven Shadow Dance or charge reset behavior in [internal/characters/nefer/charge_test.go](internal/characters/nefer/charge_test.go).

#### Validation

- Confirmed the updated package still compiles with `go test ./internal/characters/nefer ./pkg/simulation`.

### 2026-03-14

#### Phase 1. Data layer and generation

- Created [internal/characters/nefer/config.yml](internal/characters/nefer/config.yml).
- Verified `genshin_id: 10000122`.
- Validated Nefer talent mapping with `extract-talents` against the updated datamine.
- Added pipeline compatibility for obfuscated `AvatarSkillDepot` passive-open fields in [pipeline/pkg/data/dm/model.go](pipeline/pkg/data/dm/model.go) and [pipeline/pkg/data/avatar/avatar.go](pipeline/pkg/data/avatar/avatar.go).
- Ran canonical character generation successfully on the updated datamine.
- Confirmed Nefer appears in generated keys, imports, UI/DB/docs outputs.
- Switched Nefer to the canonical generated data layer by using [internal/characters/nefer/nefer_gen.go](internal/characters/nefer/nefer_gen.go) as the source of generated stats and scalings.

#### Phase 2. Character skeleton and registration

- Added Nefer registration in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go).
- Confirmed generated key registration in [pkg/core/keys/keys_char_gen.go](pkg/core/keys/keys_char_gen.go).
- Confirmed generated simulation import in [pkg/simulation/imports_char_gen.go](pkg/simulation/imports_char_gen.go).
- Set `EnergyMax = 60`, `SkillCon = 3`, `BurstCon = 5`, `NormalHitNum = 4`.
- Added `SetNumCharges(action.ActionSkill, 2)` in the character skeleton.

#### Phase 3. Base combat scaffold

- Implemented a 4-hit normal attack string in [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go).
- Implemented plunge entry points in [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go).
- Implemented basic charged attack and Phantasm branch routing in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go).
- Implemented Shadow Dance status application in [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go).
- Implemented Burst scaffold with Veil consumption in [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go).

#### Phase 4-5 partial scaffolding

- Implemented Lunar-Bloom EM bonus path in [internal/characters/nefer/asc.go](internal/characters/nefer/asc.go).
- Implemented current C6 elevation hook in [internal/characters/nefer/cons.go](internal/characters/nefer/cons.go).
