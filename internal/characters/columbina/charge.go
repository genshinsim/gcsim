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
	chargeHitmarkBloom = []int{50, 50 + 8, 50 + 8 + 22}
	verdantDewFrame    = 39
)

const (
	chargeHitmark     = 60
	chargeRadius      = 3.5
	chargeRadiusBloom = 3.3
)

func init() {
	chargeFrames = frames.InitAbilSlice(81) // CA -> Walk
	chargeFrames[action.ActionAttack] = 81
	chargeFrames[action.ActionCharge] = 81
	chargeFrames[action.ActionSkill] = verdantDewFrame
	chargeFrames[action.ActionBurst] = verdantDewFrame
	chargeFrames[action.ActionDash] = verdantDewFrame
	chargeFrames[action.ActionJump] = verdantDewFrame
	chargeFrames[action.ActionSwap] = verdantDewFrame
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.Dew() > 0 {
		return c.ChargeAttackBloom(p), nil
	}

	windup := 12
	switch c.Core.Player.CurrentState() {
	case action.ChargeAttackState:
		windup = 3
	case action.NormalAttackState:
		windup = 0
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

	// assuming actions have same frames?
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
		windup = 0
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
		IsDeployable:     true,
	}

	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, chargeRadiusBloom)
	for _, hitmark := range chargeHitmarkBloom {
		delay := hitmark + windup
		c.Core.QueueAttack(ai, ap, delay, delay)
	}

	c.QueueCharTask(func() { c.Core.Player.ConsumeDew(1) }, verdantDewFrame+windup)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}
}
