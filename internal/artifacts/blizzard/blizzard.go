package blizzard

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("blizzard strayer", New)
	core.RegisterSetFunc("blizzardstrayer", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.CryoP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "bs-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		s.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
			t := args[0].(core.Target)
			ds := args[1].(*core.Snapshot)
			if ds.ActorIndex != c.CharIndex() {
				return false
			}
			switch t.AuraType() {
			case core.Cryo:
				ds.Stats[core.CR] += .2
				s.Log.Debugw("blizzard strayer 4pc on cryo", "frame", s.F, "event", core.LogCalc, "char", c.CharIndex(), "new crit", ds.Stats[core.CR])
			case core.Frozen:
				ds.Stats[core.CR] += .4
				s.Log.Debugw("blizzard strayer 4pc on frozen", "frame", s.F, "event", core.LogCalc, "char", c.CharIndex(), "new crit", ds.Stats[core.CR])
			}
			return false
		}, fmt.Sprintf("bs4-%v", c.Name()))

	}
	//add flat stat to char
}
