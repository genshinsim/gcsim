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

	mATK := make([]float64, core.EndStatType)
	mATK[core.ATKP] = .45 + float64(r)*0.15

	mEM := make([]float64, core.EndStatType)
	mEM[core.EM] = 180 + float64(r)*60

	mDmg := make([]float64, core.EndStatType)
	dmg := .36 + float64(r)*.12
	mDmg[core.PyroP] = dmg
	mDmg[core.HydroP] = dmg
	mDmg[core.CryoP] = dmg
	mDmg[core.ElectroP] = dmg
	mDmg[core.AnemoP] = dmg
	mDmg[core.GeoP] = dmg
	mDmg[core.DendroP] = dmg

	icd := -1
	c.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		next := args[1].(int)
		//ignore if char is not the active one
		if next != char.CharIndex() {
			return false
		}
		//if char is the active one then we just came on to field
		if c.F < icd {
			return false
		}
		icd = c.F + 60*30

		//random 1 of 3
		i := c.Rand.Intn(3)
		var stat string
		var fn func() ([]float64, bool)
		switch i {
		case 0:
			stat = "em"
			fn = func() ([]float64, bool) { return mEM, true }
		case 1:
			stat = "dmg%"
			fn = func() ([]float64, bool) { return mDmg, true }
		case 2:
			stat = "atk%"
			fn = func() ([]float64, bool) { return mATK, true }
		}

		expiry := c.F + 60*10
		char.AddMod(core.CharStatMod{
			Key:    "widsith",
			Expiry: expiry,
			Amount: fn,
		})
		c.Log.NewEvent("widsith proc'd", core.LogWeaponEvent, char.CharIndex(), "stat", stat, "expiring", expiry)

		return false
	}, fmt.Sprintf("width-%v", char.Name()))

	return "thewidsith"

}
