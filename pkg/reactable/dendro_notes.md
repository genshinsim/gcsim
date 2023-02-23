This is to keep track Dendro implementation notes/thoughts:

## Gadgets

Will need to implement generic gadgets to handle proper interaction with Bloom

```go
type Gadget interface {
    OnAdded() // called when gadget gets created
    OnRemoved() // called when gadget either expires, or is destroyed
}
```

All gadgets need to implement `combat.Target` interface. This is to make sure gadgets can be added to the combat system. Gadgets should live for as long as it has `durability` or unless a `duration` is set.

### Default implementation

Should provide a default template providing common functions that specific gadgets can build upon. Something like:

```go
type Gadget struct {
    x, y int
    index int
    thinkInterval int
    onThinkInterval func() // this function to be called on think interval
}

func (g *Gadget) Index() { return g.index }
func (g *Gadget) SetIndex(idx int) { g.index = idx }
func (g *Gadget) Tick() {
    // per tick logic here
}

// any other basic stuff, etc etc
```

### Possible `combat.Target` refactor required

Not sure yet how to handle the following methods of `combat.Target`. These methods really aren't useful to gadgets.

```go
    IsAlive() bool
    MaxHP() float64
    HP() float64
    AuraContains(e ...attributes.Element) bool
```

Can maybe refactor `combat.Target` as `combat.Handler` probably doesn't actually need to know if a target is alive or how much hp it has. We can refactor `combat.Target` to the following:

```go
type Target interface {
	Index() int         //should correspond to index
	SetIndex(index int) //update the current index
	Type() TargettableType   //type of target
	Shape() Shape            // shape of target
	Pos() (float64, float64) // center of target
	SetPos(x, y float64)     // move target
	AttackWillLand(AttackPattern) bool   // hitbox check
	Attack(*AttackEvent, glog.Event) (float64, bool)
	Tick()
	Kill()
}
```

This way `AttackWillLand` can handle any dead or alive checks. Alternatively change it so that dead targets are removed from the targets array. Only problem with this is that logging needs to be changed to log not by target index but by some sort of target key and target key cannot be duplicates.

Also `AuraContains(e ...attributes.Element) bool` can be moved to a different interface in combat, and the absorb check be changed accordingly:

```go
type TargetWithAura interface {
	AuraContains(e ...attributes.Element) bool
}

func (c *Handler) AbsorbCheck(p AttackPattern, prio ...attributes.Element) attributes.Element {

	// check targets for collision first
	for _, e := range prio {
		for i, x := range c.targets {
			t, ok := x.(TargetWithAura)
			if !ok {
				continue
			}
			if WillCollide(p, t, i) && t.AuraContains(e) {
				c.Log.NewEvent(
					"infusion check picked up "+e.String(),
					glog.LogElementEvent,
					-1,
				)
				return e
			}
		}
	}
	return attributes.NoElement
}
```

## Quicken

Quicken is a new element. It behaves effectively as a `dendro` element for reactions. Can coexist with `electro`, `dendro`, and `cryo` (because there is no `dendro` + `cryo` reactions). Can coexist with burning (QUESTION: does it act as burning fuel in this case?)

### Reaction multiplier

The reaction has a 1:1 multiplier between `electro` and `dendro`

### Duration and durability

