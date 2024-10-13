package sethos

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2Key = "sethos-c2"
const c2Dur = 10 * 60

const c4Key = "sethos-c4"
const c4Dur = 10 * 60

const c6Key = "sethos-c6"
const c6IcdKey = "sethos-c6-icd"
const c6IcdDur = 15 * 60

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("sethos-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			if atk.Info.Abil != shadowPierceShotAil {
				return nil, false
			}
			return m, true
		},
	})
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	mElectro := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBase(c2Key, -1),
		Amount: func() ([]float64, bool) {
			stackCount := min(c.c2Stacks, 2.0)
			if stackCount == 0 {
				return nil, false
			}
			mElectro[attributes.ElectroP] = 0.15 * float64(stackCount)
			return mElectro, true
		},
	})
}

func (c *char) c2AddStack() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2Stacks += 1
	c.SetTag(c2Key, min(c.c2Stacks, 2))
	c.QueueCharTask(func() {
		// tags currently aren't visible in the results UI
		// the user can still access it using .char.tags.sethos-c2
		c.c2Stacks -= 1
		c.SetTag(c2Key, min(c.c2Stacks, 2))
	}, c2Dur)
}

var c4Buff []float64

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c4Buff = make([]float64, attributes.EndStatType)
	c4Buff[attributes.EM] = 80
}

func (c *char) makeC4cb() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	count := 0
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if count >= 2 {
			return
		}
		count += 1
		if count == 2 {
			for _, char := range c.Core.Player.Chars() {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(c4Key, c4Dur),
					AffectedStat: attributes.EM,
					Amount: func() ([]float64, bool) {
						return c4Buff, true
					},
				})
			}
		}
	}
}

func (c *char) makeC6cb(energy float64) combat.AttackCBFunc {
	if c.Base.Cons < 6 {
		return nil
	}
	if c.Base.Ascension < 1 {
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
		if c.StatusIsActive(c6IcdKey) {
			return
		}
		done = true
		c.AddStatus(c6IcdKey, c6IcdDur, true)
		c.AddEnergy(c6Key, energy)
	}
}
