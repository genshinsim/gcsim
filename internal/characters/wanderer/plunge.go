package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var lowPlungeFrames []int

const lowPlungeHitmark = 36

func init() {
	lowPlungeFrames = frames.InitAbilSlice(55)
	lowPlungeFrames[action.ActionDash] = 43
	lowPlungeFrames[action.ActionJump] = 50
	lowPlungeFrames[action.ActionSwap] = 50

}

func (c *char) LowPlungeAttack(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	// Not in falling state
	if delay == 0 || !(c.Core.Player.LastAction.Char == c.Index &&
		c.Core.Player.LastAction.Type == action.ActionSkill && !c.StatusIsActive(skillKey)) {
		// Nothing so far?
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge Attack",
		AttackTag:  combat.AttackTagPlunge,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}

	// TODO: check snapshot delay
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2.28),
		delay+chargeHitmark, delay+chargeHitmark)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return delay + lowPlungeFrames[next] },
		AnimationLength: delay + lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   delay + lowPlungeHitmark,
		State:           action.PlungeAttackState,
	}
}
