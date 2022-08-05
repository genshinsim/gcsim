package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/construct"
)

type ushi struct {
	src    int
	expiry int
	char   *char
}

func (c *char) newUshi(dur int) construct.Construct {
	return &ushi{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
	}
}

func (u *ushi) OnDestruct()                      { u.char.addStrStack(1) }
func (u *ushi) Key() int                         { return u.src }
func (u *ushi) Type() construct.GeoConstructType { return construct.GeoConstructIttoSkill }
func (u *ushi) Expiry() int                      { return u.expiry }
func (u *ushi) IsLimited() bool                  { return true }
func (u *ushi) Count() int                       { return 1 }
