package reactable

import (
	"fmt"
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	sswContributorMult = []float64{0.6, 0.3, 0.05, 0.05}
	sswMultStack       = []float64{0, 2, 2, 3, 3, 3, 3}
)

const (
	StellarSwirlKey    = "stellar-swirl"
	sswStackKey        = StellarSwirlKey + "-stacks"
	maxStacks          = 6
	sswContributionKey = "stellar-swirl-contribution"
	sswOwnerKey        = "stellar-swirl-owner"
	sswSrcKey          = "stellar-swirl-src"
)

type sswContribution = struct {
	dmg     float64
	isCrit  bool
	charInd int
	ae      info.AttackEvent
}

func (r *Reactable) queueStellarSwirl(charIndex int) {
	// stellar swirl triggers an aoe attack
	ai := info.AttackInfo{
		ActorIndex:       charIndex,
		DamageSrc:        r.self.Key(),
		Abil:             "Stellar Swirl",
		AttackTag:        attacks.AttackTagReactionStellarSwirl,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Anemo,
		IgnoreDefPercent: 1,
	}
	ap := combat.NewCircleHitOnTarget(r.self, nil, 5)

	var contribMap [info.MaxChars]bool

	contribMap[charIndex] = true
	for charInd, dur := range r.Durability[info.ReactionModKeyCryo] {
		if dur <= info.ZeroDur {
			continue
		}
		contribMap[charInd] = true
	}

	// TODO: Does it use frozen aura when there is cryo?
	for charInd, dur := range r.Durability[info.ReactionModKeyFrozen] {
		if dur <= info.ZeroDur {
			continue
		}
		contribMap[charInd] = true
	}

	ai, snap := r.calcStellarSwirlDmg(ai, ap, contribMap, 0.75)

	ai.ActorIndex = charIndex

	r.core.QueueAttackWithSnap(ai, snap, ap, 3)

	r.addSSwContributor(contribMap)
	r.core.Flags.Custom[sswOwnerKey] = float64(charIndex)

	if r.core.Status.Duration(StellarSwirlKey) > 0 {
		r.core.Flags.Custom[sswStackKey] += 1
		if r.core.Flags.Custom[sswStackKey] == maxStacks {
			r.detonateSSW(charIndex)
		}
	} else {
		r.core.Status.Add(StellarSwirlKey, 3*60)
		r.core.Flags.Custom[sswStackKey] = 1
		src := float64(r.core.F)
		r.core.Flags.Custom[sswSrcKey] = src
		r.core.Tasks.Add(func() {
			if r.core.Flags.Custom[sswSrcKey] != src {
				return
			}
			owner := int(r.core.Flags.Custom[sswOwnerKey])
			r.detonateSSW(owner)
		}, 3*60+1)
	}
}

