package iansan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Status = "precise-movement"
	a1ICD    = "precise-movement-icd"
	a4Status = "warming-up"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	cb := func(args ...interface{}) {
		idx := args[0].(int)
		if idx != c.Core.Player.Active() {
			return
		}
		if !c.StatModIsActive(a1Status) {
			return
		}
		if c.StatusIsActive(a1ICD) {
			return
		}
		c.AddStatus(a1ICD, 2.8*60, true)
		c.a1Increase = true
	}
	c.Core.Events.Subscribe(event.OnNightsoulGenerate, cb, "iansan-a1-generate")
	c.Core.Events.Subscribe(event.OnNightsoulConsume, cb, "iansan-a1-consume")
}

func (c *char) a1ATK() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.2
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag(a1Status, 15*60),
		Amount: func() []float64 {
			return m
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
		c.a1ATK()
		done = true
	}
}

func (c *char) a1Points() float64 {
	if c.Base.Ascension < 1 {
		return 0.0
	}
	if !c.StatModIsActive(a1Status) {
		return 0.0
	}
	if c.a1Increase {
		c.a1Increase = false
		return 4.0
	}
	return 1.0
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) {
		c.AddStatus(a4Status, 10*60, true)
		c.a4Src = c.Core.F
		c.a4Task(c.a4Src)
	}, "iansan-a4")
}

func (c *char) a4Task(src int) {
	c.QueueCharTask(func() {
		if c.a4Src != src {
			return
		}

		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index(),
			Target:  c.Core.Player.Active(),
			Message: "Warming Up",
			Src:     c.TotalAtk() * 0.6,
			Bonus:   c.Stat(attributes.Heal),
		})

		c.nightsoulState.GeneratePoints(1)
		c.a4Task(src)
	}, 2.8*60)
}
