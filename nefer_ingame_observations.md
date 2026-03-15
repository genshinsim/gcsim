# Nefer In-Game Observations

## Scope

This file records direct in-game observations that affect the current Nefer implementation and should be treated as stronger evidence than earlier assumptions.

## Confirmed Observations

### Shared Seed Or Core Limit

- There can be no more than 5 active cores or seeds on the field at the same time.
- Cores and seeds appear to share the same global field-object limit.
- When the limit is exceeded, the oldest existing object is the one that disappears or explodes.

### Seed Lifetime Reset On Conversion

- Based on in-game observation, converting an existing core into a seed appears to refresh or reset that object's lifetime.

### Seed Lifetime

- A seed exists for 12 seconds.

## Implementation Impact

- A final Seeds of Deceit implementation should not use an unlimited internal counter.
- The eventual model should enforce a shared cap of 5 for cores and seeds.
- Conversion into a seed should reset the lifetime timer.
- Seed lifetime should currently be treated as 12 seconds unless stronger contradictory evidence appears.
