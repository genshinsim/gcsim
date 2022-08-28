package collei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const (
	explosionHitmark = 25
	leapHitmark      = 68
	leapTickPeriod   = 30
	fieldStart       = 43
	burstKey         = "collei-burst"
)

func init() {
	burstFrames = frames.InitAbilSlice(67)
	burstFrames[action.ActionAttack] = 65
	burstFrames[action.ActionAim] = 65
	burstFrames[action.ActionSwap] = 66
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
		explosionHitmark,
		explosionHitmark,
	)

	c.Core.Tasks.Add(func() {
		c.AddStatus(burstKey, 360, false)
		snap := c.Snapshot(&ai)
		c.Core.Tasks.Add(func() {
			c.burstTicks(snap)
		}, leapHitmark-fieldStart)
		c.burstA4Ticks()
	}, fieldStart)

	c.c4()
	c.SetCD(action.ActionBurst, 900)
	c.ConsumeEnergy(7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
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
	}, leapTickPeriod)
}

func (c *char) burstA4Ticks() {
	// TODO: add range check within field
	if !c.StatusIsActive(burstKey) {
		return
	}
	c.Core.Player.ActiveChar().AddStatus(a4Key, 60, true)
	c.Core.Tasks.Add(func() { c.burstA4Ticks() }, 30)
}
