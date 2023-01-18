package faruzan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C4: The vortex created by Wind Realm of Nasamjnin will restore Energy to
// Faruzan based on the number of opponents hit: If it hits 1 opponent, it
// will restore 2 Energy for Faruzan. Each additional opponent hit will
// restore 0.5 more Energy for Faruzan.
// A maximum of 4 Energy can be restored to her per vortex.
func (c *char) makeC4Callback() func(combat.AttackCB) {
	if c.Base.Cons < 4 {
		return nil
	}
	count := 0
	return func(a combat.AttackCB) {
		if count > 4 {
			return
		}
		amt := 0.5
		if count == 0 {
			amt = 2
		}
		count++
		c.AddEnergy("faruzan-c4", amt)
	}
}

// C6: Characters affected by The Wind's Secret Ways' Prayerful Wind's Gift
// have 40% bonus CRIT DMG when they deal Anemo DMG. When your own active
// character deals DMG while affected by Prayerful Wind's Gift, they will fire
// another Hurricane Arrow at opponents. This effect can be triggered once
// every 2.5s.
func (c *char) c6Buff(char *character.CharWrapper) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CD] = 0.4
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

func (c *char) c6Collapse() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if dmg := args[2].(float64); dmg == 0 {
			return false
		}
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
		c.AddStatus(c6ICDKey, 180, false)
		enemy := args[0].(*enemy.Enemy)
		c.pressurizedCollapse(enemy.Pos())
		return false
	}, "faruzan-c6-hook")
}
