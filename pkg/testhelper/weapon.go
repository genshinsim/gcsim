package testhelper

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type Weapon struct {
	Index int
}

func (b *Weapon) SetIndex(idx int) { b.Index = idx }
func (b *Weapon) Init() error      { return nil }

func NewFakeWeapon(_ *core.Core, _ *character.CharWrapper, _ info.WeaponProfile) (info.Weapon, error) {
	return &Weapon{}, nil
}
