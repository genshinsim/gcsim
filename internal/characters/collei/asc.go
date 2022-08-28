package collei

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) a4() {
	dendroEvents := []event.Event{event.OnOverload} // TODO: put all dendro events here
	for _, event := range dendroEvents {
		c.Core.Events.Subscribe(event, func(args ...interface{}) bool {
			if !c.StatusIsActive(burstKey) {
				return false
			}
			if c.burstExtendCount >= 3 {
				return false
			}
			// TODO: check for increment ICD
			c.ExtendStatus(burstKey, 60)
			c.burstExtendCount++
			c.Core.Log.NewEvent("collei a4 proc", glog.LogCharacterEvent, c.Index).
				Write("extend_count", c.burstExtendCount)
			return false
		}, "collei-a4")
	}
}
