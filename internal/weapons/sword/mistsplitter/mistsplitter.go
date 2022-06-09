package mistsplitter

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("mistsplitter reforged", weapon)
	core.RegisterWeaponFunc("mistsplitterreforged", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	base := 0.09 + float64(r)*0.03
	m[core.PyroP] = base
	m[core.HydroP] = base
	m[core.CryoP] = base
	m[core.ElectroP] = base
	m[core.AnemoP] = base
	m[core.GeoP] = base
	m[core.DendroP] = base
	stack := 0.06 + float64(r)*0.02
	max := 0.03 + float64(r)*0.01
	bonus := core.EleToDmgP(char.Ele())

	normal := 0
	skill := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {

		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if atk.Info.Element == core.Physical {
			return false
		}
		normal = c.F + 300 // lasts 5 seconds
		return false
	}, fmt.Sprintf("mistsplitter-%v", char.Name()))

	c.Events.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		skill = c.F + 600
		return false

	}, fmt.Sprintf("mistsplitter-%v", char.Name()))

	char.AddMod(core.CharStatMod{
		Key: "mistsplitter",
		Amount: func() ([]float64, bool) {
			count := 0
			if char.CurrentEnergy() < char.MaxEnergy() {
				count++
			}
			if normal > c.F {
				count++
			}
			if skill > c.F {
				count++
			}
			dmg := float64(count) * stack
			if count >= 3 {
				dmg += max
			}
			//bonus for current char
			m[bonus] = base + dmg
			return m, true
		},
		Expiry: -1,
	})
	return "mistsplitterreforged"
}
