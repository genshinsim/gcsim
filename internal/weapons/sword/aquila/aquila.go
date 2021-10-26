package aquila

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("aquila favonia", weapon)
	core.RegisterWeaponFunc("aquilafavonia", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = .15 + .05*float64(r)
	char.AddMod(core.CharStatMod{
		Key: "acquila favonia",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	dmg := 1.7 + .3*float64(r)
	heal := .85 + .15*float64(r)

	last := -1

	c.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if c.F-last < 900 && last != -1 {
			return false
		}
		last = c.F
		d := char.Snapshot(
			"Aquila Favonia",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			dmg,
		)
		d.Targets = core.TargetAll
		char.QueueDmg(&d, 1)

		atk := d.BaseAtk*(1+d.Stats[core.ATKP]) + d.Stats[core.ATK]

		c.Log.Debugw("acquila heal triggered", "frame", c.F, "event", core.LogWeaponEvent, "atk", atk, "heal amount", atk*heal)
		c.Health.HealActive(atk * heal)
		return false
	}, fmt.Sprintf("aquila-%v", char.Name()))
}
