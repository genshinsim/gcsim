package chongyun

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) c4() {
	icd := 0
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(core.Reactable)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.F < icd {
			return false
		}
		if !t.AuraContains(attributes.Cryo) {
			return false
		}

		c.AddEnergy("chongyun-c4", 2)

		c.Core.Log.NewEvent("chongyun c4 recovering 2 energy", glog.LogCharacterEvent, c.Index).
			Write("final energy", c.Energy)
		icd = c.Core.F + 120

		return false
	}, "chongyun-c4")
}

func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.15
	c.AddAttackMod("chongyun-c6", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.AttackTag != combat.AttackTagElementalBurst {
			return nil, false
		}
		if t.HP()/t.MaxHP() < c.HPCurrent/c.MaxHP() {
			return m, true
		}
		return nil, false
	})
}
