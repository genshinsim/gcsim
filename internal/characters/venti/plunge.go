package venti

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var highPlungeFrames []int

const highPlungeHitmark = 58

func (c *char) HighPlungeAttack(p map[string]int) action.ActionInfo {

	// check if hold skill was used
	lastAct := c.Core.Player.LastAction
	if lastAct.Char != c.Index || lastAct.Type != action.ActionSkill || lastAct.Param["hold"] != 0 {
		c.Core.Log.NewEvent("high_plunge should be preceded by hold skill", glog.LogActionEvent, c.Index, "action", action.ActionHighPlunge)
		return action.ActionInfo{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			Post:            1200,
			State:           action.Idle,
		}
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Plunge",
		AttackTag:      combat.AttackTagPlunge,
		ICDTag:         combat.ICDTagNone,
		ICDGroup:       combat.ICDGroupDefault,
		StrikeType:     combat.StrikeTypeBlunt,
		Element:        attributes.Physical,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		IgnoreInfusion: true,
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), highPlungeHitmark, highPlungeHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(highPlungeFrames),
		AnimationLength: highPlungeFrames[action.InvalidAction],
		CanQueueAfter:   highPlungeHitmark,
		Post:            highPlungeHitmark,
		State:           action.PlungeAttackState,
	}
}
