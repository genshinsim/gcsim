package charlotte

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int

const (
	chargeHitmark = 67
	chargeRadius  = 2.5
	arkheRadius   = 2.5
	arkheIcdKeys  = "spiritbreath-thorn-icd"
)

func init() {
	chargeFrames = frames.InitAbilSlice(84) // CA -> CA
	chargeFrames[action.ActionAttack] = 79
	chargeFrames[action.ActionSkill] = 71
	chargeFrames[action.ActionBurst] = 71
	chargeFrames[action.ActionDash] = 21
	chargeFrames[action.ActionJump] = 21
	chargeFrames[action.ActionWalk] = 74
	chargeFrames[action.ActionSwap] = 60
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
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		windup = 14
	case action.SkillState:
		if c.Core.Player.LastAction.Param["hold"] == 0 {
			windup = 8
		}
	case action.BurstState:
		windup = 3
	}

	pos := c.Core.Combat.PrimaryTarget().Pos()
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(
			pos,
			nil,
			chargeRadius,
		),
			0,
			0,
		)
		c.QueueCharTask(c.arkhe(pos), 30)
	}, chargeHitmark-windup)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionDash] - windup,
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) arkhe(pos geometry.Point) func() {
	return func() {
		if c.StatusIsActive(arkheIcdKeys) {
			return
		}
		c.AddStatus(arkheIcdKeys, 6*60, true)
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Spiritbreath Thorn" + " (" + c.Base.Key.Pretty() + ")",
			AttackTag:      attacks.AttackTagExtra,
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypePierce,
			Element:        attributes.Cryo,
			Durability:     0,
			Mult:           arkhe[c.TalentLvlAttack()],
			IgnoreInfusion: true,
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				pos,
				nil,
				arkheRadius,
			),
			0,
			0,
		)
	}
}
