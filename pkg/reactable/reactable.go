package reactable

import (
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type Reactable struct {
	Durability [info.ReactionModKeyEnd]info.Durability
	DecayRate  [info.ReactionModKeyEnd]info.Durability
	// Source     []int //source frame of the aura
	self info.Target
	core *core.Core
	// ec specific
	ecAtk      info.AttackInfo // index of owner of next ec ticks
	ecSnapshot info.Snapshot
	ecTickSrc  int
	// burning specific
	burningAtk      info.AttackInfo
	burningSnapshot info.Snapshot
	burningTickSrc  int
	// freeze specific
	FreezeResist float64
	// gcd specific
	overloadGCD     int
	shatterGCD      int
	superconductGCD int
	swirlElectroGCD int
	swirlHydroGCD   int
	swirlCryoGCD    int
	swirlPyroGCD    int
	crystallizeGCD  int
}

type Enemy interface {
	QueueEnemyTask(f func(), delay int)
}

const (
	frzDelta    info.Durability = 2.5 / (60 * 60) // 2 * 1.25
	frzDecayCap info.Durability = 10.0 / 60.0
)

func (r *Reactable) Init(self info.Target, c *core.Core) *Reactable {
	r.self = self
	r.core = c
	r.DecayRate[info.ReactionModKeyFrozen] = frzDecayCap
	r.ecTickSrc = -1
	r.burningTickSrc = -1
	r.overloadGCD = -1
	r.shatterGCD = -1
	r.superconductGCD = -1
	r.swirlElectroGCD = -1
	r.swirlHydroGCD = -1
	r.swirlCryoGCD = -1
	r.swirlPyroGCD = -1
	r.crystallizeGCD = -1
	return r
}

func (r *Reactable) ActiveAuraString() []string {
	var result []string
	for i, v := range r.Durability {
		if v > info.ZeroDur {
			result = append(result, info.ReactionModKey(i).String()+": "+strconv.FormatFloat(float64(v), 'f', 3, 64))
		}
	}
	return result
}

func (r *Reactable) AuraCount() int {
	count := 0
	for _, v := range r.Durability {
		if v > info.ZeroDur {
			count++
		}
	}
	return count
}

func (r *Reactable) React(a *info.AttackEvent) {
	// TODO: double check order of reactions
	switch a.Info.Element {
	case attributes.Electro:
		// hyperbloom
		r.TryAggravate(a)
		r.TryOverload(a)
		r.TryAddEC(a)
		r.TryFrozenSuperconduct(a)
		r.TrySuperconduct(a)
		r.TryQuicken(a)
	case attributes.Pyro:
		// burgeon
		r.TryOverload(a)
		r.TryVaporize(a)
		r.TryMelt(a)
		r.TryBurning(a)
	case attributes.Cryo:
		r.TrySuperconduct(a)
		r.TryMelt(a)
		r.TryFreeze(a)
	case attributes.Hydro:
		r.TryVaporize(a)
		r.TryFreeze(a)
		r.TryBloom(a)
		r.TryAddEC(a)
	case attributes.Anemo:
		r.TrySwirlElectro(a)
		r.TrySwirlPyro(a)
		r.TrySwirlHydro(a)
		r.TrySwirlCryo(a)
		r.TrySwirlFrozen(a)
	case attributes.Geo:
		// can't double crystallize it looks like
		// freeze can trigger hydro first
		// https://docs.google.com/spreadsheets/d/1lJSY2zRIkFDyLZxIor0DVMpYXx3E_jpDrSUZvQijesc/edit#gid=0
		r.TryCrystallizeElectro(a)
		r.TryCrystallizeHydro(a)
		r.TryCrystallizeCryo(a)
		r.TryCrystallizePyro(a)
		r.TryCrystallizeFrozen(a)
	case attributes.Dendro:
		r.TrySpread(a)
		r.TryQuicken(a)
		r.TryBurning(a)
		r.TryBloom(a)
	}
}

// AttachOrRefill is called after the damage event if the attack has not reacted with anything
// and will either create a new modifier if non exist, or update according to the rules of
// each modifier
func (r *Reactable) AttachOrRefill(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	if a.Reacted {
		return false
	}
	// handle pyro, electro, hydro, cryo, dendro
	// special attachment of dendro (burning fuel) is handled in tryBurning
	var mod info.ReactionModKey
	switch a.Info.Element {
	case attributes.Pyro:
		mod = info.ReactionModKeyPyro
	case attributes.Electro:
		mod = info.ReactionModKeyElectro
	case attributes.Hydro:
		mod = info.ReactionModKeyHydro
	case attributes.Cryo:
		mod = info.ReactionModKeyCryo
	case attributes.Dendro:
		mod = info.ReactionModKeyDendro
	default:
		return false
	}
	r.attachOrRefillNormalEle(mod, a.Info.Durability)
	return true
}

func (r *Reactable) GetAuraDurability(mod info.ReactionModKey) info.Durability {
	return r.Durability[mod]
}

func (r *Reactable) GetDurability() []info.Durability {
	result := make([]info.Durability, info.ReactionModKeyEnd)
	for i := range info.ReactionModKeyEnd {
		result[i] = r.Durability[i]
	}
	return result
}

func (r *Reactable) GetAuraDecayRate(mod info.ReactionModKey) info.Durability {
	return r.DecayRate[mod]
}

func (r *Reactable) SetAuraDurability(mod info.ReactionModKey, dur info.Durability) {
	r.Durability[mod] = dur
}

func (r *Reactable) SetAuraDecayRate(mod info.ReactionModKey, decay info.Durability) {
	r.DecayRate[mod] = decay
}

func (r *Reactable) SetFreezeResist(fr float64) {
	r.FreezeResist = fr
}

// attachOrRefillNormalEle is used for pyro, electro, hydro, cryo, and dendro which don't have special attachment
// rules
func (r *Reactable) attachOrRefillNormalEle(mod info.ReactionModKey, dur info.Durability) {
	amt := 0.8 * dur
	if mod == info.ReactionModKeyPyro {
		r.attachOverlapRefreshDuration(info.ReactionModKeyPyro, amt, 6*dur+420)
	} else {
		r.attachOverlap(mod, amt, 6*dur+420)
	}
}

func (r *Reactable) attachOverlap(mod info.ReactionModKey, amt, length info.Durability) {
	if r.Durability[mod] > info.ZeroDur {
		add := max(amt-r.Durability[mod], 0)
		if add > 0 {
			r.addDurability(mod, add)
		}
	} else {
		r.Durability[mod] = amt
		if length > info.ZeroDur {
			r.DecayRate[mod] = amt / length
		}
	}
}

func (r *Reactable) attachOverlapRefreshDuration(mod info.ReactionModKey, amt, length info.Durability) {
	if amt < r.Durability[mod] {
		return
	}
	r.Durability[mod] = amt
	r.DecayRate[mod] = amt / length
}

func (r *Reactable) attachBurning() {
	r.Durability[info.ReactionModKeyBurning] = 50
	r.DecayRate[info.ReactionModKeyBurning] = 0
}

func (r *Reactable) addDurability(mod info.ReactionModKey, amt info.Durability) {
	r.Durability[mod] += amt
	r.core.Events.Emit(event.OnAuraDurabilityAdded, r.self, mod, amt)
}

// AuraCountains returns true if any element e is active on the target
func (r *Reactable) AuraContains(e ...attributes.Element) bool {
	for _, v := range e {
		for i := range info.ReactionModKeyEnd {
			if i.Element() == v && r.Durability[i] > info.ZeroDur {
				return true
			}
		}
		// TODO: not sure if this is best way to go about it? perhaps supplying frozen element is better?
		if v == attributes.Cryo && r.Durability[info.ReactionModKeyFrozen] > info.ZeroDur {
			return true
		}
	}
	return false
}

func (r *Reactable) IsBurning() bool {
	if r.Durability[info.ReactionModKeyBurningFuel] > info.ZeroDur && r.Durability[info.ReactionModKeyBurning] > info.ZeroDur {
		return true
	}
	return false
}

// reduce the requested element by dur * factor, return the amount of dur consumed
// if multiple modifier with same element are present, all of them are reduced
// the max on reduced is used for consumption purpose
func (r *Reactable) reduce(e attributes.Element, dur, factor info.Durability) info.Durability {
	m := dur * factor // maximum amount reduceable
	var reduced info.Durability

	for i := range info.ReactionModKeyEnd {
		if i.Element() != e {
			continue
		}
		if r.Durability[i] < info.ZeroDur {
			// also skip if durability already 0
			// this allows us to safely call reduce even if an element doesn't exist
			continue
		}
		// reduce by lesser of remaining and m

		red := m

		if red > r.Durability[i] {
			red = r.Durability[i]
			// reset decay rate to 0
		}

		r.Durability[i] -= red

		if red > reduced {
			reduced = red
		}
	}

	return reduced / factor
}

func (r *Reactable) deplete(m info.ReactionModKey) {
	if r.Durability[m] <= info.ZeroDur {
		r.Durability[m] = 0
		r.DecayRate[m] = 0
		r.core.Events.Emit(event.OnAuraDurabilityDepleted, r.self, attributes.Element(m))
	}
}

func (r *Reactable) Tick() {
	// duability is reduced by decay * (1 + purge)
	// where purge is 0 for anything that's not freeze
	// for freeze, purge = 0.25 * time
	// while defrosting, purge rate is reduced back down to 0 at new = old - 0.5 * time
	// where time is in seconds
	//
	// per frame then we have decay * (1 + 0.25 * (x/60))

	// anything after the delim is special decay so we ignore
	for i := range info.ReactionModKeySpecialDecayDelim {
		// skip zero decay rates i.e. modifiers that don't decay (i.e. burning)
		if r.DecayRate[i] == 0 {
			continue
		}
		if r.Durability[i] > info.ZeroDur {
			r.Durability[i] -= r.DecayRate[i]
			r.deplete(i)
		}
	}

	// check burning first since that affects dendro/quicken decay
	if r.burningTickSrc > -1 && r.Durability[info.ReactionModKeyBurningFuel] < info.ZeroDur {
		// reset src when burning fuel is gone
		r.burningTickSrc = -1
		// remove burning
		r.Durability[info.ReactionModKeyBurning] = 0
		// remove existing dendro and quicken
		r.Durability[info.ReactionModKeyDendro] = 0
		r.DecayRate[info.ReactionModKeyDendro] = 0
		r.Durability[info.ReactionModKeyQuicken] = 0
		r.DecayRate[info.ReactionModKeyQuicken] = 0
	}

	// if burning fuel is present, dendro and quicken uses burning fuel decay rate
	// otherwise it uses it's own
	for i := info.ReactionModKeyDendro; i <= info.ReactionModKeyQuicken; i++ {
		if r.Durability[i] < info.ZeroDur {
			continue
		}
		rate := r.DecayRate[i]
		if r.Durability[info.ReactionModKeyBurningFuel] > info.ZeroDur {
			rate = r.DecayRate[info.ReactionModKeyBurningFuel]
			if i == info.ReactionModKeyDendro {
				rate = max(rate, r.DecayRate[i]*2)
			}
		}
		r.Durability[i] -= rate
		r.deplete(i)
	}

	// for freeze, durability can be calculated as:
	// d_f(t) = -1.25 * (t/60)^2 - k * (t/60) + d_f(0)
	if r.Durability[info.ReactionModKeyFrozen] > info.ZeroDur {
		// ramp up decay rate first
		r.DecayRate[info.ReactionModKeyFrozen] += frzDelta
		r.Durability[info.ReactionModKeyFrozen] -= r.DecayRate[info.ReactionModKeyFrozen] / info.Durability(1.0-r.FreezeResist)

		r.checkFreeze()
	} else if r.DecayRate[info.ReactionModKeyFrozen] > frzDecayCap { // otherwise ramp down decay rate
		r.DecayRate[info.ReactionModKeyFrozen] -= frzDelta * 2

		// cap decay
		if r.DecayRate[info.ReactionModKeyFrozen] < frzDecayCap {
			r.DecayRate[info.ReactionModKeyFrozen] = frzDecayCap
		}
	}

	// for ec we need to reset src if ec is gone
	if r.ecTickSrc > -1 {
		if r.Durability[info.ReactionModKeyElectro] < info.ZeroDur || r.Durability[info.ReactionModKeyHydro] < info.ZeroDur {
			r.ecTickSrc = -1
		}
	}
}

func (r *Reactable) calcCatalyzeDmg(atk info.AttackInfo, em float64) float64 {
	char := r.core.Player.ByIndex(atk.ActorIndex)
	return combat.CalcCatalyzeDmg(char.Base.Level, char, atk, em)
}
