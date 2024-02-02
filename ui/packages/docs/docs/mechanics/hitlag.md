---
title: Hitlag
sidebar_position: 1
---

Hitlag describes the freeze frames of a character and/or enemy that occurs when an attack hits said enemy. 

:::caution
Enemy attacks can inflict hitlag on the character and gadgets, but gcsim does not simulate enemy attacks or gadgets experiencing hitlag.
:::

## Hitlag attributes

Hitlag frames are determined by five different attributes:

| Name | Description |
| --- | --- |
| `hitHaltTime` | Number of seconds that hitlag lasts at base. | 
| `hitHaltTimeScale` | The rate at which the time of hitlag-affected entities flows during hitlag. Given as a percentage in decimal format. |
| `canBeDefenseHalt` | Boolean value that determines whether `hitHaltTime` is extended by 0.06s, when the enemy hit does not have their poise broken. |
| `deployable` | Boolean value that determines whether the hitlag applies to the character or not. The hitlag will always apply to the enemy/gadget hit by the character. This value is determined experimentally. | 
| `framerate` | Number of frames that represent 1s. Used to convert hitlag in seconds into hitlag in frames. |

:::info
For simplicity, gcsim assumes that time during hitlag is frozen. 
This is because the values chosen for `hitHaltTimeScale` and the time spent in hitlag is usually so low, that it is negligible. 
Because of this, the following calculations ignore `hitHaltTimeScale`.
:::

:::caution
Enemy poise is not simulated, so if `canBeDefenseHalt` is true for an attack, then gcsim will always add 0.06s. 
Most bosses and other heavy enemies like Ruin Guards cannot have their poise broken.
:::

## Hitlag calculation

$$ 
\text{hitlag in frames(\text{framerate})} = 
\begin{cases}
    \lceil (\text{hitHaltTime} + 0.06) * \text{framerate} \rceil &\text{if } \text{canBeDefenseHalt} = \text{true} \land \text{enemy is not poise broken}\\
    \lceil \text{hitHaltTime} * \text{framerate} \rceil &\text{if } \text{canBeDefenseHalt} = \text{false} \lor (\text{canBeDefenseHalt} = \text{true} \land \text{enemy is poise broken})
\end{cases}
$$

### Example: Keqing N1
#### Calculation

For example, Keqing's N1 has the following values for the previously described attributes:

| Name | Value |
| --- | --- |
| `hitHaltTime` | 0.03 | 
| `hitHaltTimeScale` | 0.01 |
| `canBeDefenseHalt` | true |
| `deployable` | false |
| `framerate` | 60 |

Inserting the values into the function described in the previous section leads to the following result:

$$
\text{hitlag in frames(60)} = 
\begin{cases}
    \lceil (0.03 + 0.06) * 60 \rceil = 6 &\text{if } \text{canBeDefenseHalt} = \text{true} \land \text{enemy is not poise broken}\\
    \lceil 0.03 * 60 \rceil = 2 &\text{if } \text{canBeDefenseHalt} = \text{false} \lor (\text{canBeDefenseHalt} = \text{true} \land \text{enemy is poise broken})
\end{cases}
$$
 
#### Ingame

<iframe width="560" height="315" src="https://www.youtube.com/embed/kI6N3Mn5BQY?start=20" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>

By looking at Keqing N1 hitlag footage against a tree and counting the number of hitlag frames for the different cases, we can confirm that the calculation matches ingame.
Hitlag starts 1 frame after the character significantly slows down (significant camera shake/movement can also be used) and ends once the character starts moving again.

:::note
Trees have poise and they can have their poise broken.
:::

- not poise broken: $1260 - 1254 = 6$ 
- poise broken: $1522 - 1520 = 2$

## Hitlag extension

The vast majority of abilities (status/modifier in gcsim) can be hitlag extended. 
The main exception is deployable abilities like Xiangling's Pyronado and Bennett's Inspiration Field.

:::info
The current assumption for gcsim is that hitlag extension is equal to the hitlag frames. 
:::

## Where to find hitlag attributes for specific attacks

- [gcsim source code](https://github.com/genshinsim/gcsim)
- [reference section in these docs](/reference)
- [wiki page for hitlag data](https://genshin-impact.fandom.com/wiki/Hitlag/Data)
- [this sheet](https://docs.google.com/spreadsheets/d/1am9g03w1dEqHPzgBsD_-a-MqsC3fDkg-dw_wJnWrQAU/edit?usp=sharing)

:::info
You can look at:
- how much hitlag a character/enemy in gcsim experiences
- which modifiers are getting hitlag extended 

by turning on the `hitlag` log category in the Sample tab settings menu!
:::
