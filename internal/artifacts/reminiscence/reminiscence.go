package reminiscence

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("reminiscence of shime", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.ATKP] = 0.18
		c.AddMod(def.CharStatMod{
			Key: "rem-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		m := make([]float64, def.EndStatType)
		m[def.DmgP] = 0.50
		s.AddEventHook(func(s def.Sim) bool {
			if s.ActiveCharIndex() != c.CharIndex() {
				return false
			}
			if c.CurrentEnergy() > 15 {
				//consume 15 energy, increased normal/charge/plunge dmg by 50%
				c.AddEnergy(-15)
				c.AddMod(def.CharStatMod{
					Key: "rem-4pc",
					Amount: func(ds def.AttackTag) ([]float64, bool) {
						if ds != def.AttackTagNormal && ds != def.AttackTagExtra && ds != def.AttackTagPlunge {
							return nil, false
						}
						return m, true
					},
					Expiry: s.Frame() + 600,
				})

			}
			return false
		}, "rem-4pc", def.PreSkillHook)

	}
	//add flat stat to char
}
