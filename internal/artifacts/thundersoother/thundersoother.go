package thundersoother

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("thundersoother", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		log.Warnw("thundersoother 2 pc not implemented", "event", core.LogArtifactEvent, "char", c.CharIndex(), "frame", s.Frame())
	}
	if count >= 4 {
		s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			if t.AuraContains(core.Electro) {
				ds.Stats[core.DmgP] += .35
				log.Debugw("thundersoother 4pc on electro", "frame", s.Frame(), "event", core.LogCalc, "char", c.CharIndex(), "new dmg", ds.Stats[core.DmgP])
			}
		}, fmt.Sprintf("ts4-%v", c.Name()))
	}
	//add flat stat to char
}
