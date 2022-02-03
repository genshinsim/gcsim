package calamity

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("calamity queller", weapon)
	core.RegisterWeaponFunc("calamityqueller", weapon)
}

// Gain 12/15/18/21/24% All Elemental DMG Bonus.
// Obtain Consummation for 20s after using an Elemental Skill, causing ATK to increase by 3.2/4/4.8/5.6/6.4% per second.
// This ATK increase has a maximum of 6 stacks.
// When the character equipped with this weapon is not on the field, Consummation's ATK increase is doubled.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	m := make([]float64, core.EndStatType)
	dmg := .09 + float64(r)*.03
	atkbonus := .024 + float64(r)*.008
	m[core.PyroP] = dmg
	m[core.HydroP] = dmg
	m[core.CryoP] = dmg
	m[core.ElectroP] = dmg
	m[core.AnemoP] = dmg
	m[core.GeoP] = dmg
	m[core.DendroP] = dmg

	expiry := -1
	skillInitF := 0

	c.Events.Subscribe(core.PreSkill, func(args ...interface{}) bool {

		if expiry < c.F {
			skillInitF = c.F
		}
		expiry = c.F + 60*20

		return false
	}, fmt.Sprintf("calamity-queller-%v", char.Name()))

	char.AddMod(core.CharStatMod{
		Key:    "calamity-queller",
		Expiry: -1,
		Amount: func() ([]float64, bool) {

			stacks := 0
			if expiry > c.F {
				stacks = (c.F - skillInitF) / 60
			}
			if stacks > 6 {
				stacks = 6
			}

			atk := atkbonus * float64(stacks)
			if c.ActiveChar != char.CharIndex() {
				atk *= 2
			}
			m[core.ATKP] = atk

			return m, true
		},
	})

	return "calamityqueller"
}
