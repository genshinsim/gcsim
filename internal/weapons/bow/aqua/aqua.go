package aqua

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.AquaSimulacra, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	//add on hit effect to sim?
	m := make([]float64, attributes.EndStatType)
	v := make([]float64, attributes.EndStatType)
	v[attributes.HPP] = 0.12 + float64(r)*0.04
	m[attributes.DmgP] = 0.15 + float64(r)*0.05

	char.AddStatMod("aquasimulacra", -1, attributes.NoStat, func() ([]float64, bool) {
		return v, true
	})

	char.AddAttackMod("aquasimulacra", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		//TODO: need range check here
		return m, true
	})

	return w, nil
}
