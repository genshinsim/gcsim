package bloodstained

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("bloodstained chivalry", New)
	core.RegisterSetFunc("bloodstainedchivalry", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.PhyP] = 0.25
		c.AddMod(coretype.CharStatMod{
			Key: "bloodstained-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
}
