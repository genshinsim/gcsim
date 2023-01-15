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
	burstFramesNormal = frames.InitAbilSlice(101)
	burstFramesNormal[action.ActionAttack] = 94
	burstFramesNormal[action.ActionCharge] = 96
	burstFramesNormal[action.ActionSkill] = 95
	burstFramesNormal[action.ActionDash] = 97
	burstFramesNormal[action.ActionJump] = 96
	burstFramesNormal[action.ActionSwap] = 94

	// Includes Falling down for swap
	burstFramesE = frames.InitAbilSlice(145)
	burstFramesE[action.ActionAttack] = 117
	burstFramesE[action.ActionCharge] = 119
	burstFramesE[action.ActionDash] = 119
	burstFramesE[action.ActionJump] = 119
	burstFramesE[action.ActionWalk] = 117

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

	if c.StatusIsActive(skillKey) {
		// Can only occur if delay == 0, so it can be disregarded
		return c.WindfavoredBurst(p)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Kyougen: Five Ceremonial Plays",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.Tasks.Add(c.c2, delay)

	for i := 0; i < 5; i++ {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5),
			burstSnapshotDelay, burstHitmark+i*burstHitmarkDelay)
	}

	//TODO: Check CD with or without delay, check energy consume frame

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(5)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return delay + burstFramesNormal[next] },
		AnimationLength: delay + burstFramesNormal[action.InvalidAction],
		CanQueueAfter:   delay + burstFramesNormal[action.ActionAttack],
		State:           action.BurstState,
	}
}

func (c *char) WindfavoredBurst(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Kyougen: Five Ceremonial Plays (Windfavored)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.c2()

	for i := 0; i < 5; i++ {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5),
		burstSnapshotDelay, burstHitmark+i*burstHitmarkDelay)
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(5)

	// Necessary, as transitioning into the SwapState is impossible otherwise
	c.Core.Player.SwapCD = 26

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return burstFramesE[next] },
		AnimationLength: burstFramesE[action.InvalidAction],
		CanQueueAfter:   burstFramesE[action.ActionWalk],
		State:           action.BurstState,
		OnRemoved: func(next action.AnimationState) {
			c.skydwellerPoints = 0
			if next == action.SwapState {
				c.checkForSkillEnd()
			}
		},
	}
}
