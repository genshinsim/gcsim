package raiden

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 22

func init() {
	// charge -> x
	chargeFrames = frames.InitAbilSlice(37) //n1, skill, burst all at 37
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 36
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(burstKey) {
		return c.swordCharge(p)
	}

	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Charge Attack",
		AttackTag:        combat.AttackTagExtra,
		ICDTag:           combat.ICDTagExtraAttack,
		ICDGroup:         combat.ICDGroupDefault,
		Element:          attributes.Physical,
		Durability:       25,
		HitlagHaltFrames: 0.02 * 60, //all raiden normals have 0.02s hitlag
		HitlagFactor:     0.01,
		Mult:             charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), chargeHitmark, chargeHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}

var swordCAFrames []int
var swordCAHitmarks = []int{24, 31}

func init() {
	// charge (burst) -> x
	swordCAFrames = frames.InitAbilSlice(56)
	swordCAFrames[action.ActionDash] = swordCAHitmarks[len(swordCAHitmarks)-1]
	swordCAFrames[action.ActionJump] = swordCAHitmarks[len(swordCAHitmarks)-1]
}

func (c *char) swordCharge(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Musou Isshin (Charge Attack)",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Electro,
		Durability:         25,
		HitlagHaltFrames:   0.02 * 60, //all raiden normals have 0.02s hitlag
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	for i, mult := range chargeSword {
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewDefCircHit(5, false, combat.TargettableEnemy),
				0,
				0,
				c.burstRestorefunc,
				c.c6,
			)
		}, swordCAHitmarks[i])
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(swordCAFrames),
		AnimationLength: swordCAFrames[action.InvalidAction],
		CanQueueAfter:   swordCAHitmarks[len(swordCAHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
