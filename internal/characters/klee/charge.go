package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var chargeFrames []int

const chargeHitmark = 84

func init() {
	chargeFrames = frames.InitAbilSlice(84)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
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
	//stam is calculated before this func is called so it's safe to
	//set spark to 0 here
	snap := c.Snapshot(&ai)
	if c.Core.Status.Duration("kleespark") > 0 {
		snap.Stats[attributes.DmgP] += .50
		c.Core.Status.Delete("kleespark")
		c.Core.Log.NewEvent("klee consumed spark", glog.LogCharacterEvent, c.Index).
			Write("icd", c.sparkICD)
	}

	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(2, false, combat.TargettableEnemy), chargeHitmark+travel)

	c.c1(chargeHitmark + travel)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
