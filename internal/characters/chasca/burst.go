package chasca

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFramesGrounded []int
var burstFramesNS []int
var burstSecondaryHitmark = []int{103, 139, 147, 153, 157, 160}

const (
	burstHitmark = 96
	burstNSFall  = 102
)

func init() {
	burstFramesGrounded = frames.InitAbilSlice(114)
	burstFramesGrounded[action.ActionAttack] = 99 // Q -> N1
	burstFramesGrounded[action.ActionAim] = 100   // Q -> Aim
	burstFramesGrounded[action.ActionSkill] = 102 // Q -> E
	burstFramesGrounded[action.ActionDash] = 101  // Q -> D
	burstFramesGrounded[action.ActionJump] = 100  // Q -> J
	burstFramesGrounded[action.ActionSwap] = 98   // Q -> Swap

	burstFramesNS = frames.InitAbilSlice(111)
	burstFramesNS[action.ActionAttack] = 106 // Q -> N1
	burstFramesNS[action.ActionSkill] = 104  // Q -> E
	burstFramesNS[action.ActionDash] = 103   // Q -> D
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

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.0)
	c.Core.QueueAttack(
		ai,
		ap,
		burstHitmark,
		burstHitmark,
	)

	ai.ICDTag = attacks.ICDTagChascaBurst
	ai.ICDGroup = attacks.ICDGroupChascaBurst

	var c4cb combat.AttackCBFunc

	enemies := c.Core.Combat.EnemiesWithinArea(ap, nil)
	burstBullets := make([]attributes.Element, 0, 6)
	burstBullets = append(burstBullets, c.partyPHECTypes...)
	burstBullets = append(burstBullets, c.partyPHECTypes...)
	c.Core.Rand.Shuffle(len(burstBullets), func(i, j int) {
		burstBullets[i], burstBullets[j] = burstBullets[j], burstBullets[i]
	})
	burstFrame := c.Core.F
	for i := 0; i < 6; i++ {
		switch {
		case i < len(burstBullets):
			ele := burstBullets[i]
			ai.Abil = fmt.Sprintf("Radiant Soulseeker Shell (%s)", ele.String())
			ai.Mult = burstRadiant[c.TalentLvlBurst()]
			ai.Element = ele
			c4cb = c.c4cb(burstFrame)
		default:
			ai.Abil = "Soulseeker Shell"
			ai.Mult = burstSoulseeker[c.TalentLvlBurst()]
			ai.Element = attributes.Anemo
		}
		target := enemies[i%len(enemies)]
		c.Core.QueueAttack(ai,
			combat.NewSingleTargetHit(target.Key()),
			burstSecondaryHitmark[i],
			burstSecondaryHitmark[i], c4cb)
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(3)

	c.QueueCharTask(func() {
		if c.nightsoulState.HasBlessing() {
			return
		}
		c.Core.Log.NewEvent("nightsoul ended, falling", glog.LogCharacterEvent, c.Index)
		c.AddStatus(plungeAvailableKey, 26, true)
	}, burstNSFall)

	frames := frames.NewAbilFunc(burstFramesGrounded)
	if c.nightsoulState.HasBlessing() {
		// if we Q while in the air, we need to add the frames of fall down
		// TODO: set fall down animation to be "idle/skill" instead of burst?
		return action.Info{
			Frames:          c.skillNextFrames(frames, 0),
			AnimationLength: burstFramesNS[action.InvalidAction],
			CanQueueAfter:   burstNSFall, // can't start falling until frame 102
			State:           action.BurstState,
		}, nil
	}

	return action.Info{
		Frames:          frames,
		AnimationLength: burstFramesGrounded[action.InvalidAction],
		CanQueueAfter:   burstFramesGrounded[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
