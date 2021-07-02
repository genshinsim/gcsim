package bolide

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("retracing bolide", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		s.AddShieldBonus(func() float64 {
			return 0.35 //shield bonus always active
		})
	}
	if count >= 4 {
		m := make([]float64, def.EndStatType)
		m[def.DmgP] = 0.4
		c.AddMod(def.CharStatMod{
			Key: "bolide-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, s.IsShielded() && (a == def.AttackTagNormal || a == def.AttackTagExtra)
			},
			Expiry: -1,
		})
	}
}
