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

## execute_action

:::danger
**THIS FUNCTION IS EXPERIMENTAL AND SUBJECT TO CHANGE.**

**USE AT YOUR OWN RISK.**
:::

```
execute_action(char, action, params);
```

`execute_action` evaluates to null and is used by the sim to execute actions.
The intent behind this system function is to allow for proper typing/functional support in the future.
It being exposed here is an unintended side effect which can be used to implement a function that runs before every action.


:::danger
The following example is subject to breaking in the future!
:::

With that in mind it is possible to add (random) frame delays before each action:
```
fn rand_delay(mean, stddev) {
    let del = randnorm() * stddev + mean;
    if del > (mean + mean) {
        del = mean + mean;
    }
    delay(del);
}

let prev_char_id = -1;
let prev_action_id = -1;

let _execute_action = execute_action;
fn execute_action(char_id number, action_id number, p map) {
    print(prev_char_id, " ", prev_action_id, " ", char_id, " ", action_id);

    if action_id == .action.swap {
        # add delay before swap
        rand_delay(12, 3);
    } else if prev_action_id == .action.attack && action_id != .action.attack && action_id != .action.charge {
        # add delay after attack, but only if not followed by another attack or charge
        rand_delay(3, 1);
    } else if prev_action_id != .action.attack {
        # add delay to everything else
        rand_delay(3, 1);
    }

    prev_char_id = char_id;
    prev_action_id = action_id;
    return _execute_action(char_id, action_id, p);
}
```

:::danger
- `char` and `action` must be a number or an expression that evaluates to a number. 
- `params` must be a map or an expression that evaluates to a map.
:::

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

## set_swap_icd

