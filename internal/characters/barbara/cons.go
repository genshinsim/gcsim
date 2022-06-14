package barbara

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (c *char) c2() {
	val := make([]float64, attributes.EndStatType)
	val[attributes.HydroP] = 0.15
	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue
		}
		char.AddStatMod("barbara-c2",
			-1, attributes.NoStat, func() ([]float64, bool) {
				if c.Core.Status.Duration("barbskill") >= 0 {
					return val, true
				} else {
					return nil, false
				}
			})

	}
}

func (c *char) c1(delay int) {
	c.AddTask(func() {
		c.AddEnergy("barbara-c1", 1)
		c.c1(0)
	}, "barbara-c1", delay+10*60)
}

// inspired from hutao c6
//TODO: does this even work?
func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index { //trigger only when not barbara
			c.checkc6()
		}
		return false
	}, "barbara-c6")
}

func (c *char) checkc6() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Core.F < c.c6icd && c.c6icd != 0 {
		return
	}
	//if dead, revive back to 1 hp
	if c.HP() <= -1 {
		c.HPCurrent = c.MaxHP()
	}

	c.c6icd = c.Core.F + 60*60*15
}
