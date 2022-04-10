package swordofdescension

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("swordofdescension", weapon)
}

// Descension
// This weapon's effect is only applied on the following platform(s):
// "PlayStation Network"
// Hitting enemies with Normal or Charged Attacks grants a 50% chance to deal 200% ATK as DMG in a small AoE. This effect can only occur once every 10s.
// Additionally, if the Traveler equips the Sword of Descension, their ATK is increased by 66.
//  * Weapon refines do not affect this weapon
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	icd := 0
	m := make([]float64, core.EndStatType)

	if char.Key() < core.TravelerDelim {
		char.AddMod(core.CharStatMod{
			Key:    "swordofdescension",
			Expiry: -1,
			Amount: func() ([]float64, bool) {
				m[core.ATK] = 66
				return m, true
			},
		})
	}

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}

		// ignore if character not on field
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		// Ignore if neither a charged nor normal attack
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}

		// Ignore if icd is still up
		if c.F < icd {
			return false
		}

		// Ignore 50% of the time, 1:1 ratio
		if c.Rand.Float64() < 0.5 {
			return false
		}

		icd = c.F + 10*60

		ai := core.AttackInfo{
			ActorIndex: char.CharIndex(),
			Abil:       "Sword of Descension Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       2.00,
		}

		c.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 1)

		return false
	}, fmt.Sprintf("swordofdescension-%v", char.Name()))
	return "swordofdescension"
}
