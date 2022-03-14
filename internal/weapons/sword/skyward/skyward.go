package skyward

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("skyward blade", weapon)
	core.RegisterWeaponFunc("skywardblade", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	dur := -1
	c.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		dur = c.Frame + 720
		c.Log.NewEvent("Skyward Blade activated", coretype.LogWeaponEvent, char.Index(), "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-blade-%v", char.Name()))

	m := make([]float64, core.EndStatType)
	m[core.CR] = 0.03 + float64(r)*0.01

	char.AddMod(coretype.CharStatMod{
		Key: "skyward blade",
		Amount: func() ([]float64, bool) {
			m[core.AtkSpd] = 0
			if dur > c.Frame {
				m[core.AtkSpd] = 0.1 //if burst active
			}
			return m, true
		},
		Expiry: -1,
	})

	//damage procs
	atk := .15 + .05*float64(r)

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {

		ae := args[1].(*coretype.AttackEvent)

		//check if char is correct?
		if ae.Info.ActorIndex != char.Index() {
			return false
		}
		if ae.Info.AttackTag != coretype.AttackTagNormal && ae.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		//check if buff up
		if dur < c.Frame {
			return false
		}

		//add a new action that deals % dmg immediately
		ai := core.AttackInfo{
			ActorIndex: char.Index(),
			Abil:       "Skyward Blade Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       atk,
		}
		c.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, coretype.TargettableEnemy), 0, 1)

		return false

	}, fmt.Sprintf("skyward-blade-%v", char.Name()))

	return "skywardblade"
}
