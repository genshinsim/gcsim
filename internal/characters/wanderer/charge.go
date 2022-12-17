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

	relevantHitmark := chargeHitmarkNormal
	relevantFrames := chargeFramesNormal

	if c.StatusIsActive(skillKey) {
		relevantHitmark = chargeHitmarkE
		relevantFrames = chargeFramesE
	}

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
		delay+relevantHitmark, delay+relevantHitmark)
	return action.ActionInfo{
		Frames:          func(next action.Action) int { return delay +
			frames.AtkSpdAdjust(relevantFrames[next], c.Stat(attributes.AtkSpd)) },
		AnimationLength: delay + relevantFrames[action.InvalidAction],
		CanQueueAfter:   delay + relevantHitmark,
		State:           action.ChargeAttackState,
	}
}
