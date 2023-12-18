---
title: Config File
sidebar_position: 2
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
| `iteration`| Number of iterations to run gcsim for. | 1000 |
| `duration` | Duration to run gcsim for (in seconds). Fractional duration is allowed, for example: 11.5. In this case, gcsim will run until the duration has passed or there are no more actions to perform. This option is ignored if any `target` has `hp` specified. In that case, gcsim will run until all enemies are dead. | 90 |
| `swap_delay` | Number of frames it takes to swap characters. | 1 |
| `workers` | Number of workers to use. Only valid when using cli, ignored in web. | 20 |
| `hitlag` | Whether hitlag should be enabled. See the [hitlag page](/mechanics/hitlag) for more details. | true |
| `defhalt` | Whether to enable `canBeDefenseHalt` for hitlag. See the [hitlag page](/mechanics/hitlag) for more details. | true |

### Set energy generation

Example: 
```
energy every interval=480,720 amount=1;
```
This means that gcsim will generate 1 clear elemental particle every 480 to 720 frames randomly.

:::note
If multiple `energy every` lines are added, then the values specified by the final one will be used.
:::

:::info
Generating energy only at a specific frame for one time can also be specified. 
Multiple `energy once` lines can be added to spawn particles at different points in time.

Example:
```
energy once interval=300 amount=1;
```
This drops 1 clear elemental particle once at frame 300.
:::

### Set hurt generation 

Example: 
```
hurt every interval=480,720 amount=1,300 element=physical;
```
This means that gcsim will deal between 1 and 300 physical damage to the active character every 480 to 720 frames randomly.

:::note
If multiple `hurt every` lines are added, then the values specified by the final one will be used.
:::

:::info
Generating hurt only at a specific frame for one time can also be specified. 
Multiple `hurt once` lines can be added to deal damage at different points in time.

Example:
```
hurt once interval=300 amount=1,300 element=physical;
```
This will deal between 1 and 300 physical damage to the active character once at frame 300.
:::

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
| `start_hp` | Set the character's starting hp. | -1 (Character's max hp). |
| `start_hp%` | Set the character's starting hp ratio. | -1 (Character's max hp). Values should be 1 <= `start_hp%` <= 100. |
| `start_energy` | Set the character's starting energy. | Character's max energy. |

:::info
Some details about `start_hp` and `start_hp%`:
- a value <= 0 for both (manually supplied or by omission) will mean that the character's hp is set to max
- `start_hp%` can only be set as a percentage without decimal places, so 50 or 49 but not 49.5.
- the two params work additively so supplying both with a value > 0 will add them together

Example:
- `start_hp` is 10
- `start_hp%` is 49
- the sim will set the character's starting hp to be 49% of max hp + 10 flat hp
:::

:::info
Example: 
```
bennett char lvl=70/80 cons=2 talent=6,8,8 +params=[start_hp=10,start_hp%=49,start_energy=20];
```
This will set Bennett's starting hp to 49% + 10 and starting energy to 20.
:::

:::caution
There is no sanity check on these params. 
If you set this to a negative number or a really large number, behaviour is undefined. 
:::

### Add enemies

Example:
```
target lvl=88 resist=0.1 pos=0,0 radius=2 freeze_resist=0.8 hp=9999 particle_threshold=200 particle_drop_count=2;
```

| name | description | default |
| --- | --- | --- |
| `lvl` | Level of the enemy. | 0 |
| `resist` | Resistance to all types of elemental damage. Percentage represented as a decimal value. | 0 |
| `pyro`/`hydro`/`anemo`/`electro`/`dendro`/`cryo`/`geo`/`physical` | Resistance to the specified elemental damage. Percentage represented as a decimal value. | 0 |
| `pos` | Position of the enemy as x,y. | 0,0 |
| `radius` | The radius of the enemy's circle [hurtbox](https://en.wiktionary.org/wiki/hurtbox) in meters. | 1 |
| `freeze_resist` | How much freeze resistance the enemy has. `0` means no freeze resistance, `1` means immune to being frozen. The reaction still happens though. | 0 |
| `hp` | HP of the enemy. If this is set, duration in the sim options will be ignored and the sim will run until all enemies have died. If `hp` is set for at least one enemy, then it has to be set for all enemies. | - |
| `particle_threshold` | Only available if the `hp` is set. Determines after how much damage the enemy drops clear elemental particles. Example: If the enemy has 500 HP and this is set to 200, then the enemy will drop particles at 300 and 100 HP. | - |
| `particle_drop_count` | Only available if the `hp` is set. Number of clear elemental particles to drop at `particle_threshold`. | - |

:::danger
All configs must have at least one enemy specified. Otherwise you will get an error. 
:::

:::info
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

- `<condition>` can be any expression that evalutes into a number. 
- `<condition>` is considered false if it evaluates to 0 and true otherwise.

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

- A case is executed if the switch expression equals the case expression. 
- There is no `break;` at the end of each case. By default, once a case finishes evaluating, the switch statement will exit. 
- The exception to this is if a fallthrough is present. This will cause the case immediately below the current case to be executed as well.
- The `default` case is executed if none of the cases equals the switch expression. If no `default` is present, the switch will simply exit.

### `while` statement

`while` statement takes the following format:

```
while <condition> {

}
```

- `<condition>` can be any expression that evalutes into a number. 
- `<condition>` is considered false if it evaluates to 0 and true otherwise.
- `while` will repeat the block for as long as `<condition>` evalutes true.

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

### `for` statement

`for` statement takes the following format:

```
for <init>; <condition>; <post> { 
    
}
```

- `<init>` must be a variable initialization.
- `<condition>` can be any expression that evalutes into a number. 
- `<condition>` is considered false if it evaluates to 0 and true otherwise.
- `<post>` must be a variable assignment without a `;`.
- `for` will repeat the block for as long as `<condition>` evalutes true.

:::info
Example:
```
for let i = 0; i < 5; i = i + 1 {
    xiangling attack;
}
```
This will execute `xiangling attack;` 5 times.
:::

### `break` statement

`break` immediately exits the innermost enclosing `while` loop, `for` loop or `switch` statement. 

:::info
Example:
```
let i = 0;
while 1 {
    if i == 1 {
        break;
    }
    xiangling attack;
    i = i + 1;
}
```
This will execute `xiangling attack;` one time, because in the second iteration it will execute the `break;` statement and thus exit the `while` loop.
:::

### `continue` statement

`continue` skips to the next iteration of a `while` or `for` loop.

:::info
Example:
```
for let i = 0; i < 5; i = i + 1 {
    if i == 0 {
        continue;
    }
    xiangling attack;
}
```
This will skip the very first iteration of the `for` loop and execute `xiangling attack;` 4 times.
:::