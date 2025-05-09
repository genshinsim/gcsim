package citlali

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const (
	chargeHitmark = 50
	chargeRadius  = 0.8
)

func init() {
	chargeFrames = frames.InitAbilSlice(56) // CA -> Walk
	chargeFrames[action.ActionAttack] = 49
	chargeFrames[action.ActionCharge] = 52
	chargeFrames[action.ActionSkill] = 52
	chargeFrames[action.ActionBurst] = 51
	chargeFrames[action.ActionDash] = 45
	chargeFrames[action.ActionJump] = 46
	chargeFrames[action.ActionSwap] = 50
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Charge Attack",
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Cryo,
		Durability:   25,
		Mult:         charge[c.TalentLvlAttack()],
		HitlagFactor: 0.05,
	}

	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, chargeRadius)
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, ap, 0, travel)
	}, chargeHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}
