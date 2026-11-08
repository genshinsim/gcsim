package sandrone

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1MaxStacks = 10
	a4Key       = "sandrone-a4"
)

// Radiance: Stellar Glimmer: When using the Elemental Skill Formule Phenomenale: Differential
// Analysis, if Fagio's Decoding Power is greater than 50, the second Prism Shot fired by Fagio at
// nearby opponents will deal 400% of its original DMG as Fagio's Decoding Power is decreased.
// Additionally, for every 10 points of Decoding Power Fagio loses, it gains 1 stack of Refined
// Tactics for 60s (max 10 stacks). When Sandrone uses the Elemental Burst Formule Phenomenale:
// Q.E.D. while in the Radiance: Stellar Glimmer state, all stacks of Refined Tactics are cleared,
// causing the Convective Inhibition Ray to deal 100% + (number of stacks cleared) x 10% of its
// original DMG.
func (c *char) a1OnSkillBonusMult() float64 {
	if c.Base.Ascension < 1 {
		return 0.0
	}

	if c.decode <= 50 {
		return 0.0
	}

	if c.getRadiance() == radianceNone {
		return 0.0
	}

	return 3.0
}

func (c *char) a1OnDecreaseDecode(amt int) {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1DecreasedDecode += amt

	// TODO: RefinedTactics lasts 60s
	c.a1RefinedTactics += int(c.a1DecreasedDecode / 100)
	if c.a1RefinedTactics > a1MaxStacks {
		c.a1RefinedTactics = a1MaxStacks
	}
	c.a1DecreasedDecode = c.a1DecreasedDecode % 100
}

func (c *char) a1OnBurstRayStellar() float64 {
	if c.Base.Ascension < 1 {
		return 0.0
	}

	if c.getRadiance() == radianceNone {
		return 0.0
	}
	bonusMult := float64(c.a1RefinedTactics) * 0.1
	c.a1RefinedTactics = 0
	return bonusMult
}

// Increases Sandrone's Elemental Mastery based on her ATK. Every 100 ATK will increase her
// Elemental Mastery by 8. The maximum increase she can gain this way is 160.
func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(a4Key, -1),
		Extra:        true,
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
			m[attributes.EM] = min(stats.TotalATK()*0.08, 160)
			return m
		},
	})
}
