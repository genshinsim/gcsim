package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type CrystallizeShield struct {
	*shield.Tmpl
	emBonus float64
}

func (r *Reactable) TryCrystallizeElectro(a *combat.AttackEvent) bool {
	if r.Durability[Electro] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Electro, reactions.CrystallizeElectro, event.OnCrystallizeElectro)
	}
	return false
}

func (r *Reactable) TryCrystallizeHydro(a *combat.AttackEvent) bool {
	if r.Durability[Hydro] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Hydro, reactions.CrystallizeHydro, event.OnCrystallizeHydro)
	}
	return false
}

func (r *Reactable) TryCrystallizeCryo(a *combat.AttackEvent) bool {
	if r.Durability[Cryo] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Cryo, reactions.CrystallizeCryo, event.OnCrystallizeCryo)
	}
	return false
}

func (r *Reactable) TryCrystallizePyro(a *combat.AttackEvent) bool {
	if r.Durability[Pyro] > ZeroDur || r.Durability[Burning] > ZeroDur {
		reacted := r.tryCrystallizeWithEle(a, attributes.Pyro, reactions.CrystallizePyro, event.OnCrystallizePyro)
		r.burningCheck()
		return reacted
	}
	return false
}

func (r *Reactable) TryCrystallizeFrozen(a *combat.AttackEvent) bool {
	if r.Durability[Frozen] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Frozen, reactions.CrystallizeCryo, event.OnCrystallizeCryo)
	}
	return false
}

// Used to simulate booking to reduce the CD between allowable crystallize occurences without increasing the battle timer
// Because gcsim only has global frames, reduce crystallizeGCD rather than increase the cursed timer
// Param f: Number of frames to reduce the gcd
// Booking will only have an effect if crystallize has already occurred, and will reduce the timer for the next instance of crystallize only.
// Future instances of crystallize will still have the standard cd.
// "overbooking" will have no adverse effects. r.crystallizeGCD will simply go negative and will be capped at -1.
func (r *Reactable) ReduceCrystallizeGCD(f int) {
	r.crystallizeGCD -= f
	if r.crystallizeGCD < 0 {
		r.crystallizeGCD = -1
	}
	r.core.Events.Emit(event.OnTimeManip, f)
}

func (r *Reactable) tryCrystallizeWithEle(a *combat.AttackEvent, ele attributes.Element, rt reactions.ReactionType, evt event.Event) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.crystallizeGCD != -1 && r.core.F < r.crystallizeGCD {
		return false
	}
	r.crystallizeGCD = r.core.F + 60
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	r.addCrystallizeShard(char, rt, ele, r.core.F)
	// reduce
	r.reduce(ele, a.Info.Durability, 0.5)
	//TODO: confirm u can only crystallize once
	a.Info.Durability = 0
	a.Reacted = true
	// event
	r.core.Events.Emit(evt, r.self, a)
	// check freeze + ec
	switch {
	case ele == attributes.Electro && r.Durability[Hydro] > ZeroDur:
		r.checkEC()
	case ele == attributes.Frozen:
		r.checkFreeze()
	}
	return true
}

func NewCrystallizeShield(index int, typ attributes.Element, src, lvl int, em float64, expiry int) *CrystallizeShield {
	s := &CrystallizeShield{}
	s.Tmpl = &shield.Tmpl{}

	lvl--
	if lvl > 89 {
		lvl = 89
	}
	if lvl < 0 {
		lvl = 0
	}

	s.Tmpl.ActorIndex = index
	s.Tmpl.Target = -1
	s.Tmpl.Ele = typ
	s.Tmpl.ShieldType = shield.Crystallize
	s.Tmpl.Name = "Crystallize " + typ.String()
	s.Tmpl.Src = src
	s.Tmpl.HP = shieldBaseHP[lvl]
	s.Tmpl.Expires = expiry

	s.emBonus = (40.0 / 9.0) * (em / (1400 + em))

	return s
}

func (c *CrystallizeShield) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	bonus += c.emBonus
	return c.Tmpl.OnDamage(dmg, ele, bonus)
}

