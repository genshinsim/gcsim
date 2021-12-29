package skyward

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyward spine", weapon)
	core.RegisterWeaponFunc("skywardspine", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	m := make([]float64, core.EndStatType)
	m[core.CR] = 0.06 + float64(r)*0.02
	m[core.AtkSpd] = 0.12

	char.AddMod(core.CharStatMod{
		Key: "skyward spine",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	icd := 0
	atk := .25 + .15*float64(r)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		//check if char is correct?
		if ae.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		//check if cd is up
		if icd > c.F {
			return false
		}
		if c.Rand.Float64() > .5 {
			return false
		}

		//add a new action that deals % dmg immediately
		ai := core.AttackInfo{
			ActorIndex: char.CharIndex(),
			Abil:       "Skyward Spine Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       atk,
		}
		c.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, 1)

		//trigger cd
		icd = c.F + 120
		return false
	}, fmt.Sprintf("skyward-spine-%v", char.Name()))
	return "skywardspine"
}