Quicken have different formula than normal "attachment". Existing code uses [this function](https://github.com/genshinsim/gcsim/blob/40b1617647c1dbcd541e684b7f5b09d9dd424851/pkg/reactable/reactable.go#L182) to handle attachment:

```go
func (r *Reactable) attach(e attributes.Element, dur combat.Durability, m combat.Durability) {
	//calculate duration based on dur
	r.DecayRate[e] = m * dur / (6*dur + 420)
	r.addDurability(e, m*dur)
}
```

With Quicken specifically, the duration formula is changed to: `12 * dur + 360`. So the decay rate is actually:

```go
	r.DecayRate[attributes.Quicken] = m * dur / (12 * dur + 360)
```

Probably best to just not use the `attach` function call and leave as is. Realistically the `attach` function should be refactored to handle all cases.

### Aggravate/Spread

When `quicken` reacts with `electro` or `dendro`, aggravate and and spread is triggered respectively. Both reactions adds to the `AttackInfo.FlatDmg` of the triggering attack.

The amount of damage added follows the same base [reaction damage formula](https://github.com/genshinsim/gcsim/blob/40b1617647c1dbcd541e684b7f5b09d9dd424851/pkg/reactable/reactable.go#L170) but uses a different EM multiplier. Calculation should be as follows:

```go
func (r *Reactable) calcQuickenReactionDmgBase(atk combat.AttackInfo, em float64) float64 {
	char := r.core.Player.ByIndex(atk.ActorIndex)
	lvl := char.Base.Level - 1
	if lvl > 89 {
		lvl = 89
	}
	if lvl < 0 {
		lvl = 0
	}
	// use 5 * em here instead of the usual 16 * em
	return (1 + ((5 * em) / (2000 + em)) + r.core.Player.ByIndex(atk.ActorIndex).ReactBonus(atk)) * reactionLvlBase[lvl]
}

func (r *Reactable) calcAggravateDmg(atk combat.AttackInfo, em float64) float64 {
	return 1.15 * r.calcQuickenReactionDmgBase(atk, em)
}

func (r *Reactable) calcSpreadDmg(atk combat.AttackInfo, em float64) float64 {
	return 1.25 * r.calcQuickenReactionDmgBase(atk, em)
}
```

Multiplier for aggravate and spread is 1.15 and 1.25 respectively.

### Other `quicken` reactions

`quicken` behaves the same as `dendro` when reacting with `hydro` or `pyro`, triggering `bloom` and `burning` respectively

## Bloom

Bloom is the reaction triggered when `dendro` reacts with `hydro`. This one is fairly straightforward. On reaction, a "seed" is generated.

The position of the seed (dendro core) appears to be distance 1 from the target, with a random 0 to 60 degree angle. Presuming this is 0 to 60 degree from x-axis on the cartesian plane (in game uses y). The seed then has an initial velocity and a deceleration. Needs more testing on this. For simplicty probably ok to just spawn randomly within 1 radius of target.

Seeds have a lifetime of 5 seconds. From spawn it should be closer to 6 seconds?

## Bloom explode

If the dendro core does not come in contact with `pyro` or `electro` (simple attack collision check), it will explode on expiriy. Radius is 5. AttackInfo is as follows:

```go
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(combat.Bloom),
		AttackTag:        attacks.AttackTagBloom,
		ICDTag:           attacks.ICDTagBloomDamage,
		ICDGroup:         combat.ICDGroupReactionA,
		Element:          attributes.Dendro,
		Durability:       0,
		IgnoreDefPercent: 1,
	}
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	atk.FlatDmg = 2.0 * r.calcReactionDmg(atk, em)
```

### Burgeon

When the dendro core comes in contact with `pyro`, `burgeon` is triggered, dealing an attack with radius 5. AttackInfo is the same as `bloom` except with `attacks.AttackTagBurgeon`, `attacks.AttackTagBurgeon` and `attacks.ICDTagBurgeon`

### Hyperbloom

When dendro core comes in contact with `electro`, `hyperbloom` is triggered, dealing an attack with radius 1. AttackInfo is the same as `bloom` except with `attacks.AttackTagHyperbloom`, `attacks.AttackTagHyperbloom` and `attacks.ICDTagHyperbloom`

### Bloom/Burgeon/Hyperbloom self damage

All 3 explosions will trigger an additional attack that damages the player, with AttackInfo as follows:

```go
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(combat.Bloom),
		AttackTag:        attacks.AttackTagBloom, // or AttackTagBurgeon, AttackTagHyperbloom
		ICDTag:           attacks.ICDTagBloomDamage,
		ICDGroup:         combat.ICDGroupReactionA,
		Element:          attributes.Dendro,
		Durability:       0,
		IgnoreDefPercent: 1,
	}
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	atk.FlatDmg = 2.0 * r.calcReactionDmg(atk, em)
```

Radius is the same as each respective type of explosion (i.e. radius 1 for Hyperbloom)

## Burning

Probably the most complex reaction to date (way worse than EC or Freeze).... Introdues 2 new modifiers but no additional elements, which technically "breaks" the current implementation but thankfully I think should be ok implementing a hack solution. As long as Hoyoverse doesn't introduce any new reactions.

### Modifiers

Before getting into burning, first have to talk about modifiers. In the current implementation, gcsim is treating auras as some sort of array of buckets attached to each target object. In reality, auras are just modifiers, which are fairly generic properties that are attached to just about anything.

Historically for each "aura" (or existing element), you just had one modifier. Usually something like `FireModifier` for attached pyro, or `ElectroModifier` for attached electro. This made our array of whatever buckets sort of work.

However with Burning, we are introducing 2 new modifiers (`burning` and `burning_fuel`) with elements attached BUT:

1. These modifiers do not introduce new element types
2. These modifier can coexist with other modifiers having the exact same element (namely Pyro)

This throws a bit of a wrench into our implementation because gcsim is treating modifiers as having a 1 to 1 relationshp with element types. To fix this I think we have two options:

1. Gut the reaction system (again), changing it to use generic modifiers so that everything is key'd by modifier names and we can easily have coexisting elements of the same type
2. Introduce fake elements `burning` and `burning_fuel`

For now... we'll have to go with approach 2. in the interest of time. This is our hack solution. However, I am concerned about this translating to problems later down the road either in the form of hard to maintain code or unexpected bugs (involing deal more damage if x element present for example).

Hopefully it works out...

### Reaction, Duration, and Durability

When `dendro` reacts with `pyro`, you get two resulting modifiers (or "elements" in our hack solution): `burning` and `burning_fuel`.

`burning` always has 50 durability regardless of the reacted amount. It also does not have a set duration. Instead, `burning` will last as long as both:

1. `burning` durability > 0, and
2. `burning_fuel` exists still

`burning` has element of type `pyro` and will react just like normal `pyro`. However it cannot be topped up by additional pyro application (makes sense since it's a completely different modifier). Any additional `pyro` application will merely coexist.

For our hack implementation, we'll just have to go with:

```go
	r.DecayRate[attributes.Burning] = 0
	r.Durability[attributes.Burning] = 50
```

And then on tick we'll need to check for if `burning_fuel` is still present, and if not remove `burning`

`burning_fuel` durability has a 0.8x multiplier on the reacted amount. So for example:

- Existing: `dendro 20`, applied from a `dendro 25` attack
- Applying: `pyro 25`
- Result: `dendro (burning_fuel) 16`

The duration formula follows the standard `0.1 * d + 7` in seconds, or `6 * d + 420` in frames.

In addition, the decay rate is subject to a minimum threshold of at least 10 durability per second. Effectively, you have the following for decay rate:

```go
	decayRate := m * dur / (6 * dur + 420)
	if decayRate < 10.0/60.0 {
		decayRate = 10.0/60.0
	}
	r.DecayRate[attributes.BurningFuel] = decayRate
```

Note here I'm treating `burning_fuel` as it's own element i.e. `attributes.BurningFuel` when that's not actually the case here in game (see modifier discussion about).

### `burning_fuel` refresh

Unlike other element attachment, when additional Dendro is applied to existing `burning_fuel`, it does not overlap but instead refreshes. This means that the existing `burning_fuel` durability/duration will be overwritten by the applying Dendro durability/duration. However, the existing decay rate is kept.

For example:

- Existing: `dendro 40 (burning_fuel)`
- Applying: `dendro 25`
- Result: `dendro 25 (burning_fuel)`

So you can actually "shrink" the amount of `dendro (burning_fuel)` aura per say.

### Burning damage

While `burning` is active, every 0.25s it will trigger a spherical attack with radius 1, center on the target that is burning.

AttackInfo is as follows:

```go
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(combat.Burning),
		AttackTag:        attacks.AttackTagBurningDamage,
		ICDTag:           attacks.ICDTagBurningDamage,
		ICDGroup:         combat.ICDGroupBurning,
		Element:          attributes.Pyro,
		Durability:       25,
		IgnoreDefPercent: 1,
	}
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	atk.FlatDmg = 0.25 * r.calcReactionDmg(atk, em)
```

TODO: it's possible that if burning ends early before the next tick happens then the next tick happens right away? not sure

POST STILL WIP
