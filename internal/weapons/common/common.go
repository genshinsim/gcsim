package common

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type NoEffect struct {
	Index int
}

func (b *NoEffect) SetIndex(idx int) { b.Index = idx }
func (b *NoEffect) Init() error      { return nil }

func NewNoEffect(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	return &NoEffect{}, nil
}
