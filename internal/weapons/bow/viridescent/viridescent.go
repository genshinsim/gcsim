package viridescent

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the viridescent hunt", weapon)
	core.RegisterWeaponFunc("theviridescenthunt", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	cd := 900 - r*60
	icd := 0
	mult := 0.3 + float64(r)*0.1

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		//check if char is correct?
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		// Vhunt passive only applies for NAs and CAs
		// For Tartaglia this also includes melee NAs/CAs
		// See https://youtu.be/EBtOiFhrs94?t=221, Test 4 and 5
		if !((atk.Info.AttackTag == core.AttackTagNormal) || (atk.Info.AttackTag == core.AttackTagExtra)) {
			return false
		}
		//check if cd is up
		if icd > c.F {
			return false
		}
		//50% chance to proc
		if c.Rand.Float64() > 0.5 {
			return false
		}

		//add a new action that deals % dmg immediately
		ai := core.AttackInfo{
			ActorIndex: char.CharIndex(),
			Abil:       "Viridescent",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       mult,
		}

		for i := 0; i <= 240; i += 30 {
			c.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, i+1)
		}

		//trigger cd
		icd = c.F + cd

		return false
	}, fmt.Sprintf("veridescent-%v", char.Name()))

	return "theviridescenthunt"
}
