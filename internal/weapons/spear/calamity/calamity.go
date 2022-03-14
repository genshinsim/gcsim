package calamity

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("calamity", weapon)
	core.RegisterWeaponFunc("calamityqueller", weapon)
}

// Gain 12/15/18/21/24% All Elemental DMG Bonus.
// Obtain Consummation for 20s after using an Elemental Skill, causing ATK to increase by 3.2/4/4.8/5.6/6.4% per second.
// This ATK increase has a maximum of 6 stacks.
// When the character equipped with this weapon is not on the field, Consummation's ATK increase is doubled.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	dmg := .09 + float64(r)*.03
	atkbonus := .024 + float64(r)*.008

	skillInitF := -1

	c.Subscribe(core.PreSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}

		// update init frame on first cast or after 20s from last cast
		dur := 60 * 20
		if skillInitF == -1 || (skillInitF+dur) < c.Frame {
			skillInitF = c.Frame
		}

		char.AddMod(coretype.CharStatMod{
			Key:    "calamity-consummation",
			Expiry: c.Frame + dur,
			Amount: func() ([]float64, bool) {
				m := make([]float64, core.EndStatType)

				stacks := (c.Frame - skillInitF) / 60
				if stacks > 6 {
					stacks = 6
				}

				atk := atkbonus * float64(stacks)
				if c.ActiveChar != char.Index() {
					atk *= 2
				}
				m[core.ATKP] = atk

				return m, true
			},
		})

		return false
	}, fmt.Sprintf("calamity-queller-%v", char.Name()))

	m := make([]float64, core.EndStatType)
	m[core.PyroP] = dmg
	m[core.HydroP] = dmg
	m[coretype.CryoP] = dmg
	m[core.ElectroP] = dmg
	m[core.AnemoP] = dmg
	m[core.GeoP] = dmg
	m[core.DendroP] = dmg

	char.AddMod(coretype.CharStatMod{
		Key:    "calamity-queller",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	return "calamityqueller"
}
