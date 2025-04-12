package mizuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstHitmark            = 93
	burstDurability         = 50
	burstPoise              = 100
	burstDuration           = 12 * 60
	burstCdDelay            = 1
	burstEnergyDrainDelay   = 4
	burstCd                 = 15 * 60
	burstRadius             = 8
	snackInterval           = 1.5 * 60
	snackHitmark            = 93
	snackDurability         = 25
	snackPoise              = 30
	snackRadius             = 4
	snackHealTriggerHpRatio = 0.7
	burstKey                = "mizuki-burst"
)

func init() {
	burstFrames = frames.InitAbilSlice(93)
	burstFrames[action.ActionCharge] = 92 // Q -> CA
	burstFrames[action.ActionDash] = 91   // Q -> D
	burstFrames[action.ActionWalk] = 92   // Q -> Walk
	burstFrames[action.ActionSwap] = 94   // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Anraku Secret Spring Therapy",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: burstDurability,
		PoiseDMG:   burstPoise,
		Mult:       burstDMG[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, burstRadius), burstHitmark, burstHitmark)

	snackAttack := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Munen Shockwave",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: snackDurability,
		PoiseDMG:   snackPoise,
		Mult:       snackDMG[c.TalentLvlBurst()],
	}

	c.AddStatus(burstKey, burstDuration, false)

	var hitFunc func()
	hitFunc = func() {
		if !c.StatusIsActive(burstKey) {
			return
		}

		if c.Core.Player.ActiveChar().CurrentHP() > (c.Core.Player.ActiveChar().MaxHP() * snackHealTriggerHpRatio) {
			c.Core.QueueAttack(snackAttack, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, snackRadius), 0, snackHitmark)
		} else {
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "Snack Pick-Up",
				Src:     (c.Stat(attributes.EM) * snackHealEM[c.TalentLvlBurst()]) + snackHealFLat[c.TalentLvlBurst()],
				Bonus:   c.Stat(attributes.Heal),
			})
		}

		c.QueueCharTask(hitFunc, snackInterval)
	}
	c.QueueCharTask(hitFunc, snackInterval)

	c.ConsumeEnergy(burstEnergyDrainDelay)
	c.SetCDWithDelay(action.ActionBurst, burstCd, burstCdDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
