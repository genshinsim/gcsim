package amos

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
	core.RegisterWeaponFunc(keys.AmosBow, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmgpers := 0.06 + 0.02*float64(r)

	m := make([]float64, attributes.EndStatType)
	// m[attributes.DmgP] = 0.09 + 0.03*float64(r)
	flat := 0.09 + 0.03*float64(r)

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("amos", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			m[attributes.DmgP] = flat
			travel := float64(c.F-atk.Snapshot.SourceFrame) / 60
			stacks := int(travel / 0.1)
			if stacks > 5 {
				stacks = 5
			}
			m[attributes.DmgP] += dmgpers * float64(stacks)
			return m, true
		},
	})

	return w, nil
}
