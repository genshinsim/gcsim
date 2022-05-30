package razor

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

// Picking up an Elemental Orb or Particle increases Razor's DMG by 10% for 8s.
func (c *char) c1() {
	dur := 0
	c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		dur = c.Core.F + 8*60
		return false
	}, "razor-c1")

	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.1
	c.AddStatMod("c1", -1, attributes.DmgP, func() ([]float64, bool) {
		if c.Core.F > dur {
			return nil, false
		}
		return val, true
	})
}

// Increases CRIT Rate against opponents with less than 30% HP by 10%.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.1

	c.AddAttackMod("razor-c2", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if t.HP()/t.MaxHP() < 0.3 {
			return m, true
		}
		return nil, false
	})
}
