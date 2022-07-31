package xinyan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 98

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Riff Revolution",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Physical,
		Durability:         100,
		Mult:               burstDmg[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), 28, 28)

	// 7 hits
	ai = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Riff Revolution (DoT)",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagElementalBurstPyro,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               burstDot[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
	}
	for i := 0; i < 7; i++ {
		f := 63 + i*17
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), f, f)
		if i == 0 {
			ai.CanBeDefenseHalted = false // only the first DoT has hitlag
		}
	}

	if c.Base.Cons >= 2 {
		stats, _ := c.Stats()
		defFactor := c.Base.Def*(1+stats[attributes.DEFP]) + stats[attributes.DEF]
		c.updateShield(3, defFactor)
	}

	c.ConsumeEnergy(11)
	c.SetCDWithDelay(action.ActionBurst, 15*60, 3)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
