package generic

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.Rust, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	//Increases Normal Attack DMG by 40% but decreases Charged Attack DMG by 10%.
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	inc := .3 + float64(r)*0.1
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("rust", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagNormal {
				m[attributes.DmgP] = inc
				return m, true
			}
			if atk.Info.AttackTag == attacks.AttackTagExtra {
				m[attributes.DmgP] = -0.1
				return m, true
			}
			return nil, false
		},
	})

	return w, nil
}
