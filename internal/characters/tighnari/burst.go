package tighnari

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstRelease = 77

var burstHitmarks = []int{112, 117, 120, 121, 126, 128}
var burstSecondHitmarks = []int{147, 153, 160, 161, 171, 175}

func init() {
	burstFrames = frames.InitAbilSlice(118)
	burstFrames[action.ActionAttack] = 114
	burstFrames[action.ActionAim] = 114
	burstFrames[action.ActionSkill] = 117
	burstFrames[action.ActionDash] = 117
	burstFrames[action.ActionSwap] = 115
}

func (c *char) Burst(p map[string]int) action.Info {
	travel, ok := p["travel"]
	if !ok {
		travel = 0
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Tanglevine Shaft",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 1)

	for i := 0; i < 6; i++ {
		c.Core.QueueAttack(ai, ap, burstRelease, burstHitmarks[i]+travel)
	}

	ai.Abil = "Secondary Tanglevine Shaft"
	ai.Mult = burstSecond[c.TalentLvlBurst()]
	for i := 0; i < 6; i++ {
		c.Core.QueueAttack(ai, ap, burstHitmarks[i]+travel, burstSecondHitmarks[i]+travel)
	}

	c.ConsumeEnergy(7)
	c.SetCD(action.ActionBurst, 12*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}
}
