package durin

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var chargeFrames []int

const chargeHitmark = 17

func init() {
	chargeFrames = frames.InitAbilSlice(58)         // CA -> N1
	chargeFrames[action.ActionAttack] = 53          // CA -> N1
	chargeFrames[action.ActionSkill] = 52           // CA -> E
	chargeFrames[action.ActionBurst] = 52           // CA -> Q
	chargeFrames[action.ActionDash] = chargeHitmark // CA -> D
	chargeFrames[action.ActionJump] = chargeHitmark // CA -> J
	chargeFrames[action.ActionSwap] = 51            // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge",
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Mult:       charge[c.TalentLvlAttack()],
		Durability: 25,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), info.Point{Y: 0.3}, 2.2, 270),
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
