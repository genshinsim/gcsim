package zhongli

import "github.com/genshinsim/gsim/pkg/core"

type construct struct {
	src    int
	expiry int
	char   *char
}

func (c *construct) Key() int {
	return c.src
}

func (c *construct) Type() core.GeoConstructType {
	return core.GeoConstructZhongliSkill
}

func (c *construct) OnDestruct() {
	if c.char.steeleCount > 0 {
		c.char.steeleCount--
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
