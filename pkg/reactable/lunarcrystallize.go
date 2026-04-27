package reactable

import (
	"fmt"
	"slices"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	LcrKey              = "lunarcrystallize"
	LcrExtraHitOverride = "lunarcrystallize-bonus-hit-chance"
	lcrCountKey         = "lunarcrystallize-count"
	lcrDur              = 9 * 60
)

var lcrContributorMult = []float64{0.6, 0.3, 0.05, 0.05} // TODO: move to a lunar.go ?

func (r *Reactable) TryAddLCr(a *info.AttackEvent) bool {
	if r.GetAuraDurability(info.ReactionModKeyHydro) <= info.ZeroDur {
		return false
	}

	if a.Info.Durability < info.ZeroDur {
		return false
	}

	if r.core.Status.Duration(LcrKey) > 0 {
		r.extendLunarCrystallizeConstructDur()
	} else {
		// TODO: Check if constructs expiring will reset the counter
		r.core.Constructs.NewNoLimitCons(r.newLunarCrystallizeConstruct(r.self.Direction(), r.self.Pos().Add(info.Point{Y: 1, X: 0})), false)
		r.core.Constructs.NewNoLimitCons(r.newLunarCrystallizeConstruct(r.self.Direction(), r.self.Pos().Add(info.Point{Y: -0.5, X: 0.866})), false)
		r.core.Constructs.NewNoLimitCons(r.newLunarCrystallizeConstruct(r.self.Direction(), r.self.Pos().Add(info.Point{Y: -0.5, X: -0.866})), false)
	}

	r.core.Flags.Custom[lcrCountKey] += 1
	r.core.Status.Add(LcrKey, lcrDur)
	r.addLCrContributor(a)

	if r.core.Flags.Custom[lcrCountKey] >= 3 {
		// trigger three attacks
		r.core.Flags.Custom[lcrCountKey] = 0
		r.core.Events.Emit(event.OnMoondriftHarmony, r.self, &a)
		r.core.Log.NewEvent("Lunar Crystallize attack triggered", glog.LogElementEvent, a.Info.ActorIndex)
		r.DoLCrAttack()
	}

	// reduce
	consumed := r.reduce(attributes.Hydro, a.Info.Durability, 0.5)
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true
	r.lcrAtkOwner = a.Info.ActorIndex

	// event
	r.core.Events.Emit(event.OnLunarCrystallize, r.self, a)

	return true
}

// TODO this needs to be global?
func (r *Reactable) addLCrContributor(a *info.AttackEvent) {
	r.lcrContributors[a.Info.ActorIndex] = true
	for charInd, dur := range r.Durability[info.ReactionModKeyHydro] {
		if dur <= info.ZeroDur {
			continue
		}
		r.lcrContributors[charInd] = true
	}
}

func (r *Reactable) extendLunarCrystallizeConstructDur() {
	matched, _ := r.core.Constructs.ConstructsByType(construct.GeoConstructLunarCrystallize)
	for _, construct := range matched {
		c, ok := (construct).(*skillConstruct)
		if !ok {
			continue
		}
		c.expiry = r.core.F + lcrDur
	}
}

type skillConstruct struct {
	src    int
	expiry int
	react  *Reactable
	dir    info.Point
	pos    info.Point
}

func (r *Reactable) newLunarCrystallizeConstruct(dir, pos info.Point) *skillConstruct {
	return &skillConstruct{
		src:    r.core.F,
		expiry: r.core.F + lcrDur,
		react:  r,
		dir:    dir,
		pos:    pos,
	}
}

func (c *skillConstruct) OnDestruct() {}
func (c *skillConstruct) Key() int    { return c.src }
func (c *skillConstruct) Type() construct.GeoConstructType {
	return construct.GeoConstructLunarCrystallize
}
func (c *skillConstruct) Expiry() int           { return c.expiry }
func (c *skillConstruct) IsLimited() bool       { return true }
func (c *skillConstruct) Count() int            { return 1 }
func (c *skillConstruct) Direction() info.Point { return c.dir }
func (c *skillConstruct) Pos() info.Point       { return c.pos }

type lcrContribution = struct {
	dmg     float64
	isCrit  bool
	charInd int
	ae      info.AttackEvent
}

