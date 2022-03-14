package tenacity

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("tenacity of millelith", New)
	core.RegisterSetFunc("tenacity of the millelith", New)
	core.RegisterSetFunc("tenacityofthemillelith", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.HPP] = 0.2
		c.AddMod(coretype.CharStatMod{
			Key: "tom-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		icd := 0

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
			s.AddStatus("tom-proc", 180)
			icd = s.Frame + 30 //.5 second icd

			for _, char := range s.Chars {
				char.AddMod(coretype.CharStatMod{
					Key: "tom-4pc",
					Amount: func() ([]float64, bool) {
						m := make([]float64, core.EndStatType)
						m[core.ATKP] = 0.2
						if s.StatusDuration("tom-proc") == 0 {
							return nil, false
						}
						return m, true
					},
					Expiry: s.Frame + 180,
				})

			}
			s.Log.NewEvent("tom 4pc proc", coretype.LogArtifactEvent, c.Index(), "expiry", s.Frame+180, "icd", s.Frame+30)
			return false
		}, fmt.Sprintf("tom4-%v", c.Name()))

	}
	//add flat stat to char
}
