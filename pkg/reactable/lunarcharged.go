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

const (
	LcKey    = "lunarcharged-cloud"
	lcSrcKey = "lunarcharged-cloud-src"
	lcIcdKey = "lunarcharged-cloud-icd"
)

var lcContributorMult = []float64{1.0, 1.0 / 2.0, 1.0 / 12.0, 1.0 / 12.0}

func (r *Reactable) TryAddLC(a *info.AttackEvent) bool {
	if a.Info.Durability < info.ZeroDur {
		return false
	}
	// if there's still frozen left don't try to ec
	// game actively rejects lc reaction if frozen is present
	if r.GetAuraDurability(info.ReactionModKeyFrozen) > info.ZeroDur {
		return false
	}

	// adding lc or hydro just adds to durability
	switch a.Info.Element {
	case attributes.Hydro:
		// if there's no existing hydro or electro then do nothing
		if r.GetAuraDurability(info.ReactionModKeyElectro) < info.ZeroDur {
			return false
		}
		// add to hydro durability (can't add if the atk already reacted)
		// TODO: this shouldn't happen here
		if !a.Reacted {
			r.attachOrRefillNormalEle(info.ReactionModKeyHydro, a.Info.Durability, a.Info.ActorIndex)
		}
	case attributes.Electro:
		// if there's no existing hydro or electro then do nothing
		if r.GetAuraDurability(info.ReactionModKeyHydro) < info.ZeroDur {
			return false
		}
		// add to electro durability (can't add if the atk already reacted)
		if !a.Reacted {
			r.attachOrRefillNormalEle(info.ReactionModKeyElectro, a.Info.Durability, a.Info.ActorIndex)
		}
	default:
		return false
	}

	a.Reacted = true
	r.core.Events.Emit(event.OnLunarCharged, r.self, a)

	// if not LC cloud exists, create a new cloud and start ticking
	if r.core.Status.Duration(LcKey) == 0 {
		r.core.Flags.Custom[lcSrcKey] = float64(r.core.F)
		r.core.Tasks.Add(r.nextLCTick(r.core.F), 9)
	}
	// the LC cloud lasts 5.5s, or lifetime is extended to last another 5.5s
	r.core.Status.Add(LcKey, 5.5*60)
	return true
}

type lcContribution = struct {
	dmg     float64
	isCrit  bool
	charInd int
	cr      float64
	cd      float64
	em      float64
}

func (r *Reactable) DoLCAttack() {
	contributions := []lcContribution{}

	ap := combat.NewSingleTargetHit(r.self.Key())

	// Do we need to make a new one for each character?
	ai := info.AttackInfo{
		DamageSrc:        r.self.Key(),
		Abil:             string(info.ReactionTypeLunarCharged),
		AttackTag:        attacks.AttackTagReactionLunarCharge,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Electro,
		IgnoreDefPercent: 1,
	}

	for charInd, char := range r.core.Player.Chars() {
		if r.Durability[info.ReactionModKeyHydro][charInd] <= info.ZeroDur && r.Durability[info.ReactionModKeyElectro][charInd] <= info.ZeroDur {
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

		// Emit even so PreDamageMods can be applied to the individual LC contributions
		// Is there a way to collect these attackMods to show in logs?
		r.core.Events.Emit(event.OnLunarChargedReactionAttack, r.self, &ae)

		em := ae.Snapshot.Stats[attributes.EM]
		cr := ae.Snapshot.Stats[attributes.CR]
		cd := ae.Snapshot.Stats[attributes.CD]

		flatdmg := 1.8 * combat.CalcLunarChargedDmg(char.Base.Level, char, ae.Info, em)
		isCrit := false

		if r.core.Rand.Float64() <= cr {
			flatdmg *= (1 + cd)
			isCrit = true
		}

		contributions = append(contributions, lcContribution{flatdmg, isCrit, charInd, cr, cd, em})
	}

	if len(contributions) == 0 {
		return
	}

	slices.SortStableFunc(contributions, func(i, j lcContribution) int {
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

	for i, contr := range contributions {
		r.core.Combat.Log.NewEvent(fmt.Sprint("lunarcharged contributor ", (i+1)), glog.LogElementEvent, contr.charInd).
			Write("target", r.self.Key()).
			Write("damage", &contr.dmg).
			Write("crit", &contr.isCrit).
			Write("mult", lcContributorMult[i]).
			Write("cr", &contr.cr).
			Write("cd", &contr.cd).
			Write("em", &contr.em)

		ai.FlatDmg += contr.dmg * lcContributorMult[i]
	}

	// TODO: Make lunarcharged attack count as all contributor's attacks
	ai.ActorIndex = contributions[0].charInd
	snap := info.Snapshot{}
	if contributions[0].isCrit {
		snap.Stats[attributes.CR] = 1.0
	}
	r.core.QueueAttackWithSnap(
		ai,
		snap,
		ap,
		0,
		r.reduceLCAuraCB,
	)
}

func (r *Reactable) reduceLCAuraCB(a info.AttackCB) {
	var existing []string
	if r.core.Flags.LogDebug {
		existing = r.ActiveAuraString()
	}

	r.reduceMod(info.ReactionModKeyElectro, 10)
	r.reduceMod(info.ReactionModKeyHydro, 10)

	if r.core.Flags.LogDebug {
		r.core.Log.NewEvent(
			"application",
			glog.LogElementEvent,
			a.AttackEvent.Info.ActorIndex,
		).
			Write("attack_tag", a.AttackEvent.Info.AttackTag).
			Write("applied_ele", string(info.ReactionTypeLunarCharged)).
			Write("abil", a.AttackEvent.Info.Abil).
			Write("target", r.self.Key()).
			Write("existing", existing).
			Write("after", r.ActiveAuraString())
	}
}

func (r *Reactable) nextLCTick(src int) func() {
	return func() {
		if r.core.Flags.Custom[lcSrcKey] != float64(src) {
			// tick src changed
			return
		}

		if r.core.Status.Duration(LcKey) == 0 {
			// lunarcharge cloud expired
			return
		}

		// check all enemies in area
		// TODO: should be save the area when the cloud is generated or update when the player moves?
		for _, e := range r.core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(r.core.Combat.Player(), nil, 8), nil) {
			enemy, ok := e.(info.Reactable)
			if !ok {
				continue
			}
			if e.StatusIsActive(lcIcdKey) {
				continue
			}

			if enemy.GetAuraDurability(info.ReactionModKeyElectro) <= info.ZeroDur || enemy.GetAuraDurability(info.ReactionModKeyHydro) <= info.ZeroDur {
				continue
			}
			e.AddStatus(lcIcdKey, 2*60, true)
			enemy.DoLCAttack()
		}
		// queue up next tick
		r.core.Tasks.Add(r.nextLCTick(src), 6)
	}
}
