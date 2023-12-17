package charlotte

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// TODO frames & aoe
var chargeFrames []int

const (
	chargeHitmark = 0
	chargeRadius  = 0
	arkheRadius   = 0
	arkheIcdKeys  = "spiritbreath-thorn-icd"
)

func init() {
	chargeFrames = frames.InitAbilSlice(0)
	chargeFrames[action.ActionCharge] = 0
	chargeFrames[action.ActionSkill] = 0
	chargeFrames[action.ActionBurst] = 0
	chargeFrames[action.ActionDash] = 0
	chargeFrames[action.ActionJump] = 0
	chargeFrames[action.ActionSwap] = 0
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	windup := 0
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		windup = 0
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			chargeRadius,
		),
		chargeHitmark-windup,
		chargeHitmark-windup,
		c.arkheCB,
	)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionDash] - windup, // TODO earliest cancel
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) arkheCB(a combat.AttackCB) {
	if c.StatusIsActive(arkheIcdKeys) {
		return
	}

	c.AddStatus(arkheIcdKeys, 6*60, true)

	c.QueueCharTask(func() {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Spiritbreath Thorn" + " (" + c.Base.Key.Pretty() + ")",
			AttackTag:      attacks.AttackTagExtra,
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeSlash,
			Element:        attributes.Cryo,
			Durability:     0,
			Mult:           arkhe[c.TalentLvlAttack()],
			IgnoreInfusion: true,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				arkheRadius,
			),
			0,
			0,
		)
	}, 30)
}
