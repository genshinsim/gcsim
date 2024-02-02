package testhelper

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Weapon struct {
	Index int
	data  *model.WeaponData
}

func (b *Weapon) SetIndex(idx int)        { b.Index = idx }
func (b *Weapon) Init() error             { return nil }
func (b *Weapon) Data() *model.WeaponData { return b.data }

func NewFakeWeapon(_ *core.Core, _ *character.CharWrapper, _ info.WeaponProfile) (info.Weapon, error) {
	return &Weapon{data: &model.WeaponData{}}, nil
}
