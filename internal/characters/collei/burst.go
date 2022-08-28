package collei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const (
	burstTickRate = 27 // TODO: find tick rate
	burstHitmark  = 58 // TODO: actual frames
	burstKey      = "collei-burst"
)

func init() {
	burstFrames = frames.InitAbilSlice(64)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Trump-Card Kitty (Explosion)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst, // TODO: find ICD
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstExplosion[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy),
		0,
		burstHitmark,
	) // TODO: snapshot timing

	c.AddStatus(burstKey, 360, false)

	snap := c.Snapshot(&ai) // TODO: snapshot timing
	c.Core.Tasks.Add(func() {
		c.burstTicks(snap)
	}, burstHitmark+burstTickRate)

	c.c4() // TODO: figure out c4 delay
	c.SetCDWithDelay(action.ActionBurst, 900, 41) // TODO: find cooldown delay
	c.ConsumeEnergy(43)                           // TODO: find energy consumption delay

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark, // TODO: correct cancel frame
		State:           action.BurstState,
	}
}

func (c *char) burstTicks(snap combat.Snapshot) {
	if !c.StatusIsActive(burstKey) {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Trump-Card Kitty (Leap)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst, // TODO: find ICD
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstLeap[c.TalentLvlBurst()],
	}
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy),
		0,
	)
	c.Core.Tasks.Add(func() {
		c.burstTicks(snap)
	}, burstTickRate)
}
