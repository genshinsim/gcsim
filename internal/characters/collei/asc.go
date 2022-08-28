package collei

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var dendroEvents = []event.Event{event.OnOverload} // TODO: put all dendro events here

const (
	a1Key        = "collei-sprout"
	a1Hitmark    = 86
	a1TickPeriod = 89
)

func (c *char) a1Init() {
	for _, event := range dendroEvents {
		c.Core.Events.Subscribe(event, func(args ...interface{}) bool {
			if !c.StatusIsActive(skillKey) {
				return false
			}
			c.a1Triggered = true
			c.Core.Log.NewEvent("collei a1 proc", glog.LogCharacterEvent, c.Index)
			return false
		}, "collei-a1")
	}
}

func (c *char) a4() {
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

func (c *char) a1Ticks(startFrame int) {
	if !c.StatusIsActive(a1Key) {
		return
	}
	if startFrame != c.a1StartFrame {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Floral Sidewinder",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone, // TODO: find ICD
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       0.4,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy),
		0,
		0,
	)
	c.Core.Tasks.Add(func() {
		c.a1Ticks(startFrame)
	}, a1TickPeriod)
}
