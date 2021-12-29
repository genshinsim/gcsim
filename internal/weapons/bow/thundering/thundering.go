package thundering

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("thundering pulse", weapon)
	core.RegisterWeaponFunc("thunderingpulse", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.15 + float64(r)*0.05
	stack := 0.09 + float64(r)*0.03
	max := 0.3 + float64(r)*0.1

	normal := 0
	skill := 0

	key := fmt.Sprintf("thundering-pulse-%v", char.Name())

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal {
			return false
		}
		normal = c.F + 300 // lasts 5 seconds
		return false
	}, key)

	c.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		skill = c.F + 600
		return false
	}, key)

	char.AddMod(core.CharStatMod{
		Key: "thundering",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.DmgP] = 0
			if a != core.AttackTagNormal {
				return m, true
			}
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
			if count > 3 {
				count = 3 // should never happen
				dmg = max
			}
			m[core.DmgP] = dmg
			return m, true
		},
		Expiry: -1,
	})

	return "thunderingpulse"
}
