package reactable

import (
	"math/bits"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type Reactable struct {
	// Durability is only written by SetAuraDurability, addDurability, reduceMod and removeMod,
	// which keep live up to date.
	Durability [info.ReactionModKeyEnd][info.MaxChars]info.Durability
	DecayRate  [info.ReactionModKeyEnd]info.Durability
	// live has bit mod set while some source of Durability[mod] may be at or above ZeroDur
	// (or NaN). Every write that can raise a durability sets it and Tick clears it once all
	// sources are below ZeroDur, so Tick skips the mods that have nothing to decay.
	live uint16
	// Source     []int //source frame of the aura
	self info.Target
	core *core.Core
	// ec specific
	ecAtk      info.AttackInfo // index of owner of next ec ticks
	ecSnapshot info.Snapshot
	ecTickSrc  int
	lcAtkOwner int // index of owner of next LC tick
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
	// TODO: change this to be initialized after the characters are added
	// for i := range r.Durability {
	// 	r.Durability[i] = make([]info.Durability, len(c.Player.Chars()))
	// }
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
	for i := range info.ReactionModKeyEnd {
		v := r.GetAuraDurability(i)
		if v > info.ZeroDur {
			result = append(result, i.String()+": "+strconv.FormatFloat(float64(v), 'f', 3, 64))
		}
	}
	return result
}

func (r *Reactable) AuraCount() int {
	count := 0
	for i := range info.ReactionModKeyEnd {
		if r.GetAuraDurability(i) > info.ZeroDur {
			count++
		}
	}
	return count
}

func (r *Reactable) emitEvent(evt event.Event, args ...any) {
	if r.core == nil {
		return
	}
	r.core.Events.Emit(evt, args...)
}

func (r *Reactable) React(a *info.AttackEvent) {
	r.emitEvent(event.OnElementApplied, r.self, a)

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
	// fmt.Printf("AttachOrRefill a.Info.ActorIndex %v\n", a.Info.ActorIndex)
	r.attachOrRefillNormalEle(mod, a.Info.Durability, a.Info.ActorIndex)
	return true
}

// The aura helpers below and reduceMod are unrolled for 4 sources, since on wasm every loop
// iteration jumps through the function's br_table dispatch. This stops compiling if
// info.MaxChars changes.
var _ [4]info.Durability = [info.MaxChars]info.Durability{}

// GetAuraDurability returns the largest durability of mod across all sources, folding from the
// left like slices.Max.
func (r *Reactable) GetAuraDurability(mod info.ReactionModKey) info.Durability {
	d := &r.Durability[mod]
	return maxDur(maxDur(maxDur(d[0], d[1]), d[2]), d[3])
}

// auraAbove, auraBelow and auraAtMost report GetAuraDurability(mod) > x, < x and <= x without
// computing the max. A NaN source makes the max NaN and every comparison false: auraBelow and
// auraAtMost are false on NaN by themselves, auraAbove checks for it.
func (r *Reactable) auraAbove(mod info.ReactionModKey, x info.Durability) bool {
	d := &r.Durability[mod]
	// d == d is the NaN check; isNaN would push this past the inlining budget
	return (d[0] > x || d[1] > x || d[2] > x || d[3] > x) &&
		d[0] == d[0] && d[1] == d[1] && d[2] == d[2] && d[3] == d[3] //nolint:gocritic // NaN check
}

func (r *Reactable) auraBelow(mod info.ReactionModKey, x info.Durability) bool {
	d := &r.Durability[mod]
	return d[0] < x && d[1] < x && d[2] < x && d[3] < x
}

func (r *Reactable) auraAtMost(mod info.ReactionModKey, x info.Durability) bool {
	d := &r.Durability[mod]
	return d[0] <= x && d[1] <= x && d[2] <= x && d[3] <= x
}

func (r *Reactable) GetDurability() []info.Durability {
	result := make([]info.Durability, info.ReactionModKeyEnd)
	for i := range info.ReactionModKeyEnd {
		result[i] = r.GetAuraDurability(i)
	}
	return result
}

func (r *Reactable) GetAuraDecayRate(mod info.ReactionModKey) info.Durability {
	return r.DecayRate[mod]
}

func (r *Reactable) SetAuraDurability(mod info.ReactionModKey, dur info.Durability, src int) {
	r.Durability[mod][src] = dur
	r.live |= 1 << mod
	r.emitEvent(event.OnAuraDurabilityAdded, r.self, mod, dur)
}

func (r *Reactable) SetAuraDecayRate(mod info.ReactionModKey, decay info.Durability) {
	r.DecayRate[mod] = decay
}

func (r *Reactable) SetFreezeResist(fr float64) {
	r.FreezeResist = fr
}

// attachOrRefillNormalEle is used for pyro, electro, hydro, cryo, and dendro which don't have special attachment
// rules
func (r *Reactable) attachOrRefillNormalEle(mod info.ReactionModKey, dur info.Durability, src int) {
	amt := 0.8 * dur
	if mod == info.ReactionModKeyPyro {
		r.attachOverlapRefreshDuration(info.ReactionModKeyPyro, amt, 6*dur+420, src)
	} else {
		r.attachOverlap(mod, amt, 6*dur+420, src)
	}
}

func (r *Reactable) attachOverlap(mod info.ReactionModKey, amt, length info.Durability, src int) {
	if r.GetAuraDurability(mod) > info.ZeroDur {
		add := maxDur(amt-r.Durability[mod][src], 0)
		if add > 0 {
			r.addDurability(mod, add, src)
		}
	} else {
		r.SetAuraDurability(mod, amt, src)
		if length > info.ZeroDur {
			r.DecayRate[mod] = amt / length
		}
	}
}

func (r *Reactable) attachOverlapRefreshDuration(mod info.ReactionModKey, amt, length info.Durability, src int) {
	if amt < r.GetAuraDurability(mod) {
		return
	}
	r.SetAuraDurability(mod, amt, src)

	// only update decay rate is amt is greater or equal to the current durability
	if amt < r.GetAuraDurability(mod) {
		return
	}
	r.DecayRate[mod] = amt / length
}

func (r *Reactable) attachBurning(src int) {
	r.removeMod(info.ReactionModKeyBurning)
	r.SetAuraDurability(info.ReactionModKeyBurning, 50, src)
}

func (r *Reactable) addDurability(mod info.ReactionModKey, amt info.Durability, src int) {
	r.Durability[mod][src] += amt
	r.live |= 1 << mod
	r.emitEvent(event.OnAuraDurabilityAdded, r.self, mod, amt)
}

// AuraCountains returns true if any element e is active on the target
func (r *Reactable) AuraContains(e ...attributes.Element) bool {
	for _, v := range e {
		if v == attributes.NoElement {
			continue
		}
		for i := range info.ReactionModKeyEnd {
			if i.Element() == v && r.GetAuraDurability(i) > info.ZeroDur {
				return true
			}
		}
		// TODO: not sure if this is best way to go about it? perhaps supplying frozen element is better?
		if v == attributes.Cryo && r.GetAuraDurability(info.ReactionModKeyFrozen) > info.ZeroDur {
			return true
		}
	}
	return false
}

func (r *Reactable) IsBurning() bool {
	if r.GetAuraDurability(info.ReactionModKeyBurningFuel) > info.ZeroDur && r.GetAuraDurability(info.ReactionModKeyBurning) > info.ZeroDur {
		return true
	}
	return false
}

func (r *Reactable) reduceMod(mod info.ReactionModKey, amt info.Durability) {
	d := &r.Durability[mod] // unrolled, see GetAuraDurability
	d[0] -= minDur(amt, d[0])
	d[1] -= minDur(amt, d[1])
	d[2] -= minDur(amt, d[2])
	d[3] -= minDur(amt, d[3])
	// a negative or NaN amt raises durability
	r.live |= 1 << mod
	if r.auraAtMost(mod, info.ZeroDur) {
		r.emitEvent(event.OnAuraDurabilityDepleted, r.self, mod)
	}
}

func (r *Reactable) removeMod(mod info.ReactionModKey) {
	for i := range r.Durability[mod] {
		r.Durability[mod][i] = 0
	}
	r.DecayRate[mod] = 0
	r.emitEvent(event.OnAuraDurabilityDepleted, r.self, mod)
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
		if r.auraBelow(i, info.ZeroDur) {
			// also skip if durability already 0
			// this allows us to safely call reduce even if an element doesn't exist
			continue
		}
		// reduce by lesser of remaining and m

		red := minDur(m, r.GetAuraDurability(i))

		r.reduceMod(i, red)

		if red > reduced {
			reduced = red
		}
	}

	return reduced / factor
}

