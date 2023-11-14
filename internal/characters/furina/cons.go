package furina

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const c4Key = "furina-c4"
const c4IcdKey = "furina-c4-icd"

const c6Key = "furina-c6"

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
	return scaleHP
}
