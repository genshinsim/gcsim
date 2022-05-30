package razor

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
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
