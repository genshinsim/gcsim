package thundering

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("thundering pulse", weapon)
	core.RegisterWeaponFunc("thunderingpulse", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.15 + float64(r)*0.05
	stack := 0.09 + float64(r)*0.03
	max := 0.3 + float64(r)*0.1

	normal := 0
	skill := 0

	key := fmt.Sprintf("thundering-pulse-%v", char.Name())

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal {
			return false
		}
		normal = c.Frame + 300 // lasts 5 seconds
		return false
	}, key)

	c.Subscribe(core.PreSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		skill = c.Frame + 600
		return false
	}, key)

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key: "thundering",
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			m[core.DmgP] = 0
			if atk.Info.AttackTag != coretype.AttackTagNormal {
				return m, true
			}
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
			m[core.DmgP] = dmg
			return m, true
		},
		Expiry: -1,
	})

	return "thunderingpulse"
}
