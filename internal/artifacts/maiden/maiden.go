package maiden

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("maiden beloved", New)
	core.RegisterSetFunc("maidenbeloved", New)
}

// 2 piece: Character Healing Effectiveness +15%
// 4 piece: Using an Elemental Skill or Burst increases healing received by all party members by 20% for 10s.
func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.Heal] = 0.15
		c.AddMod(coretype.CharStatMod{
			Key: "maiden-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		dur := 0

		s.Subscribe(core.PostBurst, func(args ...interface{}) bool {
			// s.s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if s.Player.ActiveChar != c.Index()() {
				return false
			}
			dur = s.Frame + 600
			s.Log.NewEvent("maiden 4pc proc", coretype.LogArtifactEvent, c.Index(), "expiry", dur)
			return false
		}, fmt.Sprintf("maid 4pc - %v", c.Name()))

		// Applies to all characters, so no filters needed
		s.Health.AddIncHealBonus(func(healedCharIndex int) float64 {
			if s.Frame < dur {
				return 0.2
			}
			return 0
		})
	}
	//add flat stat to char
}
