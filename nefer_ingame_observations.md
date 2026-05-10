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

- Converting an existing core into a seed resets that object's lifetime.
- When a core is converted on the last second of its lifetime, the resulting seed still persists for a full 12 seconds.

### Seed Lifetime

- A seed exists for 12 seconds.

### Seed Behavior After P1 Window End

- When the 15s P1 conversion window ends, already existing Seeds of Deceit remain on the field.
- No special cleanup, forced conversion, or immediate detonation was observed at P1 window end.

### Slither Behavior

- While Slither is active, Nefer continuously absorbs nearby Seeds of Deceit while moving.
- Slither movement speed matches ordinary character running speed.
- Slither stamina consumption uses the sourced Charged Attack charging drain value of 18.15 stamina per second.

### Shadow Dance Swap Reset

- Swapping to another character removes Shadow Dance.
- After swapping away and back, Nefer's next ordinary Charged Attack consumes the full 50 stamina cost instead of the reduced 25 stamina Shadow Dance route.

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
