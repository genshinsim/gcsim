package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C1:
// After using Royal Descent: Behold, Itto the Evil!, Arataki Itto gains 2 stacks of Superlative Superstrength.
// After 1s, Itto will gain 1 stack of Superlative Superstrength every 0.5s for 1.5s.
// TODO: add link to itto-c1-mechanics tcl entry later
func (c *char) c1() {
	// gain 2 initial stacks around 75f after pressing Q
	c.addStrStack(2)
	// "After 1s" refers to 1s after gaining the initial 2 stacks, so queue up the stacks properly
	for i := 60; i <= 120; i += 30 {
		c.QueueCharTask(func() { c.addStrStack(1) }, i)
	}
}

// copied from raiden c4
// When the Raging Oni King state caused by Royal Descent: Behold, Itto the Evil! ends,
// all nearby party members gain 20% DEF and 20% ATK for 10s.
func (c *char) c4() {
	if !c.applyC4 {
		return
	}
	c.applyC4 = false

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.2
	m[attributes.DEFP] = 0.2
	for _, x := range c.Core.Player.Chars() {
		x.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("itto-c4", 10*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CD] = 0.7
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("itto-c6", -1),
		Amount: func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			return m, a.Info.AttackTag == combat.AttackTagExtra
		},
	})
}
