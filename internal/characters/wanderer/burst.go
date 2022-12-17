package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(72)
	burstFrames[action.ActionAttack] = 71
	burstFrames[action.ActionSkill] = 71
	burstFrames[action.ActionJump] = 70
	burstFrames[action.ActionSwap] = 69

}

// First Hitmark
const burstHitmark = 34

// Delay between each additional hit
const burstHitmarkDelay = 1

// Frames until snapshot stage is reached
const burstSnapshotDelay = 30

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

	//TODO: when does wanderer burst snapshot?
	for i := 0; i < 5; i++ {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 6),
				burstSnapshotDelay, burstHitmark+i*burstHitmarkDelay)
		}, delay)
	}

	//TODO: Check CD with or without delay, check energy consume frame
	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(3)
	c.skydwellerPoints = 0

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return delay + burstFrames[next] },
		AnimationLength: delay + burstFrames[action.InvalidAction],
		CanQueueAfter:   delay + burstFrames[action.ActionSwap],
		State:           action.BurstState,
		OnRemoved: func(next action.AnimationState) {
			if next == action.SwapState {
				c.skillEndRoutine()
			}
		},
	}
}
