package albedo

import (
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type skillConstruct struct {
	src    int
	expiry int
	char   *char
	dir    info.Point
	pos    info.Point
}

func (c *char) newConstruct(dur int, dir, pos info.Point) *skillConstruct {
	return &skillConstruct{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
		dir:    dir,
		pos:    pos,
	}
}

func (c *skillConstruct) OnDestruct()                      { c.char.skillActive = false }
func (c *skillConstruct) Key() int                         { return c.src }
func (c *skillConstruct) Type() construct.GeoConstructType { return construct.GeoConstructAlbedoSkill }
func (c *skillConstruct) Expiry() int                      { return c.expiry }
func (c *skillConstruct) IsLimited() bool                  { return true }
func (c *skillConstruct) Count() int                       { return 1 }
func (c *skillConstruct) Direction() info.Point            { return c.dir }
func (c *skillConstruct) Pos() info.Point                  { return c.pos }
