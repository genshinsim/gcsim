---
title: Fields
sidebar_position: 3
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
| `debuff` | `res`/`def` | `t0`/`t1`/`t2`/... | res/def modifier name | Evaluates to the remaining duration of the specified res/def modifier on the specified target. See the relevant character/weapon/artifact page for acceptable modifier names. |
| `element` |  `t0`/`t1`/`t2`/... | `pyro`/`hydro`/`anemo`/`electro`/`dendro`/`cryo`/`geo`/`frozen`/`quicken` | - | Evaluates to the remaining durability of the specified element on the specified target. |
| `status` | core status name | - | - | Evaluates to the remaining duration of the specified core status. See the relevant character/weapon/artifact page for acceptable status names. |
| `stam` | - | - | - | Evaluates to the player's remaining stamina. |
| `construct` | `duration`/`count` | construct name | - | Evaluates to the duration/count of the specified construct. See individual character page for acceptable construct names. |
| `gadgets` | `dendrocore` | `count` | - | Evaluates to the current number of Dendro Cores. |
| `gadgets` | `sourcewaterdroplet` | `count` | - | Evaluates to the current number of Sourcewater Droplets. Use character specific fields to get number of Sourcewater Droplets in range. |
| `gadgets` | `crystallizeshard` | `all`/`pyro`/`hydro`/`electro`/`cryo` | - | Evaluates to the current number of Crystallize Shards that can be picked up. `all` will return the total number of Crystallize Shards while the others will only count the ones of the given element. |
| `keys` | `char`/`weapon`/`artifact` | char/weapon/artifact name | - | Evaluates to the key for the specified char/weapon/artifact name. See the relevant character/weapon/artifact page for acceptable names. |
| `keys` | `element` | element name | - | Evaluates to the key for the specified element name. Element names are `electro`, `pyro`, `cryo`, `hydro`, `dendro`, `quicken`, `frozen`, `anemo`, `geo`, `physical`. |
| `action` |  `skill`/`burst`/`attack`/`charge`/`high_plunge`/`low_plunge`/`aim`/`dash`/`jump`/`swap`/`walk`/`wait` | - | Evaluates to the key for the specified action name. |
| `state` | - | - | - | Evaluates to the current state of the player. | 
| `previous-char` | - | - | - | Evaluates to the char that executed the previous action. Use `.keys.char.<char name>` for comparison. | 
| `previous-action` | - | - | - | Evaluates to the previously executed action. Use `.action.<action name>` for comparison. | 
| `airborne` | - | - | - | `1` if the player is airborne (via buffed jump from Xianyun Q for example), `0` otherwise. | 
| character name | `cons` | - | - | Evaluates to the character's constellation count. |
| character name | `energy` | - | - | Evaluates to the character's current energy. |
| character name | `energymax` | - | - | Evaluates to the character's maximum energy. |
| character name | `hp` | - | - | Evaluates to the character's current hp. |
| character name | `hpratio` | - | - | Evaluates to the character's hp ratio. |
| character name | `hpmax` | - | - | Evaluates to the character's maximum hp. |
| character name | `bol` | - | - | Evaluates to the character's current Bond of Life. |
| character name | `bolratio` | - | - | Evaluates to the character's Bond of Life ratio. |
| character name | `normal` | - | - | Evaluates to the character's next normal counter. Example: If the character is in idle or just performed their last normal attack, then the next normal counter is `1` (N1). |
| character name | `onfield` | - | - | `1` if the character is on the field, `0` otherwise. |
| character name | `weapon` | - | - | Evaluates to the character's weapon. Use `.keys.weapon.<weapon name>` for comparison purposes.
| character name | `mods`/`status` | mod/status name | - | Evaluates to the remaining duration of the mod/status on the character. See the relevant character page for acceptable mod/status names. | 
| character name | `infusion` | infusion name | - | Evaluates to the remaining duration of the weapon infusion on the character. See the relevant character page for acceptable infusion names. |
| character name | `tags` | tag name | - | Evaluates to the value of the tag on the character. See the relevant character page for acceptable tag names. |
| character name | `sets` | set name | - | Evaluates to the count of the set on the character. |
| character name | `stats` | `def%`/`def`/`hp`/`hp%`/`atk`/`atk%`/`er`/`em`/`cr`/`cd`/`heal`/`pyro%`/`hydro%`/`cryo%`/`electro%`/`anemo%`/`geo%`/`dendro%`/`phys%`/`atkspd%`/`dmg%` | - | Evaluates to the value of the stat on the character. | 
| character name | `skill`/`burst`/`attack`/`charge`/`high_plunge`/`low_plunge`/`aim`/`dash`/`jump`/`swap`/`walk`/`wait` | `cd`/`charge`/`ready` | - | Evaluates to the following things for the specified action of the character: remaining cooldown / remaining charges (example: Sucrose Skill) / `1` if the action is ready, `0` otherwise. |
| character name | `nightsoul` | `state` | - | `1` if the character is in Nightsoul's Blessing state, `0` otherwise. |
| character name | `nightsoul` | `point` | - | Evaluates to the character's Nightsoul points. |
| character name | `nightsoul` | `duration` | - | Evaluates to the duration of the character's Nightsoul Blessing state. |