func (r *Reactable) DoLCrAttack() {
	DoLCrAttackWithContrib(r.lcrContributors, r.self, r.core, r.lcAtkOwner)
	// clear contributors after last attack
	r.core.Tasks.Add(func() {
		for i := range r.lcrContributors {
			r.lcrContributors[i] = false
		}
	}, 7)
}

// Perform a Lunar Crystallize reaction 3-hit attack with the given contributors
func DoLCrAttackWithContrib(contribMap [info.MaxChars]bool, target info.Target, core *core.Core, owner int) {
	for _, delay := range []int{1, 4, 7} {
		core.Tasks.Add(func() { doSingleLCrAttack(contribMap, target, core, owner) }, delay)
		if chance, ok := core.Flags.Custom[LcrExtraHitOverride]; ok && core.Rand.Float64() < chance {
			core.Tasks.Add(func() { doSingleLCrAttack(contribMap, target, core, owner) }, delay)
		}
	}
}

func doSingleLCrAttack(contribMap [info.MaxChars]bool, target info.Target, core *core.Core, owner int) {
	contributions := []lcrContribution{}

	ap := combat.NewSingleTargetHit(target.Key())

	// Do we need to make a new one for each character?
	ai := info.AttackInfo{
		DamageSrc:        target.Key(),
		Abil:             string(info.ReactionTypeLunarCrystallize),
		AttackTag:        attacks.AttackTagReactionLunarCrystallize,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Geo,
		IgnoreDefPercent: 1,
	}

	for charInd, char := range core.Player.Chars() {
		if !contribMap[charInd] {
			continue
		}

		ai.ActorIndex = charInd
		snap := char.Snapshot(&ai)

		ae := info.AttackEvent{
			Info:        ai,
			Pattern:     ap,
			SourceFrame: core.F,
			Snapshot:    snap,
		}

		// Emit even so PreDamageMods can be applied to the individual LCr contributions
		core.Events.Emit(event.OnLunarReactionAttack, target, &ae)

		em := ae.Snapshot.Stats[attributes.EM]
		cr := ae.Snapshot.Stats[attributes.CR]
		cd := ae.Snapshot.Stats[attributes.CD]

		flatdmg := combat.CalcLunarReactionDmg(char.Base.Level, char.ReactBonus(ae.Info), ae.Info, em)
		isCrit := false

		if core.Rand.Float64() <= cr {
			flatdmg *= (1 + cd)
			isCrit = true
		}

		contributions = append(contributions, lcrContribution{flatdmg, isCrit, charInd, ae})
	}

	if len(contributions) == 0 {
		return
	}

	slices.SortStableFunc(contributions, func(i, j lcrContribution) int {
		diff := j.dmg - i.dmg
		switch {
		case diff < 0:
			return -1
		case diff > 0:
			return 1
		default:
			return 0
		}
	})

	for i := range contributions {
		contr := &contributions[i]
		core.Combat.Log.NewEvent(fmt.Sprint("lunarcrystallize contributor ", (i+1)), glog.LogElementEvent, contr.charInd).
			Write("target", target.Key()).
			Write("damage", &contr.dmg).
			Write("crit", &contr.isCrit).
			Write("mult", lcrContributorMult[i]).
			Write("contrib_map", contribMap).
			Write("cr", &contr.ae.Snapshot.Stats[attributes.CR]).
			Write("cd", &contr.ae.Snapshot.Stats[attributes.CD]).
			Write("em", &contr.ae.Snapshot.Stats[attributes.EM]).
			Write("base_damage_bonus", &contr.ae.Info.BaseDmgBonus).
			Write("react_bonus", core.Player.Chars()[contr.charInd].ReactBonus(contr.ae.Info)).
			Write("flat_dmg", &contr.ae.Info.FlatDmg).
			Write("elevation", &contr.ae.Info.Elevation)

		ai.FlatDmg += contr.dmg * lcrContributorMult[i]
	}

	snap := info.Snapshot{}
	if contributions[0].isCrit {
		snap.Stats[attributes.CR] = 1.0
	}
	// LCr is owned by the character that last triggered Lunar Crystallize for this instance of Lunar Crystallize
	ai.ActorIndex = owner
	core.QueueAttackWithSnap(
		ai,
		snap,
		ap,
		0,
	)
}
