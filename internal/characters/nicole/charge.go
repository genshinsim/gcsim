package nicole

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	chargeHitmark = 66
)

var chargeFrames []int

func init() {
	chargeFrames = frames.InitAbilSlice(99)
	chargeFrames[action.ActionSkill] = chargeHitmark
	chargeFrames[action.ActionBurst] = chargeHitmark
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = chargeHitmark
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	// TODO: Find actual windup and conditions for windup
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		if c.NormalCounter == 1 || c.NormalCounter == 2 {
			windup = 14
		}
	}

	c.Core.Tasks.Add(func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Charge",
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       charge[c.TalentLvlAttack()],
		}

		ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 3)
		c.Core.QueueAttack(
			ai,
			ap,
			0,
			0,
		)
	}, chargeHitmark)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionJump] - windup, // earliest cancel
		State:           action.ChargeAttackState,
	}, nil
}
