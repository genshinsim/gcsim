package thundersoother

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("thundersoother", New)
	core.RegisterSetFunc("thundersoother", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		s.Log.Warnw("thundersoother 2 pc not implemented", "event", core.LogArtifactEvent, "char", c.CharIndex(), "frame", s.F)
	}
	if count >= 4 {
		c.AddPreDamageMod(core.PreDamageMod{
			Key:    "4ts",
			Expiry: -1,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				m := make([]float64, core.EndStatType)
				//frozen check first so we don't mistaken coexisting cryo
				if t.AuraContains(core.Electro) {
					m[core.DmgP] = .35
					return m, true
				}
				return nil, false
			},
		})
	}
	//add flat stat to char
}
