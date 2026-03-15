# Nefer Implementation Plan

## Purpose

This is the current completion plan for the Nefer branch.

It is a live execution document, not a historical snapshot. Completed work belongs in [nefer_implementation_progress.md](nefer_implementation_progress.md). Remaining approximations and missing mechanics belong in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).

## Current Baseline

- Canonical generated data is in place.
- Character registration, imports, shortcuts, and generated docs are in place.
- Nefer has a working Normal Attack string, Charged Attack routing, Skill, Burst, seed conversion window, seed gadgets, Veil stack scaffolding, Lunar-Bloom EM bonus path, and implemented constellation logic through C6.
- Shadow Dance now supports a real Slither loop, Phantasm Performance Charges, swap reset, and self-completing standalone special charges.
- C2 should be interpreted as a `1 + 5 * 8%` multiplier on the base Phantasm Performance hit terms at 5 Veil stacks, not as a separate hidden 140%-formula problem.
- The current branch now encodes the Phantasm formula ordering explicitly: C1 raises Shades MV, and the Veil bonus is applied afterward by multiplying the base constructed Phantasm Performance hit terms directly in the Phantasm path.
- The branch is no longer blocked on core scaffolding. The remaining work is specifically charge or burst timing closure, Skill or Burst or Phantasm or seed geometry closure, explicit poise or hitlag mapping, final validation of the already assigned ICD or StrikeType or durability or particle behavior, and review packaging.

## Planning Principles

- Do not guess hidden numeric behavior where the source data is still ambiguous.
- Prefer explicit approximation notes over silent assumptions.
- Land improvements in a way that preserves the current working gameplay loop.
- Treat public-facing documentation and PR artifacts as part of the deliverable, not as post-work cleanup.

## Remaining Workstreams

### Implementation Track

Objective:

- Finish the code paths that are still missing or only partially implemented.

Work:

1. Close the remaining seed-side gaps around absorb geometry.
2. Re-check the exact Phantasm and C6 extra-hit timing and geometry now that the extra hits are implemented.
3. Perform the deferred explicit hitlag and poise mapping pass, and validate the already assigned ICD, StrikeType, durability, and particle behavior.

Exit criteria:

- Every missing mechanic that already has a sufficiently clear interpretation is implemented in code.

### Research And Verification Track

Objective:

- Reduce the remaining uncertainty that still depends on frame data, in-game checking, geometry validation, or datamine mapping.

Work:

1. Reconcile mixed frame values in [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md) with the current NA, CA, Skill, and Burst implementations.
2. Finalize or explicitly bound the remaining provisional timings in [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go), [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go), [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go), and [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go).
3. Replace approximate hitboxes with better-supported geometry for NA, Slither exit behavior, Phantasm hits, Skill, Burst, and seed absorption.
4. Re-validate the Slither to Phantasm loop against the current in-game observation notes, especially release timing and chaining behavior.
5. Re-validate seed lifetime reset assumptions, shared cap behavior, and absorb radius behavior.
6. Re-validate Skill particle behavior against hit timing and enemy-contact edge cases.
7. Re-validate exact swap timing, lockout, and any on-exit edge cases around Shadow Dance only if later gameplay evidence contradicts the now-confirmed reset behavior.

Exit criteria:

- Remaining combat-data approximations are either resolved or reduced to a short list covering only mixed frame rows, open AoE mappings, and metadata validation.

### PR Packaging And Review Readiness

Objective:

- Make the branch easy to review without overstating completeness.

Work:

1. Keep [Nefer PR checklist.md](Nefer%20PR%20checklist.md) current.
2. Keep [nefer_implementation_progress.md](nefer_implementation_progress.md) as the factual changelog of completed work.
3. Keep [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md) as the authoritative live gap register.
4. Keep [nefer_ingame_observations.md](nefer_ingame_observations.md) limited to evidence, not implementation policy.
5. Ensure [ui/packages/docs/docs/reference/characters/nefer.md](ui/packages/docs/docs/reference/characters/nefer.md) and generated data outputs remain in sync with the implementation.

Exit criteria:

- A reviewer can identify exactly which rows, radii, timings, and semantic questions are still open without reconstructing the branch history manually.

## Recommended Execution Order

1. Re-check the new P2 Verdant Dew interpretation against stronger source evidence, then close the remaining seed-system gaps before doing another C6 timing or geometry pass.
2. Run the research and verification pass immediately around those code paths: frames, geometry, Slither or Phantasm timing, seed behavior, particles, and swap semantics.
3. Perform the remaining attack-metadata cleanup for explicit hitlag and poise mapping, then validate the current ICD, StrikeType, durability, and particle assignments.
4. Regenerate and re-check any generated artifacts affected by code changes.
5. Refresh PR-facing documentation and gap summaries immediately before review.

## Explicit Non-Goals Until New Evidence Exists

- Do not hard-claim exact geometry where only partial lock-shape metadata exists.
- Do not claim final C6 timing or geometry until the extra-hit frame placement and AoE mapping are supported by stronger evidence.

## Done Definition For The Branch

The branch is ready to present as a complete Nefer implementation only when:

- core combat behavior is implemented,
- remaining frame and geometry risk is narrow and documented as specific open rows or AoE mappings,
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
