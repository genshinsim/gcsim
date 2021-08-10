package kaeya

import (
	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/shield"
)

func (c *char) c1() {
	c.Sim.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return
		}
		if t.AuraContains(core.Cryo) {
			ds.Stats[core.CR] += 0.15
			c.Log.Debugw("kaeya c1 - adding crit", "event", core.LogCalc, "char", c.Index, "frame", c.Sim.Frame(), "final cr", ds.Stats[core.CR])
		}
	}, "kaeya-c1")
}

func (c *char) c4() {
	c.Sim.AddOnHurt(func(s core.Sim) {
		if s.Frame() < c.c4icd && c.c4icd != 0 {
			return
		}
		if c.HPCurrent/c.HPMax < .2 {
			c.c4icd = c.Sim.Frame() + 3600
			c.Sim.AddShield(&shield.Tmpl{
				Src:        c.Sim.Frame(),
				ShieldType: core.ShieldKaeyaC4,
				HP:         .3 * c.HPMax,
				Ele:        core.Cryo,
				Expires:    c.Sim.Frame() + 1200,
			})
		}
	})
}
