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
	LcrExtraHitOverride    = "lunarcrystallize-bonus-hit-chance"
	lcrContributionKey     = "lunarcrystallize-contribution"
	lcrCountKey            = "lunarcrystallize-count"
	lcrDur                 = 9 * 60
	lcrTravel              = 10
	moondriftHarmonyICDKey = "moondrift-harmony-icd"
)

var (
	lcrContributorMult = []float64{0.6, 0.3, 0.05, 0.05}
	moondriftOffset    = []info.Point{{Y: 1, X: 0}, {Y: -0.5, X: 0.866}, {Y: -0.5, X: -0.866}}
	lcrHitmarks        = []int{13, 13 + 12, 13 + 12 + 12}
)

func (r *Reactable) TryAddLCr(a *info.AttackEvent) bool {
	if r.GetAuraDurability(info.ReactionModKeyHydro) <= info.ZeroDur {
		return false
	}

	if a.Info.Durability < info.ZeroDur {
		return false
	}

	r.core.Flags.Custom[lcrCountKey] = min(r.core.Flags.Custom[lcrCountKey]+1, 6)
	r.contributeToLCrWithAttack(a)
	// event
	r.core.Events.Emit(event.OnLunarCrystallize, r.self, a)

	moondriftsNearby := 0
	moondrifts, _ := r.core.Constructs.ConstructsByType(construct.GeoConstructLunarCrystallize)
	playerPos := r.core.Combat.Player().Pos()
	for _, moondrift := range moondrifts {
		if playerPos.Distance(moondrift.Pos()) < 20 {
			moondriftsNearby += 1
		}
	}
	if moondriftsNearby < 3 {
		// TODO: Set accurate location for spawned moondrifts
		for i := range 3 - moondriftsNearby {
			moondrift := r.newLunarCrystallizeConstruct(r.self.Direction(), r.self.Pos().Add(moondriftOffset[i]))
			r.core.Constructs.NewNoLimitCons(moondrift, false)
		}
	} else if r.core.Flags.Custom[lcrCountKey] >= 3 && r.core.Status.Duration(moondriftHarmonyICDKey) == 0 {
		// trigger three attacks
		r.core.Flags.Custom[lcrCountKey] -= 3
		r.core.Events.Emit(event.OnMoondriftHarmony, r.self, &a)
		r.core.Log.NewEvent("Moondrift Harmony triggered", glog.LogElementEvent, a.Info.ActorIndex)
		r.core.Status.Add(moondriftHarmonyICDKey, lcrHitmarks[len(lcrHitmarks)-1])
		r.DoLCrAttack(a.Info.ActorIndex)
		r.extendNearbyLunarCrystallizeConstructDur()
	}

	// reduce
	consumed := r.reduce(attributes.Hydro, a.Info.Durability, 0.5)
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	return true
}

func contributorKey(charInd int) string {
	return fmt.Sprintf("%v-%v", lcrContributionKey, charInd)
}

func (r *Reactable) addLCrContributor(charInd int) {
	r.core.Flags.Custom[contributorKey(charInd)] = 1
}

func (r *Reactable) removeLCrContributor(charInd int) {
	r.core.Flags.Custom[contributorKey(charInd)] = 0
}

func (r *Reactable) lcrContributors() [info.MaxChars]bool {
	var contributors [info.MaxChars]bool
	for _, char := range r.core.Player.Chars() {
		contributors[char.Index()] = r.core.Flags.Custom[contributorKey(char.Index())] == 1
	}
	return contributors
}

func (r *Reactable) contributeToLCrWithAttack(a *info.AttackEvent) {
	r.addLCrContributor(a.Info.ActorIndex)
	for charInd, dur := range r.Durability[info.ReactionModKeyHydro] {
		if dur <= info.ZeroDur {
			continue
		}
		r.addLCrContributor(charInd)
	}
}

func (r *Reactable) extendNearbyLunarCrystallizeConstructDur() {
	matched, _ := r.core.Constructs.ConstructsByType(construct.GeoConstructLunarCrystallize)
	playerPos := r.core.Combat.Player().Pos()
	for _, construct := range matched {
		c, ok := (construct).(*LunarCrystallizeConstruct)
		if !ok {
			continue
		}
		if playerPos.Distance(construct.Pos()) > 20 {
			continue
		}
		c.expiry = r.core.F + lcrDur
	}
}

type lcrContribution = struct {
	dmg     float64
	isCrit  bool
	charInd int
	ae      info.AttackEvent
}

func (r *Reactable) DoLCrAttack(owner int) {
	DoLCrAttackWithContrib(r.lcrContributors(), r.self, r.core, owner)
	// clear contributors after obtaining the contributor map
	for _, char := range r.core.Player.Chars() {
		r.removeLCrContributor(char.Index())
	}
}

// Perform a Lunar Crystallize reaction 3-hit attack with the given contributors
func DoLCrAttackWithContrib(contribMap [info.MaxChars]bool, target info.Target, core *core.Core, owner int) {
	for _, delay := range lcrHitmarks {
		core.Tasks.Add(func() { doSingleLCrAttack(contribMap, target, core, owner) }, delay)
		if chance, ok := core.Flags.Custom[LcrExtraHitOverride]; ok && core.Rand.Float64() < chance {
			core.Tasks.Add(func() { doSingleLCrAttack(contribMap, target, core, owner) }, delay)
		}
	}
}

func doSingleLCrAttack(contribMap [info.MaxChars]bool, target info.Target, core *core.Core, owner int) {
	contributions := []lcrContribution{}

	ap := combat.NewSingleTargetHit(target.Key())

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

		// Emit event so PreDamageMods can be applied to the individual LCr contributions
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
		lcrTravel,
	)
}
