package sharpshooter

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.SharpshootersOath, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	//Increases DMG against weak spots by 24%.
	w := &Weapon{}
	r := p.Refine

	dmg := 0.18 + float64(r)*0.06
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("sharpshooter", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			m := make([]float64, attributes.EndStatType)
			if atk.Info.HitWeakPoint {
				m[attributes.DmgP] = dmg
				return m, true
			}
			return nil, false
		},
	})

	return w, nil
}
