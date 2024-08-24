package clorinde

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	chargeFrames  []int
	chargeHitmark = 38
	chargeRadius  = 4.1
)

func init() {
	chargeFrames = frames.InitAbilSlice(44) // CA -> E
	chargeFrames[action.ActionAttack] = 66
	chargeFrames[action.ActionBurst] = 67
	chargeFrames[action.ActionSwap] = 45
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Charge",
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       attacks.ICDTagNormalAttack,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Physical,
		Durability:   25,
		Mult:         charge[c.TalentLvlAttack()],
		HitlagFactor: 0.02,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, chargeRadius),
		chargeHitmark,
		chargeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}
