﻿package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// A1: When Arataki Itto uses consecutive Arataki Kesagiri, he obtains the following effects:
// - Each slash increases the ATK SPD of the next slash by 10%. Max ATK SPD increase is 30%.
// - Increases his resistance to interruption.
// These effects will be cleared once he stops performing consecutive slashes.
func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("itto-a1", -1),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			if c.a1Stacks == 0 || c.Core.Player.CurrentState() != action.ChargeAttackState {
				return nil, false
			}
			m[attributes.AtkSpd] = 0.10 * float64(c.a1Stacks)
			return m, true
		},
	})
}

func (c *char) a1Update(curSlash SlashType) {
	switch curSlash {
	case SaichiSlash:
		// reset a1 stacks if we are doing a CA0
		c.a1Stacks = 0
		c.Core.Log.NewEvent("itto-a1 reset atkspd stacks", glog.LogCharacterEvent, c.Index).
			Write("a1Stacks", c.a1Stacks).
			Write("slash", curSlash.String())
	case LeftSlash, RightSlash:
		// increment a1 stacks if we are doing CA1/CA2
		c.a1Stacks++
		if c.a1Stacks > 3 {
			c.a1Stacks = 3
		}
		c.Core.Log.NewEvent("itto-a1 atkspd stacks increased", glog.LogCharacterEvent, c.Index).
			Write("a1Stacks", c.a1Stacks).
			Write("slash", curSlash.String())
	}
}