```
set_swap_icd(arg1);
```
:::caution
- This function replicates behavior not found in typical gameplay. 
- By default, characters in Genshin cannot swap more than once per second. However, by 'booking' (opening the Adventurer's Handbook mid-combat), the swap timer can continue while other in-game timers (such as the Spiral Abyss timer) remain paused.
- If you use this function, the resulting dps will not represent damage per real time, but will instead represent damage per in-game time.
- Changing the swap icd will not affect the cooldown for any swap delays in progress; only future swaps will use the new icd.
- As this is one of the "Booking" functions, a Time Manipulation warning will be displayed on the results if used.
:::

- `set_swap_icd` will set the default swap ICD for all characters equal to the number of frames in `arg1`.
- This function will cause a `cooldown` log to appear in the sample.
- If `arg1` evaluates to a number that is less than 0, an error will be returned.
- `set_swap_icd` will always evaluate to 0.

:::danger 
`arg1` must be a number or an expression that evaluates to a number. 
:::

:::info 
Example:
Reduce swap ICD to 0, allowing future swaps to be able to occur as frequently as `swap_delay` allows.
```
set_swap_icd(0);
```

:::

## reduce_swap_cd

```
reduce_swap_cd(arg1, [optional arg2 delay, optional arg3 hitlag_delay]);
```
:::caution
- This function replicates behavior not found in typical gameplay. 
- By default, characters in Genshin cannot swap more than once per second. However, by 'booking' (opening the Adventurer's Handbook mid-combat), the swap timer can continue while other in-game timers (such as the Spiral Abyss timer) remain paused.
- If you use this function, the resulting dps will not represent damage per real time, but will instead represent damage per in-game time.
- Reducing the swap cd will not affect the cooldown for any future swaps; only the current swap cd will be affected.
- As this is one of the "Booking" functions, a Time Manipulation warning will be displayed on the results if used.
:::

- `reduce_swap_cd` will reduce the current swap cd by the number of frames in `arg1`. If the cd would be reduced to less than 0, it is capped at 0 instead.
- This function will cause a `cooldown` log to appear in the sample.
- If `arg1` evaluates to a number that is less than 0, an error will be returned.
- `arg2` will cause the execution of this function to be delayed by `delay` frames. If `arg3` is set, this delay will be further extended by hitlag. This function is non-blocking and will not cause a delay within the config. 
- `reduce_swap_cd` will always evaluate to 0.

:::danger 
- `arg1` must be a non-negative number or an expression that evaluates to a non-negative number. 
- `arg2` May be excluded, or must be a non-negative number or an expression that evaluates to a non-negative number.
- `arg3` May be excluded, or must be a non-negative number or an expression that evaluates to a non-negative number. `arg3` may not be included if `arg2` is excluded.
:::

:::info 
Example:
- Reduce swap cd by 40 frames, allowing the next swap to occur after as few as 20 frames pass.
```
reduce_swap_cd(40);
```
- Reduce swap cd by 40 frames, to be executed after 10 frames have elapsed. If a swap occurs between the calling of this function and the swap cd reduction, then the swap cd reduction will apply on this most recent swap.
```
set_swap_cd(40, 10);
```
- Reduce swap cd by 40 frames, to be executed after 10 frames have elapsed. If hitlag occurs between the calling of this function and the end of the 10f delay, then the reduction of the swap cd will be further delayed by the hitlag.
```
set_swap_cd(40, 10, 1);
```
:::

## reduce_crystallize_gcd

```
reduce_crystallize_gcd(arg1, [optional arg2 delay, optional arg3 hitlag_delay]);
```
:::caution
- This function replicates behavior not found in typical gameplay. 
- By default, characters in Genshin cannot crystallize any given enemy more than once per second. However, by 'booking' (opening the Adventurer's Handbook mid-combat), the gcd timer can continue while other in-game timers (such as the Spiral Abyss timer) remain paused.
- If you use this function, the resulting dps will not represent damage per real time, but will instead represent damage per in-game time.
- As this is one of the "Booking" functions, a Time Manipulation warning will be displayed on the results if used. 
:::

- `reduce_crystallize_gcd` will reduce the current swap cd by the number of frames in `arg1`. If the cd would be reduced to less than 0, it is capped at 0 instead.
- All enemies will have their crystallize gcd reduced when this function is called.
- This function will cause a `cooldown` log to appear in the sample.
- If `arg1` evaluates to a number that is less than 0, an error will be returned.
- `arg2` will cause the internal execution of this function to be delayed by `delay` frames. If `arg3` is set, this delay will be further extended by hitlag. This function will not cause a delay within the config. 
- `reduce_crystallize_gcd` will always evaluate to 0.

:::danger 
`arg1` must be a non-negative number or an expression that evaluates to a non-negative number. 
`arg2` May be excluded, or must be a non-negative number or an expression that evaluates to a non-negative number.
`arg3` May be excluded, or must be a non-negative number or an expression that evaluates to a non-negative number. `arg3` may not be included if `arg2` is excluded.
:::

:::info 
Example:
- Reduce crystallize gcd by 40 frames, allowing the next crystallize to occur after as few as 20 frames pass.
```
reduce_crystallize_gcd(40);
```
- Reduce crystallize gcd by 40 frames, to be executed after 10 frames have elapsed. If a crystallize occurs between the calling of this function and the crystallze gcd reduction, then the crystallize gcd reduction will apply on this most recent crystallize.
```
reduce_crystallize_gcd(40, 10);
```
- Reduce crystallize gcd by 40 frames, to be executed after 10 frames have elapsed. If hitlag occurs between the calling of this function and the end of the 10f delay, then the reduction of the crystallize gcd will be further delayed by the hitlag.
```
reduce_crystallize_gcd(40, 10, 1);
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

## is_target_dead

:::danger
**THIS FUNCTION IS EXPERIMENTAL AND SUBJECT TO CHANGE.**

**USE AT YOUR OWN RISK.**
:::

```
is_target_dead(arg);
```

- `is_target_dead` will evaluates to 1 if the target with index `arg` is dead and 0 otherwise.

:::danger
`arg` must be a number or an expression that evaluates to a number.
:::

:::danger
If `arg` is an invalid target (i.e. 3 when there are only 2 targets), then gcsim will exit with an error.
:::


## pick_up_crystallize

```
pick_up_crystallize(element);
```

- `pick_up_crystallize` will pick up the oldest crystallize shard with the specified `element` supplied as a string.
- `pick_up_crystallize` will not pick up any shard if:
 - no shard with the specified `element` exists 
 - there is a shard with the specfied `element`, but it cannot be picked up yet
- `pick_up_crystallize` will return the number of crystallize shards that were picked up (either 0 or 1).

:::info
`element` can also be "any" to pick up the oldest crystallize shard of any element.
:::

:::danger
`element` must be a string or an expression that evaluates to a string.
:::

## is_even

```
is_even(arg);
```

`is_even` evaluates if a given number is even or not. If a number is a floating point, the number if floored first. 

:::danger
`arg` must be a number or an expression that evaluates to a number.
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

:::danger
**THIS FUNCTION IS EXPERIMENTAL AND SUBJECT TO CHANGE.**

**USE AT YOUR OWN RISK.**
:::

```
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