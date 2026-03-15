# Nefer Implementation Readiness Assessment

## Purpose

This document describes the current readiness of the Nefer branch as of 2026-03-15.

It is a live readiness assessment. It should be read together with [Nefer PR checklist.md](Nefer%20PR%20checklist.md), [nefer_implementation_progress.md](nefer_implementation_progress.md), and [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).

## Executive Summary

- Data generation, registration, imports, shortcuts, and generated docs are ready.
- The core gameplay loop is implemented: Skill, Shadow Dance, Slither, Phantasm routing, seed conversion, seed absorption, Veil stack scaffolding, and Burst consumption.
- The branch is ready for functional and architectural review.
- The branch is not yet ready to claim final combat-accuracy completeness.
- The main remaining risks are incomplete passives or constellations, unresolved timing or geometry details, and a still-open metadata validation pass for hitlag, poise, and some attack tagging details.
- The current Phantasm formula interpretation is now explicit in code: C1 raises Shades MV and the Veil bonus is then applied locally to the base constructed Phantasm Performance hit terms rather than using a shared generic bonus slot; later additive reaction terms such as Spread remain outside that multiplier.

## Domain Readiness

| Area | Readiness | Assessment |
| --- | --- | --- |
| Data pipeline and generated assets | High | Canonical generation is in place and the generated outputs are present on the branch. |
| Character registration and simulator integration | High | Character package, keys, imports, shortcuts, and docs outputs are wired in. |
| Normal Attack and plunge baseline | Medium | Usable and implemented, but still carrying frame and hitbox approximation risk. |
| Charged Attack, Slither, and Phantasm loop | Medium | Functional and substantially improved, but still not final on timing, geometry, and some branch-specific details. |
| Skill and Burst | Medium | Playable and integrated, but still provisional on exact geometry and hit timelines. |
| Seeds, Veil, and Lunar-Bloom integration | Medium | Core loop exists, but several semantics remain inferred rather than fully confirmed. |
| Passives and constellations | Medium | C1, C2, and C4 are implemented, and C6 remains incomplete. |
| Combat metadata polish | Low | Hitlag, poise, exact StrikeType, durability, and some ICD validation are still pending. |
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
- Recent smoke validation confirmed that swapping off Nefer clears Shadow Dance and Phantasm Performance Charges, so a charged attack after swapping back consumes the full 50 stamina normal route.

## Main Readiness Blockers

### 1. Missing Character Mechanics

- C6 only covers the Lunar-Bloom elevation hook and not the extra damage instances.
- C4 is now implemented functionally, but its nearby-opponent radius is still an approximation rather than a source-confirmed area definition.
- Veil of Falsehood now has a sourced duration model in code, but the exact refresh-target semantics at cap are still an implementation assumption.

### 2. Combat-Accuracy Gaps

- Several frame values remain mixed or provisional.
- Hitbox geometry is still approximate across multiple actions.
- Seed absorption radius is still an assumption rather than a fully verified mapping.
- Burst timeline and Phantasm finisher timing are not yet final.

### 3. Metadata And Engine-Mapping Gaps

- Final hitlag mapping is still open.
- Final poise mapping is still open.
- Some StrikeType and ICD decisions still need a deliberate validation pass.
- Durability values exist, but the final validation pass is still missing.

## Review Recommendation

This branch is ready for serious PR review if it is presented honestly as a functional but still incomplete Nefer implementation.

It is not yet in a state where it should be described as fully combat-accurate or mechanically complete.

The right review posture is:

1. review the state machine, data integration, and current mechanic modeling now,
2. review the remaining gaps against [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md),
3. decide whether the missing combat-accuracy and constellation work must land before merge.

## Merge Readiness Conclusion

- Ready for review: yes.
- Ready for merge as a partial, explicitly documented implementation: plausible, depending on project expectations.
- Ready for merge as a complete Nefer implementation: no.
