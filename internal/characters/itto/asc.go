package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Arataki Itto uses consecutive Arataki Kesagiri, he obtains the following effects:
//
// - Each slash increases the ATK SPD of the next slash by 10%. Max ATK SPD increase is 30%.
//
// TODO: - Increases his resistance to interruption.
//
// These effects will be cleared once he stops performing consecutive slashes.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	mAtkSpd := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("itto-a1", -1),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			if c.a1Stacks == 0 || c.Core.Player.CurrentState() != action.ChargeAttackState {
				return nil, false
			}
			mAtkSpd[attributes.AtkSpd] = 0.10 * float64(c.a1Stacks)
			return mAtkSpd, true
		},
	})
}

func (c *char) a1Update(curSlash SlashType) {
	if c.Base.Ascension < 1 {
		return
	}
	switch curSlash {
	case SaichiSlash:
		// reset a1 stacks if we are doing a CA0
		c.a1Stacks = 0
		c.Core.Log.NewEvent("itto-a1 reset atkspd stacks", glog.LogCharacterEvent, c.Index).
			Write("a1Stacks", c.a1Stacks).
			Write("slash", curSlash.String())
	case LeftSlash, RightSlash:
		// increment a1 stacks if we are doing CA1/CA2
		// increment stacks for A1, max is 3 stacks
		c.a1Stacks++
		if c.a1Stacks > 3 {
			c.a1Stacks = 3
		}
		c.Core.Log.NewEvent("itto-a1 atkspd stacks increased", glog.LogCharacterEvent, c.Index).
			Write("a1Stacks", c.a1Stacks).
			Write("slash", curSlash.String())
	}
	// do nothing if we are doing a CAF
}

// Arataki Kesagiri DMG is increased by 35% of Arataki Itto's DEF.
func (c *char) a4(ai *combat.AttackInfo) {
	if c.Base.Ascension < 4 {
		return
	}
	ai.FlatDmg = (c.Base.Def*(1+c.Stat(attributes.DEFP)) + c.Stat(attributes.DEF)) * 0.35
	c.Core.Log.NewEvent("itto-a4 applied", glog.LogCharacterEvent, c.Index)
}
