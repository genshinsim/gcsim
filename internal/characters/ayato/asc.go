package ayato

import (
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// Kamisato Art: Kyouka has the following properties:
//
// - After it is used, Kamisato Ayato will gain 2 Namisen stacks.
func (c *char) a1OnSkill() {
	if c.Base.Ascension < 1 {
		return
	}
	c.stacks = 2
	c.Core.Log.NewEvent("ayato a1 proc'd", glog.LogCharacterEvent, c.Index)
}

// Kamisato Art: Kyouka has the following properties:
//
// - When the water illusion explodes, Ayato will gain a Namisen effect equal to the maximum number of stacks possible.
func (c *char) a1OnExplosion() {
	if c.Base.Ascension < 1 {
		return
	}
	c.stacks = c.stacksMax
	c.Core.Log.NewEvent("ayato a1 set namisen stacks to max", glog.LogCharacterEvent, c.Index).
		Write("stacks", c.stacks)
}

// If Kamisato Ayato is not on the field and his Energy is less than 40, he will regenerate 2 Energy for himself every second.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	if c.Core.Player.Active() == c.Index {
		return
	}
	if c.Energy >= 40 {
		return
	}
	c.AddEnergy("ayato-a4", 2)
	c.Core.Tasks.Add(c.a4, 60)
}
