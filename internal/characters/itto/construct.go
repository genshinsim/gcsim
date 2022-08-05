package itto

import (
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type ushi struct {
	src    int
	expiry int
	char   *char
}

func (u *ushi) OnDestruct() {
	u.char.changeStacks(1)
	u.char.Core.Log.NewEvent("itto ushi stack gained on exit", glog.LogCharacterEvent, u.char.Index).
		Write("stacks", u.char.Tags[u.char.stackKey])
}

func (u *ushi) Key() int                         { return u.src }
func (u *ushi) Type() construct.GeoConstructType { return construct.GeoConstructIttoSkill }
func (u *ushi) Expiry() int                      { return u.expiry }
func (u *ushi) IsLimited() bool                  { return true }
func (c *ushi) Count() int                       { return 1 }
