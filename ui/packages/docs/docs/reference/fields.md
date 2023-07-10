---
title: Fields
sidebar_position: 2
---

Fields have the following structure:

```
.field1.field2.field3.field4
```

All fields evaluate to a number.

## Available Fields

:::tip
Most of the specific tags can be located in their respective character/weapon/artifact page. 
For example, if you are looking for the tag for Lisa's A4 (defense shred), then look in Lisa's character page.
:::

<!-- prettier-ignore -->
| field1 | field2 | field3 | field4 | description |
| --- | --- | --- | --- | --- |
| `debuff` | `res`/`def` | `t1`/`t2`/`t3`/... | res/def modifier name | Evaluates to the remaining duration of the specified res/def modifier on the specified target. See the relevant character/weapon/artifact page for acceptable modifier names. |
| `element` |  `t1`/`t2`/`t3`/... | `pyro`/`hydro`/`anemo`/`electro`/`dendro`/`cryo`/`geo`/`frozen`/`quicken` | - | `1` if the specified element exists on specified target, `0` otherwise. |
| `status` | core status name | - | - | Evaluates to the remaining duration of the specified core status. See the relevant character/weapon/artifact page for acceptable status names. |
| `stam` | - | - | - | Evaluates to the player's remaining stamina. |
| `construct` | `duration`/`count` | construct name | - | Evaluates to the duration/count of the specified construct. See individual character page for acceptable construct names. |
| `gadgets` | `dendrocore` | `count` | - | Evaluates to the current number of Dendro Cores. |
| `keys` | `char`/`weapon`/`artifact` | char/weapon/artifact name | - | Evaluates to the key for the specified char/weapon/artifact name. See the relevant character/weapon/artifact page for acceptable names. |
| `state` | - | - | - | Evaluates to the current state of the player. | 
| character name | `cons` | - | - | Evaluates to the character's constellation count. |
| character name | `energy` | - | - | Evaluates to the character's current energy. |
| character name | `energymax` | - | - | Evaluates to the character's maximum energy. |
| character name | `normal` | - | - | Evaluates to the character's next normal counter. Example: If the character is at N1, then the next normal counter is `1` (N2). |
| character name | `onfield` | - | - | `1` if the character is on the field, `0` otherwise. |
| character name | `weapon` | - | - | Evaluates to the character's weapon. Use `.keys.weapon.<weapon name>` for comparison purposes.
| character name | `mods`/`status` | mod/status name | - | Evaluates to the remaining duration of the mod/status on the character. See the relevant character page for acceptable mod/status names. | 
| character name | `infusion` | infusion name | - | Evaluates to the remaining duration of the weapon infusion on the character. See the relevant character page for acceptable infusion names. |
| character name | `tags` | tag name | - | Evaluates to the value of the tag on the character. See the relevant character page for acceptable tag names. |
| character name | `stats` | `def%`/`def`/`hp`/`hp%`/`atk`/`atk%`/`er`/`em`/`cr`/`cd`/`heal`/`pyro%`/`hydro%`/`cryo%`/`electro%`/`anemo%`/`geo%`/`dendro%`/`phys%`/`atkspd%`/`dmg%` | - | Evaluates to the value of the stat on the character. | 
| character name | `skill`/`burst`/`attack`/`charge`/`high_plunge`/`low_plunge`/`aim`/`dash`/`jump`/`swap`/`walk`/`wait` | `cd`/`charge`/`ready` | - | Evaluates to the following things for the specified action of the character: remaining cooldown / remaining charges (example: Sucrose Skill) / `1` if the action is ready, `0` otherwise. |
