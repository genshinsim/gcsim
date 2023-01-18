package collei

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	sproutKey        = "collei-a1"
	sproutHitmark    = 86
	sproutTickPeriod = 89
	a4Key            = "collei-a4-modcheck"
)

func (c *char) a1() {
	for _, event := range dendroEvents {
		c.Core.Events.Subscribe(event, func(args ...interface{}) bool {
			if c.sproutShouldProc {
				return false
			}
			if !c.StatusIsActive(skillKey) {
				return false
			}
			c.sproutShouldProc = true
			c.Core.Log.NewEvent("collei a1 proc", glog.LogCharacterEvent, c.Index)
			return false
		}, "collei-a1")
	}
}

func (c *char) a1AttackInfo() combat.AttackInfo {
	return combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Floral Sidewinder (A1)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagColleiSprout,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       0.4,
	}
}

func (c *char) a4() {
	for _, event := range dendroEvents {
		c.Core.Events.Subscribe(event, func(args ...interface{}) bool {
			if !c.StatusIsActive(burstKey) {
				return false
			}
			atk := args[1].(*combat.AttackEvent)
			char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
			if !char.StatusIsActive(a4Key) {
				return false
			}
			if c.burstExtendCount >= 3 {
				return false
			}
			c.ExtendStatus(burstKey, 60)
			c.burstExtendCount++
			c.Core.Log.NewEvent("collei a4 proc", glog.LogCharacterEvent, c.Index).
				Write("extend_count", c.burstExtendCount)
			return false
		}, "collei-a4")
	}
}

func (c *char) a1Ticks(startFrame int, snap combat.Snapshot) {
	if !c.StatusIsActive(sproutKey) {
		return
	}
	if startFrame != c.sproutSrc {
		c.Core.Log.NewEvent("collei a1 tick ignored, src diff", glog.LogCharacterEvent, c.Index).
			Write("src", startFrame).
			Write("new src", c.sproutSrc)
		return
	}
	c.Core.QueueAttackWithSnap(
		c.a1AttackInfo(),
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2),
		0,
	)
	c.Core.Tasks.Add(func() {
		c.a1Ticks(startFrame, snap)
	}, sproutTickPeriod)
}
