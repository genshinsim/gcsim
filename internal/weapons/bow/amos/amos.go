package amos

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("amos' bow", weapon)
	core.RegisterWeaponFunc("amosbow", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	dmgpers := 0.06 + 0.02*float64(r)

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.09 + 0.03*float64(r)
	char.AddMod(core.CharStatMod{
		Key: "amos",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, a == core.AttackTagNormal || a == core.AttackTagExtra
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if char.CharIndex() != ds.ActorIndex {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		//calculate travel time
		travel := float64(c.F-ds.SourceFrame-ds.AnimationFrames) / 60
		stacks := int(travel / 0.1)
		if stacks > 5 {
			stacks = 5
		}
		ds.Stats[core.DmgP] += dmgpers * float64(stacks)
		c.Log.Debugw("amos bow", "frame", c.F, "event", core.LogCalc, "stacks", stacks, "final dmg%", ds.Stats[core.DmgP])
		return false
	}, fmt.Sprintf("amos-%v", char.Name()))

}
