package furina

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
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

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		hi := args[0].(*info.HealInfo)
		overheal := args[3].(float64)

		if hi.Caller == c.Index {
			return false
		}

		if overheal <= 0 {
			return false
		}

		if hi.Target != c.Core.Player.Active() && hi.Target != -1 {
			return false
		}

		if !c.StatusIsActive(a1HealKey) {
			c.a1src = c.Core.F
			c.QueueCharTask(c.a1HealingOverTime(c.a1src), 2*60)
		}

		c.AddStatus(a1HealKey, 4*60, true)

		return false
	}, "furina-a1")
}

func (c *char) a1HealingOverTime(src int) func() {
	return func() {
		if !c.StatusIsActive(a1HealKey) {
			return
		}
		if c.a1src != src {
			return
		}
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Type:    info.HealTypePercent,
			Message: "Endless Waltz",
			Src:     0.02,
			Bonus:   c.Stat(attributes.Heal),
		})

		c.QueueCharTask(c.a1HealingOverTime(src), 2*60)
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

	c.a4Buff[attributes.DmgP] = min(c.MaxHP()/1000*0.007, 0.28)
	c.a4IntervalReduction = min(c.MaxHP()/1000.0*0.004, 0.16)

	c.QueueCharTask(c.a4Tick, 30)
}
