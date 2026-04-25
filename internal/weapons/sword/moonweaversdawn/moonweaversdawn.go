package moonweaversdawn

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.MoonweaversDawn, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine
	cost := char.EnergyMax

	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.15 + 0.05*float64(r)
	if cost <= 40 {
		val[attributes.DmgP] += 0.21 + 0.07*float64(r)
	} else if cost <= 60 {
		val[attributes.DmgP] += 0.12 + 0.04*float64(r)
	}

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("moonweavers-dawn", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag == attacks.AttackTagElementalBurst {
				return val
			}
			return nil
		},
	})

	return w, nil
}
