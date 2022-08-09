package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var paramitaChargeFrames []int

const chargeHitmark = 19
const paramitaChargeHitmark = 6

func init() {
	// charge -> x
	chargeFrames = frames.InitAbilSlice(62)
	chargeFrames[action.ActionAttack] = 57
	chargeFrames[action.ActionSkill] = 57
	chargeFrames[action.ActionSkill] = 60
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark

	// charge (paramita) -> x
	paramitaChargeFrames = frames.InitAbilSlice(44)
	paramitaChargeFrames[action.ActionBurst] = 35
	paramitaChargeFrames[action.ActionDash] = paramitaChargeHitmark
	paramitaChargeFrames[action.ActionJump] = paramitaChargeHitmark
	paramitaChargeFrames[action.ActionSwap] = 42
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	if c.StatModIsActive(paramitaBuff) {
		return c.ppChargeAttack(p)
	}

	//check for particles
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Charge Attack",
		AttackTag:        combat.AttackTagExtra,
		ICDTag:           combat.ICDTagExtraAttack,
		ICDGroup:         combat.ICDGroupPole,
		StrikeType:       combat.StrikeTypeSlash,
		Element:          attributes.Physical,
		Durability:       25,
		Mult:             charge[c.TalentLvlAttack()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.01 * 60,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy), 0, chargeHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}

func (c *char) ppChargeAttack(p map[string]int) action.ActionInfo {
	//TODO: currently assuming snapshot is on cast since it's a bullet and nothing implemented re "pp slide"
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupPole,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy), 0, paramitaChargeHitmark, c.ppParticles, c.applyBB)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(paramitaChargeFrames),
		AnimationLength: paramitaChargeFrames[action.InvalidAction],
		CanQueueAfter:   paramitaChargeHitmark,
		State:           action.ChargeAttackState,
	}
}
