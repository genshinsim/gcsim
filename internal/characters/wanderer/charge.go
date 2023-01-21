package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFramesNormal []int
var chargeFramesE []int

const chargeHitmarkNormal = 34
const chargeHitmarkE = 36

func init() {
	chargeFramesNormal = frames.InitAbilSlice(69)
	chargeFramesNormal[action.ActionAttack] = 51
	chargeFramesNormal[action.ActionCharge] = 50
	chargeFramesNormal[action.ActionSkill] = 49
	chargeFramesNormal[action.ActionBurst] = 49
	chargeFramesNormal[action.ActionDash] = 37
	chargeFramesNormal[action.ActionJump] = 37
	chargeFramesNormal[action.ActionSwap] = 47

	chargeFramesE = frames.InitAbilSlice(70)
	chargeFramesE[action.ActionAttack] = 49
	chargeFramesE[action.ActionCharge] = 49
	chargeFramesE[action.ActionSkill] = 38
	chargeFramesE[action.ActionBurst] = 38
	chargeFramesE[action.ActionDash] = 38
	chargeFramesE[action.ActionJump] = 38
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	if c.StatusIsActive(skillKey) {
		// Can only occur if delay == 0, so it can be disregarded
		return c.WindfavoredChargeAttack(p)
	}

	windup := c.chargeWindupNormal() + delay

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	// TODO: check snapshot delay
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 2.5),
		windup+chargeHitmarkNormal, windup+chargeHitmarkNormal,
		c.a4CB,
		c.a1ElectroCB,
		c.particleGenCB)
	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return windup +
				frames.AtkSpdAdjust(chargeFramesNormal[next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: windup + chargeFramesNormal[action.InvalidAction],
		CanQueueAfter:   windup + chargeFramesNormal[action.ActionDash],
		State:           action.ChargeAttackState,
	}
}

func (c *char) WindfavoredChargeAttack(p map[string]int) action.ActionInfo {
	windup := c.chargeWindupE()

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack (Windfavored)",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillNABonus[c.TalentLvlSkill()] * charge[c.TalentLvlAttack()],
	}

	// TODO: check snapshot delay
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5),
		windup+chargeHitmarkE, windup+chargeHitmarkE,
		c.a4CB, c.a1ElectroCB, c.particleGenCB)
	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return windup +
				frames.AtkSpdAdjust(chargeFramesE[next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: windup + chargeFramesE[action.InvalidAction],
		CanQueueAfter:   windup + chargeFramesE[action.ActionDash],
		State:           action.ChargeAttackState,
	}
}

func (c *char) chargeWindupNormal() int {
	switch c.Core.Player.LastAction.Type {
	case action.ActionAttack:
		return -2
	case action.ActionLowPlunge:
		return -4
	default:
		return 0
	}
}

func (c *char) chargeWindupE() int {
	switch c.Core.Player.LastAction.Type {
	case action.ActionAttack:
		return -4
	case action.ActionCharge, action.ActionJump:
		return -7
	case action.ActionBurst:
		return -2
	case action.ActionDash:
		return -9
	default:
		return 0
	}
}
