package xinyan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstInitialHitmark = 22
const burstShieldStart = 43
const burstDoT1Hitmark = 57

func init() {
	burstFrames = frames.InitAbilSlice(87) // Q -> E/D/J
	burstFrames[action.ActionAttack] = 86  // Q -> N1
	burstFrames[action.ActionSwap] = 86    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Riff Revolution",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Physical,
		Durability:         100,
		Mult:               burstDmg[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 3), burstInitialHitmark, burstInitialHitmark)

	// 7 hits
	ai = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Riff Revolution (DoT)",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagElementalBurstPyro,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               burstDot[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
	}
	// 1st DoT
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 4), 0, 0)
		ai.CanBeDefenseHalted = false // only the first DoT has hitlag
		// 2nd DoT onwards
		c.QueueCharTask(func() {
			for i := 0; i < 6; i++ {
				c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 4), i*17, i*17)
			}
		}, 17)
	}, burstDoT1Hitmark)

	if c.Base.Cons >= 2 {
		// TODO: snapshot timing?
		stats, _ := c.Stats()
		defFactor := c.Base.Def*(1+stats[attributes.DEFP]) + stats[attributes.DEF]
		c.QueueCharTask(func() {
			c.updateShield(3, defFactor)
		}, burstShieldStart)
	}

	c.ConsumeEnergy(5)
	c.SetCDWithDelay(action.ActionBurst, 15*60, 1)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}
}
