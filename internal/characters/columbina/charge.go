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
	chargeHitmarkBloom = []int{60, 65, 70}
	verdantDewFrame    = 57
)

const (
	chargeHitmark = 60
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
	c.QueueCharTask(func() {
		if c.Core.Player.VerdantDew() > 0 {
			c.ChargeAttackBloom(p)
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
		delay := chargeHitmark - verdantDewFrame
		c.Core.QueueAttack(ai, ap, delay, delay)
	}, verdantDewFrame)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) ChargeAttackBloom(p map[string]int) {
	ai := info.AttackInfo{
		ActorIndex:   c.Index(),
		Abil:         "Moondew Cleanse",
		AttackTag:    attacks.AttackTagDirectLunarBloom,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Dendro,
		Mult:         chargeLB[c.TalentLvlAttack()],
		UseHP:        true,
		IsDeployable: true,
	}

	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, chargeRadius)
	for _, hitmark := range chargeHitmarkBloom {
		delay := hitmark - verdantDewFrame
		c.Core.QueueAttack(ai, ap, delay, delay)
	}
	c.Core.Player.ConsumeVerdantDew(1)
}
