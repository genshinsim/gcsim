package barbara

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// The Stamina Consumption of characters within Let the Show Begin♪'s Melody Loop is reduced by 12%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	// a1 lasts for duration of barb skill which is 900 frames
	c.Core.Player.AddStamPercentMod("barb-a1-stam", skillDuration, func(a action.Action) (float64, bool) {
		return -0.12, false
	})
}

// When your active character gains an Elemental Orb/Particle, the duration of the Melody Loop of Let the Show Begin♪ is extended by 1s.
// The maximum extension is 5s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnParticleReceived, func(_ ...interface{}) bool {
		if c.Core.Status.Duration(barbSkillKey) == 0 {
			return false
		}
		if c.a4extendCount == 5 {
			return false
		}

		c.a4extendCount++
		c.Core.Status.Extend(barbSkillKey, 60)

		c.Core.Log.NewEvent("barbara skill extended from a4", glog.LogCharacterEvent, c.Index)

		return false
	}, "barbara-a4")
}
