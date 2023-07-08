---
title: System Functions
sidebar_position: 3
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

## wait

```
wait(arg);
```

- `wait` is a special function that will ask gcsim to wait a number of frames. 
- `wait` will always evaluate to 0.

:::danger
`arg` must be a number or an expression that evaluates to a number and represents the number of frames the simulator will wait for.
:::

:::caution
Due to how gcsim handles actions, the current implementation of `wait` is not intuitive.

Example:
```
keqing attack;  // Keqing N1
wait(2);        // wait for 2 frames
keqing attack;  // Keqing N2
```

Many users would expect that gcsim waits for 2 frames *after* Keqing's N1 action ends. 
This is not how `wait` works.
`wait` makes the sim wait for 2 frames *after* Keqing's N1 action has reached its specified `CanQueueAfter` value. 
The duration of `wait` counts towards the action length.

- expected: 
    - N1 starts
    - N1 ends after 15 frames
    - gcsim waits for 2 frames 
    - N2 starts after a total of 15 + 2 = 17 frames
- reality: 
    - N1 starts
    - N1 `CanQueueAfter` is reached after 11 frames
    - gcsim waits for 2 frames
    - now there are 15 - (11 + 2) = 2 frames left in the N1 animation
    - N1 continues for 2 more frames until the N1 animation is over
    - N2 starts after a total of 11 + 2 + 2 = 15 frames

To make gcsim wait for 1 frame after Keqing's N1 action ends, the user would have to insert a `wait(5);`.
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
