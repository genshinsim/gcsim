# Nefer Implementation Plan

## Purpose

This is the current completion plan for the Nefer branch.

It is a live execution document, not a historical snapshot. Completed work belongs in [nefer_implementation_progress.md](nefer_implementation_progress.md). Remaining approximations and missing mechanics belong in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).

## Current Baseline

- Canonical generated data is in place.
- Character registration, imports, shortcuts, and generated docs are in place.
- Nefer has a working Normal Attack string, Charged Attack routing, Skill, Burst, seed conversion window, seed gadgets, Veil stack scaffolding, Lunar-Bloom EM bonus path, and partial constellation coverage.
- Shadow Dance now supports a real Slither loop, Phantasm Performance Charges, swap reset, and self-completing standalone special charges.
- The branch is no longer blocked on core scaffolding. The remaining work is about accuracy, missing mechanics, and review packaging.

## Planning Principles

- Do not guess hidden numeric behavior where the source data is still ambiguous.
- Prefer explicit approximation notes over silent assumptions.
- Land improvements in a way that preserves the current working gameplay loop.
- Treat public-facing documentation and PR artifacts as part of the deliverable, not as post-work cleanup.

## Remaining Workstreams

### 1. Combat Accuracy Pass

Objective:

- Raise Nefer from functional to reviewable combat accuracy for timing, geometry, and attack metadata.

Work:

1. Reconcile mixed frame values in [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md) with the current NA, CA, Skill, and Burst implementations.
2. Finalize or explicitly bound the remaining provisional timings in [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go), [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go), [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go), and [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go).
3. Replace approximate hitboxes with better-supported geometry for NA, Slither exit behavior, Phantasm hits, Skill, Burst, and seed absorption.
4. Perform the deferred hitlag, ICD, StrikeType, durability, and poise review pass.

Exit criteria:

- Remaining combat-data approximations are either resolved or reduced to a short, reviewable list.

### 2. Passive And Constellation Completion

Objective:

- Finish the currently missing character mechanics without destabilizing the working core loop.

Work:

1. Finish Veil of Falsehood duration and refresh behavior.
2. Close the remaining passive gaps around P1 or P2 semantics where the current implementation is still partial.
3. Implement C1.
4. Complete the rest of C2 beyond the current stack-cap and stack-grant behavior.
5. Implement C4.
6. Complete C6 beyond the current Lunar-Bloom elevation hook.

Exit criteria:

- Every passive and constellation is either implemented or explicitly blocked by a source-data gap that can be defended in review.

### 3. Nefer-Specific Mechanic Validation

Objective:

- Validate the most branch-specific behavior so the remaining risk is narrow and explicit.

Work:

1. Re-validate the Slither to Phantasm loop against the current in-game observation notes, especially release timing and chaining behavior.
2. Re-validate seed lifetime reset assumptions, shared cap behavior, and absorb radius behavior.
3. Re-validate Skill particle behavior against hit timing and enemy-contact edge cases.
4. Confirm swap-reset behavior for Phantasm Performance Charges against intended gameplay semantics.

Exit criteria:

- The main Nefer gameplay loop is validated well enough that open risk is concentrated in known precision gaps rather than state-machine uncertainty.

### 4. PR Packaging And Review Readiness

Objective:

- Make the branch easy to review without overstating completeness.

Work:

1. Keep [Nefer PR checklist.md](Nefer%20PR%20checklist.md) current.
2. Keep [nefer_implementation_progress.md](nefer_implementation_progress.md) as the factual changelog of completed work.
3. Keep [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md) as the authoritative live gap register.
4. Keep [nefer_ingame_observations.md](nefer_ingame_observations.md) limited to evidence, not implementation policy.
5. Ensure [ui/packages/docs/docs/reference/characters/nefer.md](ui/packages/docs/docs/reference/characters/nefer.md) and generated data outputs remain in sync with the implementation.

Exit criteria:

- A reviewer can understand what is complete, what is partial, and what still needs evidence without reconstructing the branch history manually.

## Recommended Execution Order

1. Finish the combat accuracy pass where source data is already available.
2. Finish the highest-value missing mechanics: Veil duration or refresh rules, C1, C2, C4, and C6 extra hits.
3. Re-run focused smoke validation for Slither, Phantasm chaining, seed absorption, Burst stack consumption, and particle behavior.
4. Regenerate and re-check any generated artifacts affected by code changes.
5. Refresh PR-facing documentation and gap summaries immediately before review.

## Explicit Non-Goals Until New Evidence Exists

- Do not invent hidden Veil durations or refresh rules.
- Do not hard-claim exact geometry where only partial lock-shape metadata exists.
- Do not mark C1, C2, C4, or C6 complete until the missing behavior is implemented rather than inferred from wording.

## Done Definition For The Branch

The branch is ready to present as a complete Nefer implementation only when:

- core combat behavior is implemented,
- remaining frame and geometry risk is narrow and documented,
- passives and constellations are implemented or defensibly blocked,
- PR docs and generated artifacts are consistent with the actual code state.

## PR Readiness Criteria

A PR is ready only when all of the following are true:

- the code compiles
- the data layer and generated artifacts are aligned
- the code does not hide assumptions where source data is incomplete
- every unresolved part is explicitly documented as requiring additional analysis
- Nefer's base gameplay loop behaves correctly
- C6 elevation uses the existing engine path instead of ad hoc logic
- a reviewer can clearly distinguish exact behavior from still-provisional behavior
