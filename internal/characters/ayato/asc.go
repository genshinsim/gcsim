package ayato

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

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

func (c *char) a4() {
	c.Core.Tasks.Add(c.a4task, 60)
}

func (c *char) a4task() {
	if c.Core.Player.Active() == c.Index {
		return
	}
	if c.StatusIsActive(a4ICDKey) {
		return
	}
	if c.Energy >= 40 {
		return
	}
	c.AddEnergy("ayato-a4", 2)
	c.Core.Tasks.Add(c.a4task, 60)
	c.AddStatus(a4ICDKey, 60, true)
}
