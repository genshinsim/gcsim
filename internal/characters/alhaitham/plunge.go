package alhaitham

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var lowPlungeFrames []int

const lowPlungeHitmark = 18

func init() {
	lowPlungeFrames = frames.InitAbilSlice(50)
	lowPlungeFrames[action.ActionAttack] = 29
	lowPlungeFrames[action.ActionSkill] = 30
	lowPlungeFrames[action.ActionBurst] = 30
	lowPlungeFrames[action.ActionDash] = 20
	lowPlungeFrames[action.ActionSwap] = 38

}

func (c *char) LowPlungeAttack(p map[string]int) action.ActionInfo {
	//last action must be hold skill
	if c.Core.Player.LastAction.Type != action.ActionSkill ||
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
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1}, 3),
		lowPlungeHitmark, lowPlungeHitmark, c.a1CB, c.projectionAttack)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(lowPlungeFrames),
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}