var shieldBaseHP = []float64{
	91.1791000366211,
	98.7076644897461,
	106.236221313477,
	113.764770507813,
	121.293319702148,
	128.821884155273,
	136.35041809082,
	143.878982543945,
	151.407516479492,
	158.936080932617,
	169.991485595703,
	181.076248168945,
	192.190368652344,
	204.048202514648,
	215.938995361328,
	227.862747192383,
	247.685943603516,
	267.542114257813,
	287.431213378906,
	303.826416015625,
	320.225219726563,
	336.627624511719,
	352.319274902344,
	368.010925292969,
	383.702545166016,
	394.432373046875,
	405.181457519531,
	415.949920654297,
	426.737640380859,
	437.544708251953,
	450.600006103516,
	463.700286865234,
	476.845581054688,
	491.127502441406,
	502.554565429688,
	514.012084960938,
	531.409606933594,
	549.979614257813,
	568.584899902344,
	584.996520996094,
	605.670349121094,
	626.38623046875,
	646.052307128906,
	665.755615234375,
	685.49609375,
	700.839416503906,
	723.333129882813,
	745.865295410156,
	768.435729980469,
	786.791931152344,
	809.538818359375,
	832.329040527344,
	855.162658691406,
	878.039611816406,
	899.484802246094,
	919.361999511719,
	946.039611816406,
	974.764221191406,
	1003.57861328125,
	1030.07702636719,
	1056.63500976563,
	1085.24633789063,
	1113.92443847656,
	1149.25866699219,
	1178.06481933594,
	1200.22375488281,
	1227.66027832031,
	1257.24304199219,
	1284.91735839844,
	1314.7529296875,
	1342.66516113281,
	1372.75244140625,
	1396.32104492188,
	1427.31237792969,
	1458.37451171875,
	1482.33581542969,
	1511.91088867188,
	1541.54931640625,
	1569.15368652344,
	1596.81433105469,
	1622.41967773438,
	1648.07397460938,
	1666.37609863281,
	1684.67822265625,
	1702.98034667969,
	1726.10473632813,
	1754.67150878906,
	1785.86657714844,
	1817.13745117188,
	1851.06030273438,
}

type CrystallizeShard struct {
	*gadget.Gadget
	// earliest that a shard can be picked up after spawn
	EarliestPickup int
	// captures the shield because em snapshots
	Shield *CrystallizeShield
	// for logging purposes
	src    int
	expiry int
}

func (r *Reactable) addCrystallizeShard(char *character.CharWrapper, rt reactions.ReactionType, typ attributes.Element, src int) {
	// delay shard spawn
	r.core.Tasks.Add(func() {
		// grab current snapshot for shield
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			DamageSrc:  r.self.Key(),
			Abil:       string(rt),
		}
		snap := char.Snapshot(&ai)
		lvl := snap.CharLvl
		// shield snapshots em on shard spawn
		em := snap.Stats[attributes.EM]
		// expiry will get set properly later
		shd := NewCrystallizeShield(char.Index, typ, src, lvl, em, -1)
		cs := NewCrystallizeShard(r.core, r.self.Shape(), shd)
		r.core.Combat.AddGadget(cs)
		r.core.Log.NewEvent(
			fmt.Sprintf("%v crystallize shard spawned", cs.Shield.Ele),
			glog.LogElementEvent,
			cs.Shield.ActorIndex,
		).
			Write("src", cs.src).
			Write("expiry", cs.expiry).
			Write("earliest_pickup", cs.EarliestPickup)
	}, 23)
}

func NewCrystallizeShard(c *core.Core, shp geometry.Shape, shd *CrystallizeShield) *CrystallizeShard {
	cs := &CrystallizeShard{}

	circ, ok := shp.(*geometry.Circle)
	if !ok {
		panic("rectangle target hurtbox is not supported for crystallize shard spawning")
	}

	// for simplicity, crystallize shards spawn randomly at radius + 0.5
	r := circ.Radius() + 0.5
	// radius 2 is ok
	cs.Gadget = gadget.New(c, geometry.CalcRandomPointFromCenter(circ.Pos(), r, r, c.Rand), 2, combat.GadgetTypCrystallizeShard)

	// shard lasts for 15s from shard spawn
	cs.Gadget.Duration = 15 * 60
	// earliest shard pickup is 54f from crystallize text, so 31f from shard spawn
	cs.EarliestPickup = c.F + 31
	cs.Shield = shd
	cs.src = c.F
	cs.expiry = c.F + cs.Gadget.Duration

	return cs
}

func (cs *CrystallizeShard) AddShieldKillShard() bool {
	// don't pick up if shard is not available for pick up yet
	if cs.Core.F < cs.EarliestPickup {
		cs.Core.Log.NewEvent(
			fmt.Sprintf("%v crystallize shard could not be picked up", cs.Shield.Ele),
			glog.LogElementEvent,
			cs.Core.Player.Active(),
		).
			Write("src", cs.src).
			Write("expiry", cs.expiry).
			Write("earliest_pickup", cs.EarliestPickup)
		return false
	}
	cs.Core.Log.NewEvent(
		fmt.Sprintf("%v crystallize shard picked up", cs.Shield.Ele),
		glog.LogElementEvent,
		cs.Core.Player.Active(),
	).
		Write("src", cs.src).
		Write("expiry", cs.expiry).
		Write("earliest_pickup", cs.EarliestPickup)
	// add shield
	cs.Shield.Expires = cs.Core.F + 15.1*60 // shield lasts for 15.1s from shard pickup
	cs.Core.Player.Shields.Add(cs.Shield)
	// kill self
	cs.Kill()
	return true
}

func (cs *CrystallizeShard) HandleAttack(atk *combat.AttackEvent) float64 {
	cs.Core.Events.Emit(event.OnGadgetHit, cs, atk)
	return 0
}
func (cs *CrystallizeShard) Attack(*combat.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (cs *CrystallizeShard) SetDirection(trg geometry.Point)                        {}
func (cs *CrystallizeShard) SetDirectionToClosestEnemy()                            {}
func (cs *CrystallizeShard) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
