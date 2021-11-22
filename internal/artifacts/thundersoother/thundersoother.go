package thundersoother

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("thundersoother", New)
	core.RegisterSetFunc("thundersoother", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		s.Log.Warnw("thundersoother 2 pc not implemented", "event", core.LogArtifactEvent, "char", c.CharIndex(), "frame", s.F)
	}
	if count >= 4 {
		s.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
			ds := args[1].(*core.Snapshot)
			t := args[0].(core.Target)
			if ds.ActorIndex != c.CharIndex() {
				return false
			}
			if t.AuraContains(core.Electro) {
				ds.Stats[core.DmgP] += .35
				s.Log.Debugw("thundersoother 4pc on electro", "frame", s.F, "event", core.LogCalc, "char", c.CharIndex(), "new dmg", ds.Stats[core.DmgP])
			}
			return false
		}, fmt.Sprintf("ts4-%v", c.Name()))

	}
	//add flat stat to char
}
