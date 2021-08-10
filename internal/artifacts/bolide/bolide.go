package bolide

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("retracing bolide", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		s.AddShieldBonus(func() float64 {
			return 0.35 //shield bonus always active
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.4
		c.AddMod(core.CharStatMod{
			Key: "bolide-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, s.IsShielded() && (a == core.AttackTagNormal || a == core.AttackTagExtra)
			},
			Expiry: -1,
		})
	}
}
