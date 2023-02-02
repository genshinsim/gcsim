package beidou

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// A1 is not implemented:
// TODO: Counterattacking with Tidecaller at the precise moment when the character is hit grants the maximum DMG Bonus.

// Gain the following effects for 10s after unleashing Tidecaller with its maximum DMG Bonus:
// - DMG dealt by Normal and Charged Attacks is increased by 15%. ATK SPD of Normal and Charged Attacks is increased by 15%.
// TODO: - Greatly reduced delay before unleashing Charged Attacks.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = .15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("beidou-a4-dmg", 600),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
				return nil, false
			}
			return mDmg, true
		},
	})

	mAtkSpd := make([]float64, attributes.EndStatType)
	mAtkSpd[attributes.AtkSpd] = .15
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("beidou-a4-atkspd", 600),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			return mAtkSpd, true
		},
	})
}
