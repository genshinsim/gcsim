package prune

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const burstKey = "prune-burst"

func init() {
	burstFrames = frames.InitAbilSlice(73) // walk
	burstFrames[action.ActionAttack] = 73
	burstFrames[action.ActionCharge] = 82
	burstFrames[action.ActionSkill] = 72
	burstFrames[action.ActionDash] = 74
	burstFrames[action.ActionJump] = 73
	burstFrames[action.ActionSwap] = 71
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	duration := 813 + c.c6BurstBonusDur()

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "The Bell Tolls! The Hunt Is On!",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5),
		34,
		34,
	)

	// burst status delay based on dm calculation of burst status duration, assume to be affected by hitlag
	burstStatusDelay := 57
	c.QueueCharTask(func() {
		c.AddStatus(burstKey, duration-burstStatusDelay, false)
		c.c2OnBurst(duration - burstStatusDelay)
	}, burstStatusDelay)

	c.burstSrc = c.Core.F
	c.Core.Tasks.Add(c.burstTick(c.burstSrc), 137)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(18) // average of 6 trials

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTick(src int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Witchlure Bell",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       burstDot[c.TalentLvlBurst()],
		}

		detectionArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10)
		target := c.Core.Combat.ClosestEnemyWithinArea(detectionArea, nil)

		if target == nil {
			return
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(target, nil, 1.6),
			0,
			0,
			c.makeC2CB,
		)
		c.Core.Tasks.Add(c.burstTick(src), 117)
	}
}
