package shenhe

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4BuffKey = "shenhe-c4"

func (c *char) c2(active *character.CharWrapper, dur int) {
	active.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("shenhe-c2", dur),
		Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if ae.Info.Element != attributes.Cryo {
				return nil, false
			}
			return c.c2buff, true
		},
	})
}

func (c *char) c4() {
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("shenhe-c4-dmg", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
				return nil, false
			}
			if !c.StatusIsActive(c4BuffKey) {
				c.c4count = 0
				return nil, false
			}
			c.c4bonus[attributes.DmgP] += 0.05 * float64(c.c4count)
			c.c4count = 0
			return c.c4bonus, true
		},
	})
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
			return false
		}
		c.DeleteStatus(c4BuffKey)
		return false
	}, "shenhe-c4-reset")
}
