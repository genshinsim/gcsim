---
title: Animation Speed
sidebar_position: 3
---

## Attack Speed

Animation Speed (ATK SPD) is a character stat that increases the animation speed of a character's Normal Attacks and thus decreases the amount of frames that they take. 
In gcsim, ATK SPD starts at 0 (normal animation speed, 100% ATK SPD) and caps at 0.6 (160% ATK SPD).

In gcsim, ATK SPD is modelled like this:

$$
\text{action frames with ATK SPD(action frames, Current ATK SPD)} = \text{action frames} - \lfloor \text{min}(\text{Current ATK SPD}, 0.1+\frac{\text{Current ATK SPD}-0.1}{2}) * \text{action frames}\rfloor
$$


:::danger
The current implementation of ATK SPD in gcsim has some flaws:
- ATK SPD currently snapshots in gcsim, but it does not snapshot ingame.
- ATK SPD does not adjust the hitmark timing in gcsim, but it does ingame. 
This makes some action cancels still take the same amount of time as without any ATK SPD.
- ATK SPD does not decrease frames in an intuitive way ingame. 
The current formula is an approximation that is not accurate for all characters.
:::

### Example

Eula's N1 if followed by N2 takes 34 frames at 60 fps.
If she currently has 130% ATK SPD, then her N1 will take this amount of frames:

$$
\text{action frames with ATK SPD(34, 0.3)} = 34 - \lfloor \text{min}(0.3, 0.1+\frac{0.3-0.1}{2}) * 34 \rfloor = 34 - \lfloor \text{min}(0.3, 0.2) * 34 \rfloor = 34 - \lfloor 0.2 * 34 \rfloor = 34 - \lfloor 6.8 \rfloor = 34 - 6 = 28
$$

## Overall Animation Speed

Overall Animation Speed is a character stat that is able to increase the animation speed of all the character's actions. 
Overall Animation Speed starts at 100% (normal animation speed) and caps at 140%.
This stat is usually applied conditionally based on which actions it should affect. 
For example, Itto's Charge Attack ATK SPD and Dehya's Flame-Mane's Fist strike speed increase is implemented this way.

:::danger
The implementation of this shares the same problems as ATK SPD.
:::

## Movement Speed

Movement Speed shares the same cap as Overall Animation Speed and it only affects the character's animation speed while walking, dashing, jumping or during certain alternative movement techniques like Sayu's Skill or Yelan's Skill.

:::danger
Movement Speed is not implemented. The only way the player can move in gcsim is via teleportation.
:::
