package furina

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const c4Key = "furina-c4"
const c4IcdKey = "furina-c4-icd"

const c6Key = "center-of-attention"
const c6IcdKey = "furina-c6-icd"
const c6OusiaHealKey = "furina-c6-ousia-heal"

func (c *char) c4cb(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(c4IcdKey) {
		return
	}
	c.AddStatus(c4IcdKey, 5*60, true)
	c.AddEnergy(c4Key, 4)
}

func (c *char) c6BonusDMG() float64 {
	scaleHP := 0.18
	if c.arkhe == pneuma {
		scaleHP += 0.25
	}
	return scaleHP * c.MaxHP()
}

func (c *char) c6cb(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}

	c.AddStatus(c6IcdKey, 0.1*60, true)

	switch c.arkhe {
	case ousia:
		if !c.StatusIsActive(c6OusiaHealKey) {
			c.QueueCharTask(c.c6heal(c.Core.F), 60)
		}
		c.AddStatus(c6OusiaHealKey, 2.9*60, true)
	case pneuma:
		for _, char := range c.Core.Player.Chars() {
			hpDrain := char.CurrentHP() * 0.01
			c.Core.Player.Drain(player.DrainInfo{
				ActorIndex: char.Index,
				Abil:       "Furina C6 Pneuma Drain",
				Amount:     hpDrain,
				External:   false,
			})
		}
	}
	c.c6Count += 1
	if c.c6Count == 6 {
		c.DeleteStatus(c6Key)
	}
}

func (c *char) c6heal(src int) func() {
	return func() {
		if c.c6HealSrc != src {
			return
		}
		if !c.StatusIsActive(c6OusiaHealKey) {
			return
		}
		heal := c.Stat(attributes.Heal)
		amt := c.MaxHP() * 0.04
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Type:    player.HealTypePercent,
			Message: "Furina C6 Ousia Heal",
			Src:     amt,
			Bonus:   heal,
		})
		c.QueueCharTask(c.c6heal(src), 60)
	}
}
