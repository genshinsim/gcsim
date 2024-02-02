package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

type ushi struct {
	src    int
	expiry int
	char   *char
	dir    geometry.Point
	pos    geometry.Point
}

func (c *char) newUshi(dur int, dir, pos geometry.Point) construct.Construct {
	return &ushi{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
		dir:    dir,
		pos:    pos,
	}
}

func (u *ushi) OnDestruct()                      { u.char.addStrStack("ushi-exit", 1) }
func (u *ushi) Key() int                         { return u.src }
func (u *ushi) Type() construct.GeoConstructType { return construct.GeoConstructIttoSkill }
func (u *ushi) Expiry() int                      { return u.expiry }
func (u *ushi) IsLimited() bool                  { return true }
func (u *ushi) Count() int                       { return 1 }
func (u *ushi) Direction() geometry.Point        { return u.dir }
func (u *ushi) Pos() geometry.Point              { return u.pos }
