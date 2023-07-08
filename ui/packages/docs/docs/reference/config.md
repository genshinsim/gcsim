---
sidebar_position: 1
title: Config File
---

## Config File Reference

### Set sim options

Options can be set as follows:

```
options iteration=1000 duration=90 swap_delay=14;
```

#### Valid options

| name | description | default |
| --- | --- | --- |
| `defhalt` | Whether to enable `canBeDefenseHalt` for hitlag. See the [hitlag page](/mechanics/hitlag) for more details. | true |
| `hitlag` | Whether hitlag should be enabled. See the [hitlag page](/mechanics/hitlag) for more details. | true |
| `duration` | Number of iterations to run gcsim for. | 1000 |
| `iteration` | Duration to run gcsim for (in seconds). Fractional duration is allowed, for example: 11.5. In this case, gcsim will run until the duration has passed or there are no more actions to perform. This option is ignored if any `target` has `hp` specified. In that case, gcsim will run until all enemies are dead. | 90 |
| `workers` | Number of workers to use. Only valid when using cli, ignored in web. | 20 |
| `swap_delay` | Number of frames it takes to swap characters. | 1 |

### Set energy generation

Example: 
```
energy every interval=480,720 amount=1;
```

This means that gcsim will generate 1 clear elemental particle every 480 to 720 frames randomly.

### Perform character, weapon, artifact setup

Character data can be roughly broken into 4 parts:

- `<name> char` data such as level, cons and talents
- `<name> add weapon=<weapon name>` data such as refine and level
- `<name> add set=<set name>` data such as count 
- `<name> add stats` for any character stats

For example:

```
bennett char lvl=70/80 cons=2 talent=6,8,8 +params=[a=1];
bennett add weapon="favoniussword" refine=1 lvl=90/90 +params=[b=2];
bennett add set="noblesseoblige" count=4 +params=[c=1];
bennett add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cr=0.311; # main
bennett add stats hp=717 hp%=0.058 atk=121 atk%=0.635 def=102 em=42 er=0.156 cr=0.128 cd=0.265; # subs
```

:::danger
With the exception of the stats (i.e. `hp`, `atk`, etc...), all other fields not starting with a `+` are mandatory.
:::danger

:::info
An optional param flag may be added to the character/weapon/artifact set via the `+params` flag. This optional param is defined by each character/weapon/artifact set.
:::

#### Optional global character params

| name | description | default |
| --- | --- | --- |
| `start_hp` | Set the character's starting hp. | Character's max hp. |
| `start_energy` | Set the character's starting energy. | Character's max energy. |

:::info
Example: 
```
bennett char lvl=70/80 cons=2 talent=6,8,8 +params=[start_hp=10,start_energy=20];
```
This will set Bennett's starting hp to 10 and starting energy to 20.
:::

:::caution
There is no sanity check on these params. 
If you set this to a negative number or a really large number, behaviour is undefined. 
:::

### Add enemies

Example:
```
target lvl=88 resist=0.1 pos=0,0;
```

:::danger
All configs must have at least one enemy specified. Otherwise you will get an error. 
:::

:::info
You can also specify each resist separately:
```
target lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=0.1 cryo=0.1 dendro=0.1 pos=0,0;
```
:::

:::info
Target starting position can be specified with `pos=x,y`. 
Note that if no position is provided, the target will default to (0, 0). 
If you have multiple targets, make sure to set their starting position properly. 
Otherwise you may get unintended behaviour such as otherwise single target abilities hitting multiple targets.
:::

:::info
To add multiple enemies, simply repeat the target line. 
Each enemy does not have to have the same level/resistance/position.

For example:
```
target lvl=100 resist=0.1;
target lvl=88 resist=0.05;
```

This would add two targets (making it a multi target simulation). 
Each target has different level and resistance.
:::

### Set the active character

Example:

```
active xiangling;
```

:::danger
All configs must have an active character specified. Otherwise you will get an error. 
:::

### Add comments

Any text following either a `//` or a `#` is treated as a comment and will be ignored until the end of the line. There are no multiline comments.

## gcsl

