# Project TODO

## Missing Character Implementations

- [x] Detect missing characters programmatically
- [x] Generate pipeline configs from exact local labels
- [x] Generate character files and `zz_*.dm.go` through the community-data generator
- [x] Run gofmt
- [x] First compile and independent structural review
- [x] Automatically fix confirmed config and registration issues
- [ ] Complete character-specific passive, constellation, and runtime-state implementations
- [ ] Verify combat frame data and full simulations

Summary:

- Missing discovered: 13
- Generated: 13
- Compile-ready: 13
- Basic-simulation-ready: 8
- Partial-simulation: 5
- Core-support-blocked: 0
- Full-simulation-ready: 0
- Failed: 0

The generated baseline simulates confirmed direct talent damage. It does not claim that unknown frames, particles, hitboxes, passives, constellations, or star-reaction formulas are complete.

## Star-Reaction Core Support

- [ ] Lunar-Bloom reaction event and confirmed damage formula
  - Required by: Nefer
- [ ] Moonsign state/query and Ascendant Gleam state
  - Required by: Illuga, Jahoda, Linnea, Nefer, Zibai
- [ ] Verdant Dew team resource
  - Required by: Nefer
- [ ] Seed of Deceit field entity and conversion support
  - Required by: Nefer
- [ ] Connect public star-reaction events to character callbacks

## Imported Characters

### Alyosha — 10000148

Status: compile-ready
Simulation level: basic-simulation-ready

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000148.json`
- Lunaris fallback: https://lunaris.moe/character/10000148
- Nanoka fallback query: `Alyosha 10000148`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- None identified

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

### Iansan — 10000110

Status: compile-ready
Simulation level: basic-simulation-ready

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000110.json`
- Lunaris fallback: https://lunaris.moe/character/10000110
- Nanoka fallback query: `Iansan 10000110`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- None identified

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

### Ifa — 10000113

Status: compile-ready
Simulation level: basic-simulation-ready

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000113.json`
- Lunaris fallback: https://lunaris.moe/character/10000113
- Nanoka fallback query: `Ifa 10000113`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- None identified

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

### Illuga — 10000127

Status: compile-ready
Simulation level: partial-simulation

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000127.json`
- Lunaris fallback: https://lunaris.moe/character/10000127
- Nanoka fallback query: `Illuga 10000127`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- Moonsign state/query, Ascendant Gleam state

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

Star-reaction implementation:

- [x] Description dependency classified
- [ ] Character-side trigger logic
- [ ] Character-side stacks/status
- [ ] Character-provided reaction modifiers
- [ ] Talent and constellation interactions
- [ ] Public reaction event integration
- [ ] Confirmed base damage and EM/Moonsign/critical formulas
- [ ] Final damage validation

### Jahoda — 10000124

Status: compile-ready
Simulation level: partial-simulation

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000124.json`
- Lunaris fallback: https://lunaris.moe/character/10000124
- Nanoka fallback query: `Jahoda 10000124`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- Moonsign state/query, Ascendant Gleam state

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

Star-reaction implementation:

- [x] Description dependency classified
- [ ] Character-side trigger logic
- [ ] Character-side stacks/status
- [ ] Character-provided reaction modifiers
- [ ] Talent and constellation interactions
- [ ] Public reaction event integration
- [ ] Confirmed base damage and EM/Moonsign/critical formulas
- [ ] Final damage validation

### Kachina — 10000100

Status: compile-ready
Simulation level: basic-simulation-ready

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000100.json`
- Lunaris fallback: https://lunaris.moe/character/10000100
- Nanoka fallback query: `Kachina 10000100`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- None identified

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

### Linnea — 10000130

Status: compile-ready
Simulation level: partial-simulation

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000130.json`
- Lunaris fallback: https://lunaris.moe/character/10000130
- Nanoka fallback query: `Linnea 10000130`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- Moonsign state/query, Ascendant Gleam state

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

Star-reaction implementation:

- [x] Description dependency classified
- [ ] Character-side trigger logic
- [ ] Character-side stacks/status
- [ ] Character-provided reaction modifiers
- [ ] Talent and constellation interactions
- [ ] Public reaction event integration
- [ ] Confirmed base damage and EM/Moonsign/critical formulas
- [ ] Final damage validation

### Lohen — 10000129

Status: compile-ready
Simulation level: basic-simulation-ready

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000129.json`
- Lunaris fallback: https://lunaris.moe/character/10000129
- Nanoka fallback query: `Lohen 10000129`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- None identified

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

### Nefer — 10000122

Status: compile-ready
Simulation level: partial-simulation

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000122.json`
- Lunaris fallback: https://lunaris.moe/character/10000122
- Nanoka fallback query: `Nefer 10000122`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- Lunar-Bloom reaction event and formula, Moonsign state/query, Ascendant Gleam state, Verdant Dew team resource

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

Star-reaction implementation:

- [x] Description dependency classified
- [ ] Character-side trigger logic
- [ ] Character-side stacks/status
- [ ] Character-provided reaction modifiers
- [ ] Talent and constellation interactions
- [ ] Public reaction event integration
- [ ] Confirmed base damage and EM/Moonsign/critical formulas
- [ ] Final damage validation

### Odette — 10000150

Status: compile-ready
Simulation level: basic-simulation-ready

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000150.json`
- Lunaris fallback: https://lunaris.moe/character/10000150
- Nanoka fallback query: `Odette 10000150`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- None identified

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

### Prune — 10000132

Status: compile-ready
Simulation level: basic-simulation-ready

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000132.json`
- Lunaris fallback: https://lunaris.moe/character/10000132
- Nanoka fallback query: `Prune 10000132`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- None identified

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

### Sandrone — 10000133

Status: compile-ready
Simulation level: basic-simulation-ready

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000133.json`
- Lunaris fallback: https://lunaris.moe/character/10000133
- Nanoka fallback query: `Sandrone 10000133`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- None identified

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

### Zibai — 10000126

Status: compile-ready
Simulation level: partial-simulation

- [x] Local data
- [x] Exact-label config.yml
- [x] Pipeline-generated `zz_*.dm.go`
- [x] Base character and registration
- [x] Normal attack
- [x] Charged attack
- [x] Plunge
- [x] Skill direct damage baseline
- [x] Burst direct damage baseline
- [ ] Complete A1/A4 implementation
- [ ] Complete C1-C6 implementation
- [x] First package compile
- [x] Independent structural review
- [x] Confirmed structural fixes
- [ ] Character-specific mechanic review

Sources:

- Local: `pipeline/community/data/6.7.54/characters/10000126.json`
- Lunaris fallback: https://lunaris.moe/character/10000126
- Nanoka fallback query: `Zibai 10000126`

Combat data TODO:

- [ ] Frames and cancel windows
- [ ] Hitboxes, offsets, and targeting
- [ ] ICD, gauge, hitlag, and particles
- [ ] Snapshot and runtime state-machine timing

Core support required:

- Moonsign state/query, Ascendant Gleam state

Validation:

- gofmt: passed
- package compile: passed
- basic simulation: not run (global simulation package build did not complete in this workspace)
- full simulation: not claimed

Star-reaction implementation:

- [x] Description dependency classified
- [ ] Character-side trigger logic
- [ ] Character-side stacks/status
- [ ] Character-provided reaction modifiers
- [ ] Talent and constellation interactions
- [ ] Public reaction event integration
- [ ] Confirmed base damage and EM/Moonsign/critical formulas
- [ ] Final damage validation

