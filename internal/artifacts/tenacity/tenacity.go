package tenacity

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("tenacity of millelith", New)
	core.RegisterSetFunc("tenacity of the millelith", New)
	core.RegisterSetFunc("tenacityofthemillelith", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.HPP] = 0.2
		c.AddMod(core.CharStatMod{
			Key: "tom-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.2
		icd := 0

		s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			atk := args[1].(*core.AttackEvent)
			if atk.Info.ActorIndex != c.CharIndex() {
				return false
			}
			if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
				return false
			}
			if icd > s.F {
				return false
			}
			s.Status.AddStatus("tom-proc", 180)
			icd = s.F + 30 //.5 second icd

			for _, char := range s.Chars {
				char.AddMod(core.CharStatMod{
					Key:    "tom-4pc",
					Expiry: s.F + 180,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}
			s.Log.NewEvent("tom 4pc proc", core.LogArtifactEvent, c.CharIndex(), "expiry", s.F+180, "icd", s.F+30)
			return false
		}, fmt.Sprintf("tom4-%v", c.Name()))
	}
}