func (r *Reactable) calcStellarSwirlDmg(ai info.AttackInfo, ap info.AttackPattern, contribMap [info.MaxChars]bool, mult float64) (info.AttackInfo, info.Snapshot) {
	contributions := []sswContribution{}
	for charInd, char := range r.core.Player.Chars() {
		if !contribMap[charInd] {
			continue
		}

		ai.ActorIndex = charInd
		snap := char.Snapshot(&ai)

		ae := info.AttackEvent{
			Info:        ai,
			Pattern:     ap,
			SourceFrame: r.core.F,
			Snapshot:    snap,
		}

		// Emit event so PreDamageMods can be applied to the individual LCr contributions
		r.core.Events.Emit(event.OnLunarReactionAttack, r.self, &ae)

		em := ae.Snapshot.Stats[attributes.EM]
		cr := ae.Snapshot.Stats[attributes.CR]
		cd := ae.Snapshot.Stats[attributes.CD]

		flatdmg := mult * combat.CalcLunarReactionDmg(char.Base.Level, char.ReactBonus(ae.Info), ae.Info, em)
		isCrit := false

		if r.core.Rand.Float64() <= cr {
			flatdmg *= (1 + cd)
			isCrit = true
		}

		contributions = append(contributions, sswContribution{flatdmg, isCrit, charInd, ae})
	}

	if len(contributions) == 0 {
		return ai, info.Snapshot{}
	}

	slices.SortStableFunc(contributions, func(i, j sswContribution) int {
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
		r.core.Combat.Log.NewEvent(fmt.Sprint("stellar swirl contributor ", (i+1)), glog.LogElementEvent, contr.charInd).
			Write("target", r.self.Key()).
			Write("damage", &contr.dmg).
			Write("crit", &contr.isCrit).
			Write("mult", sswContributorMult[i]).
			Write("contrib_map", contribMap).
			Write("cr", &contr.ae.Snapshot.Stats[attributes.CR]).
			Write("cd", &contr.ae.Snapshot.Stats[attributes.CD]).
			Write("em", &contr.ae.Snapshot.Stats[attributes.EM]).
			Write("base_damage_bonus", &contr.ae.Info.BaseDmgBonus).
			Write("react_bonus", r.core.Player.Chars()[contr.charInd].ReactBonus(contr.ae.Info)).
			Write("flat_dmg", &contr.ae.Info.FlatDmg).
			Write("elevation", &contr.ae.Info.Elevation)

		ai.FlatDmg += contr.dmg * sswContributorMult[i]
	}

	snap := info.Snapshot{}
	if contributions[0].isCrit {
		snap.Stats[attributes.CR] = 1.0
	}

	return ai, snap
}

func (r *Reactable) TryStellarSwirlCryo(a *info.AttackEvent) bool {
	return r.tryStellarSwirl(a, info.ReactionModKeyCryo)
}

func (r *Reactable) TryStellarSwirlFrozen(a *info.AttackEvent) bool {
	return r.tryStellarSwirl(a, info.ReactionModKeyFrozen)
}

func (r *Reactable) tryStellarSwirl(a *info.AttackEvent, mod info.ReactionModKey) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	if r.GetAuraDurability(mod) < info.ZeroDur {
		return false
	}

	a.Reacted = true

	r.core.Events.Emit(event.OnStellarSwirl, r.self, a)

	r.queueStellarSwirl(a.Info.ActorIndex)

	rd := r.reduce(attributes.Cryo, a.Info.Durability, 0.5)
	a.Info.Durability -= rd

	return true
}

func (r *Reactable) detonateSSW(owner int) {
	detonate := func() {
		contribMap := r.sswContributors()
		r.removeSSwContributors()
		ai := info.AttackInfo{
			ActorIndex:       owner,
			DamageSrc:        r.self.Key(),
			Abil:             "Stellar Swirl Detonation",
			AttackTag:        attacks.AttackTagReactionStellarSwirl,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Cryo,
			Durability:       25,
			IgnoreDefPercent: 1,
		}
		ap := combat.NewCircleHitOnTarget(r.self, nil, 5)

		stacks := min(int(r.core.Flags.Custom[sswStackKey]), maxStacks)

		ai, snap := r.calcStellarSwirlDmg(ai, ap, contribMap, sswMultStack[stacks])

		ai.ActorIndex = owner

		r.core.QueueAttackWithSnap(ai, snap, ap, 0)

		r.core.Status.Delete(StellarSwirlKey)
		r.core.Flags.Custom[sswStackKey] = 0
	}
	r.core.Tasks.Add(detonate, 5)
}

func sswContributorKey(charInd int) string {
	return fmt.Sprintf("%v-%v", sswContributionKey, charInd)
}

func (r *Reactable) addSSwContributor(charMap [info.MaxChars]bool) {
	for charInd, contrib := range charMap {
		if !contrib {
			continue
		}
		r.core.Flags.Custom[sswContributorKey(charInd)] = 1
	}
}

func (r *Reactable) removeSSwContributors() {
	for charInd := range r.core.Player.Chars() {
		r.core.Flags.Custom[sswContributorKey(charInd)] = 0
	}
}

func (r *Reactable) sswContributors() [info.MaxChars]bool {
	var contributors [info.MaxChars]bool
	for _, char := range r.core.Player.Chars() {
		contributors[char.Index()] = r.core.Flags.Custom[sswContributorKey(char.Index())] == 1
	}
	return contributors
}
