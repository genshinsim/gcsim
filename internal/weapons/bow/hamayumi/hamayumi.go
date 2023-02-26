package hamayumi

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.Hamayumi, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	nm := .12 + .04*float64(r)
	ca := .09 + .03*float64(r)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("hamayumi", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			val := make([]float64, attributes.EndStatType)
			if atk.Info.AttackTag == attacks.AttackTagNormal {
				val[attributes.DmgP] = nm
				if char.Energy == char.EnergyMax {
					val[attributes.DmgP] = nm * 2
				}
				return val, true
			}

			if atk.Info.AttackTag == attacks.AttackTagExtra {
				val[attributes.DmgP] = ca
				if char.Energy == char.EnergyMax {
					val[attributes.DmgP] = ca * 2
				}
				return val, true
			}
			return nil, false
		},
	})

	return w, nil
}
