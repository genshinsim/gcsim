package kaeya

import (
	"github.com/genshinsim/gsim/pkg/def"
	"github.com/genshinsim/gsim/pkg/shield"
)

func (c *char) c1() {
	c.Sim.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		if t.AuraContains(def.Cryo) {
			ds.Stats[def.CR] += 0.15
			c.Log.Debugw("kaeya c1 - adding crit", "event", def.LogCalc, "char", c.Index, "frame", c.Sim.Frame(), "final cr", ds.Stats[def.CR])
		}
	}, "kaeya-c1")
}

func (c *char) c4() {
	c.Sim.AddOnHurt(func(s def.Sim) {
		if s.Frame() < c.c4icd && c.c4icd != 0 {
			return
		}
		if c.HPCurrent/c.HPMax < .2 {
			c.c4icd = c.Sim.Frame() + 3600
			c.Sim.AddShield(&shield.Tmpl{
				Src:        c.Sim.Frame(),
				ShieldType: def.ShieldKaeyaC4,
				HP:         .3 * c.HPMax,
				Ele:        def.Cryo,
				Expires:    c.Sim.Frame() + 1200,
			})
		}
	})
}
