package iansan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Key    = "precise-movement"
	a1ICDKey = "precise-movement-icd"
	a4Key    = "warming-up"
	a4ICDKey = "warming-up-icd"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1Buff = make([]float64, attributes.EndStatType)
	c.a1Buff[attributes.ATKP] = 0.2

	cb := func(args ...any) {
		idx := args[0].(int)
		if idx != c.Core.Player.Active() {
			return
		}

		if !c.StatModIsActive(a1Key) {
			return
		}

		if c.StatusIsActive(a1ICDKey) {
			return
		}

		c.AddStatus(a1ICDKey, 2.8*60, true)
		c.a1Increase = true
	}

	c.Core.Events.Subscribe(event.OnNightsoulGenerate, cb, "iansan-a1-generate")
	c.Core.Events.Subscribe(event.OnNightsoulConsume, cb, "iansan-a1-consume")
}

func (c *char) a1AddBuff() {
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag(a1Key, 15*60),
		Amount: func() []float64 {
			return c.a1Buff
		},
	})
}

func (c *char) makeA1CB() func(_ info.AttackCB) {
	if c.Base.Ascension < 1 {
		return nil
	}

	done := false
	return func(a info.AttackCB) {
		if done {
			return
		}
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		c.a1AddBuff()
		done = true
	}
}

func (c *char) a1Points() float64 {
	if c.Base.Ascension < 1 {
		return 0.0
	}
	if !c.StatModIsActive(a1Key) {
		return 0.0
	}
	if c.a1Increase {
		c.a1Increase = false
		return 4.0
	}
	return 1.0
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) {
		c.AddStatus(a4Key, 10*60, true)
	}, "iansan-a4")
}

func (c *char) a4Heal() {
	if !c.StatusIsActive(a4Key) {
		return
	}
	if c.StatusIsActive(a4ICDKey) {
		return
	}

	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  c.Core.Player.Active(),
		Message: "Warming Up",
		Src:     c.TotalAtk() * 0.6,
		Bonus:   c.Stat(attributes.Heal),
	})

	c.AddStatus(a4ICDKey, 2.8*60, true)
}
