package navia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var burstFrames []int

const burstHitmark = 100

func init() {
	burstFrames = frames.InitAbilSlice(114)
	burstFrames[action.ActionSkill] = 114
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "As the Sunlit Sky's Singing Salute",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       burst[0][c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6),
		burstHitmark,
		burstHitmark,
	)

	c.ConsumeEnergy(5)
	c.SetCD(action.ActionBurst, 15*60)

	ai.Abil = "Fire Support"
	ai.ICDTag = attacks.ICDTagElementalBurst
	ai.Durability = 25
	ai.Mult = burst[1][c.TalentLvlBurst()]

	for i := 45; i <= 12*60; i = i + 45 {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3),
			burstHitmark,
			burstHitmark+i,
			c.BurstCB(),
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) BurstCB() combat.AttackCBFunc {
	if c.StatusIsActive("navia-q-shrapnel-icd") {
		return nil
	}

	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}

		if c.shrapnel < 6 {
			c.shrapnel++
		}

		c.AddStatus("navia-q-shrapnel-icd", 2.4*60, false)
	}

}
