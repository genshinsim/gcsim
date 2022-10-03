package testhelper

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type Weapon struct {
	Index int
}

func (b *Weapon) SetIndex(idx int) { b.Index = idx }
func (b *Weapon) Init() error      { return nil }

func NewFakeWeapon(_ *core.Core, _ *character.CharWrapper, _ weapon.WeaponProfile) (weapon.Weapon, error) {
	return &Weapon{}, nil
}
