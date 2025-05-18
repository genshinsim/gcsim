package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const (
	chargeHitmark  = 76
	chargeSnapshot = 29 + 32
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

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
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
		PoiseDMG:   180,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		if c.NormalCounter == 1 || c.NormalCounter == 2 {
			windup = 14
		}
	case action.SkillState:
		windup = 14
	}

	c.Core.Tasks.Add(func() {
		// apply and clear spark
		snap := c.Snapshot(&ai)
		if c.StatusIsActive(a1SparkKey) {
			snap.Stats[attributes.DmgP] += .50
			c.DeleteStatus(a1SparkKey)
		}

		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 3),
			(chargeHitmark-chargeSnapshot)+travel,
			c.makeA4CB(),
		)
	}, chargeSnapshot-windup)

	c.c1(chargeHitmark - windup + travel)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionJump] - windup, // earliest cancel
		State:           action.ChargeAttackState,
	}, nil
}
