---
title: Understanding Config Files
sidebar_position: 2
---

:::info
For more information, visit the [config file page under references](/reference/config).
:::

The gcsim config file contains all the information necessary for the simulator to run a simulation. The config file can be roughly broken down into the following parts:

- simulator options
- character setting
- enemy, energy, and starting character setting
- gcsim script (gcsl)

:::note
There is no requirement that these parts need to appear in the above order in a config file (or even stay together for that matter). 

The simulator will read each line sequentialy and interpret their meaning without regards to where they are in the file. So you can very well have the script at the top and everything else at the bottom if you wish.
:::

Let's break them down part by part.

## Simulator Options

Simulator options set some global settings for gcsim to run. This includes things such as how long to run each simulation for, the number of iterations, and other global settings.

At the very minimum, your config should have `iteration` and `duration` set. this tells the simulator how long to run each simulation for and how many simulations to run. This can be done by adding the following line to the config file:

```
options iteration=1000 duration=90;
```

This line will tell the simulator to run 1000 simulations, and each simulation should last 90s. 

:::tip
Duration can be decimals such as `duration=90.5` to simulate fractional seconds.
:::

More commonly, you will see options that looks like the following:

```
options iteration=1000 duration=90 swap_delay=12;
```

The `swap_delay` flag tells the simulator how long character swapping should take in frames (60 frames = 1 second). If this flag is not set, then by default the simulator will assume that character swapping takes 1 frame.

:::tip
`swap_delay=12` is one of the more commonly used option for swap frames among the gcsim community. However, since in reality actual swap in game is heavily ping dependent, you can feel free to change this to fit your needs.

A swap delay of 12 corresponds to roughly 100ms ping in game. You can calculate delay in frames from ping roughly as follows:

$$
\text{delay in frames} = \frac{ \text{ping in ms} * 2 * 60 }{1000}
$$
:::

:::note
You may also see older sims with options such as `debug=true` or `mode=sl`. These correspond to older simulator options that no longer apply in the current version and are ignored by the gcsim parser.
:::

## Character Settings

Character settings consist of various syntax to provide basic character information to the simulator such as character level, constellation, talents, weapon, artifact set, and artifact stats.

Character data consists of four parts:

- Basic settings such as level, cons, talents
- Weapon settings
- Artifact set settings
- Artifact stats settings

Let's break them down one by one.

### Basic Settings

Character basic setting can be commonly set with the following line:

```
bennett char lvl=70/80 cons=6 talent=6,8,8;
```

Here we tell the simulator to add Bennett to our team, with the given level, constellation, and talent levels. We do not need to specify the base stats for the character as that's calculated automatically based on the character's level and ascension (specified by `lvl=70/80`)

### Weapon Settings

Weapon settings tell the simulator what weapon to use for a given character:

```
bennett add weapon="thealleyflash" refine=1 lvl=90/90;
```

Weapons `refine` and `lvl` are required. The simulator will automatically calculate all base stats and bonuses of the weapon.

:::note
Despite the syntax being `add weapon`, you cannot add more than 1 weapon per character.
:::

### Artifact Sets

Artifact sets are specified as follows:

```
bennett add set="noblesseoblige" count=4;
```

You can specify any combination of sets as long as the total count does not exceed 5. The simulator will automatically handle the set bonuses based on the count.


### Artifact Stats

Artifact stats can be added as follows:

```
bennett add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cr=0.311;
bennett add stats hp=717 hp%=0.058 atk=121 atk%=0.635 def=102 em=42 er=0.156 cr=0.128 cd=0.265 ; 
```

There is no limit on how many lines of `add stats` you can have for each character. You can have as many as you like. The above example splits the stats between main stats from the artifacts and a total of the sub stats. You'll see this format commonly among theorycrafters that use theoritical artifact stats in their simulations.

:::important
You only need to add artifact stats (i.e the numbers that shows up on each artifact including main stat and sub stats) to each character. Do not add any character/weapon base stats or any other bonuses such as set bonuses. Those are calculated automatically.

Unfortunately, you cannot use the attack from the artifact summary page because it does not break down stats between ATK% and flat ATK. You must add the Attack stats individually from each artifact.
:::

:::tip
Contrary to what was stated above, you can use stats from the artifact summary page for simulating your only team **IF AND ONLY IF you do not change the character's weapon**.

To do so, simply add up the Attack shown on the artifact summary page.

