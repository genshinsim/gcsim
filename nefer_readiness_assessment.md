# Nefer Implementation Readiness Assessment

- Prepared: 2026-03-14

## Inputs Used

- Lunaris character page: [Lunaris character page](https://lunaris.moe/character/10000122)
- Lunaris JSON endpoint: [Lunaris JSON endpoint](https://api.lunaris.moe/data/latest/en/char/10000122.json)
- Google Sheet frame source 1: [Workbook 1](https://docs.google.com/spreadsheets/d/1JMqNGKIXK8KM_rpP3eu8lYFnjYbIuX45f_QqxSBsTYY/edit?gid=0#gid=0)
- Google Sheet frame source 2: [Workbook 2](https://docs.google.com/spreadsheets/d/1VpqRMv7FHPayKm9PBbazJs6gL2pr4Ka2F-EI0VgkpDQ/edit?gid=0#gid=0)
- Existing repository state on branch wfpsim-custom
- Existing remote branches nefer/main, nefer/nefer-implementation-01, nefer/nefer-pipeline-run

## Available Source Data

### Character And Data Basics

- Character ID 10000122
- Name, element, weapon type, rarity, base stats, ascension stats, constellation name, birthday, short bio
- Full talent multiplier tables for Normal Attack, Charged Attack, Plunge, Skill, and Burst
- Skill cooldown, Burst cooldown, Burst cost, Skill charge count, and Shadow Dance duration
- Ascension and talent material requirements
- Particle generation metadata from Lunaris API
- Combat metadata from Lunaris API for several attack objects: element, gauge, ICD rule/source, and poise

### Gameplay And Mechanic Descriptions

- Full text descriptions of Normal Attack, Charged Attack, Plunge, Skill, and Burst
- Full passive descriptions P1-P4
- Full constellation descriptions C1-C6
- Explicit mechanical statements including:
  - Bloom converts to Lunar-Bloom through P3
  - P3 scales Lunar-Bloom base damage by 0.0175% per EM, capped at 14%
  - Skill gives Shadow Dance, has 2 charges, and lets Charged Attack become Phantasm Performance when the Verdant Dew condition is met
  - P1 creates a 15s replacement window where Dendro Cores and future Lunar-Bloom outputs become Seeds of Deceit
  - Burst consumes Veil of Falsehood stacks for bonus damage
  - C4 applies 20% Dendro RES shred with 4.5s linger logic
  - C6 adds extra Lunar-Bloom-tagged damage instances to Phantasm Performance

### Frame And Cancel Data

- Hitmark measurements for the Normal Attack string, Normal CA, and Phantasm CA
- Cancel timing measurements for several attack chains
- Secondary measurements for Yelan/Xingqiu-like timing behavior
- Source videos and notes from the counters

### Repository And Engine Context Already Present

- The current branch already has Lunar-Bloom and Lunar-Charged infrastructure in the engine
- The repository already has examples of lunar-mechanic characters such as Lauma, Flins, and Ineffa
- The repository already has attack tags, events, and damage handling for lunar reactions
- The repository already has support for special damage elevation through AttackInfo.Elevation, including existing character-side examples
- The repository already has a full Verdant Dew implementation in the player handler, including cap, gain, consumption, expiry window, script hook, conditional exposure, and a production integration example in Lauma

## Missing Or Incomplete Data

### Frame And Animation Gaps

- Definitive values for every mixed frame measurement from the Google sheets
- Skill startup hitmark and cancel timings that still depend on mixed measurements
- Burst startup hitmarks and cancel timings that still depend on mixed measurements
- Dash, Jump, Swap, Skill, and Burst cancels for most NA and CA states where the sheet currently has no value
- Final confirmed value for N3 -> N4, because the measured range is still mixed

### Combat Geometry And Collision Gaps

- Exact hitbox shapes, radii, widths, lengths, and offsets for NA hits, CA, Phantasm hits, Skill, Burst, and C6 extra hits
- Absorption radius for Seeds of Deceit
- Whether attacks are single-target locked, centered AoE, forward box, or circle beyond what is implied textually

### Mechanics Requiring Hidden Numeric Data

- Base duration of Veil of Falsehood stacks
- Exact refresh rules for Veil of Falsehood beyond the textual 3rd and 5th stack refresh behavior
- Exact quantitative effect of P2 additional Verdant Dew provision
- Exact numeric meaning of C2 wording that causes Phantasm Performance to deal up to 140% of its original damage
- Exact timing and spacing of the multiple hits inside Burst
- Exact timing and application details for C6 extra damage instances

### Engine Mapping Gaps

- Exact mapping from Lunaris attack metadata entries to specific hits inside Nefer's NA and Phantasm strings
- Final ICD tag and group mapping for each hit in gcsim terms
- Final StrikeType mapping for each hit in gcsim terms where text alone is insufficient

## What Can Be Implemented Unambiguously

### Data And Registration Layer

- Add Nefer to character keys, imports, generated character data, docs, and avatar data
- Create config and generated multiplier tables from Lunaris data
- Set basic stats, element, weapon type, rarity, ascension stats, talent scaling tables, Burst cost, Skill charges, and cooldowns

### Behavior That Is Textually Explicit Enough To Implement Directly

- Four-hit Normal Attack string with Dendro damage
- Charged Attack entering Slither, draining stamina while active, and spending additional stamina on exit hit
- Shadow Dance state granted by Skill
- Phantasm Performance replacing Charged Attack only when Shadow Dance is active and Verdant Dew is available
- Integration with the existing Verdant Dew resource through the current player handler API, following the Lauma pattern where applicable
- Skill initial charge count of 2 and skill cooldown of 9s
- Burst consuming all Veil of Falsehood stacks and applying Burst damage bonus per stack based on talent level
- P3 conversion from Bloom to Lunar-Bloom
- P3 Lunar-Bloom base damage scaling from Nefer EM: 0.0175% per EM, capped at 14%
- C3 and C5 talent level increases
- C4 Dendro RES reduction of 20% while on field in Shadow Dance, with removal 4.5s after exit or distancing
- C6 Lunar-Bloom damage elevation of 15% through the existing AttackInfo.Elevation path
- Particle generation probabilities from the Lunaris energy metadata

### Frame Data Safe To Hard-Code Now

- Any entry marked Confirmed in [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md)
- Yelan timing summary 10 only if the final implementation hook matches the same semantics as the measurement sheet; otherwise it must be re-validated at the integration point

## What Still Requires Additional Analysis

- Any frame or cancel value marked Mixed in [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md)
- Any frame or cancel value marked Missing in [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md)
- The exact semantic meaning of Yelan 10 and Xingqiu 11 from the second Google sheet relative to gcsim AnimationStartDelay or N0 trigger logic
- The exact intended implementation of Seeds of Deceit as a gadget or other entity, including hitbox, lifetime, absorption radius, and ownership behavior beyond the textual description
- The exact stack duration of Veil of Falsehood and its per-stack refresh model
- The exact quantity and formula behind P2 additional Verdant Dew provision
- The exact meaning of C2 up to 140% of original DMG
- The full Burst hit timeline and the detailed C6 extra-hit timeline
- The exact hitbox geometry for attacks that currently only have textual or partial metadata support
- The exact ICD and StrikeType mapping for hits where Lunaris text does not map one-to-one to gcsim attack events

## Source-Confirmation Gaps That Should Not Be Guessed

- The corrected N3 -> N4 cancel value still needs confirmation because current measurements remain mixed at 48-53
- Whether Phantasm CA hit 5 should be modeled as frame 0 or frame 1 in the engine
- Whether N1 and N2 hitmarks should use the majority measurement or require a new recording pass
- Whether N4 hitmark is 22 or 24, because current trials conflict

## Practical Conclusion

### Current Readiness

- High readiness for data entry, registration, scalar tables, basic state machine, and the parts of Nefer that are fully specified by text or confirmed measurements
- Medium readiness for Normal, Charged, and Phantasm timing because several values are confirmed but some important values still conflict
- Low readiness for exact combat geometry, exact Veil hidden numeric behavior, and constellation parts that rely on quantities not fully exposed by current sources, excluding C6 elevation and the base Verdant Dew resource path which are already supported by the engine

### Recommended Execution Order

1. Implement the fully specified data layer and state scaffolding first.
2. Implement the mechanics that are textually explicit and numerically defined.
3. Hard-code only confirmed frame values from the sheets.
4. Leave every mixed or missing timing or geometry item explicitly marked for additional analysis.
5. Do not guess hidden quantities such as Veil duration, the P2 Verdant Dew provision formula, or C2 140% handling without another source.
