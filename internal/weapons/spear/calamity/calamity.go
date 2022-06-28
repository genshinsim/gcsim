package calamity

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.CalamityQueller, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Gain 12% All Elemental DMG Bonus. Obtain Consummation for 20s after using
	//an Elemental Skill, causing ATK to increase by 3.2% per second. This ATK
	//increase has a maximum of 6 stacks. When the character equipped with this
	//weapon is not on the field, Consummation's ATK increase is doubled.
	w := &Weapon{}
	r := p.Refine

	//fixed elemental dmg bonus
	dmg := .09 + float64(r)*.03
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = dmg
	m[attributes.HydroP] = dmg
	m[attributes.CryoP] = dmg
	m[attributes.ElectroP] = dmg
	m[attributes.AnemoP] = dmg
	m[attributes.GeoP] = dmg
	m[attributes.DendroP] = dmg
	char.AddStatMod(character.StatMod{Base: modifier.NewBase("calamity-queller", -1), AffectedStat: attributes.NoStat, Amount: func() ([]float64, bool) {
		return m, true
	}})

	//atk increase per stack after using skill
	//double bonus if not on field
	atkbonus := .024 + float64(r)*.008
	skillInitF := -1
	skillPressBonus := make([]float64, attributes.EndStatType)
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}

		dur := 60 * 20
		if skillInitF == -1 || (skillInitF+dur) < c.F {
			skillInitF = c.F
		}
		char.AddStatMod(character.StatMod{Base: modifier.NewBase("calamity-consummation", dur), AffectedStat: attributes.NoStat, Amount: func() ([]float64, bool) {
			stacks := (c.F - skillInitF) / 60
			if stacks > 6 {
				stacks = 6
			}
			atk := atkbonus * float64(stacks)
			if c.Player.Active() != char.Index {
				atk *= 2
			}
			skillPressBonus[attributes.ATKP] = atk

			return skillPressBonus, true
		}})

		return false
	}, fmt.Sprintf("calamity-queller-%v", char.Base.Key.String()))

	return w, nil
}
