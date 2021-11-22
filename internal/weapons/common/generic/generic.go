package generic

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("generic catalyst", weapon)
	core.RegisterWeaponFunc("generic bow", weapon)
	core.RegisterWeaponFunc("generic claymore", weapon)
	core.RegisterWeaponFunc("generic sword", weapon)
	core.RegisterWeaponFunc("generic spear", weapon)
	core.RegisterWeaponFunc("generic", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

}

/**

t := args[0].(core.Target)
ds := args[1].(*core.Snapshot)
crit := args[3].(bool)

c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
t := args[0].(core.Target)
ds := args[1].(*core.Snapshot)
crit := args[3].(bool)

}, )

c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
t := args[0].(core.Target)
ds := args[1].(*core.Snapshot)
crit := args[3].(bool)

}, )

c.Events.Subscribe(core.OnReactionOccured, func(args ...interface{}) bool {

}, )


**/
