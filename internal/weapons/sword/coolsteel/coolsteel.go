package coolsteel

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func init() {
	core.RegisterWeaponFunc(keys.CoolSteel, NewWeapon)
}

//Increases DMG against enemies affected by Hydro or Cryo by 12-24%.
type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.09 + float64(r)*0.03

	char.AddAttackMod("coolsteel", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		x, ok := t.(*enemy.Enemy)
		if !ok {
			return nil, false
		}
		if x.AuraContains(attributes.Hydro, attributes.Cryo) {
			return m, true
		}
		return nil, false
	})

	return w, nil
}
