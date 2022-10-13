package cyno

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 24

func init() {
	chargeFrames = frames.InitAbilSlice(63)
	chargeFrames[action.ActionBurst] = 62
	chargeFrames[action.ActionDash] = 24
	chargeFrames[action.ActionJump] = 24
	chargeFrames[action.ActionSwap] = 61
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(burstKey) {
		return c.chargeB(p)
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge Attack",
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagExtraAttack,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		Mult:               charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy),
		chargeHitmark,
		chargeHitmark,
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}

var (
	chargeBFrames   []int
	chargeBHitmarks = 27
)

func init() {
	// charge (burst) -> x
	chargeBFrames = frames.InitAbilSlice(65)
	chargeBFrames[action.ActionSkill] = 63
	chargeBFrames[action.ActionDash] = 26
	chargeBFrames[action.ActionJump] = 26
	chargeBFrames[action.ActionSwap] = 63
}

func (c *char) chargeB(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Pactsworn Pathclearer Charge",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               chargeB[c.TalentLvlBurst()],
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy),
			0,
			0,
		)
	}, chargeBHitmarks)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeBFrames),
		AnimationLength: chargeBFrames[action.InvalidAction],
		CanQueueAfter:   chargeBHitmarks,
		State:           action.ChargeAttackState,
	}
}
