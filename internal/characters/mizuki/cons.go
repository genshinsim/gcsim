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
	c1Key        = "mizuki-c1"
	c1Interval   = 3.5 * 60
	c1Duration   = 3 * 60
	c1Multiplier = 11.0
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

		e.DeleteStatus(c1Key)

		return false
	}, c1Key)
}
