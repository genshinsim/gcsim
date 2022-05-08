package ningguang

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) newScreen(dur int) core.Construct {
	return &construct{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
	}
}

type construct struct {
	src    int
	expiry int
	char   *char
}

func (c *construct) Key() int {
	return c.src
}

func (c *construct) Type() core.GeoConstructType {
	return core.GeoConstructNingSkill
}

func (c *construct) OnDestruct() {
	if c.char.Base.Cons >= 2 {
		//make sure last reset is more than 6 seconds ago
		if c.char.c2reset <= c.char.Core.F-360 && c.char.Cooldown(core.ActionSkill) > 0 {
			//reset cd
			c.char.ResetActionCooldown(core.ActionSkill)
			c.char.c2reset = c.char.Core.F
		}
	}
}

func (c *construct) Expiry() int {
	return c.expiry
}

func (c *construct) IsLimited() bool {
	return true
}

func (c *construct) Count() int {
	return 1
}
