# Nefer Implementation Progress

## Purpose

This file is a dated changelog of completed Nefer work on the branch.

It should not be used as the source of current readiness or open-gap status. For that, use [Nefer PR checklist.md](Nefer%20PR%20checklist.md), [nefer_readiness_assessment.md](nefer_readiness_assessment.md), and [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).

## Accepted Implementation Decisions

- `C1` is implemented as an additive `+0.6` increase to the EM-scaling multiplier of each Phantasm Performance Shades hit, not as a separate downstream damage bonus.
- The Phantasm Veil bonus is applied locally to the base constructed Phantasm Performance hit terms in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go); later additive terms such as Spread remain outside that multiplier.
- `C2` is interpreted as a Phantasm base-term multiplier ceiling of `1 + 5 * 8%` at 5 Veil stacks, not as a separate hidden `140%` formula lane.
- When adding a Veil stack at cap, the branch refreshes the oldest active stack first.
- `P1` only opens a 15s conversion window. When that window ends, existing Seeds of Deceit remain on field; there is no special cleanup or forced detonation.
- Converting an existing core into a Seed of Deceit resets lifetime by replacing the old core gadget with a newly created seed gadget.
- Seeds and ordinary Dendro Cores share a field cap of `5`, and the branch intentionally models that through the shared `GadgetTypDendroCore` engine path.
- `P2` is modeled as a Slither-only Verdant Dew generation-rate bonus refreshed by party-triggered Lunar-Bloom for `5s`, scaling from `+0%` at `EM <= 500` by `+10%` per `100` EM above `500`, up to `+50%` total bonus.
- `C4` nearby-opponent RES shred is modeled as an on-field Shadow Dance refresh loop with a `4.5s` linger, while the exact nearby-opponent area definition remains open.
- `C6` is modeled as three separate pieces: the existing `+15%` Lunar-Bloom elevation hook, conversion of the second Nefer Phantasm hit into an `85% EM` Lunar-Bloom hit, and an additional `120% EM` Lunar-Bloom hit when Phantasm Performance ends.
- `AnimationYelanN0StartDelay = 10` is the current best-fit branch choice for Nefer's Yelan N0 integration, but it remains an open integration assumption rather than a closed source-backed mapping.
- `ChargeAttack(p)` now follows a four-scenario route split with explicit hold-frame semantics. `hold=0` means tap. `hold=1..150` means hold for exactly that many frames, capped to `150`. With tap, Nefer performs the fixed single-Phantasm route if Phantasm is already available, otherwise she performs the no-hold ordinary CA route. With hold, Nefer first checks the same immediate-Phantasm condition; if it is true, hold behaves identically to tap. If it is false, hold enters Slither as a prefix state and then resolves into either Phantasm CA as soon as the Phantasm condition becomes true or ordinary Charged Attack when the hold route ends.
- Workbook 1 ordinary Charged Attack rows are treated as button-to-button initiation intervals for the full no-hold Slither-release non-Phantasm route. Within the accepted branch model, ordinary CA is not separable from the Slither-entry segment that immediately precedes it, so the code intentionally maps `Charge Attack -> Normal 1` directly to `29` and `Charge Attack -> Swap` directly to `28` without introducing a second internal boundary for that embedded entry segment.
- Ordinary Nefer Charged Attack is not allowed to start unless Nefer can first enter Slither. The branch now models that through a dedicated Slither-entry stamina threshold in the shared charge readiness gate, using the engine stamina-spec contract to separate readiness requirement from actual upfront consumption: Nefer must have roughly one second worth of Slither stamina available to begin the non-Phantasm route, but that threshold is not consumed up front. Actual stamina consumption still occurs only through Slither ticks once Slither has non-zero duration.

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
- Corrected [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go) so Skill cooldown now starts on the confirmed `22f` row and `Skill -> Walk` now uses the confirmed `34f` row.
- Reinterpreted workbook 1 ordinary Charged Attack rows in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go) as button-to-button timings for the full no-hold Slither-release route rather than as post-windup release timings, and updated `Charge Attack -> Normal 1` to `29f` plus `Charge Attack -> Swap` to `28f` under the current `0f` embedded Slither-entry assumption.
- Corrected Nefer action queue metadata so earliest-cancel values now line up with the actual per-action frame tables in [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go) and [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go): normal attacks now report the earliest skill-cancel row as `CanQueueAfter`, and Burst now reports the sourced `Swap = 118f` earliest cancel instead of the later attack row.
- Added an engine stamina-spec contract across [pkg/core/action/action.go](pkg/core/action/action.go), [pkg/core/player/player.go](pkg/core/player/player.go), and [pkg/core/player/exec.go](pkg/core/player/exec.go) so abilities can express a readiness threshold separately from actual stamina consumption timing; [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go) now uses that contract to enforce the Slither-entry requirement without consuming it up front, and [pkg/core/player/dew.go](pkg/core/player/dew.go) now correctly preserves infinite Verdant Dew rate mods.
- Corrected the hold-route Slither -> Phantasm transition in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go) so updating the active action length no longer delays `CanQueueAfter` until the end of Phantasm; the held route now keeps the same early queue-window semantics as the direct Phantasm charge path while still reporting later action-specific `Frames(next)` values.
- Reworked [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go) so `hold` is now interpreted as an explicit hold duration in frames rather than a boolean. `charge` means tap, while `charge[hold=n]` means hold for `n` frames up to `150`. The hold route now acts as a Slither prefix that resolves into either the same Phantasm route or the same ordinary Charged Attack route depending on which condition wins first.
- Added in-code mechanic comments in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go) documenting the four-scenario charge split, the exact `hold` contract, and the current in-game stamina rule that ordinary CA cannot begin unless Nefer can enter Slither first.
- Collapsed the scattered charge-phase fields in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go) into an explicit `chargeRouteState` struct and renamed the route source field accordingly, so the Nefer charge path now carries one coherent internal state object instead of several loosely related top-level fields.
- Reworked [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go) cancel timing so Burst no longer uses the old `76f` scaffold: `Burst -> Skill` now uses the confirmed `119f` row, the shared action-end timing now tracks the late post-Burst cancel window, and `Burst -> Charge` now distinguishes the ordinary CA route (`131f`) from the Phantasm CA route (`124f`).
- Corrected Slither stamina drain to the sourced 18.15/s value and split Charged Attack stamina handling into the three observed modes in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go): normal 50 stamina, Shadow Dance normal 25 stamina, and Phantasm Performance 0 stamina.
- Corrected swap handling so leaving the field clears Shadow Dance itself in addition to resetting Phantasm Performance Charges in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go).
- Implemented C4 Verdant Dew gain-rate acceleration through a reusable player Verdant Dew rate modifier path in [pkg/core/player/dew.go](pkg/core/player/dew.go) and [pkg/core/player/player.go](pkg/core/player/player.go).
- Implemented C4 on-field Shadow Dance Dendro RES shred with a 4.5s linger refresh model in [internal/characters/nefer/cons.go](internal/characters/nefer/cons.go).
- Corrected P1 gating in [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go) so the Seed of Deceit passive now follows the sourced `Moonsign: Ascendant Gleam` requirement instead of unlocking from Ascension alone.
- Implemented the current branch P2 model in [internal/characters/nefer/asc.go](internal/characters/nefer/asc.go): party-triggered Lunar-Bloom refreshes a 5s window, and while Nefer is in Slither during Shadow Dance she gains an additional Verdant Dew generation-rate bonus in the same mechanical lane used by C4, scaled by EM above 500 from +0% up to +50% total bonus.
- Implemented C1 by interpreting it as an additive +0.6 increase to the EM scaling multiplier of each Phantasm Performance Shades hit in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go), and applied the Veil bonus to the base constructed Phantasm Performance hit terms in that same path so the final ordering is `(base Phantasm term) * (1 + Veil bonus)`, with C1 only enlarging the Shades MV portion before that multiplier and later additive reaction terms remaining outside it.
- Implemented the current high-confidence C6 extra-hit model in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go): the second Nefer Phantasm hit is converted into a Lunar-Bloom EM hit equal to 85% EM, and an additional 120% EM Lunar-Bloom hit is queued when Phantasm Performance ends, while [internal/characters/nefer/cons.go](internal/characters/nefer/cons.go) continues to provide the separate 15% Lunar-Bloom elevation hook under Ascendant Gleam.

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
