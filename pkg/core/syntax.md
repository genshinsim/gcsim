# Team data syntax (Proposed)

Following is the proposed syntax for the team data portion of the config file

```
xiangling char lvl=80/90 cons=4 talent=6,9,9;
xiangling add weapon="staff of homa" lvl=80/90 refine=3;
xiangling add set="seal of insulation" count=4;
xiangling add stats hp=4780 atk=311 er=.518 pyro%=0.466 cr=0.311;
xiangling add stats atk%=.0992 cr=.1655 cd=.7282 em=39.64 er=.5510 hp%=.0992 hp=507.88 atk=33.08 def%=.124 def=39.36;
```

# Option syntaxs

`options debug=true iteration=5000 duration=90 workers=24`

# Action List Syntax (Proposed)

Following is proposed action list syntax for gcsim

## general syntax

Each line typically follows the following structure

`<char or command> <list of abilities> <flags>`

For example

`bennett skill,attack,burst,attack /if=.xx.xx.xx>1 /swap=xiangling /onfield /limit=1 /try=1 /timeout=100 /noswap=50`

## char or commands

Each line needs to start either with a character name or a command

List of commands as follows:

- `chain`
- `wait`
- `reset_limit`

### `chain`

`chain` allow you to execute a list of macros.

Macro are a way to shortcut a certain sequence of actions. Macros can only be used in `+chain`.

All macros must be declared at the top of the action list before any action lines. Otherwise the config will fail to parse.

Macros starts with an `identifier` representing the name of the macro as is followed by `:`. Examples:

```
xq_seq:xingqiu skill[orbital=1],burst[orbital=1],attack;
bt_seq:bennett burst,skill;
```

Macros can also be commands

```
wait_for_bt:wait /mod=xingqiu.btburst
asdf:reset_limit
```

Not sure why you would want to chain a reset_limit but it's acceptable.

The following flags are not acceptable in a macro:

- `swap`
- `onfield`
- `needs`
- `timeout`
- `try`

The chain command itself is subject to the same format as any regular line.

`chain bt_seq,wait_for_bt,xq_seq,asdf /if=.xx.xx.xx>1 /swap=xiangling /limit=1 /try=1`

By default, if the `/try` flag is not set then every ability in each macro must be ready before this chain will be executed. Other wise only the first action of the first ability in the chain needs to be ready.

Note that the following flags cannot be assigned to `chain` as they wouldn't make sense:

- `onfield`

### `wait`

`wait` tells the simulator to stop executing any actions. The length of wait depends on the flags set. Available flags are

- `/if`: should be equal to `.char.modkey`
- `/particle=<src>`: where `<src>` is either a character name or a particle source such as `favonius`
- `/max`: max number of frames to wait for. If used with any other option than it'll specify when to give up. Can be used alone with `wait` to wait for a number of frames no other conditions.

Not sure what other waits we would have ?

### `reset_limit`

`reset_limit` resets all lines with `/limit` flag (i.e. set their usage count to 0)

## list of abilities

This should be a comma separated list of abilities. Params can be passed into each ability via square brackets. For example:

`skill[orbital=1],attack,burst[orbital=1]`

Can use `:x` as a short form to repeat certain abilitiy. For example:

`attack[travel=50]:4`

This will repeat attack, with params `[travel=50]` 4 times. Common usage would be something like:

`raidenshongun attack:4,charge,attack:4,charge`

Which is basically N4C N4C

Following is valid list of abilities

- `skill`
- `burst`
- `attack`
- `charge`
- `high_plunge`
- `low_plunge`
- `proc`
- `aim`
- `dash`
- `jump`

## flags

These are additional key words that can be added after the list of abilities that modify how the action is executed

- `/if`: condition that must be fulfiled before this line can be executed
- `/swap`: forcefully swap to another char after this line has finished
- `/swap_lock`: force the sim to stay on this character for x frames
- `/onfield`: this line can only be used if this character is already on field (does not work in `chain`)
- `/label`: a name for this line; can't have duplicated labels
- `/needs`: this line can only be executed if the previous action is the referred to label
- `/limit`: number of times this line can be executed (replaces once)
- `/timeout`: this line cannot be executed again for x number of frames
- `/try`: if this flag is set, then this line will execute as long as first ability in the list is executable. If flag is set to 1, then if the sim will keep trying to execute the next ability in the list even if it's not ready (even if it means waiting forever). If flag is set to 0, then if the next ability is not ready immediately after the previous ability, the next ability along with the rest of the sequence will be dropped
- `/char_mod`: flag specific to wait command
- `/particle`: flag specific to wait command
- `/max`: flag specific to wait command

## CALC MODE

This syntax is for calc mode which is super simplified as there are no conditionals and is just a list of actions to execute (just like a calculator)

Calc mode syntax is **NOT** compatible with regular priority based action list. Trying to use this in sim mode will result in an error.

The syntax can be demonstrated as follows:

```
bennett attack, skill, attack, burst, attack
dash
swap raiden
wait 10
raiden burst, attack:4, charge, attack:4, charge, attack:4, charge
wait_until 600
```

Some explanations:

- `dash` or `jump` will cause the calculator to use either of those
- `swap <char>` will cause the calculator to swap to the specified character. note that if you swap to say `raiden` but your very next action is say `xingqiu attack`, the calculator is going to execute the swap to raiden, realize it needs to be on xingqiu next, and therefore wait there doing nothing for the swap cooldown until it can swap to xingqiu
- `wait x` will wait x number of frames doing nothing
- `wait_until x` will wait until frame number x
