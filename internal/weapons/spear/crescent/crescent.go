package crescent

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("crescent pike", weapon)
	core.RegisterWeaponFunc("crescentpike", weapon)
}

//After defeating an enemy, ATK is increased by 12/15/18/21/24% for 30s.
//This effect has a maximum of 3 stacks, and the duration of each stack is independent of the others.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	atk := .15 + float64(r)*.05
	active := 0

	c.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		c.Log.NewEvent("crescent pike active", coretype.LogWeaponEvent, char.Index(), "expiry", c.Frame+300)
		active = c.Frame + 300

		return false
	}, fmt.Sprintf("cp-%v", char.Name()))

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*coretype.AttackEvent)
		//check if char is correct?
		if ae.Info.ActorIndex != char.Index() {
			return false
		}
		if ae.Info.AttackTag != coretype.AttackTagNormal && ae.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if c.Frame < active {
			//add a new action that deals % dmg immediately
			ai := core.AttackInfo{
				ActorIndex: char.Index(),
				Abil:       "Crescent Pike Proc",
				AttackTag:  core.AttackTagWeaponSkill,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.Physical,
				Durability: 100,
				Mult:       atk,
			}
			c.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, coretype.TargettableEnemy), 0, 1)

		}
		return false
	}, fmt.Sprintf("cpp-%v", char.Name()))
	return "crescentpike"
}
