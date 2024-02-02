package zhongli

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int

const burstHitmark = 101

func init() {
	burstFrames = frames.InitAbilSlice(139) // Q -> N1/E
	burstFrames[action.ActionDash] = 123    // Q -> D
	burstFrames[action.ActionJump] = 123    // Q -> J
	burstFrames[action.ActionSwap] = 138    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// deal damage when created
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Planet Befall",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   500,
		Element:    attributes.Geo,
		Durability: 100,
		Mult:       burst[c.TalentLvlBurst()],
		FlatDmg:    c.a4Burst(),
	}
	r := 7.5
	if c.Base.Cons >= 4 {
		r = 9
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 5}, r),
		burstHitmark,
		burstHitmark,
	)

	if c.Base.Cons >= 2 {
		c.addJadeShield()
	}

	c.SetCD(action.ActionBurst, 720)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
