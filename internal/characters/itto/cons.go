package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C1:
// After using Royal Descent: Behold, Itto the Evil!,
// 	Arataki Itto gains 2 stacks of Superlative Superstrength.
// After 1s, Itto will gain 1 stack of Superlative Superstrength every 0.5s for 1.5s.
// TODO: add link to itto-c1-mechanics tcl entry later
func (c *char) c1() func() {
	return func() {
		// gain 2 initial stacks around 75f after pressing Q
		c.changeStacks(2)
		c.Core.Log.NewEvent("itto-c1 initial stacks added", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.Tags[c.stackKey])
		// "After 1s" refers to 1s after gaining the initial 2 stacks, so queue up the stacks properly
		for i := 60; i <= 120; i += 30 {
			c.QueueCharTask(func() {
				c.changeStacks(1)
				c.Core.Log.NewEvent("itto-c1 later stack added", glog.LogCharacterEvent, c.Index).
					Write("stacks", c.Tags[c.stackKey])
			}, i)
		}
	}
}

// C2:
// After using Royal Descent: Behold, Itto the Evil!,
// 	each party member whose Element is Geo will decrease that skill's CD by 1.5s
// 	and restore 6 Energy to Arataki Itto.
// CD can be decreased by up to 4.5s in this manner.
// Max 18 Energy can be restored in this manner.
func (c *char) c2() func() {
	return func() {
		energyGain := float64(c.c1GeoMemberCount) * 6
		cdDecrease := c.c1GeoMemberCount * (1.5 * 60)
		c.AddEnergy("itto-c2", energyGain)
		c.ReduceActionCooldown(action.ActionBurst, cdDecrease)

		c.Core.Log.NewEvent("itto-c2 applied", glog.LogCharacterEvent, c.Index).
			Write("energy gained", energyGain).
			Write("new energy", c.Energy).
			Write("q cooldown decrease", cdDecrease)
	}
}

// C4:
// When the Raging Oni King state caused by Royal Descent: Behold, Itto the Evil! ends,
// 	all nearby party members gain 20% DEF and 20% ATK for 10s.
func (c *char) c4() func() {
	return func() {
		if c.c4Applied {
			return
		}
		c.c4Applied = true
		m := make([]float64, attributes.EndStatType)
		m[attributes.DEFP] = 0.2
		m[attributes.ATKP] = 0.2
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag("itto-c4", 600),
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
	}
}

// First Part of C6:
// Arataki Itto's Charged Attacks deal +70% Crit DMG.
func (c *char) c6ChargedCritDMG() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CD] = 0.7
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("itto-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagExtra {
				return nil, false
			}
			return m, true
		},
	})
}

// Second Part of C6:
// Additionally, when he uses Arataki Kesagiri, he has a 50% chance to not consume stacks of Superlative Superstrength.
func (c *char) c6StackHandler() {
	if c.Core.Rand.Float64() < 0.5 {
		c.changeStacks(-1)
		// only update if a stack was actually consumed
		c.stacksConsumed++
	} else {
		c.Core.Log.NewEvent("itto-c6 proc'd", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.Tags[c.stackKey])
	}
}
