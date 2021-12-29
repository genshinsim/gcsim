package paleflame

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("pale flame", New)
	core.RegisterSetFunc("paleflame", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.PhyP] = 0.25
		c.AddMod(core.CharStatMod{
			Key: "maiden-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		stacks := 0
		icd := 0
		dur := 0
		m := make([]float64, core.EndStatType)

		s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			atk := args[1].(*core.AttackEvent)
			if atk.Info.ActorIndex != c.CharIndex() {
				return false
			}
			if atk.Info.AttackTag != core.AttackTagElementalArt {
				return false
			}
			if icd > s.F {
				return false
			}
			stacks++
			if stacks > 2 {
				stacks = 2
				m[core.PhyP] = 0.25
			}
			m[core.ATKP] = 0.09 * float64(stacks)

			s.Log.Debugw("pale flame 4pc proc", "frame", s.F, "event", core.LogArtifactEvent, "stacks", stacks, "expiry", s.F+420, "icd", s.F+18)
			icd = s.F + 18
			dur = s.F + 420
			return false
		}, fmt.Sprintf("pf4-%v", c.Name()))

		c.AddMod(core.CharStatMod{
			Key: "pf-4pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if dur < s.F {
					m[core.ATKP] = 0
					m[core.PhyP] = 0
					return nil, false
				}

				return m, true
			},
			Expiry: -1,
		})

	}
	//add flat stat to char
}
