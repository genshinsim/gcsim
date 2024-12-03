package chasca

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int
var burstSecondaryHitmark = []int{104, 130, 150, 155, 160, 165} // 2nd, 4th, and 6th hitmarks are unknown

const (
	burstHitmark = 99
)

func init() {
	burstFrames = frames.InitAbilSlice(71)
	burstFrames[action.ActionAttack] = 61 // Q -> N1
	burstFrames[action.ActionAim] = 61    // Q -> Aim
	burstFrames[action.ActionSkill] = 61  // Q -> E
	burstFrames[action.ActionDash] = 61   // Q -> D
	burstFrames[action.ActionJump] = 63   // Q -> J
	burstFrames[action.ActionSwap] = 60   // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Galesplitting Soulseeker Shell",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           burstGalesplitting[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.0),
		burstHitmark,
		burstHitmark,
	)

	ai.ICDTag = attacks.ICDTagChascaBurst
	ai.ICDGroup = attacks.ICDGroupChascaBurst

	// TODO: Is it the anemo ones first/last or is it random

	// the anemo ones
	for i := 0; i < 6; i++ {
		switch {
		case i < c.phecCount*2:
			ele := c.partyTypes[c.Core.Rand.Intn(len(c.partyTypes))]
			ai.Abil = fmt.Sprintf("Shining Soulseeker Shell (%s)", ele.String())
			ai.Mult = burstSoulseeker[c.TalentLvlBurst()]
			ai.Element = ele
			c.Core.QueueAttack(ai,
				combat.NewSingleTargetHit(c.Core.Combat.DefaultTarget),
				burstSecondaryHitmark[i],
				burstSecondaryHitmark[i])
		default:
			ai.Abil = "Soulseeker Shell"
			ai.Mult = burstSoulseeker[c.TalentLvlBurst()]
			ai.Element = attributes.Anemo
			c.Core.QueueAttack(ai,
				combat.NewSingleTargetHit(c.Core.Combat.DefaultTarget),
				burstSecondaryHitmark[i],
				burstSecondaryHitmark[i])
		}

	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(22)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