func (r *Reactable) deplete(m info.ReactionModKey) {
	if r.auraAtMost(m, info.ZeroDur) {
		r.SetAuraDecayRate(m, 0)
		r.emitEvent(event.OnAuraDurabilityDepleted, r.self, attributes.Element(m))
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

	// Only live mods can have anything to decay (see live). The next live mod is looked up
	// after each step because the depleted events can change the auras.

	// anything after the delim is special decay so we ignore
	for i := r.nextLive(0); i < info.ReactionModKeySpecialDecayDelim; i = r.nextLive(i + 1) {
		// skip zero decay rates i.e. modifiers that don't decay (i.e. burning)
		if r.DecayRate[i] != 0 && r.auraAbove(i, info.ZeroDur) {
			r.reduceMod(i, r.DecayRate[i])
			r.deplete(i)
		} else {
			r.clearIfBelowZero(i)
		}
	}

	// check burning first since that affects dendro/quicken decay

	if r.burningTickSrc > -1 && r.auraBelow(info.ReactionModKeyBurningFuel, info.ZeroDur) {
		// reset src when burning fuel is gone
		r.burningTickSrc = -1
		// remove burning
		r.removeMod(info.ReactionModKeyBurning)
		// remove existing dendro and quicken
		r.removeMod(info.ReactionModKeyDendro)
		r.removeMod(info.ReactionModKeyQuicken)
	}

	// if burning fuel is present, dendro and quicken uses burning fuel decay rate
	// otherwise it uses it's own
	for i := r.nextLive(info.ReactionModKeyDendro); i <= info.ReactionModKeyQuicken; i = r.nextLive(i + 1) {
		if r.clearIfBelowZero(i) {
			continue
		}
		rate := r.DecayRate[i]
		if r.live&(1<<info.ReactionModKeyBurningFuel) != 0 && r.auraAbove(info.ReactionModKeyBurningFuel, info.ZeroDur) {
			rate = r.DecayRate[info.ReactionModKeyBurningFuel]
			if i == info.ReactionModKeyDendro {
				rate = maxDur(rate, r.DecayRate[i]*2)
			}
		}
		r.reduceMod(i, rate)
		r.deplete(i)
	}

	// for freeze, durability can be calculated as:
	// d_f(t) = -1.25 * (t/60)^2 - k * (t/60) + d_f(0)
	if r.live&(1<<info.ReactionModKeyFrozen) != 0 && r.auraAbove(info.ReactionModKeyFrozen, info.ZeroDur) {
		// ramp up decay rate first
		r.DecayRate[info.ReactionModKeyFrozen] += frzDelta
		r.reduceMod(info.ReactionModKeyFrozen, r.DecayRate[info.ReactionModKeyFrozen]/info.Durability(1.0-r.FreezeResist))

		r.checkFreeze()
	} else if r.DecayRate[info.ReactionModKeyFrozen] > frzDecayCap { // otherwise ramp down decay rate
		r.clearIfBelowZero(info.ReactionModKeyFrozen)
		r.DecayRate[info.ReactionModKeyFrozen] -= frzDelta * 2

		// cap decay
		if r.DecayRate[info.ReactionModKeyFrozen] < frzDecayCap {
			r.DecayRate[info.ReactionModKeyFrozen] = frzDecayCap
		}
	}

	// for ec we need to reset src if ec is gone
	if r.ecTickSrc > -1 {
		if r.auraBelow(info.ReactionModKeyElectro, info.ZeroDur) || r.auraBelow(info.ReactionModKeyHydro, info.ZeroDur) {
			r.ecTickSrc = -1
		}
	}
}

// nextLive returns the first live mod at or after from, or a key past ReactionModKeyEnd if
// there is none.
func (r *Reactable) nextLive(from info.ReactionModKey) info.ReactionModKey {
	return from + info.ReactionModKey(bits.TrailingZeros16(r.live>>from))
}

// clearIfBelowZero clears the live bit of mod and returns true if every source of mod is
// below ZeroDur.
func (r *Reactable) clearIfBelowZero(mod info.ReactionModKey) bool {
	if r.auraBelow(mod, info.ZeroDur) {
		r.live &^= 1 << mod
		return true
	}
	return false
}

func (r *Reactable) calcCatalyzeDmg(atk info.AttackInfo, em float64) float64 {
	char := r.core.Player.ByIndex(atk.ActorIndex)
	return combat.CalcCatalyzeDmg(char.Base.Level, char, atk, em)
}
