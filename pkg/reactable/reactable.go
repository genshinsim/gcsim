package reactable

import (
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

type Reactable struct {
	Durability []combat.Durability
	DecayRate  []combat.Durability
	// Source     []int //source frame of the aura
	self combat.Target
	core *core.Core
	//ec specific
	ecSnapshot combat.AttackInfo //index of owner of next ec ticks
	ecTickSrc  int
	//hitlag
	// animationFreeze float64
}

const frzDelta combat.Durability = 2.5 / (60 * 60) // 2 * 1.25
const frzDecayCap combat.Durability = 10.0 / 60.0

const ZeroDur combat.Durability = 0.00000000001

func (r *Reactable) Init(self combat.Target, c *core.Core) *Reactable {
	r.self = self
	r.core = c
	r.Durability = make([]combat.Durability, attributes.ElementDelimAttachable)
	r.DecayRate = make([]combat.Durability, attributes.ElementDelimAttachable)
	r.DecayRate[attributes.Frozen] = frzDecayCap
	r.ecTickSrc = -1
	return r
}

func (r *Reactable) ActiveAuraString() []string {
	var result []string
	for i, v := range r.Durability {
		if v > ZeroDur {
			result = append(result, attributes.ElementString[i]+": "+strconv.FormatFloat(float64(v), 'f', 3, 64))
		}
	}

	return result
}

func (r *Reactable) React(a *combat.AttackEvent) {
	//before all else, check for shatter first
	switch count := r.auraCount(); count {
	case 0:
		r.tryAttach(a.Info.Element, &a.Info.Durability)
	case 1:
		r.tryRefill(a.Info.Element, &a.Info.Durability)
		//check if refilled ok; return if so
		if a.Info.Durability < ZeroDur {
			return
		}
		fallthrough
	default:
		//TODO: double check order of reactions
		switch a.Info.Element {
		case attributes.Electro:
			r.tryAggravate(a)
			r.tryOverload(a)
			r.tryFrozenSuperconduct(a)
			r.trySuperconduct(a)
			r.tryQuicken(a)
			r.tryAddEC(a)
		case attributes.Pyro:
			r.tryOverload(a)
			r.tryMelt(a)
			r.tryVaporize(a)
			r.tryMeltFrozen(a)
		case attributes.Cryo:
			r.trySuperconduct(a)
			r.tryFreeze(a)
			r.tryMelt(a)
		case attributes.Hydro:
			r.tryAddEC(a)
			r.tryVaporize(a)
			r.tryFreeze(a)
		case attributes.Anemo:
			r.trySwirlElectro(a)
			r.trySwirlHydro(a)
			r.trySwirlCryo(a)
			r.trySwirlPyro(a)
			r.trySwirlFrozen(a)
		case attributes.Geo:
			r.tryCrystallize(a)
		case attributes.Dendro:
			r.trySpread(a)
			r.tryQuicken(a)
		default:
			//do nothing
			return
		}
	}

}

func (r *Reactable) AuraContains(e ...attributes.Element) bool {
	for _, v := range e {
		if r.Durability[v] > ZeroDur {
			return true
		}
		if v == attributes.Cryo && r.Durability[attributes.Frozen] > ZeroDur {
			return true
		}
	}
	return false
}

func (r *Reactable) AuraType() attributes.Element {
	if r.Durability[attributes.Frozen] > ZeroDur {
		return attributes.Frozen
	}
	if r.Durability[attributes.Electro] > ZeroDur && r.Durability[attributes.Hydro] > ZeroDur {
		return attributes.EC
	}

	for i, v := range r.Durability {
		if v > 0 {
			return attributes.Element(i)
		}
	}

	return attributes.NoElement
}

func (r *Reactable) auraCount() int {
	count := 0
	for _, v := range r.Durability {
		if v > ZeroDur {
			count++
		}
	}
	return count
}

func (r *Reactable) tryAttach(ele attributes.Element, dur *combat.Durability) {
	//can't attach >= frozen
	if ele >= attributes.Frozen {
		return
	}
	if *dur < ZeroDur {
		return
	}
	r.attach(ele, *dur, 0.8)
	*dur = 0
}

func (r *Reactable) tryRefill(ele attributes.Element, dur *combat.Durability) {
	//shouldn't be >= frozen
	if ele >= attributes.Frozen {
		return
	}
	if *dur < ZeroDur {
		return
	}
	//must already have existing element
	if r.Durability[ele] < ZeroDur {
		return
	}
	r.refill(ele, *dur, 0.8)
	*dur = 0
}

func (r *Reactable) calcReactionDmg(atk combat.AttackInfo, em float64) float64 {
	char := r.core.Player.ByIndex(atk.ActorIndex)
	lvl := char.Base.Level - 1
	if lvl > 89 {
		lvl = 89
	}
	if lvl < 0 {
		lvl = 0
	}
	return (1 + ((16 * em) / (2000 + em)) + r.core.Player.ByIndex(atk.ActorIndex).ReactBonus(atk)) * reactionLvlBase[lvl]
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

func (r *Reactable) attach(e attributes.Element, dur combat.Durability, m combat.Durability) {
	//calculate duration based on dur
	r.DecayRate[e] = m * dur / (6*dur + 420)
	r.addDurability(e, m*dur)
}

func (r *Reactable) refill(e attributes.Element, dur combat.Durability, m combat.Durability) {
	add := max(dur*m-r.Durability[e], 0)
	if add > 0 {
		r.addDurability(e, add)
	}
}

//reduce the requested element by dur * factor, return the amount of dur consumed
func (r *Reactable) reduce(e attributes.Element, dur combat.Durability, factor combat.Durability) (consumed combat.Durability) {
	t := dur * factor
	if t > r.Durability[e] {
		t = r.Durability[e]
	}
	r.Durability[e] -= t
	consumed = t / factor

	return

	//if dur * factor > amount of existing element, then set amont of existing element to
	//0; and consumed is equal to dur / facotr
	if dur*factor >= r.Durability[e] {
		consumed = r.Durability[e] / factor
		//aura is depleted
		r.Durability[e] = 0
		return
	}
	//otherwise consumed = dur
	//TODO: this is wrong. should be just = dur
	consumed = dur / factor
	r.Durability[e] -= dur * factor
	return
}

func (r *Reactable) addDurability(e attributes.Element, dur combat.Durability) {
	r.Durability[e] += dur
	r.core.Events.Emit(event.OnAuraDurabilityAdded, r.self, e, dur)
}

// func (r *Reactable) ApplyHitlag(factor, dur float64) {
// 	//freeze for dur * (1-factor)
// 	r.animationFreeze += dur * (1 - factor)
// }

func (r *Reactable) Tick() {
	//TODO: further testing/in game check this would be nice
	// if r.animationFreeze > 0 {
	// 	r.animationFreeze -= 1
	// 	//reset back to 0 if it was fractional
	// 	if r.animationFreeze < 0 {
	// 		r.animationFreeze = 0
	// 	}
	// 	//skip this tick (aka no decay)
	// 	r.core.Log.NewEvent("reactable skipping tick", glog.LogHitlagEvent, -1).
	// 		Write("animationFreeze", r.animationFreeze)
	// 	//in reality this is not quite accurate. the ticks aren't frozen but instead
	// 	//the decay still happens but at 0.01x instead of the normal speed
	// 	return
	// }

	//duability is reduced by decay * (1 + purge)
	//where purge is 0 for anything that's not freeze
	//for freeze, purge = 0.25 * time
	//while defrosting, purge rate is reduced back down to 0 at new = old - 0.5 * time
	//where time is in seconds
	//
	//per frame then we have decay * (1 + 0.25 * (x/60))

	for i := attributes.Element(0); i < attributes.Frozen; i += 1 {
		if r.Durability[i] > ZeroDur {
			r.Durability[i] -= r.DecayRate[i]
			if r.Durability[i] <= ZeroDur {
				r.Durability[i] = 0
				r.DecayRate[i] = 0
				// log.Println(r.self)
				// log.Println("ele", attributes.Element(i).String())
				// log.Println("core", r.core)
				// log.Println("frame", r.core.F)
				r.core.Events.Emit(event.OnAuraDurabilityDepleted, r.self, attributes.Element(i))
			}
		}
	}

	//for freeze, durability can be calculated as:
	//d_f(t) = -1.25 * (t/60)^2 - k * (t/60) + d_f(0)
	if r.Durability[attributes.Frozen] > ZeroDur {
		//ramp up decay rate first
		r.DecayRate[attributes.Frozen] += frzDelta
		r.Durability[attributes.Frozen] -= r.DecayRate[attributes.Frozen]

		r.checkFreeze()
	} else if r.DecayRate[attributes.Frozen] > frzDecayCap { //otherwise ramp down decay rate
		r.DecayRate[attributes.Frozen] -= frzDelta * 2

		//cap decay
		if r.DecayRate[attributes.Frozen] < frzDecayCap {
			r.DecayRate[attributes.Frozen] = frzDecayCap
		}
	}

	//for ec we need to reset src if ec is gone
	if r.ecTickSrc > -1 {
		if r.Durability[attributes.Electro] < ZeroDur || r.Durability[attributes.Hydro] < ZeroDur {
			r.ecTickSrc = -1
		}
	}
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
