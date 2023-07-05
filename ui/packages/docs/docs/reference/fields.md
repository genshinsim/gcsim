---
sidebar_position: 3
title: Fields
---

Fields have the following structure:

```
.field1.field2.field3
```

All fields evaluate to a number.

The following are available fields.

:::tip

Most of the specific tags can be located in their respective character/weapon/artifact page. For example, if you are looking for the tag for Lisa's A4 (defense shred), look in Lisa's character page.

:::

<!-- prettier-ignore -->
| field1 | field2 | field3 | field4 | description |
| --- | --- | --- | --- | --- |
| `energy` | `character` | - | - | character's current energy |
| `cd`| `character` | `skill` `burst`|- | cooldown remaining on skill/burst in frames |
| `stam`| - | - | - |player's current stamina |
| `status` | see each character | - | - |these are character specialized statuses, usually used to keep track of buffs |
| `tags` | `character` | see each character - | - | these are character specialized tags as defined by each character|
| `debuff` | `res` or `def` | `t1` or `t2` etc...| `tag` | Evaluates to the remaining duration of a specified default `tag`. See the relevant character/weapon/artifact page for acceptable `tag` |
| `element` |  `t1` or `t2` etc... | `pryo` `hydro` `cryo` `electro` `frozen` `ec`| - | `1` if the specified element exists, `0` otherwise|
| `ready`| `character` | `skill` `burst`| - | shorthand to check both cooldown and energy is ready |
| `mods` | `character` | `tag` | - | Evaluate to remaining duration of stat mods active on the specified character. See individual character page for acceptable `tag` |
| `infusion` | `character` | `tag` | - | Evaluate to remaining duration of infusion on the specified character. See individual character page for acceptable `tag` |
| `construct` | `duration` or `count` | `tag` | - | Evalues to the duration or count of the specified `tag`. See individual character page for acceptable `tag` |
| `normal` | `character` | - | - | Evalutes to the next normal attack number in a combo for the specified character |
| `character` | `tag` | - | - | These are character specific fields. See individual character page for details |
