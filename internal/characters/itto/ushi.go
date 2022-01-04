package itto

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

type ushi struct {
	src    int
	expiry int
	char   *char
}

func (u *ushi) Key() int {
	return u.src
}

func (u *ushi) Type() core.GeoConstructType {
	return core.GeoConstructIttoSkill
}

func (u *ushi) OnDestruct() {
	if u.char.skillStacks < 5 {
		u.char.skillStacks++
	}
}

func (u *ushi) Expiry() int {
	return u.expiry
}

func (u *ushi) IsLimited() bool {
	return true
}

func (u *ushi) Count() int {
	return 1
}

func (c *char) newCow(dur int) core.Construct {
	return &ushi{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
	}
}
