package vodyanitsa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var chargeFrames []int

const (
	idleWindup    = 15
	burstWindup   = idleWindup
	chargeHitmark = 52 - idleWindup
)

func init() {
	chargeFrames = frames.InitAbilSlice(75 - idleWindup)
	chargeFrames[action.ActionCharge] = 78 - 32
	chargeFrames[action.ActionSkill] = 56 - idleWindup
	chargeFrames[action.ActionBurst] = 54 - idleWindup
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionWalk] = 68 - idleWindup
	chargeFrames[action.ActionSwap] = 45
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.Idle:
		windup = idleWindup
	case action.BurstState:
		windup = burstWindup
	case action.SkillState:
		windup = 17
	case action.ChargeAttackState:
		windup = 32
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 2),
		chargeHitmark+windup,
		chargeHitmark+windup+travel,
	)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] + windup },
		AnimationLength: chargeFrames[action.InvalidAction] + windup,
		CanQueueAfter:   chargeHitmark + windup,
		State:           action.ChargeAttackState,
	}, nil
}
