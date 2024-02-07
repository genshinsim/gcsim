package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
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
	ppChargeFrames = frames.InitAbilSlice(42)
	ppChargeFrames[action.ActionBurst] = 33
	ppChargeFrames[action.ActionDash] = ppChargeHitmark
	ppChargeFrames[action.ActionJump] = ppChargeHitmark
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatModIsActive(paramitaBuff) {
		return c.ppChargeAttack(), nil
	}

	// check for particles
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge Attack",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagExtraAttack,
		ICDGroup:           attacks.ICDGroupPoleExtraAttack,
		StrikeType:         attacks.StrikeTypeSpear,
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

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) ppChargeAttack() action.Info {
	// pp slide: add 1.8s to paramita on charge attack start which gets removed once the charge attack ends
	c.ExtendStatus(paramitaBuff, 1.8*60)

	//TODO: currently assuming snapshot is on cast since it's a bullet and nothing implemented re "pp slide"
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge Attack",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagExtraAttack,
		ICDGroup:           attacks.ICDGroupPoleExtraAttack,
		StrikeType:         attacks.StrikeTypeSlash,
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
		c.particleCB,
		c.applyBB,
	)

	// frames changes if previous action is normal
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
			}
			return 2 // N1J
		case 1: // N2
			if next == action.ActionDash {
				return 4 // N2D
			}
			return 5 // N2J
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

	return action.Info{
		Frames:          ff,
		AnimationLength: ppChargeFrames[action.InvalidAction],
		CanQueueAfter:   1,
		State:           action.ChargeAttackState,
		OnRemoved: func(next action.AnimationState) {
			if next != action.BurstState {
				c.ExtendStatus(paramitaBuff, -1.8*60)
			}
		},
	}
}
