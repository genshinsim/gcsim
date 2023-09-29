package kuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstStart = 50

func init() {
	burstFrames = frames.InitAbilSlice(63) // Q -> D/J
	burstFrames[action.ActionAttack] = 62  // Q -> N1
	burstFrames[action.ActionSkill] = 62   // Q -> E
	burstFrames[action.ActionSwap] = 62    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gyoei Narukami Kariyama Rite",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    c.MaxHP() * burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	count := 7 // can be 11 at low HP
	if c.CurrentHPRatio() <= 0.5 {
		count = 12
	}
	interval := 2 * 60 / 7

	// C1: Gyoei Narukami Kariyama Rite's AoE is increased by 50%.
	r := 4.0
	if c.Base.Cons >= 1 {
		r = 6
	}

	// assume that the target is close enough to make her Q center on it
	for i := burstStart; i < count*interval+burstStart; i += interval {
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, r),
			i,
		)
	}

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 900) // 15s * 60

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}, nil
}
