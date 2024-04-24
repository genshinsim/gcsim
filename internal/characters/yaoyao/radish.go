package yaoyao

import "github.com/genshinsim/gcsim/pkg/core/info"

func (c *char) radishHeal(hi info.HealInfo) {
	c.Core.Player.Heal(hi)
	// c1 and a4 should not proc on c6 radish
	if hi.Message == c6HealMsg {
		return
	}
	if c.Base.Cons >= 1 {
		c.c1()
	}

	if c.Base.Ascension >= 4 {
		active := c.Core.Player.ActiveChar()
		active.AddStatus(a4Status, 5*60, true)
		c.a4Srcs[active.Index] = c.Core.F
		c.QueueCharTask(c.a4(active.Index, c.Core.F), 60)
	}
}
