package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var ppChargeFrames []int

const (
	chargeHitmark   = 19
	ppChargeHitmark = 3
)

func init() {
	// charge -> x
	chargeFrames = frames.InitAbilSlice(62)
	chargeFrames[action.ActionAttack] = 57
	chargeFrames[action.ActionSkill] = 57
	chargeFrames[action.ActionSkill] = 60
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark

	// charge (paramita) -> x
	ppChargeFrames = frames.InitAbilSlice(44)
	ppChargeFrames[action.ActionBurst] = 33
	ppChargeFrames[action.ActionDash] = ppChargeHitmark
	ppChargeFrames[action.ActionJump] = ppChargeHitmark
	ppChargeFrames[action.ActionSwap] = 42
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	if c.StatModIsActive(paramitaStatus) {
		return c.ppChargeAttack(p)
	}

	//check for particles
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

func (c *char) ppChargeAttack(p map[string]int) action.ActionInfo {

	//TODO: currently assuming snapshot is on cast since it's a bullet and nothing implemented re "pp slide"
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge Attack",
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagExtraAttack,
		ICDGroup:           combat.ICDGroupPoleExtraAttack,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               charge[c.TalentLvlAttack()],
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
		ppChargeHitmark,
		c.ppParticles,
		c.applyBB,
	)

	//frames changes if previous action is normal
	prevState := -1
	if c.Core.Player.LastAction.Char == c.Index && c.Core.Player.LastAction.Type == action.ActionAttack {
		prevState = c.NormalCounter - 1
		if prevState < 0 {
			prevState = c.NormalHitNum - 1
		}
	}
	ff := func(next action.Action) int {
		if prevState == -1 {
			return ppChargeFrames[next]
		}
		switch next {
		case action.ActionDash, action.ActionJump:
		default:
			return ppChargeFrames[next]
		}
		switch prevState {
		case 0: // N1
			if next == action.ActionDash {
				return 1 // N1D
			} else {
				return 2 // N1J
			}
		case 1: // N2
			if next == action.ActionDash {
				return 4 // N2D
			} else {
				return 5 // N2J
			}
		case 2: // N3
			return 2
		case 3: // N4
			return 3
		case 4: // N5
			return 3
		default:
			return 500 //TODO: this action is illegal; need better way to handle it
		}
	}

	return action.ActionInfo{
		Frames:          ff,
		AnimationLength: ppChargeFrames[action.InvalidAction],
		CanQueueAfter:   1,
		State:           action.ChargeAttackState,
	}
}
