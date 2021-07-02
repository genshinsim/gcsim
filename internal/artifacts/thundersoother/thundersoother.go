package thundersoother

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("thundersoother", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		log.Warnw("thundersoother 2 pc not implemented", "event", def.LogArtifactEvent, "char", c.CharIndex(), "frame", s.Frame())
	}
	if count >= 4 {
		s.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			switch t.AuraType() {
			case def.Electro:
				ds.Stats[def.DmgP] += .35
				log.Debugw("thundersoother 4pc on electro", "frame", s.Frame(), "event", def.LogCalc, "char", c.CharIndex(), "new dmg", ds.Stats[def.DmgP])
			case def.EC:
				ds.Stats[def.DmgP] += .35
				log.Debugw("thundersoother 4pc on electro", "frame", s.Frame(), "event", def.LogCalc, "char", c.CharIndex(), "new dmg", ds.Stats[def.DmgP])
			}
		}, fmt.Sprintf("ts4-%v", c.Name()))
	}
	//add flat stat to char
}
