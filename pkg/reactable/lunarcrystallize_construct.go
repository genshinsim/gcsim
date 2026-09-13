package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type LunarCrystallizeConstruct struct {
	src    int
	expiry int
	dir    info.Point
	pos    info.Point
}

func (r *Reactable) newLunarCrystallizeConstruct(dir, pos info.Point) *LunarCrystallizeConstruct {
	return &LunarCrystallizeConstruct{
		src:    r.core.F,
		expiry: r.core.F + lcrDur,
		dir:    dir,
		pos:    pos,
	}
}

func (c *LunarCrystallizeConstruct) OnDestruct() {}
func (c *LunarCrystallizeConstruct) Key() int    { return c.src }
func (c *LunarCrystallizeConstruct) Type() construct.GeoConstructType {
	return construct.GeoConstructLunarCrystallize
}
func (c *LunarCrystallizeConstruct) Expiry() int           { return c.expiry }
func (c *LunarCrystallizeConstruct) IsLimited() bool       { return true }
func (c *LunarCrystallizeConstruct) Count() int            { return 1 }
func (c *LunarCrystallizeConstruct) Direction() info.Point { return c.dir }
func (c *LunarCrystallizeConstruct) Pos() info.Point       { return c.pos }
