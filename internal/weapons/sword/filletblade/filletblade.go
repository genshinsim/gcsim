package filletblade

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("fillet blade", weapon)
	core.RegisterWeaponFunc("filletblade", weapon)
}

// On hit, has 50% chance to deal 240/280/320/360/400% ATK DMG to a single enemy.
// Can only occur once every 15/14/13/12/11s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	icd := 0
	cd := 960 - 60*r

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if icd > c.F {
			return false
		}
		if c.Rand.Float64() > 0.5 {
			return false
		}
		// add a new action that deals % dmg immediately
		// superconduct attack
		ai := core.AttackInfo{
			ActorIndex: char.CharIndex(),
			Abil:       "Fillet Blade Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       2.0 + 0.4*float64(r),
		}
		c.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, 1)

		// trigger cd
		icd = c.F + cd

		return false
	}, fmt.Sprintf("fillet-blade-%v", char.Name()))
	return "filletblade"
}
