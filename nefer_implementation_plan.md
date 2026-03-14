# Nefer Implementation Plan

## Objective

Implement Nefer on the current branch at a quality level that is suitable for normal project review, without relying on hidden assumptions where the source data is incomplete.

This plan is based on:

- the contributor checklist and contributor guide
- Lunaris data
- Google Sheets frame data
- the current wfpsim-custom infrastructure
- branch archaeology for nefer/main, nefer/nefer-implementation-01, and nefer/nefer-pipeline-run

## Execution Principles

- Do not merge nefer/nefer-pipeline-run as-is.
- Use nefer/nefer-pipeline-run only as a rough structural reference, not as a finished implementation.
- Implement only what is supported by source data or existing engine infrastructure.
- Keep every unsupported or ambiguous part explicitly documented as requiring additional analysis.
- Build the data layer and state scaffolding first, then combat logic, then timing accuracy, then documentation and verification.

## Working Artifacts

- [nefer_lunaris_10000122.md](nefer_lunaris_10000122.md)
- [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md)
- [nefer_readiness_assessment.md](nefer_readiness_assessment.md)
- [nefer_implementation_progress.md](nefer_implementation_progress.md)
- [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md)

## Phase Layout

1. Data layer and generation
2. Character scaffold and registration
3. Base attacks and state machine
4. Nefer-specific mechanics
5. Passives and constellations
6. Combat accuracy
7. Generated artifacts and docs
8. Verification and cleanup

## Phase 1: Data Layer And Generation

Goal:

- Prepare all source-driven inputs so Nefer enters the character generation pipeline cleanly.

Work:

1. Create the character package under internal/characters/nefer.
2. Prepare config.yml for Nefer.
3. Ensure genshin_id is 10000122.
4. Validate skill_data_mapping against Lunaris and extract-talents output.
5. Run character generation.
6. Verify base stats, talent scaling tables, generated character data, and generated UI, DB, and docs artifacts.

Done when:

- internal/characters/nefer/config.yml is valid.
- Character generation succeeds without manual edits to generated scaling tables.

Open questions:

- If extract-talents and Lunaris disagree on naming for Phantasm Performance, Verdant Dew, or Burst bonus fields, that mismatch must be resolved before combat logic is finalized.

## Phase 2: Character Scaffold And Registration

Goal:

- Build the minimal compilable Nefer package and wire it into the engine.

Work:

1. Create nefer.go.
2. Register the character through core.RegisterCharFunc(keys.Nefer, NewChar).
3. Set EnergyMax = 60, SkillCon = 3, BurstCon = 5, NormalHitNum = 4, and SetNumCharges(action.ActionSkill, 2).
4. Add the character key to generated keys.
5. Add the character to generated imports.
6. Add any required shortcut registration.
7. Update mode_gcsim.js if shortcut highlighting depends on it.
8. Verify that Nefer appears in generated character data.

Done when:

- The Nefer package compiles as a working skeleton.
- The simulator recognizes the character.

## Phase 3: Base Attacks And State Machine

Goal:

- Implement a functional combat foundation: NA, CA, plunge, Shadow Dance, and CA replacement.

Work:

1. Implement the 4-hit Striking Serpent Normal Attack string.
2. Bind the Normal Attack, Charged Attack, and plunge talent multipliers.
3. Use only confirmed frame values from [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md) where they exist.
4. Keep mixed frame values explicitly provisional until they are re-validated.
5. Implement Slither entry, stamina drain, exit stamina cost, and exit hit.
6. Respect the reduced extra stamina cost while Shadow Dance is active.
7. Preserve the rule that Skill and sprint do not forcibly break Slither.
8. Implement plunge, low plunge, and high plunge.
9. Create the Shadow Dance state, grant it on Skill, and set its duration to 9s.
10. Replace Charged Attack with Phantasm Performance only while Shadow Dance is active and Verdant Dew is available.

Done when:

- NA, CA, plunge, and Shadow Dance switching work as a coherent state machine.

Open questions:

