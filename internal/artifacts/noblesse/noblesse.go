package noblesse

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("noblesse oblige", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.DmgP] = 0.2
		c.AddMod(def.CharStatMod{
			Key: "nob-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, a == def.AttackTagElementalBurst
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		s.AddEventHook(func(s def.Sim) bool {
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

			log.Debugw("noblesse 4pc proc", "frame", s.Frame(), "event", def.LogArtifactEvent, "expiry", s.Status("nob-4pc"))
			return false
		}, fmt.Sprintf("no 4pc - %v", c.Name()), def.PostBurstHook)

		m := make([]float64, def.EndStatType)
		m[def.ATKP] = 0.2
		c.AddMod(def.CharStatMod{
			Key: "nob-4pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
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
