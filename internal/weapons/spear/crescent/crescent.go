package crescent

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("crescent pike", weapon)
	core.RegisterWeaponFunc("crescentpike", weapon)
}

//After defeating an enemy, ATK is increased by 12/15/18/21/24% for 30s.
//This effect has a maximum of 3 stacks, and the duration of each stack is independent of the others.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	atk := .15 + float64(r)*.05
	active := 0

	c.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		c.Log.Debugw("crescent pike active", "event", core.LogWeaponEvent, "frame", c.F, "char", char.CharIndex(), "expiry", c.F+300)
		active = c.F + 300

		return false
	}, fmt.Sprintf("cp-%v", char.Name()))

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		//check if char is correct?
		if ae.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.F < active {
			//add a new action that deals % dmg immediately
			ai := core.AttackInfo{
				ActorIndex: char.CharIndex(),
				Abil:       "Crescent Pike Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       atk,
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, 1)

		}
		return false
	}, fmt.Sprintf("cpp-%v", char.Name()))

}
