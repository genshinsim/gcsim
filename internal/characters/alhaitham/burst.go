package alhaitham

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 56 //todo:scuff

func init() {
	burstFrames = frames.InitAbilSlice(90) // Q -> D/J
	burstFrames[action.ActionAttack] = 90  // Q -> N1
	burstFrames[action.ActionSkill] = 90   // Q -> E
	burstFrames[action.ActionSwap] = 90    // Q -> Swap
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
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 7.1}, 6.8), burstHitmark+i*5, burstHitmark+i*5)

	}
	c.ConsumeEnergy(55)                             //TODO can we prefunnel?
	c.SetCDWithDelay(action.ActionBurst, 18*60, 52) //TODO: Delay needed here?

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
	}, 2*60+burstFrames[0]) //TODO:exact refund timing
}
