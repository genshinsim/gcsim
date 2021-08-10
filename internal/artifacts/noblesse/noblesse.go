package noblesse

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("noblesse oblige", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.2
		c.AddMod(core.CharStatMod{
			Key: "nob-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, a == core.AttackTagElementalBurst
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		s.AddEventHook(func(s core.Sim) bool {
			// s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if s.ActiveCharIndex() != c.CharIndex() {
				return false
			}

			nob, ok := s.GetCustomFlag("nob-4pc")
			//only activate if none existing
			if s.Status("nob-4pc") == 0 || (nob == c.CharIndex() && ok) {
				s.AddStatus("nob-4pc", 720)
				s.SetCustomFlag("nob-4pc", c.CharIndex())
			}

			log.Debugw("noblesse 4pc proc", "frame", s.Frame(), "event", core.LogArtifactEvent, "expiry", s.Status("nob-4pc"))
			return false
		}, fmt.Sprintf("no 4pc - %v", c.Name()), core.PostBurstHook)

		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.2
		c.AddMod(core.CharStatMod{
			Key: "nob-4pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if s.Status("nob-4pc") > 0 {
					return m, true
				}
				return nil, false
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
