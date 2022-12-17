package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFramesNormal []int
var burstFramesE []int

func init() {
	burstFramesNormal = frames.InitAbilSlice(94)
	burstFramesNormal[action.ActionCharge] = 96
	burstFramesNormal[action.ActionSkill] = 95
	burstFramesNormal[action.ActionDash] = 97
	burstFramesNormal[action.ActionJump] = 96
	burstFramesNormal[action.ActionWalk] = 101
	burstFramesNormal[action.ActionSwap] = 94

	burstFramesE = frames.InitAbilSlice(117)
	burstFramesE[action.ActionCharge] = 119
	burstFramesE[action.ActionDash] = 119
	burstFramesE[action.ActionJump] = 119
	// Includes Falling down
	burstFramesE[action.ActionSwap] = 145
}

// First Hitmark
const burstHitmark = 92

// Delay between each additional hit
const burstHitmarkDelay = 6

// Frames until snapshot stage is reached
// TODO: Determine correct Frame
const burstSnapshotDelay = 55

func (c *char) Burst(p map[string]int) action.ActionInfo {

	delay := c.checkForSkillEnd()

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Kyougen: Five Ceremonial Plays",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagElementalBurst,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Anemo,
		Durability:         25,
		Mult:               burst[c.TalentLvlBurst()],
		CanBeDefenseHalted: false,
	}

	c.Core.Tasks.Add(c.c2, delay)

	for i := 0; i < 5; i++ {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 6),
				burstSnapshotDelay, burstHitmark+i*burstHitmarkDelay)
		}, delay)
	}

	//TODO: Check CD with or without delay, check energy consume frame
	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(5)

	// End Windfavored State after burst
	c.Core.Tasks.Add(func() {
		c.skydwellerPoints = 0
	}, 89)

	relevantFrames := burstFramesNormal
	if delay == 0 {
		relevantFrames = burstFramesE
	}

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return delay + relevantFrames[next] },
		AnimationLength: delay + relevantFrames[action.InvalidAction],
		CanQueueAfter:   delay + relevantFrames[action.ActionAttack],
		State:           action.BurstState,
		OnRemoved: func(next action.AnimationState) {
			if next == action.SwapState {
				c.skillEndRoutine()
			}
		},
	}
}
