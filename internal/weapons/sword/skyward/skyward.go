package skyward

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyward blade", weapon)
	core.RegisterWeaponFunc("skywardblade", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	dur := -1
	c.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		dur = c.F + 720
		c.Log.Debugw("Skyward Blade activated", "frame", c.F, "event", core.LogWeaponEvent, "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-blade-%v", char.Name()))

	m := make([]float64, core.EndStatType)
	m[core.CR] = 0.03 + float64(r)*0.01

	char.AddMod(core.CharStatMod{
		Key: "skyward blade",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.AtkSpd] = 0
			if dur > c.F {
				m[core.AtkSpd] = 0.1 //if burst active
			}
			return m, true
		},
		Expiry: -1,
	})

	//damage procs
	atk := .15 + .05*float64(r)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {

		ae := args[1].(*core.AttackEvent)

		//check if char is correct?
		if ae.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		//check if buff up
		if dur < c.F {
			return false
		}

		//add a new action that deals % dmg immediately
		ai := core.AttackInfo{
			ActorIndex: char.CharIndex(),
			Abil:       "Skyward Blade Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       atk,
		}
		c.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, 1)

		return false

	}, fmt.Sprintf("skyward-blade-%v", char.Name()))

	return "skywardblade"
}
