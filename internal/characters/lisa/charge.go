package lisa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var chargeFramesNoWindup []int
var chargeFrames []int

const (
	// hitmark frame, includes CA windup
	chargeHitmark = 70
	// TODO: stacks technically only last 15s and each stack has its own timer
	conductiveTag = "lisa-conductive-stacks"
)

func init() {
	chargeFrames = frames.InitAbilSlice(93)
	chargeFrames[action.ActionAttack] = 91
	chargeFrames[action.ActionCharge] = 90
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 90
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	// skip CA windup if we're in NA or Skill (Hold) animation
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		windup = 14
	case action.SkillState:
		if c.Core.Player.LastAction.Param["hold"] != 0 {
			windup = 14
		}
	}

	cb := func(a combat.AttackCB) {
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		count := t.GetTag(conductiveTag)
		if count < 3 {
			t.SetTag(conductiveTag, count+1)
		}
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.1), chargeHitmark-windup, chargeHitmark-windup, cb)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeHitmark - windup,
		State:           action.ChargeAttackState,
	}
}
