# Nefer PR Checklist

Status legend:

- Done: implemented and present on the branch.
- Partial: implemented in a usable form, but still approximate or incomplete.
- Not done: missing or intentionally deferred.
- N/A: checklist item does not apply to the current repository layout.

| Item | Status | Notes |
| --- | --- | --- |
| New character package | Done | [internal/characters/nefer](internal/characters/nefer) exists and is wired into the build. |
| Config in character package | Done | [internal/characters/nefer/config.yml](internal/characters/nefer/config.yml) exists. |
| Run pipeline with added config (generates character curve, talent stats, `.generated.json` files) | Done | Generated outputs are present in [internal/characters/nefer/nefer_gen.go](internal/characters/nefer/nefer_gen.go), [ui/packages/db/src/Data/char_data.generated.json](ui/packages/db/src/Data/char_data.generated.json), and [ui/packages/ui/src/Data/char_data.generated.json](ui/packages/ui/src/Data/char_data.generated.json). |
| Character key | Done | Generated key exists in [pkg/core/keys/keys_char_gen.go](pkg/core/keys/keys_char_gen.go). |
| Shortcuts for character key | Done | Shortcut registration includes Nefer in [pkg/shortcut/characters.go](pkg/shortcut/characters.go). |
| Update `mode_gcsim.js` with shortcuts for syntax highlighting | N/A | No `mode_gcsim.js` file exists in the current repository layout. |
| Add Character package to imports | Done | Generated sim import exists in [pkg/simulation/imports_char_gen.go](pkg/simulation/imports_char_gen.go). |
| Normal Attack | Done | Implemented in [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go), with timing and geometry still tracked as approximate in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| Charge Attack / Aimed Shot | Partial | Charged Attack, Slither, and Phantasm routing are implemented in [internal/characters/nefer/charge.go](internal/characters/nefer/charge.go), but exact finisher timing, geometry, and constellation-sensitive branching remain incomplete. |
| Skill | Partial | Implemented in [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go), but final geometry, startup validation, and full C2 interaction are not complete. |
| Burst | Partial | Implemented in [internal/characters/nefer/burst.go](internal/characters/nefer/burst.go), but hit timeline and geometry remain provisional. |
| A1 | Partial | Seed conversion window and core replacement path are implemented in [internal/characters/nefer/seeds.go](internal/characters/nefer/seeds.go), but the passive still carries approximation risk. |
| A4 | Partial | Lunar-Bloom EM bonus path is implemented in [internal/characters/nefer/asc.go](internal/characters/nefer/asc.go), but full source-specific validation is still open. |
| C1 | Not done | Still missing; tracked in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| C2 | Partial | Stack-cap increase and initial stack grant exist in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go) and [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go), but the rest of C2 is still missing. |
| C3 | Done | Covered by `SkillCon = 3` in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go). |
| C4 | Not done | Still missing; tracked in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| C5 | Done | Covered by `BurstCon = 5` in [internal/characters/nefer/nefer.go](internal/characters/nefer/nefer.go). |
| C6 | Partial | Lunar-Bloom elevation hook exists in [internal/characters/nefer/cons.go](internal/characters/nefer/cons.go), but the extra damage instances are still missing. |
| Other necessary talents (custom dash/jump, low/high plunge, ...) | Partial | Low and high plunge are implemented in [internal/characters/nefer/attack.go](internal/characters/nefer/attack.go); no custom dash or jump behavior is implemented. |
| Hitlag | Partial | Functional defaults exist, but a final Nefer-specific hitlag pass is still open per [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md). |
| ICD | Partial | ICD tags are assigned across the current attacks, but final validation is still open. |
| StrikeType | Partial | Strike types are assigned, but they are still provisional in several places. |
| PoiseDMG (blunt attacks only for now) | Not done | No final Nefer-specific poise mapping has been completed yet. |
| Hitboxes | Partial | Current hitboxes are usable but still approximate across NA, CA, Skill, Burst, and seeds. |
| Attack durability | Partial | Durability values are assigned on current attacks, but the final validation pass is still open. |
| Particles | Partial | Skill particle generation exists in [internal/characters/nefer/skill.go](internal/characters/nefer/skill.go), but gameplay validation is still pending. |
| Frames | Partial | Core frame data is populated, but several NA, CA, Skill, and Burst timings remain provisional or mixed-source. |
| Update documentation | Done | Nefer implementation, observations, progress, register, and generated docs are present and updated. |
| Xingqiu/Yelan N0 (optional) | Not done | No Nefer-specific follow-up adjustment has been made here. |
| Xianyun Plunge (optional) | Not done | No Nefer-specific follow-up adjustment has been made here. |

## Summary

- Branch status is good enough for functional review.
- The remaining work is concentrated in combat accuracy, hidden-mechanic completion, and unfinished constellations.
- The authoritative live gap list is [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).
