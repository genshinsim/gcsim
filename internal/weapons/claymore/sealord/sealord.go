package sealord

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
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
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.F < effectLastProc+15*60 {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagElementalBurst {
			return false
		}
		effectLastProc = c.F
		ai := core.AttackInfo{
			ActorIndex: char.CharIndex(),
			Abil:       "Luxurious Sea-Lord Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       tunaDmg,
		}
		c.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, 1)

		return false
	}, fmt.Sprintf("sealord-%v", char.Name()))
}
