package sacrificial

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("sacrificial bow", weapon)
	core.RegisterWeaponFunc("sacrificial fragments", weapon)
	core.RegisterWeaponFunc("sacrificial greatsword", weapon)
	core.RegisterWeaponFunc("sacrificial sword", weapon)
	core.RegisterWeaponFunc("sacrificialbow", weapon)
	core.RegisterWeaponFunc("sacrificialfragments", weapon)
	core.RegisterWeaponFunc("sacrificialgreatsword", weapon)
	core.RegisterWeaponFunc("sacrificialsword", weapon)
}

//After damaging an opponent with an Elemental Skill, the skill has a 40/50/60/70/80%
//chance to end its own CD. Can only occur once every 30/26/22/19/16s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	last := 0
	prob := 0.3 + float64(r)*0.1
	cd := (34 - r*4) * 60
	//add on crit effect
	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagElementalArt {
			return false
		}
		if last != 0 && c.F-last < cd {
			return false
		}
		if char.Cooldown(core.ActionSkill) == 0 {
			return false
		}
		if c.Rand.Float64() < prob {
			char.ResetActionCooldown(core.ActionSkill)
			last = c.F + cd
			c.Log.Debugw("sacrificial proc'd", "frame", c.F, "event", core.LogWeaponEvent, "char", char.CharIndex())
		}
		return false
	}, fmt.Sprintf("sac-%v", char.Name()))

}
