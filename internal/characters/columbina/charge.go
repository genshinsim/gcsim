package columbina

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	chargeFrames       []int
	chargeLBFrames     []int
	chargeHitmarkBloom = []int{50, 50 + 8, 50 + 8 + 22}
)

const (
	chargeHitmark     = 57
	verdantDewFrame   = 43
	chargeRadius      = 3.5
	chargeRadiusBloom = 3.3
)

func init() {
	chargeFrames = frames.InitAbilSlice(101) // CA -> CA
	chargeFrames[action.ActionAttack] = 94
	chargeFrames[action.ActionSkill] = chargeHitmark
	chargeFrames[action.ActionBurst] = chargeHitmark
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionWalk] = 99
	chargeFrames[action.ActionSwap] = chargeHitmark

	chargeLBFrames = frames.InitAbilSlice(124) // CA-LB -> Walk
	chargeLBFrames[action.ActionAttack] = 82
	chargeLBFrames[action.ActionCharge] = 80 // Need to set to 89 if the next CA isn't CA-LB
	chargeLBFrames[action.ActionSkill] = verdantDewFrame
	chargeLBFrames[action.ActionBurst] = verdantDewFrame
	chargeLBFrames[action.ActionDash] = verdantDewFrame
	chargeLBFrames[action.ActionJump] = verdantDewFrame
	chargeLBFrames[action.ActionSwap] = verdantDewFrame
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.Dew() > 0 {
		return c.ChargeAttackBloom(p), nil
	}

	windup := 15
	switch c.Core.Player.CurrentState() {
	case action.ChargeAttackState:
		windup = 3
	case action.NormalAttackState:
		windup = -2 // CA is faster when done out of NA
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

	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, chargeRadius)
	delay := chargeHitmark + windup
	c.Core.QueueAttack(ai, ap, delay, delay)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) ChargeAttackBloom(p map[string]int) action.Info {
	windup := 12
	switch c.Core.Player.CurrentState() {
	case action.ChargeAttackState:
		windup = 3
	case action.NormalAttackState:
		windup = -5
	}

	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Moondew Cleanse",
		AttackTag:        attacks.AttackTagDirectLunarBloom,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		Mult:             chargeLB[c.TalentLvlAttack()],
		UseHP:            true,
		IgnoreDefPercent: 1.0,
	}

	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, chargeRadiusBloom)
	for _, hitmark := range chargeHitmarkBloom {
		delay := hitmark + windup
		c.Core.QueueAttack(ai, ap, delay, delay)
	}

	c.QueueCharTask(func() { c.Core.Player.ConsumeDew(1) }, verdantDewFrame+windup)

	return action.Info{
		Frames: func(next action.Action) int {
			if next == action.ActionCharge && c.Core.Player.Dew() <= 0 {
				return 89 // CAb -> CA is 89
			}
			return chargeLBFrames[next]
		},
		AnimationLength: chargeLBFrames[action.InvalidAction],
		CanQueueAfter:   chargeLBFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}
}
