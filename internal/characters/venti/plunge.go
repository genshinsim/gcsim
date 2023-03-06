package venti

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var highPlungeFrames []int

const highPlungeHitmark = 58

func init() {
	// TODO: missing counts for plunge cancels?
	// using hitmark as placeholder for now
	highPlungeFrames = frames.InitAbilSlice(highPlungeHitmark)
}

func (c *char) HighPlungeAttack(p map[string]int) action.ActionInfo {

	// check if hold skill was used
	lastAct := c.Core.Player.LastAction
	if lastAct.Char != c.Index || lastAct.Type != action.ActionSkill || lastAct.Param["hold"] != 0 {
		c.Core.Log.NewEvent("high_plunge should be preceded by hold skill", glog.LogActionEvent, c.Index).
			Write("action", action.ActionHighPlunge)
		return action.ActionInfo{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			State:           action.Idle,
		}
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "High Plunge",
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypePierce,
		Element:        attributes.Physical,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		IgnoreInfusion: true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5),
		highPlungeHitmark,
		highPlungeHitmark,
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(highPlungeFrames),
		AnimationLength: highPlungeFrames[action.InvalidAction],
		CanQueueAfter:   highPlungeHitmark,
		State:           action.PlungeAttackState,
	}
}
