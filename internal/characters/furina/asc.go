package furina

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1HealKey = "furina-a1"
	a4BuffKey = "furina-a4"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1HealsStopFrameMap = make([]int, len(c.Core.Player.Chars()))
	c.a1HealsFlagMap = make([]bool, len(c.Core.Player.Chars()))

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		hi := args[0].(*player.HealInfo)
		overheal := args[3].(float64)

		if hi.Caller == c.Index {
			return false
		}

		if overheal <= 0 {
			return false
		}

		if !c.StatusIsActive(a1HealKey) {
			c.Core.Tasks.Add(c.a1HealingOverTime(), 120)
		}

		c.AddStatus(a1HealKey, 240, false)

		return false
	}, "furina-a1")
}

func (c *char) a1HealingOverTime() func() {
	return func() {
		if !c.StatusIsActive(a1HealKey) {
			return
		}
		heal := c.Stat(attributes.Heal)
		for target := range c.Core.Player.Chars() {
			amt := c.Core.Player.Chars()[target].MaxHP() * 0.02
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  target,
				Type:    player.HealTypePercent,
				Message: "Endless Waltz",
				Src:     amt,
				Bonus:   heal,
			})
		}

		c.Core.Tasks.Add(c.a1HealingOverTime(), 120)
	}
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4Buff = make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a4BuffKey, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil, false
			}

			if !strings.Contains(atk.Info.Abil, salonMemberKey) {
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

	var dmgBuff = c.MaxHP() / 1000 * 0.007

	if dmgBuff > 0.28 {
		dmgBuff = 0.28
	}

	c.a4Buff[attributes.DmgP] = dmgBuff

	c.Core.Tasks.Add(c.a4Tick, 30)
}
