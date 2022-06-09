package travelershandysword

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("travelershandysword", weapon)
}

// Each Elemental Orb or Particle collected restores 1/1.25/1.5/1.75/2% HP.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	c.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		// ignore if character not on field
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		c.Health.Heal(core.HealInfo{
			Type:    core.HealTypePercent,
			Message: "Traveler's Handy Sword (Proc)",
			Src:     0.0075 + float64(r)*0.0025,
		})

		return false
	}, fmt.Sprintf("travelershandysword-%v", char.Name()))

	return "travelershandysword"
}
