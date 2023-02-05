package xiao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "xiao-a1"

// While under the effects of Bane of All Evil, all DMG dealt by Xiao increases
// by 5%. DMG increases by a further 5% for every 3s the ability persists. The
// maximum DMG Bonus is 25%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a1Key, 900+burstStart),
		AffectedStat: attributes.DmgP,
		Amount: func() ([]float64, bool) {
			stacks := 1 + int((c.Core.F-c.qStarted)/180)
			if stacks > 5 {
				stacks = 5
			}
			m[attributes.DmgP] = float64(stacks) * 0.05
			return m, true
		},
	})
}

// Using Lemniscatic Wind Cycling increases the DMG of subsequent uses of Lemniscatic Wind Cycling by 15%.
// This effect lasts for 7s, and has a maximum of 3 stacks. Gaining a new stack refreshes the effect's duration.
//
// - checks for ascension level in skill.go to avoid queuing this up only to fail the ascension level check
func (c *char) a4() {
	// reset stacks if buff expired
	if !c.StatModIsActive(a4BuffKey) {
		c.a4stacks = 0
	}
	// Text is not explicit, but assume that gaining a stack while at max still refreshes duration
	c.a4stacks++
	if c.a4stacks > 3 {
		c.a4stacks = 3
	}
	c.a4buff[attributes.DmgP] = float64(c.a4stacks) * 0.15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(a4BuffKey, 420),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			return c.a4buff, atk.Info.AttackTag == combat.AttackTagElementalArt
		},
	})
}
