package ningguang

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 50

func init() {
	chargeFrames = frames.InitAbilSlice(50)
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
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy), chargeHitmark, chargeHitmark+travel)

	ai = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge (Gems)",
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagExtraAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Geo,
		Durability:         50,
		Mult:               jade[c.TalentLvlAttack()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	for i := 0; i < c.Tags["jade"]; i++ {
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy), chargeHitmark, chargeHitmark+travel)
	}
	c.Tags["jade"] = 0

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
