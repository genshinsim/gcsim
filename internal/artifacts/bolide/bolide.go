package bolide

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("retracing bolide", New)
	core.RegisterSetFunc("retracingbolide", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		s.Player.AddShieldBonus(func() float64 {
			return 0.35 //shield bonus always active
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.4
		c.AddPreDamageMod(coretype.PreDamageMod{
			Key: "bolide-2pc",
			Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				return m, s.Player.IsCharShielded(c.Index()) && (atk.Info.AttackTag == coretype.AttackTagNormal || atk.Info.AttackTag == coretype.AttackTagExtra)
			},
			Expiry: -1,
		})
	}
}
