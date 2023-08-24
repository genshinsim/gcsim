package bloodtainted

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.BloodtaintedGreatsword, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Increases DMG against opponents affected by Pyro or Electro by 12/15/18/21/24%.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmg := 0.09 + float64(r)*0.03
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = dmg
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("bloodtaintedgreatsword", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			x, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if x.AuraContains(attributes.Pyro, attributes.Electro) {
				return m, true
			}
			return nil, false
		},
	})

	return w, nil
}
