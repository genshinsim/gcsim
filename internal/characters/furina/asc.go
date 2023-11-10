package furina

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a4BuffKey = "furina-a4"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		hi := args[0].(*player.HealInfo)
		target := hi.Target
		overheal := args[3].(float64)

		if hi.Caller == c.Index {
			return false
		}

		if overheal <= 0 {
			return false
		}

		for _, char := range c.Core.Player.Chars() {
			if char.Index == target {
				continue
			}

			c.a1HealsStopFrameMap[char.Index] = c.Core.F + 240

			if !c.a1HealsFlagMap[char.Index] {
				c.Core.Tasks.Add(c.a1HealingOverTime(char.Index), 120)
				c.a1HealsFlagMap[char.Index] = true
			}
		}

		return false
	}, "furina-a1")
}

func (c *char) a1HealingOverTime(target int) func() {
	return func() {
		if c.a1HealsStopFrameMap[target] <= c.Core.F {
			c.a1HealsFlagMap[target] = false
			return
		}

		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  target,
			Type:    player.HealTypePercent,
			Message: "Endless Waltz",
			Src:     0.02,
			Bonus:   c.Stat(attributes.Heal),
		})

		c.Core.Tasks.Add(c.a1HealingOverTime(target), 120)
	}
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a4BuffKey, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil, false
			}

			return c.a4Buff, true
		},
	})
}

func (c *char) a4Tick() {
	if c.Base.Ascension < 4 {
		return
	}

	maxHP := c.MaxHP()
	var dmgBuff = maxHP / 1000 * 0.007

	if dmgBuff > 0.28 {
		dmgBuff = 0.28
	}

	c.a4Buff[attributes.DmgP] = dmgBuff

	c.Core.Tasks.Add(c.a4Tick, 30)
}
