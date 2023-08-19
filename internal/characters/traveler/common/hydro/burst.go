package hydro

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// TODO: anemomc based frames
var burstHitmarks = []int{96, 94}
var burstFrames [][]int

func init() {
	burstFrames = make([][]int, 2)

	// Male
	burstFrames[0] = frames.InitAbilSlice(110) // Q -> N1
	burstFrames[0][action.ActionSkill] = 109   // Q -> E
	burstFrames[0][action.ActionDash] = 96     // Q -> D
	burstFrames[0][action.ActionJump] = 96     // Q -> J
	burstFrames[0][action.ActionSwap] = 100    // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(105) // Q -> N1
	burstFrames[1][action.ActionSkill] = 104   // Q -> E
	burstFrames[1][action.ActionDash] = 90     // Q -> D
	burstFrames[1][action.ActionJump] = 90     // Q -> J
	burstFrames[1][action.ActionSwap] = 95     // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	duration := burstHitmarks[c.gender] + 4*60

	c.Core.Status.Add("hmcburst", duration)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rising Waters",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst, // or the new one?
		ICDGroup:   attacks.ICDGroupTravelerBurst,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 0.15)
	snap := c.Snapshot(&ai)

	for i := 0; i < 8; i++ {
		c.Core.QueueAttackWithSnap(ai, snap, ap, 94+30*i)
	}

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(3)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
