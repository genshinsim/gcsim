package ayato

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// A1:
// Kamisato Art: Kyouka has the following properties:
// - After it is used, Kamisato Ayato will gain 2 Namisen stacks. (handled here)
// - When the water illusion explodes, Ayato will gain a Namisen effect equal to the maximum number of stacks possible. (handled in skill.go)
func (c *char) a1() {
	//TODO: this used to be PostSkill; check if working correctly still
	c.Core.Events.Subscribe(event.OnSkill, func(_ ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}
		c.stacks = 2
		c.Core.Log.NewEvent("ayato a1 proc'd", glog.LogCharacterEvent, c.Index)
		return false
	}, "ayato-a1")
}

// A4:
// If Kamisato Ayato is not on the field and his Energy is less than 40, he will regenerate 2 Energy for himself every second.
func (c *char) a4() {
	if c.Core.Player.Active() == c.Index {
		return
	}
	if c.Energy >= 40 {
		return
	}
	c.AddEnergy("ayato-a4", 2)
	c.Core.Tasks.Add(c.a4, 60)
}
