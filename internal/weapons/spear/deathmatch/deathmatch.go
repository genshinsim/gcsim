package deathmatch

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.Deathmatch, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	multiple := make([]float64, attributes.EndStatType)
	multiple[attributes.ATKP] = .12 + .04*float64(r)
	multiple[attributes.DEFP] = .12 + .04*float64(r)

	single := make([]float64, attributes.EndStatType)
	single[attributes.ATKP] = .18 + .06*float64(r)
	char.AddStatMod("deathmatch", -1, attributes.NoStat, func() ([]float64, bool) {
		if len(c.Combat.Targets()) > 2 {
			return multiple, true
		}
		return single, true
	})

	return w, nil
}
