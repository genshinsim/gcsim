package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 34

func init() {
	chargeFrames = frames.InitAbilSlice(36)
	chargeFrames[action.ActionAttack] = 51
	chargeFrames[action.ActionCharge] = 50
	chargeFrames[action.ActionSkill] = 49
	chargeFrames[action.ActionBurst] = 49
	chargeFrames[action.ActionDash] = 37
	chargeFrames[action.ActionJump] = 37
	chargeFrames[action.ActionWalk] = 69
	chargeFrames[action.ActionSwap] = 47
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

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
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2.28),
		delay+chargeHitmark, delay+chargeHitmark)
	return action.ActionInfo{
		Frames:          func(next action.Action) int { return delay +
			frames.AtkSpdAdjust(chargeFrames[next], c.Stat(attributes.AtkSpd)) },
		AnimationLength: delay + chargeFrames[action.InvalidAction],
		CanQueueAfter:   delay + chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
