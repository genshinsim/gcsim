package dragonbane

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("lion's roar", weapon)
	core.RegisterWeaponFunc("lionsroar", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}

		if ds.IsReactionDamage {
			return false
		}

		if t.AuraContains(core.Electro, core.Pyro) {
			ds.Stats[core.DmgP] += dmg
			c.Log.Debugw("lion's roar", "frame", c.F, "event", core.LogWeaponEvent, "char", char.CharIndex(), "final dmg%", ds.Stats[core.DmgP])
		}
		return false
	}, fmt.Sprintf("lion-%v", char.Name()))

}
