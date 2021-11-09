package kaeya

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
)

func (c *char) c1() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		if ds.ActorIndex != c.Index {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if t.AuraContains(core.Cryo) {
			ds.Stats[core.CR] += 0.15
			c.Core.Log.Debugw("kaeya c1 - adding crit", "event", core.LogCalc, "char", c.Index, "frame", c.Core.F, "final cr", ds.Stats[core.CR])
		}
		return false
	}, "kaeya-c1")
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
