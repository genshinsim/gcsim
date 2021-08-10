package lavawalker

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("lavawalker", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		log.Warnw("lavawalker 2 pc not implemented", "event", core.LogArtifactEvent, "char", c.CharIndex(), "frame", s.Frame())
	}
	if count >= 4 {
		s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			if t.AuraContains(core.Pyro) {
				ds.Stats[core.DmgP] += .35
				log.Debugw("lavawalker 4pc on pyro", "frame", s.Frame(), "event", core.LogCalc, "char", c.CharIndex(), "new dmg", ds.Stats[core.DmgP])
			}
		}, fmt.Sprintf("lw4-%v", c.Name()))
	}
	//add flat stat to char
}
