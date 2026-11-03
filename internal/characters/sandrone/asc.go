package sandrone

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	refinedTacticsKey = "sandrone-refined-tactics"
	refinedTacticsDur = 60 * 60
	refinedTacticsMax = 10
	a1PrismThreshold  = 50
	a4Key             = "sandrone-a4"
)

// Radiance: Stellar Glimmer: When using the Elemental Skill, if Fagio's Decoding Power is greater
// than 50, the second Prism Shot fired by Fagio at nearby opponents will deal 400% of its original
// DMG as Fagio's Decoding Power is decreased.
func (c *char) a1PrismActive() bool {
	if c.Base.Ascension < 1 {
		return false
	}
	return c.decodingPower > a1PrismThreshold
}

// For every 10 points of Decoding Power Fagio loses, it gains 1 stack of Refined Tactics for 60s
// (max 10 stacks).
func (c *char) a1AddTactics() {
	if c.Base.Ascension < 1 {
		return
	}
	c.tactics = min(c.refinedTactics()+1, refinedTacticsMax)
	c.AddStatus(refinedTacticsKey, refinedTacticsDur, true)
	c.Core.Log.NewEvent("refined tactics stack gained", glog.LogCharacterEvent, c.Index()).
		Write("stacks", c.tactics)
}

func (c *char) refinedTactics() int {
	if !c.StatusIsActive(refinedTacticsKey) {
		return 0
	}
	return c.tactics
}

// When Sandrone uses the Elemental Burst while in the Radiance: Stellar Glimmer state, all stacks
// of Refined Tactics are cleared, causing the Convective Inhibition Ray to deal 100% + (number of
// stacks cleared) x 10% of its original DMG.
func (c *char) a1ConsumeTactics() float64 {
	if c.Base.Ascension < 1 {
		return 1
	}
	stacks := c.refinedTactics()
	c.tactics = 0
	c.DeleteStatus(refinedTacticsKey)
	if stacks > 0 {
		c.Core.Log.NewEvent("refined tactics consumed", glog.LogCharacterEvent, c.Index()).
			Write("stacks", stacks)
	}
	return 1 + float64(stacks)*a1TacticsMult
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
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			m[attributes.EM] = min(c.TotalAtk()*a4EMPerATK, a4EMMax)
			return m
		},
	})
}
