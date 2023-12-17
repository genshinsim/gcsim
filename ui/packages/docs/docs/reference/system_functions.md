---
title: System Functions
sidebar_position: 4
---

The following system functions are available:

## print

```
print(arg1, arg2, arg3, ...);
```

- The `print` command allows the user to print any arbitrary expression. 
- The printed message can be viewed in the Sample tab of the viewer under the `user` log category. 
- There is no limit on the number of arguments that can be passed to `print`.
- `print` will always evaluate to 0.

:::info
As an exception, `print` can take both number and string as arguments. 
The following is valid:
```
print("this is a number: ", 1); // note the space after the :
```
`print` will evaluate each argument and concatenate them all into one string, then displaying it on the debug.
:::

:::info
Note that any valid expression that evaluates into a number can also be printed. For example:
```
print("bennett's current energy is: ", .bennett.energy);
```
:::

## sleep/wait

:::caution
The syntax `wait(arg)`, while valid, is deprecated in favour of `sleep(arg)`
:::

```
sleep(arg);
wait(arg); //deprecated, use sleep(arg)
```

- `sleep` is a special function that will ask gcsim to wait a number of frames. 
- `sleep` will always evaluate to 0.

:::danger
`arg` must be a number or an expression that evaluates to a number and represents the number of frames the simulator will wait for.
:::

:::caution
Due to how gcsim handles actions, the current implementation of `sleep` is not intuitive when trying to extend the duration of actions.
Please use `delay` for this purpose instead.

Example:
```
keqing attack;  // Keqing N1
sleep(2);       // sleeps for 2 frames
keqing attack;  // Keqing N2
```

Many users would expect that gcsim sleeps for 2 frames *after* Keqing's N1 action ends. 
This is not how `sleep` works.
`sleep` makes the sim sleep for 2 frames *after* Keqing's N1 action has reached its specified `CanQueueAfter` value. 
The duration of `sleep` counts towards the action length.

- expected: 
    - N1 starts
    - N1 ends after 15 frames
    - gcsim sleeps for 2 frames 
    - N2 starts after a total of 15 + 2 = 17 frames
- reality: 
    - N1 starts
    - N1 `CanQueueAfter` is reached after 11 frames
    - gcsim sleeps for 2 frames
    - now there are 15 - (11 + 2) = 2 frames left in the N1 animation
    - N1 continues for 2 more frames until the N1 animation is over
    - N2 starts after a total of 11 + 2 + 2 = 15 frames

To make gcsim sleep for 1 frame after Keqing's N1 action ends, the user would have to insert a `sleep(5);`.
:::

## delay

```
delay(arg);
```

- `delay` is a special function that will ask gcsim to delay the start of the following action by a number of frames. 
- `delay` will always evaluate to 0.

:::danger
`arg` must be a number or an expression that evaluates to a number and represents the number of frames the simulator will wait for.
:::

:::caution
`delay` is executed before the sim checks if the next action is ready. 

Example:
```
keqing burst;
delay(5);
keqing burst;
```

In this case, the sim would do the following:
- Keqing's 1st Burst is executed
- gcsim executes a `delay` for 5 frames at the end of the previous action
- Once the delay is over, gcsim checks if Keqing's 2nd Burst can be executed
- Since there is not enough energy, the sim will be stuck waiting for energy
- After enough particles were collected from energy drops, Keqing's 2nd Burst is executed
:::

:::caution
If the active character is affected by hitlag during the execution of `delay`, then it will last longer than specified.

Example:
```
noelle skill;
sleep(700);
delay(50);
noelle attack;
```

This example uses C4 Noelle to show a source of hitlag that can occur during `delay`.
The `sleep` is used so that the C4 shield explosion happens during `delay`.

- Noelle's Skill is executed
- gcsim will sleep for 700 frames after the `CanQueueAfter` of the previous action 
- gcsim starts executing a `delay` that should last 50 frames
- A few frames after `delay` starts, C4 Noelle applies 13 frames of hitlag
- Noelle's Attack is executed 50 + 13 = 63 frames after the start of `delay`
:::


## f

```
f();
```

`f` is function that takes no argument and will evaluate to the current frame number gcsim is on.

## rand

```
rand();
```

`rand` evaluates to an uniformly distributed random number between 0 and 1.

## randnorm

```
randnorm();
```

`randnorm` evaluates to a normally distributed random number with mean 0 and std dev of 1.

## type

```
type(arg);
```

`type` evaluates to the name of the gcsl type of `arg`.

## set_particle_delay

```
set_particle_delay(arg1, arg2);
```

- `set_particle_delay` will set the default particle delay for the character supplied in `arg1` to the value in `arg2`. 
- If `arg2` evaluates to a number that is less than 0, 0 will be used.
- `set_particle_delay` will always evaluate to 0.

:::danger
`arg1` must be a string (wrapped in double quotes) and `arg2` must be a number or an expression that evaluates to a number. 
:::

:::info
Example:
```
set_particle_delay("xingqiu", 100);
```
:::

## set_default_target

```
set_default_target(arg);
```

- `set_default_target` will set the default target to the index supplied by `arg`. 
- `set_default_target` will always evaluate to 0.

:::danger
`arg` must be a number or an expression that evaluates to a number. 
:::

:::danger
If `arg` is an invalid target (i.e. 3 when there are only 2 targets), then gcsim will exit with an error.
:::

:::info
For example, if there are 2 targets, then `set_default_target(2)` will set the default target to the 2nd one. 
Note that it starts at 1 and not 0 because 0 is a special case (target 0 represents the player).
:::

## set_player_pos

```
set_player_pos(x, y);
```

- `set_player_pos` will set the player's current position to the supplied `x` and `y` coordinate. 
- `set_player_pos` will always evaluate to 0.

:::danger
`x` and `y` must be a number or an expression that evaluates to a number.
:::

## set_target_pos

```
set_target_pos(arg, x, y);
```

- `set_target_pos` will set the target with index `arg` to the supplied `x` and `y` coordinates. 
- `set_target_pos` will always evaluate to 0.

:::danger
All arguments must be a number or an expression that evaluates to a number.
:::

:::danger
If `arg` is an invalid target (i.e. 3 when there are only 2 targets), then gcsim will exit with an error.
:::

## kill_target

```
kill_target(arg);
```

- `kill_target` will kill the target with index `arg`.
- `kill_target` will always evaluate to 0.

:::danger
`arg` must be a number or an expression that evaluates to a number.
:::

:::danger
If `arg` is an invalid target (i.e. 3 when there are only 2 targets), then gcsim will exit with an error.
:::

## sin

```
sin(arg);
```

`sin` evaluates to the sine of the given `arg`.

:::danger
`arg` must be a number or an expression that evaluates to a number.
:::

## cos

```
cos(arg);
```

`cos` evaluates to the cosine of the given `arg`.

:::danger
`arg` must be a number or an expression that evaluates to a number.
:::

## asin

```
asin(arg);
```

`asin` evaluates to the arcsine of the given `arg`.

:::danger
`arg` must be a number or an expression that evaluates to a number.
:::

## acos

```
acos(arg);
```

`acos` evaluates to the arccos of the given `arg`.

:::danger
`arg` must be a number or an expression that evaluates to a number.
:::

## set_on_tick

```nginx
set_on_tick(func);
```

`set_on_tick` evaluates to null and is a way to make the sim execute a user-defined function every frame.

In the following example, the player's stamina will be printed every frame:
```
fn stam() {
    print(.stam);
}
set_on_tick(stam);
```

:::danger
`func` must be a function or an expression that evaluates to a function.
:::
