package yaemiko

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 100
const burstThunderbolt1Hitmark = 154

func init() {
	burstFrames = frames.InitAbilSlice(114) // Q -> CA
	burstFrames[action.ActionAttack] = 112  // Q -> N1
	burstFrames[action.ActionSkill] = 113   // Q -> E
	burstFrames[action.ActionDash] = 103    // Q -> D
	burstFrames[action.ActionJump] = 104    // Q -> J
	burstFrames[action.ActionSwap] = 101    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Great Secret Art: Tenko Kenshin",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[0][c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 7),
		burstHitmark,
		burstHitmark,
	)

	ai.Abil = "Tenko Thunderbolt"
	ai.Mult = burst[1][c.TalentLvlBurst()]
	c.kitsuneBurst(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 7))

	c.ConsumeEnergy(2)
	c.SetCD(action.ActionBurst, 22*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
