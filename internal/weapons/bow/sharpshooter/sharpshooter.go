package sharpshooter

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {

}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Increases DMG against weak spots by 24%.
	w := &Weapon{}
	r := p.Refine

	dmg := 0.18 + float64(r)*0.06
	char.AddAttackMod("sharpshooter", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		m := make([]float64, attributes.EndStatType)
		if atk.Info.HitWeakPoint {
			m[attributes.DmgP] = dmg
			return m, true
		}
		return nil, false
	})

	return w, nil
}
