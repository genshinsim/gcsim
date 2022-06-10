package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type Royal struct {
	Index int
}

func (b *Royal) SetIndex(idx int)	{ b.Index = idx }
func (b *Royal) Init() error		{ return nil }
func NewRoyal(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Royal{}
	r := p.Refine

	stacks := 0

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if crit {
			stacks = 0
		} else {
			stacks++
			if stacks > 5 {
				stacks = 5
			}
		}
		return false
	}, fmt.Sprintf("royal-%v", char.Base.Name))

	rate := 0.06 + float64(r)*0.02
	m := make([]float64, attributes.EndStatType)
	char.AddStatMod("royal", -1, attributes.NoStat, func() ([]float64, bool) {
		m[attributes.CR] = float64(stacks) * rate
		return m, true
	})

	return w, nil
}
