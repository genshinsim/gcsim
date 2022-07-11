package beidou

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a4() {
	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = .15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("beidou-a4", 600),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
				return nil, false
			}
			return mDmg, true
		},
	})

	mAtkSpd := make([]float64, attributes.EndStatType)
	mAtkSpd[attributes.AtkSpd] = .15
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("beidou-a4", 600),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			return mAtkSpd, true
		},
	})
}
