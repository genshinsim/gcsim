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
	burstFrames = frames.InitAbilSlice(123)
	burstFrames[action.ActionAttack] = 102 // Q -> N1
	burstFrames[action.ActionAim] = 103    // Q -> Aim
	burstFrames[action.ActionSkill] = 101  // Q -> E
	burstFrames[action.ActionDash] = 104   // Q -> D
	burstFrames[action.ActionJump] = 103   // Q -> J
	burstFrames[action.ActionWalk] = 118   // Q -> Walk
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

	var c4cb combat.AttackCBFunc

	burstBullets := make([]attributes.Element, 0, 6)
	burstBullets = append(burstBullets, c.partyPHECTypes...)
	burstBullets = append(burstBullets, c.partyPHECTypes...)
	c.Core.Rand.Shuffle(len(burstBullets), func(i, j int) {
		burstBullets[i], burstBullets[j] = burstBullets[j], burstBullets[i]
	})
	for i := 0; i < 6; i++ {
		switch {
		case i < len(burstBullets):
			ele := burstBullets[i]
			ai.Abil = fmt.Sprintf("Radiant Soulseeker Shell (%s)", ele.String())
			ai.Mult = burstRadiant[c.TalentLvlBurst()]
			ai.Element = ele
			c4cb = c.c4cb(c.Core.F)
		default:
			ai.Abil = "Soulseeker Shell"
			ai.Mult = burstSoulseeker[c.TalentLvlBurst()]
			ai.Element = attributes.Anemo
		}
		c.Core.QueueAttack(ai,
			combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()),
			burstSecondaryHitmark[i],
			burstSecondaryHitmark[i], c4cb)
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(22)

	c.AddStatus(SkillActionKey, SkillActionKeyDur, true)

	frames := frames.NewAbilFunc(burstFrames)
	if c.nightsoulState.HasBlessing() {
		frames = c.skillNextFrames(frames)
	}
	return action.Info{
		Frames:          frames,
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSkill], // earliest cancel
		State:           action.BurstState,
	}, nil
}
