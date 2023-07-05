---
sidebar_position: 4
title: System Functions
---

The following system functions are available

## print

```
print(arg1, arg2, arg3, ...);
```

The print command allow the user to print any arbitrary expression. The printed message can be viewed in the debug tab of the viewer under users setting. There is no limit on the number of arguments that can be passed to print.

As an exception, print can take both number and string as arguments. The following is valid

```
print("this is a number: ", 1); // note the space after the :
```

Print will evaluate each argument and concatenate them all into one string, then displaying it on the debug.

Note that any valid expression that evaluates into a number can also be printed. For example:

```
print("bennett's current energy is: ", .energy.bennett);
```

print will always evaluate to 0.

## wait

```
wait(arg);
```

wait is a special function that will ask the simulator to wait a number of frames. arg must be a number or an expression that evaluates to a number and represents the number of frames the simulator will wait for.

wait will always evaluate to 0.

## f

```
f();
```

f is function that takes no argument and will evaluate to the current frame number the simulation is on.

## rand

```
rand();
```

rand evalutes to an uniformly distributed random number between 0 and 1.

## randnorm

```
randnorm();
```

randnorm evalutes to a normally distributed random number with mean 0 and std dev of 1.

## set_particle_delay

```
set_particle_delay(arg1, arg2);
```

set_particle_delay will set the default particle delay for the character supplied in arg1 to value in arg2. arg1 must be a string (wrapped in double quotes) and arg2 must be a number or an expression that evalutes to a number. If arg2 evaluates to a number that is less than 0, 0 will be used.

set_particle_delay will always evaluate to 0.

Example:

```
set_particle_delay("xingqiu", 100);
```

## set_default_target

```
set_default_target(arg);
```

set_default_target will set the default target to the index supplied by arg. arg must be a number or an expression that evalutes to a number. For example, if there are 2 targets, then `set_default_target(2)` will set the default target to the 2nd one. Note that it starts at 1 and not 0 because 0 is a special case (target 0 represents the player).

Also, if arg is an invalid target (i.e. 3 when there are only 2 targets), the simulation will exit with an error.

set_default_target will always evaluate to 0.

## set_player_pos

```
set_player_pos(x, y);
```

set_player_pos will set the player's current position to the supplied x and y coordinate. x and y must be a number or an expression that evalutes to a number.

set_player_pos will always evaluate to 0.

## set_target_pos

```
set_target_pos(arg, x, y);
```

set_target_pos will set the target with index arg to the supplied x and y coordinates. All arguments must be a number or an expression that evalutes to a number.

Also, if arg is an invalid target (i.e. 3 when there are only 2 targets), the simulation will exit with an error.

set_target_pos will alwasy evaluate to 0.
