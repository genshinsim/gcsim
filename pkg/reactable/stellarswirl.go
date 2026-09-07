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

var sswContributorMult = []float64{0.6, 0.3, 0.05, 0.05}

const (
	StellarSwirlKey    = "stellar-swirl"
	sswStackKey        = StellarSwirlKey + "-stacks"
	sswMaxStacks       = 6
	sswContributionKey = "stellar-swirl-contribution"
	sswOwnerKey        = "stellar-swirl-owner"
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
		ICDTag:           attacks.ICDTagSwirlCryo,   // TODO: use stellar swirl ICD tag
		ICDGroup:         attacks.ICDGroupReactionA, // TODO: use stellar swirl ICD group
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

		// Based on testing, only Cryo characters that trigger frozen can contribute
		// even if a non Cryo character triggers frozen with a cryo attack via Chongyun,
		// they won't contribute to the reaction
		if r.core.Player.Chars()[charInd].Base.Element != attributes.Cryo {
			continue
		}
		contribMap[charInd] = true
	}

	for _, e := range r.core.Combat.Enemies() {
		if willHit, _ := e.AttackWillLand(ap); !willHit {
			continue
		}
		ai, snap := r.calcStellarSwirlDmg(e, ai, ap, contribMap, 0.75)
		ai.ActorIndex = charIndex
		r.core.QueueAttackWithSnap(ai, snap, combat.NewSingleTargetHit(e.Key()), 3)
	}

	r.addSSwContributor(contribMap)
	r.setSSwOwner(charIndex)

	r.addSSwStack()

	vortex := r.nearbySSwVortex()
	if vortex == nil {
		vortex = r.newStellarVortex()
	}

	if r.sswStacks() >= sswMaxStacks {
		vortex.Kill()
	}
}

func (r *Reactable) nearbySSwVortex() info.Gadget {
	for _, g := range r.core.Combat.Gadgets() {
		if g.GadgetTyp() == info.GadgetTypStellarVortex {
			return g
		}
	}
	return nil
}

func (r *Reactable) calcStellarSwirlDmg(target info.Target, ai info.AttackInfo, ap info.AttackPattern, contribMap [info.MaxChars]bool, mult float64) (info.AttackInfo, info.Snapshot) {
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

		// Emit event so PreDamageMods can be applied to the individual contributions
		r.core.Events.Emit(event.OnSpecialReactionAttack, target, &ae)

		em := ae.Snapshot.Stats[attributes.EM]
		cr := ae.Snapshot.Stats[attributes.CR]
		cd := ae.Snapshot.Stats[attributes.CD]

		flatdmg := mult * combat.CalcSpecialReactionDmg(char.Base.Level, char.ReactBonus(ae.Info), ae.Info, em)
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
			Write("target", target.Key()).
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

func (r *Reactable) addSSwStack() {
	r.core.Flags.Custom[sswStackKey] += min(r.core.Flags.Custom[sswStackKey]+1, sswMaxStacks)
}

func (r *Reactable) sswStacks() int {
	return min(int(r.core.Flags.Custom[sswStackKey]), sswMaxStacks)
}

func (r *Reactable) setSSwOwner(char int) {
	r.core.Flags.Custom[sswOwnerKey] = float64(char)
}

func (r *Reactable) sswOwner() int {
	return int(r.core.Flags.Custom[sswOwnerKey])
}
