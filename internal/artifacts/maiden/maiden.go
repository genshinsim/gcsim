package maiden

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("maiden beloved", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.Heal] = 0.15
		c.AddMod(def.CharStatMod{
			Key: "maiden-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		dur := 0

		s.AddEventHook(func(s def.Sim) bool {
			// s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if s.ActiveCharIndex() != c.CharIndex() {
				return false
			}
			dur = s.Frame() + 600
			log.Debugw("maiden 4pc proc", "frame", s.Frame(), "event", def.LogArtifactEvent, "char", c.CharIndex(), "expiry", dur)
			return false
		}, fmt.Sprintf("maid 4pc - %v", c.Name()), def.PostBurstHook)

		s.AddIncHealBonus(func() float64 {
			if s.Frame() < dur {
				return 0.2
			}
			return 0
		})
	}
	//add flat stat to char
}
