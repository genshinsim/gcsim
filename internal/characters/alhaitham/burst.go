package alhaitham

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 94

func init() {
	burstFrames = frames.InitAbilSlice(91) // Q -> J
	burstFrames[action.ActionAttack] = 88  // Q -> N1
	burstFrames[action.ActionSkill] = 89   // Q -> E
	burstFrames[action.ActionDash] = 89    // Q -> Dash
	burstFrames[action.ActionWalk] = 89    // Q -> Walk
	burstFrames[action.ActionSwap] = 87    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Particular Field: Fetters of Phenomena",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstAtk[c.TalentLvlBurst()],
		FlatDmg:    burstEm[c.TalentLvlSkill()] * c.Stat(attributes.EM),
	}

	//X number of hits depending on mirrors when casted
	for i := 0; i < 4+2*c.mirrorCount; i++ {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 7.1}, 6.8), 67, burstHitmark+i*21)

	}
	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 18*60)

	consumed := c.mirrorCount
	generated := 3 - c.mirrorCount
	hasC4 := c.Base.Cons >= 4
	if c.Base.Cons >= 6 {
		generated = 3
	}

	c.mirrorLoss(c.lastInfusionSrc, consumed)() // consume mirrors right away
	if hasC4 {
		c.c4Loss(consumed)
	}

	c.QueueCharTask(func() {
		if c.Core.Player.Active() != c.Index {
			return
		}
		c.mirrorGain(generated)
		if hasC4 {
			c.c4Gain(generated)
		}
	}, 190)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
