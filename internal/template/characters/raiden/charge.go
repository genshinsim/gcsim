package raiden

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeAttackHitmark = 22

func (c *char) chargeAttackFrameFunc(next action.Action) int {
	return chargeFrames[next]
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	if c.Core.Status.Duration("raidenburst") > 0 {
		return c.swordCharge(p)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewDefCircHit(0.5, false, combat.TargettableEnemy),
		chargeAttackHitmark,
		chargeAttackHitmark,
	)

	return action.ActionInfo{
		Frames:          c.chargeAttackFrameFunc,
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeAttackHitmark,
		Post:            chargeAttackHitmark,
		State:           action.ChargeAttackState,
	}

}

var swordCAFrames []int

const swordCAHitmark = 22

func (c *char) swordCAFrameFunc(next action.Action) int {
	return swordCAFrames[next]
}

func (c *char) swordCharge(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Isshin (Charge Attack)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
	}

	for _, mult := range chargeSword {
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(5, false, combat.TargettableEnemy),
			swordCAHitmark,
			swordCAHitmark,
			c.burstRestorefunc,
			c.c6(),
		)
	}

	return action.ActionInfo{
		Frames:          c.swordCAFrameFunc,
		AnimationLength: swordCAFrames[action.InvalidAction],
		CanQueueAfter:   swordCAHitmark,
		Post:            swordCAHitmark,
		State:           action.ChargeAttackState,
	}
}
