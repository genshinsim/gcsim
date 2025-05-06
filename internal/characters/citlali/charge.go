package citlali

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const (
	chargeHitmark = 51
	chargeRadius  = 1
)

func init() {
	chargeFrames = frames.InitAbilSlice(52) // CA -> CA
	chargeFrames[action.ActionAttack] = 49
	chargeFrames[action.ActionBurst] = 51
	chargeFrames[action.ActionDash] = 45
	chargeFrames[action.ActionJump] = 46
	chargeFrames[action.ActionWalk] = 56
	chargeFrames[action.ActionSwap] = 50
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

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
	}, chargeHitmark+travel)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] },
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash] + travel,
		State:           action.ChargeAttackState,
	}, nil
}
