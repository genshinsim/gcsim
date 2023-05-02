package baizhu

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Baizhu gains different effects according to the current HP of your current active character:
// ·When their HP is less than 50%, Baizhu gains 20% Healing Bonus.
// ·When their HP is equal to or more than 50%, Baizhu gains 25% Dendro DMG Bonus.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	//Healing part
	mHeal := make([]float64, attributes.EndStatType)
	mHeal[attributes.Heal] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("baizhu-a1-heal-bonus", -1),
		AffectedStat: attributes.Heal,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.ActiveChar().HPCurrent/c.Core.Player.ActiveChar().MaxHP() < 0.5 {
				return mHeal, true
			}
			return nil, false
		},
	})

	//Dendro DMG part
	mDendroP := make([]float64, attributes.EndStatType)
	mDendroP[attributes.DendroP] = 0.25
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("baizhu-a1-dendro-dmg", -1),
		AffectedStat: attributes.DendroP,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.ActiveChar().HPCurrent/c.Core.Player.ActiveChar().MaxHP() >= 0.5 {
				return mDendroP, true
			}
			return nil, false
		},
	})

}

// Characters who are healed by Seamless Shields will gain the Year of Verdant Favor effect:
// Each 1,000 Max HP that Baizhu possesses that does not exceed 50,000 will increase the Burning, Bloom, Hyperbloom, and Burgeon reaction
// DMG dealt by these characters by 2%, while the Aggravate and Spread reaction DMG dealt by these characters will be increased by 0.8%.
//
//	This effect lasts 6s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Player.ActiveChar().AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("baizhu-a4", 6*60),
		Amount: func(ai combat.AttackInfo) (float64, bool) {
			limitHP := c.MaxHP() / 1000.0
			if limitHP > 50 {
				limitHP = 50
			}
			if ai.Catalyzed && (ai.CatalyzedType == reactions.Aggravate || ai.CatalyzedType == reactions.Spread) {
				return limitHP * 0.008, false
			}
			switch ai.AttackTag {
			case attacks.AttackTagBloom:
			case attacks.AttackTagHyperbloom:
			case attacks.AttackTagBurgeon:
			default:
				return 0, false
			}

			return limitHP * 0.02, false
		},
	})

}
