package kuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstStart = 78

func init() {
	burstFrames = frames.InitAbilSlice(78)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gyoei Narukami Kariyama Rite",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    c.MaxHP() * burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	count := 7 //can be 11 at low HP
	if (c.HPCurrent / c.MaxHP()) <= 0.5 {
		count = 12
	}
	interval := 2 * 60 / 7

	// C1: Gyoei Narukami Kariyama Rite's AoE is increased by 50%.
	var r float64 = 2
	if c.Base.Cons >= 1 {
		r = 3.5
	}

	for i := 0; i < count*interval; i += interval {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(r, false, combat.TargettableEnemy), i)
	}

	c.ConsumeEnergy(55) //TODO: Check if she can be pre-funneled
	c.SetCDWithDelay(action.ActionBurst, 900, 55)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.InvalidAction],
		State:           action.BurstState,
	}
}
