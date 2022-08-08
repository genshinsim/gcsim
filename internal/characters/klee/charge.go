package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
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
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
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

	adjustedHitmark := chargeHitmark
	adjustedFrames := chargeFrames
	lastAction := &c.Core.Player.LastAction
	if lastAction.Char == c.Index {
		if (lastAction.Type == action.ActionAttack && (c.NormalCounter == 1 || c.NormalCounter == 2)) ||
			lastAction.Type == action.ActionSkill { // if Klee uses any of these, the windup is removed
			adjustedHitmark -= 14
			adjustedFrames = make([]int, len(chargeFrames))
			copy(adjustedFrames, chargeFrames)
			for i := range adjustedFrames {
				adjustedFrames[i] -= 14
			}
		}
	}
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
		adjustedHitmark+travel,
	)

	c.c1(adjustedHitmark + travel)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(adjustedFrames),
		AnimationLength: adjustedFrames[action.InvalidAction],
		CanQueueAfter:   0,
		State:           action.ChargeAttackState,
	}
}
