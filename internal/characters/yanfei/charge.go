package yanfei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var chargeFrames []int

const chargeHitmark = 66

func init() {
	chargeFrames = frames.InitAbilSlice(66)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
}

// Charge attack function - handles seal use
func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	//check for seal stacks
	if c.Core.F > c.sealExpiry {
		c.Tags["seal"] = 0
	}
	stacks := c.Tags["seal"]

	// apply a1
	c.a1(stacks)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       charge[stacks][c.TalentLvlAttack()],
	}
	// TODO: Not sure of snapshot timing
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), chargeHitmark, chargeHitmark)

	c.Core.Log.NewEvent("yanfei charge attack consumed seals", glog.LogCharacterEvent, c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)

	// Clear the seals next frame just in case for some reason we call stam check late
	c.Core.Tasks.Add(func() {
		c.Tags["seal"] = 0
		c.sealExpiry = c.Core.F - 1
	}, 1)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,

		State: action.ChargeAttackState,
	}
}
