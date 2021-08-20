package bloodstained

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("bloodstained chivalry", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.PhyP] = 0.25
		c.AddMod(core.CharStatMod{
			Key: "bloodstained-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
}
