package calamity

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("calamity queller", weapon)
	core.RegisterWeaponFunc("calamityqueller", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	stacks := 0
	var ctick func(core.Character, *core.Core, int) func()
	ctick = func(char core.Character, c *core.Core, skillInitF int) func() {
		return func() {
			if c.F != skillInitF {
				return
			}
			if c.Status.Duration("consummation-"+char.Name()) <= 0 {
				stacks = 0
				return
			}
			if stacks < 6 {
				stacks++
			}
			c.Log.Debugw("consummation-"+char.Name()+" buff ticking", "frame", c.F, "event", core.LogCharacterEvent, "stacks", stacks)
			if stacks == 6 {
				char.AddTask(ctick(char, c, skillInitF), "consummation-"+char.Name()+"-tick", c.Status.Duration("consummation-"+char.Name()))
			} else {
				char.AddTask(ctick(char, c, skillInitF), "consummation-"+char.Name()+"-tick", 60)
			}
		}
	}

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
			return false
		}
		c.Status.AddStatus("consummation-"+char.Name(), 20*60)
		ctick(char, c, c.F)()
		return false

	}, fmt.Sprintf("consummation-%v", char.Name()))

	val := make([]float64, core.EndStatType)
	mod := float64(r - 1)

	char.AddMod(core.CharStatMod{
		Key:    "calamityqueller",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			val[core.PyroP] = 0.12 + 0.03*mod
			val[core.HydroP] = 0.12 + 0.03*mod
			val[core.CryoP] = 0.12 + 0.03*mod
			val[core.ElectroP] = 0.12 + 0.03*mod
			val[core.AnemoP] = 0.12 + 0.03*mod
			val[core.GeoP] = 0.12 + 0.03*mod
			val[core.PhyP] = 0.12 + 0.03*mod
			val[core.DendroP] = 0.12 + 0.03*mod
			if c.ActiveChar == char.CharIndex() {
				val[core.ATKP] = (0.032 + mod*0.008) * 2.0 * float64(stacks)

			} else {
				val[core.ATKP] = (0.032 + mod*0.008) * float64(stacks)
			}
			return val, true
		},
	})

	return "calamityqueller"
}