However, as soon as you change weapons (causing character base stat to change), you must re-enter everything.

This also applies to DEF and HP if you change the Character's Level.
:::


A more practical way of adding artifact stats for simulating your own character is to split them up between each artifact such as follows:

```
yoimiya add stats hp=4780 atk=14 cr=0.039 cd=0.35 em=21;
yoimiya add stats atk=311 atk%=0.14 cd=0.062 cr=0.105 def=16;
yoimiya add stats atk%=0.466 er=0.058 cd=0.249 hp=299 hp=0.099;
yoimiya add stats pyro%=0.466 atk%=0.146 cd=0.257 er=0.058 def=23;
yoimiya add stats cd=0.622 atk=45 er=0.058 hp=508 cr=0.093;
```

You can even split up main and substats and make use of comments to quickly test out the effects of different main stats.


### Optional Params

Character, weapon, and sets can all have an optional `+params` setting. This is used by some character/weapons/artifact to customize certain settings. For example:

```
noelle add weapon="serpentspine" lvl=90/90 refine=5 +params=[stacks=5]
```

This would allow Serpent Spine to start with 5 stacks. For details, see [reference section](/reference) under each character/weapon/artifact.

### Example

Following is an example of a team of 4 characters:


```
bennett char lvl=90/90 cons=6 talent=9,9,9; 
bennett add weapon="thealleyflash" refine=1 lvl=90/90;
bennett add set="crimsonwitchofflames" count=4;
bennett add stats hp=4780 atk=311 em=187 pyro%=0.466 cr=0.311; # main
bennett add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944;

xiangling char lvl=90/90 cons=6 talent=9,9,9; 
xiangling add weapon="thecatch" refine=5 lvl=90/90;
#xiangling add weapon="deathmatch" refine=1 lvl=90/90;
#xiangling add weapon="dragonsbane" refine=3 lvl=90/90;
xiangling add set="emblemofseveredfate" count=4;
xiangling add stats hp=4780 atk=311 em=187 pyro%=0.466 cr=0.311; # main
xiangling add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=79.28 cr=0.331 cd=0.7944;

yelan char lvl=90/90 cons=0 talent=9,9,9; 
yelan add weapon="favoniuswarbow" refine=3 lvl=90/90;
yelan add set="noblesseoblige" count=4;
yelan add stats hp=4780 atk=311 hp%=0.466 hydro%=0.466 cr=0.311; # main
yelan add stats def=39.36 def%=0.124 hp=507.88 hp%=0.1984 atk=33.08 atk%=0.0992 er=0.1102 em=39.64 cr=0.331 cd=0.7944;																															

xingqiu char lvl=90/90 cons=6 talent=9,9,9; 
#xingqiu add weapon="amenomakageuchi" refine=5 lvl=90/90;
xingqiu add weapon="harbingerofdawn" refine=5 lvl=90/90;
xingqiu add set="emblemofseveredfate" count=4;
xingqiu add stats hp=4780 atk=311 atk%=0.466 hydro%=0.466 cr=0.311; # main
xingqiu add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944;
```

:::tip
Note that under the Xiangling block, we have a couple weapons commented out

```
xiangling add weapon="thecatch" refine=5 lvl=90/90;
#xiangling add weapon="deathmatch" refine=1 lvl=90/90;
#xiangling add weapon="dragonsbane" refine=3 lvl=90/90;
```

You can make use of this to quickly switch between different weapons to test out their effect on the overall team dps.
:::

## Target/Enemy Setting

Every configuration must have at least one `target` (enemy) set. Otherwise the simulator will not run. Targets can be added as follows:

```
target lvl=100 resist=0.1;
```

This will add a target that's level 100 with 10% resistance across the board. You can customize the resistance (see [here](/reference/config#enemy)).

To simulate multi-target simulation, simply repeat the target line as many times as you desire.

```
target lvl=100 resist=0.1 pos=0,0;
target lvl=100 resist=0.1 pos=1,0;
```

:::important
Make sure to specify the position of the target when simulating multiple targets, otherwise they will all stack on top of each other at (0, 0), causing certain abilities with very small AoEs (e.g. Xingqiu Rainswords) to hit multiple targets when that's not realistically going to happen.
:::

## Starting Character

Starting character can be set as follows:

```
active raiden;
```

This setting must be present in every config. This should only be set once and tells the simulator who is the starting character on the team.

