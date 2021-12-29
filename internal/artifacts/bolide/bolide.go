package bolide

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("retracing bolide", New)
	core.RegisterSetFunc("retracingbolide", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		s.Shields.AddBonus(func() float64 {
			return 0.35 //shield bonus always active
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.4
		c.AddMod(core.CharStatMod{
			Key: "bolide-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, s.Shields.IsShielded() && (a == core.AttackTagNormal || a == core.AttackTagExtra)
			},
			Expiry: -1,
		})
	}
}
