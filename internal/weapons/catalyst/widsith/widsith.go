package widsith

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the widsith", weapon)
	core.RegisterWeaponFunc("thewidsith", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	last := 0
	expiry := 0

	atk := .45 + float64(r)*0.15
	em := 180 + float64(r)*60
	dmg := .36 + float64(r)*.12

	m := make([]float64, core.EndStatType)

	char.AddMod(core.CharStatMod{
		Key: "widsith",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, expiry > c.F
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		next := args[1].(int)
		//ignore if char is not the active one
		if next != char.CharIndex() {
			return false
		}
		//if char is the active one then we just came on to field
		if last != 0 && c.F-last < 1800 { //30 sec icd
			return false
		}
		last = c.F
		expiry = c.F + 600 //10 sec duration
		//random 1 of 3
		i := c.Rand.Intn(3)

		switch i {
		case 0:
			m[core.EM] = em
			m[core.PyroP] = 0
			m[core.HydroP] = 0
			m[core.CryoP] = 0
			m[core.ElectroP] = 0
			m[core.AnemoP] = 0
			m[core.GeoP] = 0
			m[core.DendroP] = 0
			m[core.ATKP] = 0
			c.Log.Debugw("widsith proc'd", "frame", c.F, "event", core.LogWeaponEvent, "char", char.CharIndex(), "stat", "em", "expiring", expiry)
		case 1:
			m[core.EM] = 0
			m[core.PyroP] = dmg
			m[core.HydroP] = dmg
			m[core.CryoP] = dmg
			m[core.ElectroP] = dmg
			m[core.AnemoP] = dmg
			m[core.GeoP] = dmg
			m[core.DendroP] = dmg
			m[core.ATKP] = 0
			c.Log.Debugw("widsith proc'd", "frame", c.F, "event", core.LogWeaponEvent, "char", char.CharIndex(), "stat", "dmg%", "expiring", expiry)
		default:
			m[core.EM] = 0
			m[core.PyroP] = 0
			m[core.HydroP] = 0
			m[core.CryoP] = 0
			m[core.ElectroP] = 0
			m[core.AnemoP] = 0
			m[core.GeoP] = 0
			m[core.DendroP] = 0
			m[core.ATKP] = atk
			c.Log.Debugw("widsith proc'd", "frame", c.F, "event", core.LogWeaponEvent, "char", char.CharIndex(), "stat", "atk%", "expiring", expiry)
		}

		return false
	}, fmt.Sprintf("width-%v", char.Name()))

	return "thewidsith"

}
