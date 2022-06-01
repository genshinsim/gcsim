package maiden

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("maiden beloved", New)
	core.RegisterSetFunc("maidenbeloved", New)
}

// 2 piece: Character Healing Effectiveness +15%
// 4 piece: Using an Elemental Skill or Burst increases healing received by all party members by 20% for 10s.
func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.Heal] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "maiden-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		dur := 0

		f := func(args ...interface{}) bool {
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			dur = s.F + 600
			s.Log.NewEvent("maiden 4pc proc", core.LogArtifactEvent, c.CharIndex(), "expiry", dur)
			return false
		}
		s.Events.Subscribe(core.PreBurst, f, fmt.Sprintf("maid 4pc - %v", c.Name()))
		s.Events.Subscribe(core.PreSkill, f, fmt.Sprintf("maid 4pc - %v", c.Name()))

		// Applies to all characters, so no filters needed
		s.Health.AddIncHealBonus(func(healedCharIndex int) float64 {
			if s.F < dur {
				return 0.2
			}
			return 0
		})
	}
	//add flat stat to char
}
