package kaeya

import (
	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (c *char) c1() {
	c.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "kaeya-c1",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
				return nil, false
			}
			if !t.AuraContains(coretype.Cryo, coretype.Frozen) {
				return nil, false
			}
			val[core.CR] = 0.15
			return val, true
		},
	})
}

func (c *char) c4() {
	c.Core.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.Frame < c.c4icd && c.c4icd != 0 {
			return false
		}
		if c.HPCurrent/c.HPMax < .2 {
			c.c4icd = c.Core.Frame + 3600
			c.Core.Shields.Add(&shield.Tmpl{
				Src:        c.Core.Frame,
				ShieldType: core.ShieldKaeyaC4,
				Name:       "Kaeya C4",
				HP:         .3 * c.HPMax,
				Ele:        coretype.Cryo,
				Expires:    c.Core.Frame + 1200,
			})
		}
		return false
	}, "kaeya-c4")

}
