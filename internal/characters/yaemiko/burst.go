package yaemiko

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 94

func init() {
	burstFrames = frames.InitAbilSlice(111)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Great Secret Art: Tenko Kenshin",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[0][c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy, combat.TargettableObject), burstHitmark, burstHitmark)

	ai.Abil = "Tenko Thunderbolt"
	ai.Mult = burst[1][c.TalentLvlSkill()]
	c.kitsuneBurst(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy))

	c.ConsumeEnergy(7)
	c.SetCD(action.ActionBurst, 22*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		Post:            burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