- Mixed hitmarks remain unresolved until video confirmation exists.
- Animation-accurate Skill and sprint interaction inside Slither may still need another pass.

## Phase 4: Nefer-Specific Mechanics

Goal:

- Implement the mechanics that define Nefer's gameplay loop.

Work:

1. Reuse the player-level Verdant Dew system instead of adding a Nefer-only resource.
2. Use Lauma as the main integration reference.
3. Implement Phantasm Performance as the special Charged Attack route.
4. Split Nefer-owned hits and Lunar-Bloom-owned hits with the correct ownership and tagging.
5. Implement the 15s replacement window that converts relevant Dendro Core outcomes into Seeds of Deceit.
6. Ensure Seeds of Deceit do not explode and cannot trigger Hyperbloom or Burgeon.
7. Implement seed absorption through Charged Attack and Phantasm Performance.
8. Grant one Veil of Falsehood stack per absorbed seed.
9. Implement Veil stacks, the 3-stack refresh behavior, the +100 EM for 8s effect, and Burst consumption.

Confirmed engine references:

- Base resource model: [pkg/core/player/player.go](pkg/core/player/player.go)
- Conditional exposure: [pkg/conditional/conditional.go](pkg/conditional/conditional.go)
- Starting-value script hook: [pkg/gcs/eval/sysfunc.go](pkg/gcs/eval/sysfunc.go)
- Character-side Verdant Dew usage: [internal/characters/lauma/skill.go](internal/characters/lauma/skill.go)
- Character integration setup: [internal/characters/lauma/lauma.go](internal/characters/lauma/lauma.go)

Done when:

- The main gameplay loop works: Skill -> Shadow Dance -> Phantasm -> Seeds -> Veil -> Burst.

Open questions:

- P2 and C4 interactions with Verdant Dew generation still need a dedicated design pass.
- Seed geometry, lifetime details, and absorption radius still need better source support.
- Veil base duration and exact refresh rules are still incomplete.

## Phase 5: Passives And Constellations

Goal:

- Cover passives and constellations while separating well-defined logic from still-ambiguous logic.

Work:

1. Implement P1 as the main combat passive.
2. Implement P2 as Verdant Dew enhancement logic once the numeric model is clear enough.
3. Implement P3 as Bloom -> Lunar-Bloom conversion plus Lunar-Bloom base damage scaling.
4. Implement P4 as the non-combat expedition passive.
5. Implement C1 Lunar-Bloom base damage scaling from 60% of EM for the Phantasm route, while keeping any unresolved Veil interaction documented.
6. Implement C2 duration extension, stack cap increase, instant 2-stack Skill grant, and 5-stack +200 EM behavior.
7. Keep the up to 140% original damage wording explicitly unresolved until the exact formula is known.
8. Implement C3 and C5 as talent level increases.
9. Implement C4 Dendro RES shred and the 4.5s linger logic.
10. Implement C6 extra Lunar-Bloom-tagged hits as far as the available data supports, and keep the unresolved timeline documented.
11. Implement the C6 15% Lunar-Bloom elevation through AttackInfo.Elevation.

Done when:

- Every passive and constellation is either implemented or explicitly documented as blocked by missing data.

Open questions:

- P2 exact additional Verdant Dew provision formula
- C1 interaction details with Veil
- C2 exact meaning of up to 140% of original damage
- C4 exact gain-rate interaction with the Verdant Dew pipeline
- C6 exact extra-hit timeline and spacing

## Phase 6: Combat Accuracy

Goal:

- Raise Nefer to the repository's accepted combat-accuracy standard.

Work:

1. Apply confirmed frame values from [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md).
2. Re-validate or leave mixed values explicitly unresolved.
3. Add Skill and Burst frame values only where confirmation exists.
4. Add hitlag only where it is supported by evidence.
5. Finalize attack tags and ICD groups for NA, CA, plunge, Skill, Burst, Phantasm, and extra hits.
6. Finalize StrikeType values only where the attack form is clear.
7. Use Lunaris poise metadata where the hit mapping is known.
8. Define hitbox geometry for NA, CA exit hit, Phantasm, Skill, Burst, and C6 extra hits only where data is strong enough.
9. Apply Lunaris gauge data only where the mapping to a specific hit event is reliable.
10. Implement particle generation from Lunaris energy metadata.
11. Reconcile workbook 2 Yelan and Xingqiu timing data with the actual gcsim hook point before hard-coding it.

