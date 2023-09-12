package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
)

type CrystallizeShield struct {
	*shield.Tmpl
	emBonus float64
}

func (r *Reactable) TryCrystallizeElectro(a *combat.AttackEvent) bool {
	if r.Durability[ModifierElectro] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Electro, reactions.CrystallizeElectro, event.OnCrystallizeElectro)
	}
	return false
}

func (r *Reactable) TryCrystallizeHydro(a *combat.AttackEvent) bool {
	if r.Durability[ModifierHydro] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Hydro, reactions.CrystallizeHydro, event.OnCrystallizeHydro)
	}
	return false
}

func (r *Reactable) TryCrystallizeCryo(a *combat.AttackEvent) bool {
	if r.Durability[ModifierCryo] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Cryo, reactions.CrystallizeCryo, event.OnCrystallizeCryo)
	}
	return false
}

func (r *Reactable) TryCrystallizePyro(a *combat.AttackEvent) bool {
	if r.Durability[ModifierPyro] > ZeroDur || r.Durability[ModifierBurning] > ZeroDur {
		reacted := r.tryCrystallizeWithEle(a, attributes.Pyro, reactions.CrystallizePyro, event.OnCrystallizePyro)
		r.burningCheck()
		return reacted
	}
	return false
}

func (r *Reactable) TryCrystallizeFrozen(a *combat.AttackEvent) bool {
	if r.Durability[ModifierFrozen] > ZeroDur {
		return r.tryCrystallizeWithEle(a, attributes.Frozen, reactions.CrystallizeCryo, event.OnCrystallizeCryo)
	}
	return false
}

func (r *Reactable) tryCrystallizeWithEle(a *combat.AttackEvent, ele attributes.Element, rt reactions.ReactionType, evt event.Event) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	// grab current snapshot for shield
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	ai := combat.AttackInfo{
		ActorIndex: a.Info.ActorIndex,
		DamageSrc:  r.self.Key(),
		Abil:       string(rt),
	}
	snap := char.Snapshot(&ai)
	shd := NewCrystallizeShield(char.Index, ele, r.core.F, snap.CharLvl, snap.Stats[attributes.EM], r.core.F+900)
	r.core.Player.Shields.Add(shd)
	// reduce
	r.reduce(ele, a.Info.Durability, 0.5)
	//TODO: confirm u can only crystallize once
	a.Info.Durability = 0
	a.Reacted = true
	// event
	r.core.Events.Emit(evt, r.self, a)
	// check freeze + ec
	switch {
	case ele == attributes.Electro && r.Durability[ModifierHydro] > ZeroDur:
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
	s.Tmpl.Ele = typ
	s.Tmpl.ShieldType = shield.ShieldCrystallize
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
