package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var chargeFrames []int

const (
	chargeHitmark = 76
)

func init() {
	chargeFrames = frames.InitAbilSlice(113)
	chargeFrames[action.ActionAttack] = 59
	chargeFrames[action.ActionCharge] = 59
	chargeFrames[action.ActionSkill] = 59
	chargeFrames[action.ActionBurst] = 59
	chargeFrames[action.ActionDash] = 31
	chargeFrames[action.ActionJump] = 30
	chargeFrames[action.ActionSwap] = 104
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	// stam is calculated before this func is called so it's safe to
	// set spark to 0 here
	snap := c.Snapshot(&ai)
	if c.Core.Status.Duration("kleespark") > 0 {
		snap.Stats[attributes.DmgP] += .50
		c.Core.Status.Delete("kleespark")
		c.Core.Log.NewEvent("klee consumed spark", glog.LogCharacterEvent, c.Index).
			Write("icd", c.sparkICD)
	}

	windup := 0
	if (c.Core.Player.CurrentState() == action.NormalAttackState && (c.NormalCounter == 1 || c.NormalCounter == 2)) ||
		c.Core.Player.CurrentState() == action.SkillState {
		windup = 14
	}
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 3),
		chargeHitmark-windup+travel,
		c.makeA4CB(),
	)

	c.c1(chargeHitmark - windup + travel)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionJump] - windup, // earliest cancel
		State:           action.ChargeAttackState,
	}
}
