package tighnari

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

var burstHitmarks = []int{112, 117, 120, 121, 126, 128}
var burstSecondHitmarks = []int{147, 153, 160, 161, 171, 175}

func init() {
	burstFrames = frames.InitAbilSlice(77)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 0
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Tanglevine Shaft",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	for i := 0; i < 6; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
			0,
			burstHitmarks[i]+travel,
		)
	}

	ai.Abil = "Secondary Tanglevine Shaft"
	ai.Mult = burstSecond[c.TalentLvlBurst()]
	for i := 0; i < 6; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
			0,
			burstSecondHitmarks[i]+travel,
		)
	}

	c.ConsumeEnergy(7)
	c.SetCD(action.ActionBurst, 12*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}
}
