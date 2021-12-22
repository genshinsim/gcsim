package kaeya

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
)

func (c *char) c1() {
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "kaeya-c1",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
				return nil, false
			}
			if !t.AuraContains(core.Cryo, core.Frozen) {
				return nil, false
			}
			val[core.CR] = 0.15
			return val, true
		},
	})
}

func (c *char) c4() {
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.F < c.c4icd && c.c4icd != 0 {
			return false
		}
		if c.HPCurrent/c.HPMax < .2 {
			c.c4icd = c.Core.F + 3600
			c.Core.Shields.Add(&shield.Tmpl{
				Src:        c.Core.F,
				ShieldType: core.ShieldKaeyaC4,
				HP:         .3 * c.HPMax,
				Ele:        core.Cryo,
				Expires:    c.Core.F + 1200,
			})
		}
		return false
	}, "kaeya-c4")

}
