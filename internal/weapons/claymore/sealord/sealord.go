package sealord

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("luxurious sea-lord", weapon)
	core.RegisterWeaponFunc("luxurious sealord", weapon)
	core.RegisterWeaponFunc("luxurioussealord", weapon)
}

// Increases Elemental Burst DMG by 12~24%. When Elemental Burst hits opponents, there is a 100% chance of summoning a huge onrush of tuna that charges and deals 100~200% ATK as AoE DMG. This effect can occur once every 15s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	burstDmgIncrease := .09 + float64(r)*0.03
	tunaDmg := .75 + float64(r)*0.25
	effectLastProc := -9999

	val := make([]float64, core.EndStatType)
	val[core.DmgP] = burstDmgIncrease
	char.AddMod(core.CharStatMod{
		Expiry: -1,
		Key:    "luxurious-sea-lord",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a == core.AttackTagElementalBurst {
				return val, true
			}
			return nil, false
		},
	})

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if c.F < effectLastProc+15*60 {
			return false
		}
		if ds.AttackTag != core.AttackTagElementalBurst {
			return false
		}
		effectLastProc = c.F
		d := char.Snapshot(
			"Luxurious Sea-Lord Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			tunaDmg,
		)
		d.Targets = core.TargetAll
		char.QueueDmg(&d, 1)

		return false
	}, fmt.Sprintf("sealord-%v", char.Name()))
}
