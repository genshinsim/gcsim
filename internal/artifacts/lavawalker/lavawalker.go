package lavawalker

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("lavawalker", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		s.Log.Warnw("lavawalker 2 pc not implemented", "event", core.LogArtifactEvent, "char", c.CharIndex(), "frame", s.F)
	}
	if count >= 4 {
		s.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
			t := args[0].(core.Target)
			ds := args[1].(*core.Snapshot)
			if ds.ActorIndex != c.CharIndex() {
				return false
			}
			if t.AuraContains(core.Pyro) {
				ds.Stats[core.DmgP] += .35
				s.Log.Debugw("lavawalker 4pc on pyro", "frame", s.F, "event", core.LogCalc, "char", c.CharIndex(), "new dmg", ds.Stats[core.DmgP])
			}
			return false
		}, fmt.Sprintf("lw4-%v", c.Name()))

	}
	//add flat stat to char
}
