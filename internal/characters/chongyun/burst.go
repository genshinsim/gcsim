package chongyun

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

var burstHitmarks = []int{50, 59, 67}

const burstHitmarkC6 = 77

func init() {
	burstFrames = frames.InitAbilSlice(79) // Q -> Swap
	burstFrames[action.ActionAttack] = 64  // Q -> N1
	burstFrames[action.ActionSkill] = 64   // Q -> E
	burstFrames[action.ActionDash] = 64    // Q -> D
	burstFrames[action.ActionJump] = 66    // Q -> J
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Blade: Cloud-Parting Star",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   100,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c4CB := c.makeC4Callback()

	// Spirit Blade 1-3
	for _, hitmark := range burstHitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5),
			hitmark,
			hitmark,
			c4CB,
		)
	}

	// extra Spirit Blade at C6
	if c.Base.Cons >= 6 {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5),
			burstHitmarkC6,
			burstHitmarkC6,
			c4CB,
		)
	}

	c.SetCD(action.ActionBurst, 720)
	c.ConsumeEnergy(6)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
