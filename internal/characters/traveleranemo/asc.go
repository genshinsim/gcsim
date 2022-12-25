package traveleranemo

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const a4ICDKey = "traveleranemo-a4-icd"

func (c *char) a4() {
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalArt {
			return false
		}
		if c.StatusIsActive(a4ICDKey) {
			return false
		}

		c.AddStatus(a4ICDKey, 300, true)

		for i := 0; i < 5; i = i + 1 {
			c.QueueCharTask(func() {
				c.Core.Player.Heal(player.HealInfo{
					Caller:  c.Index,
					Target:  c.Index,
					Message: "Second Wind",
					Type:    player.HealTypePercent,
					Src:     0.02,
				})
			}, (i+1)*60) // healing starts 1s after death
		}

		return false
	}, "traveleranemo-a4")
}
