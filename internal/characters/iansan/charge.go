package iansan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	chargeFrames []int
	swiftFrames  []int
)

const (
	chargeHitmark = 24
	swiftHitmark  = 14
)

func init() {
	chargeFrames = frames.InitAbilSlice(55)
	chargeFrames[action.ActionAttack] = 49
	chargeFrames[action.ActionSkill] = 49
	chargeFrames[action.ActionBurst] = 50
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 49

	swiftFrames = frames.InitAbilSlice(52)
	swiftFrames[action.ActionAttack] = 51
	swiftFrames[action.ActionSkill] = 51
	swiftFrames[action.ActionJump] = 50
	swiftFrames[action.ActionSwap] = 49
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() || c.StatusIsActive(fastSkill) {
		c.DeleteStatus(fastSkill)
		return c.chargedSwift(), nil
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charged Attack",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagExtraAttack,
		ICDGroup:           attacks.ICDGroupPoleExtraAttack,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
		Mult:               charged[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.8,
		),
		chargeHitmark,
		chargeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) chargedSwift() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Swift Stormflight",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagExtra,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           swift[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6),
		swiftHitmark,
		swiftHitmark,
		c.makeA1CB(),
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(swiftFrames),
		AnimationLength: swiftFrames[action.InvalidAction],
		CanQueueAfter:   swiftFrames[action.ActionSwap],
		State:           action.ChargeAttackState,
	}
}