gcsim language or gcsl is the script that the simulator will run. This script tells the simulator what actions to execute. The scripting language includes some basic functionalities such as variables, conditionals, and loops.

### Character actions

### Primitive types and operators

The only primitive types are numbers. For purpose of boolean type conditionals, 0 is considered false and anything other than 0 is true.

The following are valid operators:

Math operations

- `+`: plus
- `-`: minus
- `*`: multiply
- `/`: divided

Comparison operators

- `==`: equal
- `!=`: does not equal
- `<`: less than
- `<=`: less than or equal to
- `>`: greater than
- `>=`: greater than or equal to

Logic operators

- `&&`: and
- `||`: or
- `!`: not

The one exception to this is that some system functions maybe take strings as inputs. However, these strings cannot be assigned to variables or manipulated otherwise.

### Fields

Fields are special syntax for accessing in-sim data during an iteration. Fields always starts with a `.`. For example, `.energy.xiangling` will evaluate to the current energy for `xiangling`.

All fields must evaluate to a number. For fields that are true or false, they will evaluate to 1 if true and 0 if false.

For more details on the available fields, see: [Fields](fields)

### Variables

Variables can only store numbers. Variables must be first declared before they can be used. Declaration starts with a `let` statement:

```
let x = 1;
```

Note that declared variables must also be initialized (i.e. assigned a value). The following is NOT valid:

```
let x; // this will error
```

Once declared, variables can be assigned other values with the assignment operator.

```
let x = 1;
x = 5; // x now has the value 5
```

Variables can be used in expressions:

```
let y = 2;
let x = 3;
let z = x + y; //z is now 5
```

Variable is subject to scoping in blocks such as `if`, `while`, `switch`, `apl` and functions:

```
let x = 1;
if 1 {
    let x = 2;
    let y = x + 1; //y will be 3 here
}
```

### `if` statement

`if` statement takes the following format:

```
if <condition> {

} else {

}
```

`<condition>` can be any expression that evalutes into a number. `condition` is considered false if it evaluates to 0 and true otherwise.

### `while` statement

`while` statement takes the following format:

```
while <condition> {

}
```

`<condition>` can be any expression that evalutes into a number. `condition` is considered false if it evaluates to 0 and true otherwise.

`while` will repeat the block for as long as condition evalutes true.

:::caution Infinite loops

Be careful when using infinite loops. gcsim does not have a way to detect infinite loops. An infinite loop that never exits will cause the simulation to hang with no noticeable error.

The exception to this is if there is an action (or `wait(x)`) inside an infinite loop, for example:

```
while 1 {
    xiangling attack;
}
```

This is ok because any time a character action (or wait) is executed, the evaluation of the script is paused and the simulation takes over. Since the simulation itself has an exit condition (i.e. duration), this infinite loop will properly be terminated once the simulation reaches its exit condition.

However, the following will cause the simulation to hang forever because script execution is never paused, so the simulator never gets a chance to check its exit conditions:

```
while 1 {
    print("hi");
}
```

`wait` is a special case in that it behaves just like a character action. So the following will exit properly according to the simulations exit conditions:

```
while 1 {
    wait(1);
}
```

Also be careful of infinite loops that seems like it contains a character action but may not actually ever evaluate it. For example:

```
while 1 {
    if 0 {
        //this block will never be reached!!
        xiangling attack;
    }
}
```

In this example the `xiangling attack;` can never be reached, causing the script to never actually pause and therefore the simulation will never reach its exit conditon.

:::

### `switch` statement

`switch` statement takes the following format:

```
switch <expr> {
    case <expr>:
        //do action here
    case <expr>:
        //do action here
        fallthrough;
    case <expr>:
        //will continue from above
    default:
        //default case
}
```

A case is executed if the switch expression equals the case expression. There is no `break;` at the end of each case. By default, once a case finishes evaluating, the switch statement will exit. The exception to this is if a fallthrough is present. This will cause the case immediately below the current case to be executed as well.

The `default` case is executed if none of the cases equals the switch expression. If no `default` is present, the switch will simply exit.
