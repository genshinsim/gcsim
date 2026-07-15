package illuga

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var chargeFrames []int

const chargeHitmark = 50

func init() {
	chargeFrames = frames.InitAbilSlice(100)
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Charge",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupPoleExtraAttack,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               charge[c.TalentLvlAttack()],
		HitlagHaltFrames:   0,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.8,
		),
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
