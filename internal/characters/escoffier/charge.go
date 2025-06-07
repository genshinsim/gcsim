package escoffier

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int

const chargeHitmark = 20

func init() {
	chargeFrames = frames.InitAbilSlice(65)
	chargeFrames[action.ActionAttack] = 53
	chargeFrames[action.ActionSkill] = 52
	chargeFrames[action.ActionBurst] = 52
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 52
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagExtraAttack,
		ICDGroup:           attacks.ICDGroupPoleExtraAttack,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               charge[c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -0.1},
			2.8,
			4.8,
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
