package chongyun

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 79

func init() {
	burstFrames = frames.InitAbilSlice(79)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Blade: Cloud-Parting Star",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 50, 50)
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 57, 57)
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 65, 65)

	if c.Base.Cons >= 6 {
		c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 76, 76)
	}

	c.SetCDWithDelay(action.ActionBurst, 720, 10)
	c.ConsumeEnergy(10)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,

		State: action.BurstState,
	}
}
