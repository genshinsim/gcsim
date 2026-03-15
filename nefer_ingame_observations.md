# Nefer In-Game Observations

## Scope

This file records direct in-game observations that affect the current Nefer implementation and should be treated as stronger evidence than earlier assumptions.

It is intended to capture observed behavior only.

Implementation choices, approximations, and remaining gaps should be tracked in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md) instead of here.

## Confirmed Observations

### Shared Seed Or Core Limit

- There can be no more than 5 active cores or seeds on the field at the same time.
- Cores and seeds appear to share the same global field-object limit.
- When the limit is exceeded, the oldest existing object is the one that disappears or explodes.

### Seed Lifetime Reset On Conversion

- Based on in-game observation, converting an existing core into a seed appears to refresh or reset that object's lifetime.

### Seed Lifetime

- A seed exists for 12 seconds.

### Slither Behavior

- While Slither is active, Nefer continuously absorbs nearby Seeds of Deceit while moving.
- Slither movement speed matches ordinary character running speed.
- Slither stamina consumption matches ordinary character running consumption.

### Phantasm Charge Reset

- Phantasm charges are reset when switching to another character.

### Phantasm Performance Charges

- The special Shadow Dance charged attack is gated by a separate Phantasm Performance Charge count, not only by Verdant Dew.
- Skill refreshes this Phantasm Performance Charge count instead of stacking it.
- The use count is 3 per Skill activation.
- If Nefer has a Phantasm Performance Charge available but does not yet have enough Verdant Dew, she enters Slither and stays there until Verdant Dew becomes sufficient, if it does.
- As soon as Verdant Dew becomes sufficient during Slither, Phantasm Performance is triggered immediately.
- If the charged input is still held after Phantasm Performance ends, Nefer immediately chains into the next Phantasm Performance if conditions are met, or into a new Slither otherwise.

## Notes

- Use this file as an evidence log.
- If an observation is later contradicted by stronger evidence, update or remove the observation here and then re-evaluate the corresponding implementation note in [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md).
