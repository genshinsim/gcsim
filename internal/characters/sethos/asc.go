package sethos

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const a1Key = "sethos-a1"

// returns the amount of time to save, and the amount of energy considered
func (c *char) a1Calc() (int, float64) {
	if c.Base.Ascension < 1 {
		return 0, 0
	}
	energy := min(c.Energy, 20)
	// floor or round the skip?
	return int(0.285 * energy * 60), energy
}

func (c *char) a1Consume(energy float64, holdLevel int) {
	switch holdLevel {
	default:
		return
	case attacks.AimParamLv1:
		c.AddEnergy(a1Key, -energy*0.5)
	case attacks.AimParamLv2:
		c.AddEnergy(a1Key, -energy)
	}
	c.c2AddStack()
}

const a4Key = "sethos-a4"
const a4IcdKey = "sethos-a4-icd"
const a4Icd = 15 * 60

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	// buff stays active until a4 is proc'd
	c.AddStatus(a4Key, 9999999, true)
	c.a4Count = 0
}

func (c *char) makeA4cb() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}

	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		if !c.StatusIsActive(a4Key) {
			return
		}
		done = true
		if c.a4Count == 0 {
			// overwrite the expiry of the a4 buff to be 5s after
			c.AddStatus(a4Key, 5*60, true)
			c.startA4Icd()
		}
		c.a4Count += 1
		if c.a4Count >= 4 {
			c.DeleteStatus(a4Key)
		}
	}
}

func (c *char) startA4Icd() {
	if c.StatusIsActive(a4IcdKey) {
		return
	}
	c.AddStatus(a4IcdKey, a4Icd, true)
	c.QueueCharTask(c.a4, a4Icd)
}
