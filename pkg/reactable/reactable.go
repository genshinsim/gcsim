package reactable

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
)

type Modifier int

const (
	Invalid Modifier = iota
	Electro
	Pyro
	Cryo
	Hydro
	BurningFuel
	SpecialDecayDelim
	Dendro
	Quicken
	Frozen
	Anemo
	Geo
	Burning
	EndModifier
)

var modifierElement = []attributes.Element{
	attributes.UnknownElement,
	attributes.Electro,
	attributes.Pyro,
	attributes.Cryo,
	attributes.Hydro,
	attributes.Dendro,
	attributes.UnknownElement,
	attributes.Dendro,
	attributes.Quicken,
	attributes.Frozen,
	attributes.Anemo,
	attributes.Geo,
	attributes.Pyro,
	attributes.UnknownElement,
}

var ModifierString = []string{
	"",
	"electro",
	"pyro",
	"cryo",
	"hydro",
	"dendro-fuel",
	"",
	"dendro",
	"quicken",
	"frozen",
	"anemo",
	"geo",
	"burning",
	"",
}

var elementToModifier = map[attributes.Element]Modifier{
	attributes.Electro: Electro,
	attributes.Pyro:    Pyro,
	attributes.Cryo:    Cryo,
	attributes.Hydro:   Hydro,
	attributes.Dendro:  Dendro,
}

func (r Modifier) Element() attributes.Element { return modifierElement[r] }
func (r Modifier) String() string              { return ModifierString[r] }

func (r Modifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(ModifierString[r])
}

func (r *Modifier) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range ModifierString {
		if v == s {
			*r = Modifier(i)
			return nil
		}
	}
	return errors.New("unrecognized ReactableModifier")
}

type Reactable struct {
	Durability [EndModifier]reactions.Durability
	DecayRate  [EndModifier]reactions.Durability
	Mutable    [EndModifier]bool
	// Source     []int //source frame of the aura
	self combat.Target
	core *core.Core
	// ec specific
	ecAtk      combat.AttackInfo // index of owner of next ec ticks
	ecSnapshot combat.Snapshot
	ecTickSrc  int
	// burning specific
	burningAtk      combat.AttackInfo
	burningSnapshot combat.Snapshot
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

const frzDelta reactions.Durability = 2.5 / (60 * 60) // 2 * 1.25
const frzDecayCap reactions.Durability = 10.0 / 60.0

const ZeroDur reactions.Durability = 0.00000000001

func (r *Reactable) Init(self combat.Target, c *core.Core) *Reactable {
	for i := Invalid; i < EndModifier; i++ {
		r.Mutable[i] = true
	}

	r.self = self
	r.core = c
	r.DecayRate[Frozen] = frzDecayCap
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
		if v > ZeroDur {
			result = append(result, Modifier(i).String()+": "+strconv.FormatFloat(float64(v), 'f', 3, 64))
		}
	}
	return result
}

func (r *Reactable) AuraCount() int {
	count := 0
	for _, v := range r.Durability {
		if v > ZeroDur {
			count++
		}
	}
	return count
}

