package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	stellarConductText     = " (Stellar-Conduct)"
	radianceSwirlKey       = "radiance-stellar-swirl"
	dreamDrifterBuffKey    = "dreamdrifter-buff"
	dreamDrifterBuffICDKey = "dreamdrifter-buff-icd"
)

// type radianceState int

// const (
// 	radianceNone radianceState = iota
// 	radianceStellarConduct
// 	radianceStellarSwirl
// )

// func (c *char) getRadiance() radianceState {
// 	if !c.revelation {
// 		return radianceNone
// 	}

// 	if c.StatusIsActive(reactable.PolestarFieldKey) {
// 		return radianceStellarConduct
// 	}

// 	if c.StatusIsActive(radianceSwirlKey) {
// 		return radianceStellarSwirl
// 	}

// 	return radianceNone
// }

func (c *char) revelationInit() {
	if !c.revelation {
		return
	}

	// empty := make([]float64, attributes.EndStatType)

	onSkillHit := func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagElementalArt {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if !c.StatusIsActive(dreamDrifterBuffKey) {
			return
		}
		c.QueueCharTask(func() {
			c.DeleteStatus(dreamDrifterBuffKey)
		}, 0.25*60)

		atk.Info.FlatDmg += c.Stat(attributes.EM) * 10

		if !c.StatusIsActive(radianceSwirlKey) {
			return
		}

		t := args[0].(info.Target)

		c.stellarHit(t)
	}

	onSwirl := func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if c.StatusIsActive(dreamDrifterBuffICDKey) {
			return
		}
		c.AddStatus(dreamDrifterBuffICDKey, 2.5*60, true)
		c.QueueCharTask(func() {
			c.AddStatus(dreamDrifterBuffKey, -1, false)
			c.Core.Log.NewEvent("Dream Drifter state enhanced", glog.LogCharacterEvent, c.Index())
		}, 0.25*60)
	}

	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		c.AddStatus(radianceSwirlKey, 8*60, false)

		onSwirl(args...)
	}, "mizuki-ssw")

	c.Core.Events.Subscribe(event.OnSwirlCryo, onSwirl, "mizuki-swirl")
	c.Core.Events.Subscribe(event.OnSwirlPyro, onSwirl, "mizuki-swirl")
	c.Core.Events.Subscribe(event.OnSwirlHydro, onSwirl, "mizuki-swirl")
	c.Core.Events.Subscribe(event.OnSwirlElectro, onSwirl, "mizuki-swirl")

	c.Core.Events.Subscribe(event.OnEnemyHit, onSkillHit, "mizuki-swirl")
}

func (c *char) stellarHit(t info.Target) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Dreamdrifter Stellar Swirl",
		AttackTag:        attacks.AttackTagDirectStellarSwirl,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		PoiseDMG:         skillActivatePoise,
		Element:          attributes.Anemo,
		Durability:       0,
		UseEM:            true,
		Mult:             10,
		IgnoreDefPercent: 1,
		HitlagFactor:     0.05,
	}

	ap := combat.NewSingleTargetHit(t.Key())

	c.Core.QueueAttack(ai, ap, 0.2*60, 0.2*60) // datamine says this should be 12 frames while practice says the difference in damage number rendering is 20 frames ¯\_(ツ)_/¯
}
