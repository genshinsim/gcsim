package tenacity

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("tenacity of millelith", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.HPP] = 0.2
		c.AddMod(core.CharStatMod{
			Key: "tom-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		icd := 0
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.2

		s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			ds := args[1].(*core.Snapshot)
			if ds.ActorIndex != c.CharIndex() {
				return false
			}
			if ds.AttackTag != core.AttackTagElementalArt {
				return false
			}
			if icd > s.F {
				return false
			}
			s.Status.AddStatus("tom-proc", 180)
			icd = s.F + 30 //.5 second icd

			s.Log.Debugw("tom 4pc proc", "frame", s.F, "event", core.LogArtifactEvent, "expiry", s.F+180, "icd", s.F+30)
			return false
		}, fmt.Sprintf("pf4-%v", c.Name()))

		c.AddMod(core.CharStatMod{
			Key: "pf-4pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if s.Status.Duration("tom-proc") == 0 {
					return nil, false
				}
				return m, true
			},
			Expiry: -1,
		})

	}
	//add flat stat to char
}
