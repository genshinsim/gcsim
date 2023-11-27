package furina

import (
	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2BuffKey = "furina-c2-hp"

const c4Key = "furina-c4"
const c4IcdKey = "furina-c4-icd"

const c6Key = "center-of-attention"
const c6IcdKey = "furina-c6-icd"
const c6OusiaHealKey = "furina-c6-ousia-heal"

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c2BuffKey, -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			m[attributes.HPP] = common.Max(c.curFanfare-c.maxQFanfare, 0) * 0.0035
			return m, true
		},
	})
}

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

func (c *char) c6BonusDMGArkhe() float64 {
	scaleHP := 0.18
	return scaleHP * c.MaxHP()
}

func (c *char) c6cb(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}

	if c.StatusIsActive(c6IcdKey) {
		return
	}

	c.AddStatus(c6IcdKey, 0.1*60, true)

	switch c.arkhe {
	case ousia:
		if !c.StatusIsActive(c6OusiaHealKey) {
			c.c6HealSrc = c.Core.F
			for _, char := range c.Core.Player.Chars() {
				char.QueueCharTask(c.c6heal(char, c.Core.F), 60)
				char.AddStatus(c6OusiaHealKey, 2.9*60, true)
			}
		} else {
			for _, char := range c.Core.Player.Chars() {
				char.ExtendStatus(c6OusiaHealKey, 2.9*60)
			}
		}

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

func (c *char) c6heal(char *character.CharWrapper, src int) func() {
	return func() {
		if c.c6HealSrc != src {
			return
		}
		if !c.StatusIsActive(c6OusiaHealKey) {
			return
		}
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  char.Index,
			Type:    player.HealTypeAbsolute,
			Message: "Furina C6 Ousia Heal",
			Src:     0.04 * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
		char.QueueCharTask(c.c6heal(char, src), 60)
	}
}
