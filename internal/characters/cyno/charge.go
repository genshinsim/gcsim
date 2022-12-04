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
		ICDGroup:           combat.ICDGroupPoleExtraAttack,
		StrikeType:         combat.StrikeTypeSpear,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               charge[c.TalentLvlAttack()],
		HitlagHaltFrames:   0,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.8,
		),
		0,
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
	chargeBFrames  []int
	chargeBHitmark = 27
)

func init() {
	// charge (burst) -> x
	chargeBFrames = frames.InitAbilSlice(65)
	chargeBFrames[action.ActionSkill] = 63
	chargeBFrames[action.ActionDash] = 27
	chargeBFrames[action.ActionJump] = 27
	chargeBFrames[action.ActionSwap] = 63
}

func (c *char) chargeB(p map[string]int) action.ActionInfo {
	c.tryBurstPPSlide(chargeBHitmark)

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Pactsworn Pathclearer Charge",
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagExtraAttack,
		ICDGroup:           combat.ICDGroupPoleExtraAttack,
		StrikeType:         combat.StrikeTypeSpear,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               chargeB[c.TalentLvlBurst()],
		HitlagHaltFrames:   0,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.8,
		),
		0,
		chargeBHitmark,
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeBFrames),
		AnimationLength: chargeBFrames[action.InvalidAction],
		CanQueueAfter:   chargeBHitmark,
		State:           action.ChargeAttackState,
	}
}
