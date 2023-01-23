package alhaitham

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 94

func init() {
	burstFrames = frames.InitAbilSlice(89) // Q -> D/W
	burstFrames[action.ActionAttack] = 88  // Q -> N1
	burstFrames[action.ActionSkill] = 89   // Q -> E
	burstFrames[action.ActionJump] = 91    // Q -> Jump
	burstFrames[action.ActionSwap] = 87    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Particular Field: Fetters of Phenomena",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstAtk[c.TalentLvlBurst()],
		FlatDmg:    burstEm[c.TalentLvlSkill()] * c.Stat(attributes.EM)}

	//X number of hits depending on mirrors when casted
	//TODO: does the number of mirror affects the length of the attacking animations?
	for i := 0; i < 4+2*c.mirrorCount; i++ { //TODO:frame counting for dis
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 7.1}, 6.8), burstHitmark+i*21, burstHitmark+i*21)

	}
	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 18*60)

	for i := 0; i < 3; i++ {
		if c.mirrorCount <= i {

			c.burstRefundMirrors()
			if c.Base.Cons >= 4 {
				c.c4("gain", i) //TODO: exact timing of c4 buff application
			}

		} else {
			c.Core.Tasks.Add(c.mirrorLoss(c.lastInfusionSrc), 0)
			if c.Base.Cons >= 4 { //TODO: Execution on cast or posburst?
				c.c4("loss", i)
			}
		}
		if c.Base.Cons >= 6 {
			c.burstRefundMirrors()
		}

	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstRefundMirrors() {
	c.Core.Tasks.Add(func() {
		if c.Core.Player.Active() == c.Index { //stacks are refunded as long as he is on field
			c.mirrorGain()
		}
	}, 191+burstFrames[0]) //TODO:exact refund timing
}
