package yaoyao

import "github.com/genshinsim/gcsim/pkg/core/player"

func (c *char) radishHeal(hi player.HealInfo) {

	c.Core.Player.Heal(hi)
	if c.Base.Cons >= 1 {
		c.c1()
	}

	if c.Base.Ascension >= 4 {
		c.a4()
	}

}

func (c *char) radishHitCB(hi player.HealInfo) {
	if c.Base.Cons >= 2 {
		c.c2()
	}
}
