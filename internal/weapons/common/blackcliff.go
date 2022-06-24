package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type Blackcliff struct {
	Index int
}

func (b *Blackcliff) SetIndex(idx int) { b.Index = idx }
func (b *Blackcliff) Init() error      { return nil }

func NewBlackcliff(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {

	b := &Blackcliff{}

	atk := 0.09 + float64(p.Refine)*0.03
	index := 0
	stacks := []int{-1, -1, -1}

	m := make([]float64, attributes.EndStatType)
	char.AddStatMod("blackcliff", -1, attributes.ATKP, func() ([]float64, bool) {
		count := 0
		for _, v := range stacks {
			if v > c.F {
				count++
			}
		}
		m[attributes.ATKP] = atk * float64(count)
		return m, true
	})

	c.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		stacks[index] = c.F + 1800
		index++
		if index == 3 {
			index = 0
		}
		return false
	}, fmt.Sprintf("blackcliff-%v", char.Base.Name))

	return b, nil
}
