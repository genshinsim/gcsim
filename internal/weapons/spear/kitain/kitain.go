package kitain

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.KitainCrossSpear, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Increases Elemental Skill DMG by 6%. After Elemental Skill hits an
	//opponent, the character loses 3 Energy but regenerates 3 Energy every 2s
	//for the next 6s. This effect can occur once every 10s. Can be triggered
	//even when the character is not on the field.
	w := &Weapon{}
	r := p.Refine

	//permanent increase
	m := make([]float64, attributes.EndStatType)
	base := 0.045 + float64(r)*0.015
	m[attributes.DmgP] = base
	char.AddAttackMod("kitain-skill-dmg-buff", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.AttackTag == combat.AttackTagElementalArt || atk.Info.AttackTag == combat.AttackTagElementalArtHold {
			return m, true
		}
		return nil, false
	})

	regen := 2.5 + float64(r)*0.5
	icd := 0
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 600
		char.AddEnergy("kitain", -3)
		for i := 120; i <= 360; i += 120 {
			c.Tasks.Add(func() {
				char.AddEnergy("kitain", regen)
			}, i)
		}
		return false
	}, fmt.Sprintf("kitain-%v", char.Base.Key.String()))
	return w, nil
}
