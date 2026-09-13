package testhelper

import (
	_ "embed"

	"github.com/genshinsim/gcsim/pkg/catalog"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

// TODO: insert a custom weapon at runtime
const TestWeaponKey = keys.InvalidWeapon

func RegisterTestWeapon() {
	// TODO: this should be part of registration
	catalog.WeaponMap[TestWeaponKey] = catalog.WeaponMap[keys.DullBlade]
	core.RegisterWeaponFunc(TestWeaponKey, NewFakeWeapon)
}

type Weapon struct {
	Index int
}

func (b *Weapon) SetIndex(idx int) { b.Index = idx }
func (b *Weapon) Init() error      { return nil }

func NewFakeWeapon(_ *core.Core, _ *character.CharWrapper, _ info.WeaponProfile) (info.Weapon, error) {
	return &Weapon{}, nil
}
