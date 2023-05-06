package yaemiko

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int

const chargeHitmark = 90

func init() {
	chargeFrames = frames.InitAbilSlice(96) // CA -> N1/E/Q
	chargeFrames[action.ActionCharge] = 95  // CA -> CA
	chargeFrames[action.ActionDash] = 46    // CA -> D
	chargeFrames[action.ActionJump] = 47    // CA -> J
	chargeFrames[action.ActionSwap] = 94    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   attacks.ICDGroupYaeCharged,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	// skip CA windup if we're in NA animation
	windup := 0
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		windup = 14
	}

	// starts at recorded hitmark and +1.65m in target direction
	// moves at 11m/s, one attack every 0.15s (9f), so it moves at 11 * 0.15 = 1.65m per attack
	// gets gated by special damage sequence (once every 0.5s)
	initialPos := c.Core.Combat.PrimaryTarget().Pos()
	initialDirection := c.Core.Combat.Player().Direction()
	for i := 0; i < 5; i++ {
		nextPos := geometry.CalcOffsetPoint(initialPos, geometry.Point{Y: 1.65 * float64(i+1)}, initialDirection)
		c.Core.QueueAttack(
			ai,
			// direction should stay the same because primary target pos can't change during this loop
			combat.NewBoxHit(c.Core.Combat.Player(), nextPos, nil, 2, 2),
			0, // TODO: check snapshot delay
			chargeHitmark+(i*9)-windup,
		)
	}

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionDash] - windup, // earliest cancel is before hitmark
		State:           action.ChargeAttackState,
	}
}
