package paleflame

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("pale flame", New)
	core.RegisterSetFunc("paleflame", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.PhyP] = 0.25
		c.AddMod(coretype.CharStatMod{
			Key: "pf-2pc",
			Amount: func() ([]float64, bool) {
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

		s.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
			atk := args[1].(*coretype.AttackEvent)
			if atk.Info.ActorIndex != c.Index() {
				return false
			}
			if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
				return false
			}
			if icd > s.Frame {
				return false
			}
			// reset stacks if expired
			if dur < s.Frame {
				stacks = 0
			}
			stacks++
			if stacks >= 2 {
				stacks = 2
				m[core.PhyP] = 0.25
			}
			m[core.ATKP] = 0.09 * float64(stacks)

			s.Log.NewEvent("pale flame 4pc proc", coretype.LogArtifactEvent, c.Index(), "stacks", stacks, "expiry", s.Frame+420, "icd", s.Frame+18)
			icd = s.Frame + 18
			dur = s.Frame + 420
			return false
		}, fmt.Sprintf("pf4-%v", c.Name()))

		c.AddMod(coretype.CharStatMod{
			Key: "pf-4pc",
			Amount: func() ([]float64, bool) {
				if dur < s.Frame {
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
