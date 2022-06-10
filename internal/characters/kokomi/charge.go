package kokomi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

// hitmark frame, includes CA windup
const chargeHitmark = 48

func init() {
	chargeFrames = frames.InitAbilSlice(76)
	chargeFrames[action.ActionAttack] = 62
	chargeFrames[action.ActionCharge] = 62
	chargeFrames[action.ActionSkill] = 62
	chargeFrames[action.ActionBurst] = 62
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 62
}

// Standard charge attack
// CA has no travel time
func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	ai.FlatDmg = c.burstDmgBonus(ai.AttackTag)

	// skip CA windup if we're in NA animation
	windup := 0
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		windup = 14
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), chargeHitmark-windup, chargeHitmark-windup)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeHitmark - windup,
		Post:            chargeHitmark - windup,
		State:           action.ChargeAttackState,
	}
}
