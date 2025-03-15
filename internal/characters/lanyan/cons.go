package lanyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

const c2Icd = "lanyan-c2-icd"

var c1Hitmarks = []int{37, 64, 90}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		if !c.hasShield() {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}

		if c.StatusIsActive(c2Icd) {
			return false
		}
		c.AddStatus(c2Icd, 2*60, true)

		c.restoreShield(0.4)
		return false
	}, "lanyan-c2")
}
