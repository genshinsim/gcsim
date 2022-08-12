package aloy

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 100

func init() {
	burstFrames = frames.InitAbilSlice(117) // Q -> Swap
	burstFrames[action.ActionAttack] = 102  // Q -> N1
	burstFrames[action.ActionAim] = 102     // Q -> Aim, assumed because it's most likely not 117
	burstFrames[action.ActionSkill] = 102   // Q -> E
	burstFrames[action.ActionDash] = 101    // Q -> D
	burstFrames[action.ActionJump] = 101    // Q -> J
}

// Burst - doesn't do much other than damage, so fairly straightforward
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// snapshots before or during Burst Animation
	// https://library.keqingmains.com/evidence/characters/cryo/aloy#burst-mechanics
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
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), burstHitmark)

	c.SetCD(action.ActionBurst, 12*60)
	c.ConsumeEnergy(2)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}
