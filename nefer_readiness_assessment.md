# Nefer Implementation Readiness Assessment

## Purpose

This document describes the current readiness of the Nefer branch as of 2026-03-15.

It is a live readiness assessment. It should be read together with [Nefer PR checklist.md](Nefer%20PR%20checklist.md), [nefer_implementation_progress.md](nefer_implementation_progress.md), and [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).

## Executive Summary

- Data generation, registration, imports, shortcuts, and generated docs are ready.
- The core gameplay loop is implemented: Skill, Shadow Dance, Slither, explicit-hold Charged Attack routing, the no-hold Slither-release tap non-Phantasm route, Phantasm routing, seed conversion, seed absorption, Veil stack scaffolding, and Burst consumption.
- The branch is ready for functional and architectural review.
- The branch is not yet ready to claim final combat-accuracy completeness.
- The main remaining risks are the specific unresolved rows already itemized in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md): `attack.go: unresolved frame rows`, `charge.go: unresolved frame rows`, `skill.go: unresolved frame rows`, `burst.go: unresolved frame rows`, the corresponding per-family hitbox or area sections, `seeds.go and seed_gadget.go`, `cons.go`, and `internal/characters/nefer/*.go: explicit poise and hitlag mappings still missing`.
- The current Phantasm formula interpretation is now explicit in code: C1 raises Shades MV and the Veil bonus is then applied locally to the base constructed Phantasm Performance hit terms rather than using a shared generic bonus slot; later additive reaction terms such as Spread remain outside that multiplier.

## Domain Readiness

| Area | Readiness | Assessment |
| --- | --- | --- |
| Data pipeline and generated assets | High | Canonical generation is in place and the generated outputs are present on the branch. |
| Character registration and simulator integration | High | Character package, keys, imports, shortcuts, and docs outputs are wired in. |
| Normal Attack and plunge baseline | Medium | Remaining gaps are the exact rows listed in the `attack.go: unresolved frame rows`, `attack.go: unresolved hitboxes and areas`, and `internal/characters/nefer/*.go: explicit poise and hitlag mappings still missing` sections of [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| Charged Attack, Slither, and Phantasm loop | Medium | Remaining gaps are the exact rows listed in the `charge.go: unresolved frame rows`, `charge.go: unresolved hitboxes and areas`, and `seeds.go and seed_gadget.go` sections of [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md), including the still-open source validation for the current button-to-button ordinary-CA row interpretation, the temporary `0f` embedded Slither-entry assumption inside those rows, and the source value for the current Slither-entry stamina threshold. |
| Skill and Burst | Medium | Remaining gaps are the exact rows listed in the `skill.go: unresolved frame rows`, `skill.go: unresolved hitboxes, areas, and particles`, `burst.go: unresolved frame rows`, and `burst.go: unresolved hitboxes and areas` sections of [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| Seeds, Veil, and Lunar-Bloom integration | Medium | Remaining gaps are the exact rows listed in the `seeds.go and seed_gadget.go` section of [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| Passives, constellations, and auxiliary action hooks | Medium | A1, A4, C1, C2, C4, C5, and C6 are implemented, and plunge coverage is present; the remaining open items are the specific `cons.go` area row plus any dependent timing or geometry rows already listed under `charge.go`, `seeds.go and seed_gadget.go`, and `attack.go` in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| Combat metadata polish | Low | The still-open metadata work is the exact list in `internal/characters/nefer/*.go: explicit poise and hitlag mappings still missing`, plus validation rows under `skill.go: unresolved hitboxes, areas, and particles` in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| Documentation and review packaging | High | Progress log, live gap register, evidence notes, generated docs, and PR checklist exist and are aligned. |

## What Is Already Strong Enough For Review

- The branch is past the scaffolding stage.
- The gameplay state machine is coherent enough to review in code.
- The seed gadget model is materially closer to observed behavior than the earlier counter-based placeholder.
- The recent charge-loop cleanup reduced accidental complexity without changing the verified standalone special-charge behavior.
- Documentation now separates completed work, live gaps, evidence, and PR status instead of mixing them.

## Validation Already On Record

- [nefer_implementation_progress.md](nefer_implementation_progress.md) records a successful compile check with `go test ./internal/characters/nefer ./pkg/simulation`.
- Recent smoke validation confirmed that a standalone Shadow Dance special charge now self-completes without requiring a following scripted action.
- Recent smoke validation also confirmed the preserved Slither-to-Phantasm transition after the implementation simplification pass.
- Recent smoke validation confirmed that swapping off Nefer clears Shadow Dance and Phantasm Performance Charges, so a charged attack after swapping back no longer takes the immediate Phantasm route and must satisfy the current non-Phantasm entry model.

## Main Readiness Blockers

### 1. Remaining Mechanic Semantics

- `Seed absorption area` and `Seed absorption target-selection rule` remain open in the `seeds.go and seed_gadget.go` section of [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).
- `Slither entry stamina threshold` remains open as a semantics-validation item because the current branch still lacks a source-backed activation threshold row even though the engine-side readiness-versus-consumption behavior is now modeled explicitly.

### 2. Combat-Accuracy Gaps

- Open timing rows are the specific entries listed under `attack.go: unresolved frame rows`, `charge.go: unresolved frame rows`, `skill.go: unresolved frame rows`, and `burst.go: unresolved frame rows` in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).
- Open geometry rows are the specific entries listed under `attack.go: unresolved hitboxes and areas`, `charge.go: unresolved hitboxes and areas`, `skill.go: unresolved hitboxes, areas, and particles`, `burst.go: unresolved hitboxes and areas`, `seeds.go and seed_gadget.go`, and `cons.go` in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).

### 3. Metadata And Engine-Mapping Gaps

- Explicit hitlag and poise rows are still open for the attacks listed in `internal/characters/nefer/*.go: explicit poise and hitlag mappings still missing` in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).
- Validation of the current particle implementation is still open for the rows listed in `skill.go: unresolved hitboxes, areas, and particles` in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).
- ICD, StrikeType, and durability are implemented, but they still await final package-wide validation.

## Review Recommendation

This branch is ready for serious PR review if it is presented honestly as a functionally complete but still combat-inexact Nefer implementation.

It is not yet in a state where it should be described as fully combat-accurate or mechanically complete.

The right review posture is:

1. review the state machine, data integration, and current mechanic modeling now,
2. review the remaining gaps against [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md),
3. decide whether the missing combat-accuracy and constellation work must land before merge.

## Merge Readiness Conclusion

- Ready for review: yes.
- Ready for merge as a partial, explicitly documented implementation: plausible, depending on project expectations.
- Ready for merge as a complete Nefer implementation: no.
