package ororon

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1Key = "ororon-c1"
const c2Key = "ororon-c2"
const c4Key = "ororon-c4"
const c6Key = "ororon-c6"

func (c *char) c1Init() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.5

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c1Key, -1),
		Amount: func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			trg, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if ae.Info.Abil != a1Abil {
				return nil, false
			}
			if !trg.StatusIsActive(c1Key) {
				return nil, false
			}
			return m, true
		},
	})
}

func (c *char) c1ExtraBounce() int {
	if c.Base.Cons < 1 {
		return 0
	}
	return 2
}

func (c *char) makeC1cb() func(combat.AttackCB) {
	if c.Base.Cons < 1 {
		return nil
	}
	if c.Base.Ascension < 1 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		trg, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		trg.AddStatus(c1Key, 12*60, true)
	}
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2Bonus = make([]float64, attributes.EndStatType)
}

func (c *char) c2OnBurst() {
	if c.Base.Cons < 2 {
		return
	}
	c.AddStatus(c2Key, 9*60, true)
	c.SetTag(c2Key, 1)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c2Key, 10*60),
		AffectedStat: attributes.ElectroP,
		Amount: func() ([]float64, bool) {
			c.c2Bonus[attributes.ElectroP] = min(0.08*float64(c.Tags[c2Key]), 0.32)
			return c.c2Bonus, true
		},
	})
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6bonus = make([]float64, attributes.EndStatType)
	c.c6stacks = make([]*stacks.MultipleRefreshNoRemove, len(c.Core.Player.Chars()))
	for i := range c.c6stacks {
		c.c6stacks[i] = stacks.NewMultipleRefreshNoRemove(3, c.QueueCharTask, &c.Core.F)
	}
}

func (c *char) makeC2cb() func(combat.AttackCB) {
	if c.Base.Cons < 2 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		c.Tags[c2Key]++
	}
}

func (c *char) c4BurstInterval() float64 {
	if c.Base.Cons < 4 {
		return 1.0
	}
	return 0.75
}

func (c *char) c4EnergyRestore() {
	if c.Base.Cons < 4 {
		return
	}
	c.AddEnergy(c4Key, 8)
}

func (c *char) c6OnBurst() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Base.Ascension < 1 {
		return
	}
	c.hypersense(3.2, "Hypersense (C6)", c.Core.Combat.PrimaryTarget().Pos())
}

func (c *char) c6onHypersense() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Base.Ascension < 1 {
		return
	}
	c.c6stacks[c.Core.Player.Active()].Add(9 * 60)

	c.Core.Player.ActiveChar().AddStatMod(character.StatMod{
		Base:   modifier.NewBaseWithHitlag(c6Key, 9*60),
		Amount: c.c6Amount(c.Core.Player.Active()),
	})
}

func (c *char) c6Amount(ind int) func() ([]float64, bool) {
	return func() ([]float64, bool) {
		c.c6bonus[attributes.ATKP] = float64(c.c6stacks[ind].Count()) * 0.1
		return c.c6bonus, true
	}
}
