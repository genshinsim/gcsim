package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 104

func init() {
	burstFrames = frames.InitAbilSlice(125) // Q -> D
	burstFrames[action.ActionAttack] = 124  // Q -> N1
	burstFrames[action.ActionSkill] = 124   // Q -> E
	burstFrames[action.ActionJump] = 113    // Q -> J
	burstFrames[action.ActionSwap] = 123    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Soumetsu",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 25,
	}

	//5 second, 20 ticks, so once every 15 frames, bloom after 5 seconds
	ai.Mult = burstBloom[c.TalentLvlBurst()]
	ai.StrikeType = attacks.StrikeTypeDefault
	ai.Abil = "Soumetsu (Bloom)"
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			5,
		),
		burstHitmark,
		burstHitmark+300,
		c.c4,
	)

	// C2 mini-frostflake bloom
	var aiC2 combat.AttackInfo
	if c.Base.Cons >= 2 {
		aiC2 = ai
		aiC2.Mult = burstBloom[c.TalentLvlBurst()] * .2
		aiC2.Abil = "C2 Mini-Frostflake Seki no To (Bloom)"
		// TODO: Not sure about the positioning/size...
		for i := 0; i < 2; i++ {
			c.Core.QueueAttack(
				aiC2,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					3,
				),
				burstHitmark,
				burstHitmark+300,
				c.c4,
			)
		}
	}

	for i := 0; i < 19; i++ {
		ai.Mult = burstCut[c.TalentLvlBurst()]
		ai.StrikeType = attacks.StrikeTypeSlash
		ai.Abil = "Soumetsu (Cutting)"
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				combat.Point{Y: 0.3},
				3,
			),
			burstHitmark,
			burstHitmark+i*15,
			c.c4,
		)

		// C2 mini-frostflake cutting
		if c.Base.Cons >= 2 {
			aiC2.Mult = burstCut[c.TalentLvlBurst()] * .2
			aiC2.StrikeType = attacks.StrikeTypeSlash
			aiC2.Abil = "C2 Mini-Frostflake Seki no To (Cutting)"
			// TODO: Not sure about the positioning/size...
			for j := 0; j < 2; j++ {
				c.Core.QueueAttack(
					aiC2,
					combat.NewCircleHit(
						c.Core.Combat.Player(),
						c.Core.Combat.PrimaryTarget(),
						combat.Point{Y: 0.3},
						1.5,
					),
					burstHitmark,
					burstHitmark+i*15,
					c.c4,
				)
			}
		}
	}

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 20*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}
