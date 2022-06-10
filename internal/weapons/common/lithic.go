package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type Lithic struct {
	Index int
}

func (b *Lithic) SetIndex(idx int) { b.Index = idx }
func (b *Lithic) Init() error      { return nil }
func NewLithic(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	l := &Lithic{}
	r := p.Refine

	stacks := 0
	val := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnInitialize, func(args ...interface{}) bool {
		for _, char := range c.Player.Chars() {
			if char.CharZone == character.ZoneLiyue {
				stacks++
			}
		}
		val[attributes.CR] = (0.02 + float64(r)*0.01) * float64(stacks)
		val[attributes.ATKP] = (0.06 + float64(r)*0.01) * float64(stacks)
		return true
	}, fmt.Sprintf("lithic-%v", char.Base.Name))
	char.AddStatMod("lithic", -1, attributes.NoStat, func() ([]float64, bool) {
		return val, true
	})

	return l, nil
}