Done when:

- The combat layer is accurate enough for normal usage and later code review.

Open questions:

- Mixed frame values
- Missing cancels
- Exact hitbox geometry for most attacks
- Exact Burst timeline
- Exact C6 extra-hit timing
- Final mapping of some poise and gauge entries to attack events

## Phase 7: Generated Artifacts And Docs

Goal:

- Bring generated assets and public documentation into sync with the implementation.

Work:

1. Regenerate character data.
2. Regenerate UI, DB, and docs data.
3. Review the generated character documentation page.
4. Record known issues where parts of the behavior still require additional analysis.
5. Verify char_data.generated.json and related generated asset files.

Done when:

- Nefer is visible in generated docs and data artifacts.

## Phase 8: Verification And Cleanup

Goal:

- Finish with a reviewable, testable, and maintainable implementation.

Work:

1. Ensure internal/characters/nefer builds cleanly.
2. Run relevant tests.
3. Add or validate smoke scenarios for:
   - Skill -> Shadow Dance -> Phantasm
   - Skill -> Seeds of Deceit -> Veil
   - Burst consuming Veil
   - P3 Bloom -> Lunar-Bloom conversion
   - C4 Dendro RES shred
   - C6 extra hits and elevation
4. Compare observed behavior against the gathered sources.
5. Produce a residual-gaps list.

Done when:

- The implementation is ready for review.
- Remaining uncertainties are isolated in documentation instead of hidden in code.

## Status Checklist

### Can Be Closed Immediately

- New character package
- Character config
- Character generation with the new config
- Character key and imports
- Registration layer
- Normal Attack skeleton using confirmed timing only
- Charged Attack and Slither skeleton
- Skill skeleton
- Burst skeleton
- P1 core logic without hidden numeric assumptions
- P3 core logic
- C3
- C5
- C6 elevation through existing infrastructure
- Particle generation
- Documentation scaffolding

### Can Be Implemented But Still Need A Design Pass

- Phantasm Performance
- Seeds of Deceit
- Veil of Falsehood
- C1
- C2
- C4
- C6 extra hits
- Final frame pass
- ICD
- StrikeType
- PoiseDMG mapping
- Hitboxes
- Attack durability mapping
- Optional Xingqiu and Yelan N0 integration

### Cannot Be Finalized Yet Without Additional Analysis

- Exact Verdant Dew generation and provision model for Nefer-specific modifiers
- Veil duration and refresh model
- Mixed frame values
- Missing Skill, Burst, Dash, Jump, and Swap cancels
- Full Burst hit timeline
- Exact geometry for most hitboxes
- Exact C2 formula for up to 140% of original damage
- Exact C6 extra-hit timeline model

## Recommended Real Implementation Order

1. Build config.yml, generation, and the character skeleton.
2. Connect the registration layer: keys, imports, and generated data.
3. Implement NA, base CA, plunge, and Skill state.
4. Implement P3 and base Lunar-Bloom integration.
5. Implement Seeds of Deceit, Veil, and the Phantasm route.
6. Implement the Burst consumption route.
7. Implement C3 and C5, then C4, then C6 elevation.
8. After that, implement C1, C2, and C6 extra-hit logic.
9. Then improve frames, hitlag, ICD, and geometry.
10. Finish with regeneration, documentation pass, smoke scenarios, and cleanup.

## PR Readiness Criteria

A PR is ready only when all of the following are true:

- the code compiles
- the data layer and generated artifacts are aligned
- the code does not hide assumptions where source data is incomplete
- every unresolved part is explicitly documented as requiring additional analysis
- Nefer's base gameplay loop behaves correctly
- C6 elevation uses the existing engine path instead of ad hoc logic
- a reviewer can clearly distinguish exact behavior from still-provisional behavior
