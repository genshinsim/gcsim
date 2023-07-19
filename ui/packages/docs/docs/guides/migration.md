---
title: Migrating from v0 to v1
sidebar_position: 8
---

# Migrate from v0 to v1 (Rewrite)

## TL;DR

From an end user's perspective, there are a few major changes that will prevent and existing sim from running

- `restart` is removed. Instead, to repeat the rotation, wrap the rotation in a `while` loop like the below example. (This also allows for abiltities that are only used in the first rotation to be used while still telling the sim to repeat):

```
while 1 {
    # rotation goes in here

}
```

- `wait 10;` has been changed to a function call. They should be replaced with `wait(10);` for the same functionality.

- `mode=apl` has been removed. APL based sims will no longer run. This will be added back as a functionality at a later date. In the mean time, its functionality can be replicated with the use of combined `if` and `else` statments.

## New Features

- The entire config language (now known as gcsl) has been completely revamped to introduce new scripting features such as variables, arithmetic and logic operations, loops and control statement, and functions.  
**Note that while these language features are available, the are not required for you to write sims.**
- Hitlag has been implemented.
- Animation system has been reworked to remove the previous bug that caused all character actions to take 2 frames longer than intended.
- Almost all characters have had their frames completely recounted and include features such as dash and jump frames, as well as cancels.


## Actually Migrating

### SL sims (The Majority of Sims)
While the optimisation of sims will have to be redone, the migration itself is a very simple process. Only the action list will have to be changed. Simply delete `restart` and nest the rotation in a while loop. 

For example, consider the below sim.

```
options swap_delay=12 debug=true iteration=1000 duration=107 workers=30 mode=sl;
target lvl=100 resist=.1;
energy every=3 amount=1;

raiden char lvl=90/90 cons=0 talent=9,9,9;
raiden add weapon="thecatch" refine=5 lvl=90/90;
raiden add set="tenacityofthemillelith" count=4;
raiden add stats hp=4780 atk=311.0 er=0.518 cr=0.3110 electro%=0.4660;
raiden add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944 ;
			
xingqiu char lvl=90/90 cons=6 talent=9,9,9;
xingqiu add weapon="amenomakageuchi" refine=5 lvl=90/90;
xingqiu add set="emblemofseveredfate" count=4;
xingqiu add stats hp=4780 atk=311.0 atk%=0.4660 cr=0.3110 hydro%=0.4660;
xingqiu add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944 ;

bennett char lvl=90/90 cons=6 talent=9,9,9;
bennett add weapon="thealleyflash" refine=1 lvl=90/90;
bennett add set="noblesseoblige" count=4;
bennett add stats hp=4870 atk=311 er=0.518 cr=0.311 pyro%=0.466;
bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.2204 em=39.64 cr=0.331 cd=0.7944 ;

xiangling char lvl=90/90 cons=6 talent=9,9,9;
xiangling add weapon="dragonsbane" refine=3 lvl=90/90;
xiangling add set="emblemofseveredfate" count=4;
xiangling add stats hp=4780 atk=311.0 er=0.518 cr=0.3110 pyro%=0.4660;
xiangling add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=79.28 cr=0.331 cd=0.7944 ;

active raiden;

raiden attack, skill;
xingqiu skill, burst, attack;
bennett burst, attack, skill;
xiangling burst, attack, skill, attack;
raiden burst, attack:14;
bennett skill;
xiangling attack:3;
restart;
```

Only the last bit needs to be changed:
```
raiden attack, skill;
xingqiu skill, burst, attack;
bennett burst, attack, skill;
xiangling burst, attack, skill, attack;
raiden burst, attack:14, skill;
bennett skill;
xiangling attack:3;
restart;
```
It should be replaced with the following. Remember to delete the `restart`.
```
while 1 {
raiden attack, skill;
xingqiu skill, burst, attack;
bennett burst, attack, skill;
xiangling burst, attack, skill, attack;
raiden burst, attack:14;
bennett skill;
xiangling attack:3;
}
```

*Done!*

Of course there are some easy refinements to this. For example, with the ability to have a unique first rotation, we can take Raiden's Skill out and have her refresh at the end of her Burst instead. 

```
raiden attack, skill;
while 1 {
xingqiu skill, burst, attack;
bennett burst, attack, skill;
xiangling burst, attack, skill, attack;
raiden burst, attack:14, skill;
bennett skill;
xiangling attack:3;
}
```

### APL Migration
As of now, there is no clean way to migrate APL to the new format. But as APL is really just a bit `else if` statement, we can recreate it with `else` and `if` functions. 

For example, if we had a simple apl like below:

```
fischl attack, skill +if=.status.fischloz==0;
fischl attack, burst +if=.status.fischloz==0;
sucrose attack;
```

Then would replace it with the below:

```
while 1 {
    if .status.fischloz==0  && (.ready.fischl.skill || .ready.fischl.burst) {
        fischl attack
        if .ready.fischl.skill {
            fischl skill;
        } else {
            fischl burst;
        }
    } else {
        sucrose attack;
    }
}
```

Have fun. 