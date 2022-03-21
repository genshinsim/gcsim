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
		c.AddPreDamageMod(core.PreDamageMod{
			Key: "bolide-4pc",
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				return m, s.Shields.IsShielded(c.CharIndex()) && (atk.Info.AttackTag == core.AttackTagNormal || atk.Info.AttackTag == core.AttackTagExtra)
			},
			Expiry: -1,
		})
	}
}
