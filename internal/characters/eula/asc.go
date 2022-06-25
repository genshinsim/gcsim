package eula

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// A4: When Glacial Illumination is cast, the CD of Icetide Vortex is reset and Eula gains 1 stack of Grimheart.
func (c *char) a4() {
	c.Core.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}

		v := c.Tags["grimheart"]
		if v < 2 {
			v++
		}
		c.Tags["grimheart"] = v

		c.ResetActionCooldown(action.ActionSkill)
		c.Core.Log.NewEvent("eula a4 reset skill cd", glog.LogCharacterEvent, c.Index)

		return false
	}, "eula-a4")
}
