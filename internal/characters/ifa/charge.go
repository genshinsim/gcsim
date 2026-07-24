package ifa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var chargeFrames []int

const (
	chargeHitmark       = 45
	chargeRadius        = 1.5
	chargeSkillInterval = 42
)

func init() {
	chargeFrames = frames.InitAbilSlice(86) // CA -> Walk
	chargeFrames[action.ActionAttack] = 64
	chargeFrames[action.ActionCharge] = 64
	chargeFrames[action.ActionSkill] = 64
	chargeFrames[action.ActionBurst] = 64
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 62

	skillAttackHoldFrames = frames.InitNormalCancelSlice(0, chargeSkillInterval)
	skillAttackHoldFrames[action.ActionAttack] = attackSkillInterval
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.attackTapSkillState(p), nil
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	travel, ok := p["travel"]

	if !ok {
		travel = 0
	}

	pos := c.Core.Combat.PrimaryTarget().Pos()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(
			pos,
			nil,
			chargeRadius,
		),
		chargeHitmark+travel,
		chargeHitmark+travel,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) attackHoldSkillState(p map[string]int) action.Info {
	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Tonic Shot",
		AttackTag:      attacks.AttackTagNormal,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagIfaSkill,
		ICDGroup:       attacks.ICDGroupIfaSkillHit,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           skill_dmg[c.TalentLvlSkill()],
	}

	travel, ok := p["travel"]

	if !ok {
		travel = 0
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.PrimaryTarget(),
		nil,
		3,
	)

	c.QueueCharTask(func() {
		if !c.nightsoulState.HasBlessing() {
			return
		}
		c.Core.QueueAttack(
			ai,
			ap,
			travel,
			travel,
			c.particleCB,
			c.healHoldCB,
			c.c1CB,
		)
	}, 1)

	c.c6OnHoldAttackSkill()

	atkspd := c.Stat(attributes.AtkSpd)

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.AtkSpdAdjust(skillAttackHoldFrames[next], atkspd)
		},
		AnimationLength: attackSkillInterval,
		CanQueueAfter:   0, // can run out of nightsoul and start falling earlier
		State:           action.NormalAttackState,
	}
}
