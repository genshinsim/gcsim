package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var (
	chargeFrames []int
	windupFrames []int
)

const (
	chargeHitmark = 76
	windupHitmark = chargeHitmark - 14
)

func init() {
	windupFrames = frames.InitAbilSlice(113)
	windupFrames[action.ActionAttack] = 59
	windupFrames[action.ActionCharge] = 59
	windupFrames[action.ActionSkill] = 59
	windupFrames[action.ActionBurst] = 59
	windupFrames[action.ActionDash] = 31
	windupFrames[action.ActionJump] = 30
	windupFrames[action.ActionSwap] = 104
	chargeFrames = make([]int, len(windupFrames))
	for i := range windupFrames {
		chargeFrames[i] = windupFrames[i] - 14
	}
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

	adjustedHitmark := windupHitmark
	adjustedFrames := windupFrames
	lastAction := &c.Core.Player.LastAction
	if lastAction.Char == c.Index {
		switch lastAction.Type {
		case action.ActionAttack,
			action.ActionCharge,
			action.ActionSkill:
			adjustedHitmark = chargeHitmark
			adjustedFrames = chargeFrames
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
