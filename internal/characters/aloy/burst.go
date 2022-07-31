package aloy

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 118

func init() {
	burstFrames = frames.InitAbilSlice(118)
}

// Burst - doesn't do much other than damage, so fairly straightforward
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// TODO: Assuming dynamic
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Prophecies of Dawn",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	c.SetCDWithDelay(action.ActionBurst, 12*60, 8)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.InvalidAction],
		State:           action.BurstState,
	}
}
