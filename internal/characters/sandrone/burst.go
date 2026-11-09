package sandrone

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames   []int
	burstHitmarks = []int{123, 123 + 12, 123 + 12 + 12}
	rayHitmark    = 123 + 12 + 12 + 30
)

func init() {
	burstFrames = frames.InitAbilSlice(128) // Q -> CA
	burstFrames[action.ActionAttack] = 101  // Q -> N1
	burstFrames[action.ActionSkill] = 100   // Q -> E
	burstFrames[action.ActionDash] = 103    // Q -> D
	burstFrames[action.ActionJump] = 103    // Q -> J
	burstFrames[action.ActionWalk] = 105    // Q -> W
	burstFrames[action.ActionSwap] = 102    // Q -> Swap
}

// In line with established ladylike behavior, Sandrone summons a large number of Prismatic
// Resonance Cannons for a frontal bombardment before firing on her foes with a Convective
// Inhibition Ray, dealing AoE Cryo DMG.
// Radiance: Stellar Glimmer: The Convective Inhibition Ray instead deals AoE Cryo DMG that is
// considered the corresponding Stellar Glimmer reaction DMG.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Bombardment",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.PrimaryTarget(),
		nil,
		7.3,
	)

	for _, delay := range burstHitmarks {
		c.QueueCharTask(func() { c.Core.QueueAttack(ai, ap, 0, 0) }, delay)
	}

	c.QueueCharTask(func() {
		apRay := combat.NewCircleHitOnTarget(
			c.Core.Combat.PrimaryTarget(),
			nil,
			7.5,
		)
		aiRay := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Convective Inhibition Ray",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   70,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       burstRay[c.TalentLvlBurst()],
		}

		bonusMult := 1 + c.a1OnBurstRayStellar()
		dmgSSC := burstRaySSC[c.TalentLvlBurst()] * bonusMult
		dmgSSw := burstRaySSw[c.TalentLvlBurst()] * bonusMult
		c.toStellar(&aiRay, dmgSSC, dmgSSw)
		if aiRay.AttackTag != attacks.AttackTagElementalBurst {
			aiRay.PoiseDMG = 100
		}

		c.Core.QueueAttack(aiRay, apRay, 0, 0)
	}, rayHitmark)

	c.ConsumeEnergy(7)
	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSkill], // earliest cancel
		State:           action.BurstState,
	}, nil
}
