package keqing

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 56

func init() {
	burstFrames = frames.InitAbilSlice(124)
	burstFrames[action.ActionDash] = 122
	burstFrames[action.ActionSwap] = 123
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//first hit 56 frame
	//first tick 82 frame
	//last tick 162
	//last hit 197

	// trigger a4
	c.a4()

	//initial
	ai := combat.AttackInfo{
		Abil:       "Starward Sword (Initial)",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstInitial[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	//8 hits
	ai.Abil = "Starward Sword (Consecutive Slash)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	for i := 82; i < 162; i += 11 {
		c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), i, i)
	}

	//final
	ai.Abil = "Starward Sword (Last Attack)"
	ai.Mult = burstFinal[c.TalentLvlBurst()]
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 197, 197)

	if c.Base.Cons >= 6 {
		c.c6("burst")
	}

	c.ConsumeEnergy(55)
	c.SetCDWithDelay(action.ActionBurst, 720, 52)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		Post:            burstHitmark,
		State:           action.BurstState,
	}
}
