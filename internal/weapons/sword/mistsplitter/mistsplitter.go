package mistsplitter

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("mistsplitter reforged", weapon)
	core.RegisterWeaponFunc("mistsplitterreforged", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	base := 0.09 + float64(r)*0.03
	m[core.PyroP] = base
	m[core.HydroP] = base
	m[coretype.CryoP] = base
	m[core.ElectroP] = base
	m[core.AnemoP] = base
	m[core.GeoP] = base
	m[core.DendroP] = base
	stack := 0.06 + float64(r)*0.02
	max := 0.03 + float64(r)*0.01
	bonus := coretype.EleToDmgP(char.Ele())

	normal := 0
	skill := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {

		atk := args[1].(*coretype.AttackEvent)

		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal {
			return false
		}
		if atk.Info.Element == core.Physical {
			return false
		}
		normal = c.Frame + 300 // lasts 5 seconds
		return false
	}, fmt.Sprintf("mistsplitter-%v", char.Name()))

	c.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		skill = c.Frame + 600
		return false

	}, fmt.Sprintf("mistsplitter-%v", char.Name()))

	char.AddMod(coretype.CharStatMod{
		Key: "mistsplitter",
		Amount: func() ([]float64, bool) {
			count := 0
			if char.CurrentEnergy() < char.MaxEnergy() {
				count++
			}
			if normal > c.Frame {
				count++
			}
			if skill > c.Frame {
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
