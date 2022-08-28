package collei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 58 // TODO: actual frames

func init() {
	burstFrames = frames.InitAbilSlice(64)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// Initial Hit
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

	ai.Abil = "Trump-Card Kitty (Leap)"
	ai.Mult = burstLeap[c.TalentLvlBurst()]
	for i := 0; i < 12; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy),
			0,       // TODO: snapshot timing
			27+i*27, // TODO: burst hitmarks
		)
	}

	c.SetCDWithDelay(action.ActionBurst, 900, 41) // TODO: find cooldown delay
	c.ConsumeEnergy(43)                           // TODO: find energy consumption delay

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark, // TODO: correct cancel frame
		State:           action.BurstState,
	}
}
