package bloodstained

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("bloodstained chivalry", New)
	core.RegisterSetFunc("bloodstainedchivalry", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		var m [core.EndStatType]float64
		m[core.PhyP] = 0.25
		c.AddMod(core.CharStatMod{
			Key: "bloodstained-2pc",
			Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
}
