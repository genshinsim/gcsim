package faruzan

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C6: Characters affected by The Wind's Secret Ways' Prayerful Wind's Gift
// have 40% bonus CRIT DMG when they deal Anemo DMG. When your own active
// character deals DMG while affected by Prayerful Wind's Gift, they will fire
// another Hurricane Arrow at opponents. This effect can be triggered once
// every 2.5s.
func (c *char) c6Buff(char *character.CharWrapper) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CD] = burstBuff[c.TalentLvlBurst()]
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("faruzan-c6", 240),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.Element != attributes.Anemo {
				return nil, false
			}
			return m, true
		},
	})
}

const c6ICDKey = "faruzan-c6-icd"

func (c *char) c6Arrows() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		char := c.Core.Player.ActiveChar()
		if char.Index != atk.Info.ActorIndex {
			return false
		}
		if !char.StatusIsActive(burstBuffKey) {
			return false
		}
		if c.StatusIsActive(c6ICDKey) {
			return false
		}
		c.AddStatus(c6ICDKey, 150, false)
		c.hurricaneArrow(10, false)
		return false
	}, "faruzan-c6-hook")
}