:::note
Repeating `active <char>` multiple times or setting `active <char>` to different characters throughout the script will not cause the simulator to execute swaps. If multiple `active <char>;` lines are present, the simulator will simply use the last one to determine who is the starting character.

To execute a swap manually, see [here](#swaps).
:::

## Energy setting

Optionally, you can add a method of obtaining energy from enemies. The first, and easier way to to simply make energy drop every so often with the below syntax.

```
energy every interval=480,720 amount=1;
```

`interval` tells the sim how often to drop a particle. In this case, it drops a particles a random number between 8 and 12s after the last time it dropped a particle. (480 and 720 frames respectivele) `amount` tells it to drop one clear particle. 

The second way to HP based drops. When you declare an enemy, you can also tell it to drop energy after it takes a certain amount of damage. 

```
target lvl=100 resist=0.1 particle_threshold=250000 particle_drop_count=1;
```

`particle_threshold` tells the sim to drop a particle after the enemy takes 250000 damage while `particle_drop_count` tells it to drop 1 clear particle. 

For reference, a level 100 Maguu Kenki in Abyss would be

```
target lvl=100 resist=0.1 particle_threshold=460000 particle_drop_count=3;
```

## gcsim Script (`gcsl`)

The gcsim script is the core of the config file. It tells the simulator what actions to execute. At the very basic level, it's simply a written down instruction of the buttons you would press in game. For example:

```
active raiden;

raiden skill;
xingqiu skill, burst;
bennett skill, burst;
xiangling burst, skill;
raiden burst;
raiden attack:5, jump;
raiden attack:5, jump;
raiden attack:5, jump;
```

This represents a very basic rotation that started with Raiden as the active character. Raiden will use her skill, followed by Xingqiu skill, burst, etc.

:::tip
We used the short form `raiden attack:5, jump;` in the above example. This is simply a shortcut for `raiden attack, attack, attack, attack, attack, jump`.

This short form will work for any actions, such as `raiden jump:2;` if you wish to jump twice for whatever reason.
:::

Commands written in gcsl are executed sequentially in the order they appear in the config file. Each line must finish executing before the next can be executed.

:::important
When attempting to execute a character's ability and that ability is not ready (due to cooldown, energy, or stamina), the simulator will wait on that line and keep trying to execute that action until it has succeeded before moving onto the next action.

For example, if you were to do `xingqiu burst` and Xingqiu does not have energy, the simulator will keep trying to use burst (and failing) until Xingqiu finally has enough energy to use burst. In situations where Xingqiu does not have enough energy and there is no additional energy coming in the form of particles, this can cause the simulation to stall for the remaining duration.
:::


### Swaps

By default, you do not need to input swap commands between characters. In our original example, we had:

```
active raiden;

raiden skill;
xingqiu skill, burst;
```

After executing `raiden skill`, the simulator will recognize that the next action to be executed requires us to be on Xingqiu but Raiden is currently active. It will then automatically insert a swap action before executing `xingqiu skill`.

You can however manually force swaps with `<char> swap`. Here the `<char>` should represent the character you want to swap to.

### Repeating a rotation

A common usage is to repeat certain actions (or a rotation) multiple times throughout the duration of a simulation. 

#### Naive approach
One way to do this is to simply copy and paste the same block of actions (a rotation) the number of times you wish to repeat the rotation.

#### While loop
Another way is to wrap the entire block in a while loop:
```
while 1 {
    raiden skill;
    xingqiu skill, burst;
    bennett skill, burst;
    xiangling burst, skill;
    raiden burst;
    raiden attack:5, jump;
    raiden attack:5, jump;
    raiden attack:5, jump;
}
```

This will execute the block until gcsim stops running. 

:::info
See the [config page in the reference section](/reference/config) for information on when gcsim stops running.
:::

#### For loop
The entire block can also be wrapped in a for loop to repeat it a specific number of times:
```
for let i = 0; i < 5; i = i + 1 {
    raiden skill;
    xingqiu skill, burst;
    bennett skill, burst;
    xiangling burst, skill;
    raiden burst;
    raiden attack:5, jump;
    raiden attack:5, jump;
    raiden attack:5, jump;
}
```

This loop will execute the block 5 times.

### Advanced

gcsl features most of the basic functionality any scripting language has, including:

- variables
- arithmetics
- logical operations
- functions
- loops
- conditional statements

:::info
For more details, check out the [config page in the reference section](/reference/config#gcsl).
:::
