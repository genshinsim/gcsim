# Sim structure

Sim is structured into a couple core packages, briefly described as follows:

- `core`: This package is essentially a controller that brings all the parts together. See details below
- `monster`: This package handles the damage calculation, reactions, and any target debuffs.
- `parse`: This package handle reading the config files
- `character`: This package handles the basic/universal character logic and is then overriden by each character to implement their own actions

# Core

The `core` package contains all the necessary

# Parse

`https://github.com/genshinsim/gsim/tree/main/pkg/parse`

This package contains the necessary code to lex/parse the custom config file syntax into the config data structure that's used by the core. The core logic is based on Rob Pike's [talk](https://talks.golang.org/2011/lex.slide#1) as well as go's template parsing [implementation](https://cs.opensource.google/go/go/+/refs/tags/go1.17:src/text/template/parse/)

# Monster

The `monster` package handles all the logic relating to:

- hitbox resolution (although not currently implemented)
- damage calculation
- ICD
- aura tracking and reactions
- target resistance and defense mods

All of this is implemented in a `Target` struct, which implements the `core.Target` interface. Thus multi target is simply having multiple copies of this `Target` struct. We'll call this the "target" for simplicity.

## Hitbox

When an attack is generated, whether or not a target will be hit/damaged is resolved by each target independently. Each of the attack [snapshot](https://github.com/genshinsim/gsim/blob/main/pkg/core/snapshot.go) contains the information necessary for the target to determine if it will be hit or not.

For now the implementation is relatively simple. There is a `Targets` field in each snapshot. If this field is equal to the index of the current target or if this field is equal to -1 (representing all targets), then the current target will take damage.

In the future, this implementation can be changed to include 2D geometry.

## Damage



## Auras

## Reactions


# Characters



# Weapons


# Artifacts


