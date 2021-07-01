package favonius

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("favonius sword", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	p := 0.50 + float64(r)*0.1
	cd := 810 - r*90
	icd := 0
	//add on crit effect
	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.Actor != c.Name() {
			return
		}
		if s.ActiveCharIndex() != c.CharIndex() {
			return
		}
		if icd > s.Frame() {
			return
		}

		if s.Rand().Float64() > p {
			return
		}
		log.Debugw("favonius sword proc'd", "frame", s.Frame(), "event", def.LogWeaponEvent)

		c.QueueParticle("favonius sword", 3, def.NoElement, 150)

		icd = s.Frame() + cd

	}, fmt.Sprintf("favo-sword-%v", c.Name()))

}
