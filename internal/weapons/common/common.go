package common

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

type NoEffect struct {
	Index int
	data  *model.WeaponData
}

func (n *NoEffect) SetIndex(idx int)        { n.Index = idx }
func (n *NoEffect) Init() error             { return nil }
func (n *NoEffect) Data() *model.WeaponData { return n.data }

func NewNoEffect(data *model.WeaponData) core.NewWeaponFunc {
	n := &NoEffect{data: data}
	return n.NewWeapon
}

func (n *NoEffect) NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	return n, nil
}
