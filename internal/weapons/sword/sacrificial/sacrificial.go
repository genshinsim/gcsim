package sacrificial

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("sacrificial sword", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	last := 0
	prob := 0.3 + float64(r)*0.1
	cd := (34 - r*4) * 60
	//add on crit effect
	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.Actor != c.Name() {
			return
		}
		if ds.AttackTag != def.AttackTagElementalArt {
			return
		}
		if last != 0 && s.Frame()-last < cd {
			return
		}
		if c.Cooldown(def.ActionSkill) == 0 {
			return
		}
		if s.Rand().Float64() < prob {
			c.ResetActionCooldown(def.ActionSkill)
			last = s.Frame() + cd
			log.Debugw("sacrificial sword proc'd", "frame", s.Frame(), "event", def.LogWeaponEvent, "char", c.CharIndex())
		}

	}, fmt.Sprintf("sac-sword-%v", c.Name()))

}
