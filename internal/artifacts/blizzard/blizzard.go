package blizzard

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("blizzard strayer", New)
	core.RegisterSetFunc("blizzardstrayer", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[coretype.CryoP] = 0.15
		c.AddMod(coretype.CharStatMod{
			Key: "bs-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		c.AddPreDamageMod(coretype.PreDamageMod{
			Key:    "4bs",
			Expiry: -1,
			Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				m := make([]float64, core.EndStatType)
				//frozen check first so we don't mistaken coexisting cryo
				if t.AuraContains(coretype.Frozen) {
					m[core.CR] = 0.4
					return m, true
				}
				if t.AuraContains(coretype.Cryo) {
					m[core.CR] = 0.2
					return m, true
				}
				return nil, false
			},
		})

	}
	//add flat stat to char
}
