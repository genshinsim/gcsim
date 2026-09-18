package iansan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	chargeFrames []int
	swiftFrames  []int
)

const (
	chargeHitmark = 21
	swiftHitmark  = 25
)

func init() {
	chargeFrames = frames.InitAbilSlice(55) // E -> Walk
	chargeFrames[action.ActionAttack] = 52
	chargeFrames[action.ActionSkill] = 51
	chargeFrames[action.ActionBurst] = 50
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 49

	swiftFrames = frames.InitAbilSlice(33) // E -> Walk
	swiftFrames[action.ActionAttack] = 32
	swiftFrames[action.ActionBurst] = 24
	swiftFrames[action.ActionDash] = 30
	swiftFrames[action.ActionJump] = 27
	swiftFrames[action.ActionSwap] = 31
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() || c.StatusIsActive(fastSkill) {
		c.DeleteStatus(fastSkill)
		return c.chargedSwift(), nil
	}

	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
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
		Mult:               charge[c.TalentLvlAttack()],
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
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Swift Stormflight",
		AdditionalTags:   []attacks.AttackTag{attacks.AttackTagNightsoul},
		PoiseDMG:         120,
		AttackTag:        attacks.AttackTagExtra,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupPoleExtraAttack,
		StrikeType:       attacks.StrikeTypeBlunt,
		IgnoreInfusion:   true,
		Element:          attributes.Electro,
		Durability:       25,
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.05 * 60,
		Mult:             swift[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.PrimaryTarget(), info.Point{Y: 0.5}, 2.5, 360),
		swiftHitmark,
		swiftHitmark,
		c.makeA1CB(),
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(swiftFrames),
		AnimationLength: swiftFrames[action.InvalidAction],
		CanQueueAfter:   swiftFrames[action.ActionBurst],
		State:           action.ChargeAttackState,
	}
}
