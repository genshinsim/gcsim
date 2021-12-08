package blizzard

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("blizzard strayer", New)
	core.RegisterSetFunc("blizzardstrayer", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		var m [core.EndStatType]float64
		m[core.CryoP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "bs-2pc",
			Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		c.AddPreDamageMod(core.PreDamageMod{
			Key:    "4bs",
			Expiry: -1,
			Amount: func(atk *core.AttackEvent, t core.Target) ([core.EndStatType]float64, bool) {
				var m [core.EndStatType]float64
				//frozen check first so we don't mistaken coexisting cryo
				if t.AuraContains(core.Frozen) {
					m[core.CR] = 0.4
					return m, true
				}
				if t.AuraContains(core.Cryo) {
					m[core.CR] = 0.2
					return m, true
				}
				return m, false
			},
		})

	}
	//add flat stat to char
}
