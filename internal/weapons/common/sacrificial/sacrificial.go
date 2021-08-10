package sacrificial

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("sacrificial bow", weapon)
	combat.RegisterWeaponFunc("sacrificial fragments", weapon)
	combat.RegisterWeaponFunc("sacrificial greatsword", weapon)
	combat.RegisterWeaponFunc("sacrificial sword", weapon)
}

//After damaging an opponent with an Elemental Skill, the skill has a 40/50/60/70/80%
//chance to end its own CD. Can only occur once every 30/26/22/19/16s.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	last := 0
	prob := 0.3 + float64(r)*0.1
	cd := (34 - r*4) * 60
	//add on crit effect
	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.Actor != c.Name() {
			return
		}
		if ds.AttackTag != core.AttackTagElementalArt {
			return
		}
		if last != 0 && s.Frame()-last < cd {
			return
		}
		if c.Cooldown(core.ActionSkill) == 0 {
			return
		}
		if s.Rand().Float64() < prob {
			c.ResetActionCooldown(core.ActionSkill)
			last = s.Frame() + cd
			log.Debugw("sacrificial proc'd", "frame", s.Frame(), "event", core.LogWeaponEvent, "char", c.CharIndex())
		}

	}, fmt.Sprintf("sac-%v", c.Name()))

}
