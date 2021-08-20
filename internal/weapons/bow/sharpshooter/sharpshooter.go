package sharpshooter

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("sharpshooter's oath", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	dmg := 0.18 + float64(r)*0.06
	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.HitWeakPoint {
			ds.Stats[core.DmgP] += dmg
			c.Log.Debugw("sharpshooter", "frame", c.F, "event", core.LogWeaponEvent, "final dmg%", ds.Stats[core.DmgP])
		}
		return false
	}, fmt.Sprintf("sharpshooter-%v", char.Name()))
}
