package twilight

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("fading twilight", weapon)
	core.RegisterWeaponFunc("fadingtwilight", weapon)
	core.RegisterWeaponFunc("twilight", weapon)
}

//Has three states, Evengleam, Afterglow, and Dawnblaze, which increase DMG dealt by 7.5%/12.5%/17.5% respectively.
// When attacks hit opponents, this weapon will switch to the next state. This weapon can change states once every 7s. 
//The character equipping this weapon can still trigger the state switch while not on the field.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	cycle:=0
	base := 0.0

	m[core.DmgP] = base
	//buff is assumed to be dynamic but idk
	char.AddPreDamageMod(core.PreDamageMod{
		Expiry: -1,
		Key:    "twilight-bonus-dmg",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			switch cycle {
			case 2:
				base = 0.14 + float64(r)*0.035
			case 1:
				base = 0.10 + float64(r)*0.025
			default:
				base = 0.06 + float64(r)*0.015
			}
	

			m[core.DmgP] = base
			return nil, false
		},
	})

	icd := 0
	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}

		if icd > c.F {
			return false
		}
		icd = c.F + 420 //once every 7 seconds #smokeweedeveryday
		cycle++
		cycle=cycle%3

		return false
	}, fmt.Sprintf("fadingtwilight-%v", char.Name()))

	return "fadingtwilight"
}
