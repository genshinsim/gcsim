package lavawalker

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("lavawalker", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		s.Log.NewEvent("lavawalker 2 pc not implemented", coretype.LogArtifactEvent, c.Index())
	}
	if count >= 4 {
		c.AddPreDamageMod(coretype.PreDamageMod{
			Key:    "4lw",
			Expiry: -1,
			Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				m := make([]float64, core.EndStatType)
				//frozen check first so we don't mistaken coexisting cryo
				if t.AuraContains(core.Pyro) {
					m[core.DmgP] = .35
					return m, true
				}
				return nil, false
			},
		})
	}
	//add flat stat to char
}