func (r *Reactable) React(a *combat.AttackEvent) {
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
		//https://docs.google.com/spreadsheets/d/1lJSY2zRIkFDyLZxIor0DVMpYXx3E_jpDrSUZvQijesc/edit#gid=0
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
func (r *Reactable) AttachOrRefill(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if a.Reacted {
		return false
	}
	// handle pyro, electro, hydro, cryo, dendro
	// special attachment of dendro (burning fuel) is handled in tryBurning
	if mod, ok := elementToModifier[a.Info.Element]; ok {
		r.attachOrRefillNormalEle(mod, a.Info.Durability)
		return true
	}
	return false
}

// attachOrRefillNormalEle is used for pyro, electro, hydro, cryo, and dendro which don't have special attachment
// rules
func (r *Reactable) attachOrRefillNormalEle(mod Modifier, dur reactions.Durability) {
	amt := 0.8 * dur
	if mod == Pyro {
		r.attachOverlapRefreshDuration(Pyro, amt, 6*dur+420)
	} else {
		r.attachOverlap(mod, amt, 6*dur+420)
	}
}

func (r *Reactable) attachOverlap(mod Modifier, amt, length reactions.Durability) {
	if !r.Mutable[mod] {
		return
	}
	if r.Durability[mod] > ZeroDur {
		add := max(amt-r.Durability[mod], 0)
		if add > 0 {
			r.addDurability(mod, add)
		}
	} else {
		r.Durability[mod] = amt
		if length > ZeroDur {
			r.DecayRate[mod] = amt / length
		}
	}
}

func (r *Reactable) attachOverlapRefreshDuration(mod Modifier, amt, length reactions.Durability) {
	if !r.Mutable[mod] {
		return
	}
	if amt < r.Durability[mod] {
		return
	}
	r.Durability[mod] = amt
	r.DecayRate[mod] = amt / length
}

func (r *Reactable) attachBurning() {
	r.Durability[Burning] = 50
	r.DecayRate[Burning] = 0
}

func (r *Reactable) addDurability(mod Modifier, amt reactions.Durability) {
	if !r.Mutable[mod] {
		return
	}
	r.Durability[mod] += amt
	r.core.Events.Emit(event.OnAuraDurabilityAdded, r.self, mod, amt)
}

// AuraCountains returns true if any element e is active on the target
func (r *Reactable) AuraContains(e ...attributes.Element) bool {
	for _, v := range e {
		for i := Invalid; i < EndModifier; i++ {
			if i.Element() == v && r.Durability[i] > ZeroDur {
				return true
			}
		}
		//TODO: not sure if this is best way to go about it? perhaps supplying frozen element is better?
		if v == attributes.Cryo && r.Durability[Frozen] > ZeroDur {
			return true
		}
	}
	return false
}

func (r *Reactable) IsBurning() bool {
	if r.Durability[BurningFuel] > ZeroDur && r.Durability[Burning] > ZeroDur {
		return true
	}
	return false
}

// reduce the requested element by dur * factor, return the amount of dur consumed
// if multiple modifier with same element are present, all of them are reduced
// the max on reduced is used for consumption purpose
func (r *Reactable) reduce(e attributes.Element, dur, factor reactions.Durability) reactions.Durability {
	m := dur * factor // maximum amount reduceable
	var reduced reactions.Durability

	for i := Invalid; i < EndModifier; i++ {
		if i.Element() != e {
			continue
		}
		if r.Durability[i] < ZeroDur {
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

		if r.Mutable[i] {
			r.Durability[i] -= red
		}

		if red > reduced {
			reduced = red
		}
	}

	return reduced / factor
}

func (r *Reactable) deplete(m Modifier) {
	if r.Durability[m] <= ZeroDur {
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
	for i := Invalid; i < SpecialDecayDelim; i++ {
		// skip zero decay rates i.e. modifiers that don't decay (i.e. burning)
		if r.DecayRate[i] == 0 {
			continue
		}
		if r.Durability[i] > ZeroDur && r.Mutable[i] {
			r.Durability[i] -= r.DecayRate[i]
			r.deplete(i)
		}
	}

	// check burning first since that affects dendro/quicken decay
	if r.burningTickSrc > -1 && r.Durability[BurningFuel] < ZeroDur {
		// reset src when burning fuel is gone
		r.burningTickSrc = -1
		// remove burning
		r.Durability[Burning] = 0
		// remove existing dendro and quicken
		r.Durability[Dendro] = 0
		r.DecayRate[Dendro] = 0
		r.Durability[Quicken] = 0
		r.DecayRate[Quicken] = 0
	}

	// if burning fuel is present, dendro and quicken uses burning fuel decay rate
	// otherwise it uses it's own
	for i := Dendro; i <= Quicken; i++ {
		if r.Durability[i] < ZeroDur {
			continue
		}
		if !r.Mutable[i] {
			continue
		}
		rate := r.DecayRate[i]
		if r.Durability[BurningFuel] > ZeroDur {
			rate = r.DecayRate[BurningFuel]
			if i == Dendro {
				rate = max(rate, r.DecayRate[i]*2)
			}
		}
		r.Durability[i] -= rate
		r.deplete(i)
	}

	// for freeze, durability can be calculated as:
	// d_f(t) = -1.25 * (t/60)^2 - k * (t/60) + d_f(0)
	if r.Durability[Frozen] > ZeroDur {
		// ramp up decay rate first
		r.DecayRate[Frozen] += frzDelta
		r.Durability[Frozen] -= r.DecayRate[Frozen] / reactions.Durability(1.0-r.FreezeResist)

		r.checkFreeze()
	} else if r.DecayRate[Frozen] > frzDecayCap { // otherwise ramp down decay rate
		r.DecayRate[Frozen] -= frzDelta * 2

		// cap decay
		if r.DecayRate[Frozen] < frzDecayCap {
			r.DecayRate[Frozen] = frzDecayCap
		}
	}

	// for ec we need to reset src if ec is gone
	if r.ecTickSrc > -1 {
		if r.Durability[Electro] < ZeroDur || r.Durability[Hydro] < ZeroDur {
			r.ecTickSrc = -1
		}
	}
}

func calcReactionDmg(char *character.CharWrapper, atk combat.AttackInfo, em float64) (float64, combat.Snapshot) {
	lvl := char.Base.Level - 1
	if lvl > 89 {
		lvl = 89
	}
	if lvl < 0 {
		lvl = 0
	}
	snap := combat.Snapshot{
		CharLvl:  char.Base.Level,
		ActorEle: char.Base.Element,
	}
	snap.Stats[attributes.EM] = em
	return (1 + ((16 * em) / (2000 + em)) + char.ReactBonus(atk)) * reactionLvlBase[lvl], snap
}

func (r *Reactable) calcCatalyzeDmg(atk combat.AttackInfo, em float64) float64 {
	char := r.core.Player.ByIndex(atk.ActorIndex)
	lvl := char.Base.Level - 1
	if lvl > 89 {
		lvl = 89
	}
	if lvl < 0 {
		lvl = 0
	}
	return (1 + ((5 * em) / (1200 + em)) + r.core.Player.ByIndex(atk.ActorIndex).ReactBonus(atk)) * reactionLvlBase[lvl]
}

var reactionLvlBase = []float64{
	17.1656055450439,
	18.5350475311279,
	19.9048538208007,
	21.27490234375,
	22.6453990936279,
	24.6496124267578,
	26.6406421661376,
	28.8685874938964,
	31.3676795959472,
	34.1433448791503,
	37.201000213623,
	40.6599998474121,
	44.4466667175292,
	48.5635185241699,
	53.7484817504882,
	59.0818977355957,
	64.4200439453125,
	69.7244567871093,
	75.1231384277343,
	80.5847778320312,
	86.1120300292968,
	91.703742980957,
	97.24462890625,
	102.812644958496,
	108.409561157226,
	113.201690673828,
	118.102905273437,
	122.979316711425,
	129.727325439453,
	136.292907714843,
	142.670852661132,
	149.029022216796,
	155.4169921875,
	161.825500488281,
	169.106307983398,
	176.518081665039,
	184.07273864746,
	191.709518432617,
	199.556915283203,
	207.382049560546,
	215.398895263671,
	224.165664672851,
	233.502166748046,
	243.35057067871,
	256.063079833984,
	268.543487548828,
	281.526062011718,
	295.013641357421,
	309.067199707031,
	323.601593017578,
	336.757537841796,
	350.530303955078,
	364.482696533203,
	378.619171142578,
	398.600402832031,
	416.398254394531,
	434.386993408203,
	452.951049804687,
	472.606231689453,
	492.884887695312,
	513.568542480468,
	539.103210449218,
	565.510559082031,
	592.538757324218,
	624.443420410156,
	651.470153808593,
	679.496826171875,
	707.794067382812,
	736.671447753906,
	765.640258789062,
	794.773376464843,
	824.677368164062,
	851.157775878906,
	877.742065429687,
	914.229125976562,
	946.746765136718,
	979.411376953125,
	1011.22302246093,
	1044.79174804687,
	1077.44372558593,
	1109.99755859375,
	1142.9765625,
	1176.36950683593,
	1210.18444824218,
	1253.83569335937,
	1288.95275878906,
	1325.48413085937,
	1363.45690917968,
	1405.09741210937,
	1446.853515625,
}
