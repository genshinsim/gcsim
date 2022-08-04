package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/construct"
)

type ushi struct {
	src    int
	expiry int
	char   *char
}

func (u *ushi) OnDestruct() {
	u.char.Tags[u.char.stackKey] += 1
	if u.char.Tags[u.char.stackKey] > 5 {
		u.char.Tags[u.char.stackKey] = 5
	}
}

func (u *ushi) Key() int                         { return u.src }
func (u *ushi) Type() construct.GeoConstructType { return construct.GeoConstructIttoSkill }
func (u *ushi) Expiry() int                      { return u.expiry }
func (u *ushi) IsLimited() bool                  { return true }
func (c *ushi) Count() int                       { return 1 }
