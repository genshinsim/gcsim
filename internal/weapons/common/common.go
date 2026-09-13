package common

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type NoEffect struct {
	Index int
}

func (n *NoEffect) SetIndex(idx int) { n.Index = idx }
func (n *NoEffect) Init() error      { return nil }

func NewNoEffect() *NoEffect {
	return &NoEffect{}
}

func (n *NoEffect) NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	return n, nil
}
