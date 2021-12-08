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
	//equipment with no effect
	core.RegisterWeaponFunc("dullblade", weapon)
	core.RegisterWeaponFunc("silversword", weapon)
	core.RegisterWeaponFunc("wastergreatsword", weapon)
	core.RegisterWeaponFunc("oldmercspal", weapon)
	core.RegisterWeaponFunc("huntersbow", weapon)
	core.RegisterWeaponFunc("seasonedhuntersbow", weapon)
	core.RegisterWeaponFunc("apprenticesnotes", weapon)
	core.RegisterWeaponFunc("pocketgrimoire", weapon)
	core.RegisterWeaponFunc("blacktassel", weapon)
	core.RegisterWeaponFunc("ironpoint", weapon)
	core.RegisterWeaponFunc("beginnersprotector", weapon)

}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

}

/**

t := args[0].(core.Target)
atk := args[1].(*core.AttackEvent)
crit := args[3].(bool)

c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
t := args[0].(core.Target)
atk := args[1].(*core.AttackEvent)
crit := args[3].(bool)

}, )

c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
t := args[0].(core.Target)
atk := args[1].(*core.AttackEvent)
crit := args[3].(bool)

}, )

c.Events.Subscribe(core.OnReactionOccured, func(args ...interface{}) bool {

}, )


**/
