package faruzan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	a4Key    = "faruzan-a4"
	a4ICDKey = "faruzan-a4-icd"
)

func (c *char) a4() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.Element != attributes.Anemo {
			return false
		}

		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if char.StatusIsActive(burstBuffKey) && !char.StatusIsActive(a4ICDKey) {
			char.AddStatus(a4Key, 6, true)
			char.AddStatus(a4ICDKey, 60, true)
		}

		if char.StatusIsActive(a4Key) {
			stats, _ := c.Stats()
			amt := 0.574 * ((c.Base.Atk+c.Weapon.Atk)*(1+stats[attributes.ATKP]) + stats[attributes.ATK])
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("faruzan a4 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt)
			}
			atk.Info.FlatDmg += amt
		}

		return false
	}, "faruzan-a4-hook")
}
