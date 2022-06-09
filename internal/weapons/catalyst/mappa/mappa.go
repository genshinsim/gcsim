package mappa

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("mappa mare", weapon)
	core.RegisterWeaponFunc("mappamare", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	stacks := 0
	dur := 0

	addStack := func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		if c.F > dur {
			stacks = 1
			dur = c.F + 600
			c.Log.NewEvent("mappa proc'd", core.LogWeaponEvent, char.CharIndex(), "stacks", stacks, "expiry", dur)
		} else if stacks < 2 {
			stacks++
			c.Log.NewEvent("mappa proc'd", core.LogWeaponEvent, char.CharIndex(), "stacks", stacks, "expiry", dur)
		}
		return false
	}

	for i := core.EventType(core.ReactionEventStartDelim + 1); i < core.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, addStack, "mappa"+char.Name())
	}

	dmg := 0.06 + float64(r)*0.02
	m := make([]float64, core.EndStatType)

	char.AddMod(core.CharStatMod{
		Key: "mappa",
		Amount: func() ([]float64, bool) {
			if c.F > dur {
				return nil, false
			}

			m[core.PyroP] = dmg * float64(stacks)
			m[core.HydroP] = dmg * float64(stacks)
			m[core.CryoP] = dmg * float64(stacks)
			m[core.ElectroP] = dmg * float64(stacks)
			m[core.AnemoP] = dmg * float64(stacks)
			m[core.GeoP] = dmg * float64(stacks)
			m[core.DendroP] = dmg * float64(stacks)
			return m, true
		},
		Expiry: -1,
	})

	return "mappamare"
}
