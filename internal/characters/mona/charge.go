package mona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var chargeFrames []int

const chargeHitmark = 66

func init() {
	chargeFrames = frames.InitAbilSlice(113) // CA -> N1
	chargeFrames[action.ActionCharge] = 105  // CA -> CA
	chargeFrames[action.ActionSkill] = 92    // CA -> E
	chargeFrames[action.ActionBurst] = 68    // CA -> Q
	chargeFrames[action.ActionDash] = 72     // CA -> D
	chargeFrames[action.ActionJump] = 64     // CA -> J
	chargeFrames[action.ActionSwap] = 54     // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	// add windup if we're in idle or swap only
	// TODO: this ignores N4 -> CA (which should be illegal anyways)
	windup := 14
	if c.Core.Player.CurrentState() == action.Idle || c.Core.Player.CurrentState() == action.SwapState {
		windup = 0
	}

	c.doCA(c.Core.Combat.PrimaryTarget(), chargeHitmark-windup)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionSwap] - windup, // earliest cancel is before hitmark
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) doCA(target info.Target, delay int) {
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
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			target,
			nil,
			3,
		),
		delay,
		delay,
		c.makeC6CAResetCB(),
		c.astralGlowGainCB,
		c.omenRefreshCB,
		c.c2HexereiCB,
	)
}
