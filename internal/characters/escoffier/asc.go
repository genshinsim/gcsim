package escoffier

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Scaling   = 1.3824
	a1key       = "escoffier-a1"
	a1Count     = 9
	a1FirstTick = 157
	a1Interval  = 58.5
	a4Dur       = 12 * 60
)

var a4Shred = []float64{0.0, 0.05, 0.10, 0.15, 0.55}

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.a1Src = c.Core.F
	ticks := a1Count + c.c4ExtraCount()
	for i := range ticks {
		c.QueueCharTask(c.a1Tick(c.a1Src), a1FirstTick+ceil(a1Interval*float64(i)))
	}
	// this status is purely cosmetic and doesn't do anything right now
	c.AddStatus(a1key, a1FirstTick+ceil(float64(ticks-1)*a1Interval), true)
}

func (c *char) a1Tick(src int) func() {
	return func() {
		if src != c.a1Src {
			return
		}
		scale := a1Scaling + c.c4ExtraHeal()
		heal := scale * c.TotalAtk()
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index(),
			Target:  c.Core.Player.Active(),
			Message: "Rehab Diet",
			Src:     heal,
			Bonus:   c.Stat(attributes.Heal),
		})
	}
}

func (c *char) a4Init() {
	c.a4HydroCryoCount = 0
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Hydro, attributes.Cryo:
			c.a4HydroCryoCount++
		default:
		}
	}
}

func (c *char) makeA4CB() info.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		shred := a4Shred[c.a4HydroCryoCount]
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("escoffier-a4-shred-cryo", a4Dur),
			Ele:   attributes.Cryo,
			Value: -shred,
		})
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("escoffier-a4-shred-hydro", a4Dur),
			Ele:   attributes.Hydro,
			Value: -shred,
		})
	}
}
