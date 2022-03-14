package skyward

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("skyward pride", weapon)
	core.RegisterWeaponFunc("skywardpride", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.06 + float64(r)*0.02
	char.AddMod(coretype.CharStatMod{
		Key: "skyward pride",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	counter := 0
	dur := 0

	dmg := 0.6 + float64(r)*0.2

	c.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		dur = c.Frame + 1200
		counter = 0
		c.Log.NewEvent("Skyward Pride activated", coretype.LogWeaponEvent, char.Index(), "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-pride-%v", char.Name()))

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		//check if char is correct?
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		//check if cd is up
		if c.Frame > dur {
			return false
		}
		if counter >= 8 {
			return false
		}

		counter++
		ai := core.AttackInfo{
			ActorIndex: char.Index(),
			Abil:       "Skyward Pride Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		c.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, coretype.TargettableEnemy), 0, 1)
		return false
	}, fmt.Sprintf("skyward-pride-%v", char.Name()))
	return "skywardpride"
}
