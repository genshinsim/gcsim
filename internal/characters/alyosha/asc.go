package alyosha

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a4Key = "alyosha-a4"

// When Tugarin attacks, it will also restore HP to nearby active characters at 120% of Alyosha's ATK.
func (c *char) a1OnTugarin() {
	if c.Base.Ascension < 1 {
		return
	}

	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  c.Core.Player.Active(),
		Message: "Alyosha (A1)",
		Src:     1.2 * c.TotalAtk(),
		Bonus:   c.Stat(attributes.Heal),
	})
}

// Increases the DMG Alyosha deals with his Elemental Skill and Elemental Burst by 0.35% for every
// 1% of his Energy Recharge. Up to a 70% increase can be obtained in this way.
func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a4Key, -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt:
			case attacks.AttackTagElementalArtHold:
			case attacks.AttackTagElementalBurst:
			default:
				return nil
			}
			// This doesn't use NonExtra ER
			m[attributes.DmgP] = min(c.Stat(attributes.ER)/0.01*0.0035, 0.7)
			return m
		},
	})
}
