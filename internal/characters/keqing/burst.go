package keqing

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 56

func init() {
	burstFrames = frames.InitAbilSlice(124) // Q -> D/J
	burstFrames[action.ActionAttack] = 123  // Q -> N1
	burstFrames[action.ActionSkill] = 123   // Q -> E
	burstFrames[action.ActionSwap] = 122    // Q -> Swap
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
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstInitial[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 8)
	c.Core.QueueAttack(ai, ap, burstHitmark, burstHitmark)

	//8 hits
	ai.Abil = "Starward Sword (Consecutive Slash)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	for i := 82; i < 162; i += 11 {
		c.Core.QueueAttack(ai, ap, i, i)
	}

	//final
	ai.Abil = "Starward Sword (Last Attack)"
	ai.Mult = burstFinal[c.TalentLvlBurst()]
	c.Core.QueueAttack(ai, ap, 197, 197)

	if c.Base.Cons >= 6 {
		c.c6("burst")
	}

	c.ConsumeEnergy(55)
	c.SetCDWithDelay(action.ActionBurst, 720, 52)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
