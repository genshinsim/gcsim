package itto

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) newUshi(dur int) core.Construct {
	return &construct{
		src:    c.Core.Frame,
		expiry: c.Core.Frame + dur,
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
	return core.GeoConstructIttoSkill
}

func (c *construct) OnDestruct() {
	c.char.Tags["strStack"] += 1
	if c.char.Tags["strStack"] > 5 {
		c.char.Tags["strStack"] = 5
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
