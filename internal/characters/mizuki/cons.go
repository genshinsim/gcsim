package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	c1Key               = "mizuki-c1"
	c1Interval          = 3.5 * 60
	c1Duration          = 3 * 60
	c1Multiplier        = 11.0
	c2Key               = "mizuki-c2"
	c2EMMultiplier      = 0.0004
	c4EnergyGenerations = 4
	c4Key               = "mizuki-c4"
	c4Energy            = 5
	c6Key               = "mizuki-c6"
	c6CR                = 0.3
	c6CD                = 1.0
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		// Only when dream drifter is active
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return false
		}

		e := args[0].(*enemy.Enemy)
		if !e.StatusIsActive(c1Key) {
			return false
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagSwirlCryo:
		case attacks.AttackTagSwirlElectro:
		case attacks.AttackTagSwirlHydro:
		case attacks.AttackTagSwirlPyro:
		default:
			return false
		}

		additionalDmg := c1Multiplier * c.Stat(attributes.EM)

		c.Core.Log.NewEvent("mizuki c1 proc", glog.LogPreDamageMod, atk.Info.ActorIndex).
			Write("before", atk.Info.FlatDmg).
			Write("addition", additionalDmg).
			Write("final", atk.Info.FlatDmg+additionalDmg).
			Write("em", c.Stat(attributes.EM))

		atk.Info.FlatDmg += additionalDmg
		atk.Info.Abil += " (Mizuki C1)"

		e.DeleteStatus(c1Key)

		return false
	}, c1Key)
}

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		ae := args[1].(*combat.AttackEvent)

		switch ae.Info.AttackTag {
		case attacks.AttackTagSwirlPyro:
		case attacks.AttackTagSwirlCryo:
		case attacks.AttackTagSwirlHydro:
		case attacks.AttackTagSwirlElectro:
		default:
			return false
		}

		if !c.StatusIsActive(dreamDrifterStateKey) {
			return false
		}

		ae.Snapshot.Stats[attributes.CR] += c6CR
		ae.Snapshot.Stats[attributes.CD] += c6CD

		c.Core.Log.NewEvent("mizuki c6 buff", glog.LogCharacterEvent, ae.Info.ActorIndex).
			Write("final_crit", ae.Snapshot.Stats[attributes.CR])

		return false
	}, c6Key)
}
