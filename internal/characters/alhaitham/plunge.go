package alhaitham

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var lowPlungeFrames []int

const lowPlungeHitmark = 19

func init() {
	lowPlungeFrames = frames.InitAbilSlice(50)
	lowPlungeFrames[action.ActionAttack] = 29
	lowPlungeFrames[action.ActionCharge] = 30
	lowPlungeFrames[action.ActionBurst] = 31
	lowPlungeFrames[action.ActionDash] = 20
	lowPlungeFrames[action.ActionSwap] = 38

}

func (c *char) LowPlungeAttack(p map[string]int) action.ActionInfo {
	//last action must be hold skill
	if c.Core.Player.LastAction.Char != c.Index ||
		c.Core.Player.LastAction.Type != action.ActionSkill ||
		c.Core.Player.LastAction.Param["hold"] != 1 {
		c.Core.Log.NewEvent("only plunge after hold skill ends", glog.LogActionEvent, c.Index).
			Write("action", action.ActionLowPlunge)
		return action.ActionInfo{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			State:           action.Idle,
		}
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
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		lowPlungeHitmark, lowPlungeHitmark, c.a1CB)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return lowPlungeFrames[next] },
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeHitmark,
		State:           action.PlungeAttackState,
	}
}
