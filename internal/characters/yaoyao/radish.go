package yaoyao

import "github.com/genshinsim/gcsim/pkg/core/player"

func (c *char) radishHeal(hi player.HealInfo) {
	c.Core.Player.Heal(hi)
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
