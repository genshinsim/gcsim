package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	chargeFrames     []int
	chargeBikeFrames []int
)

// 22, 39, 42, 45, 43, 46, 44, ... - 43.1 avg
const (
	chargeHitmark             = 52
	bikeChargeHitmarkInitial  = 22
	bikeChargeHitmarkInterval = 43
	bikeChargeFinalHitmark    = 83
)

func init() {
	// imaginary numbers
	chargeFrames = frames.InitAbilSlice(60)
	chargeBikeFrames = frames.InitAbilSlice(80)
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() && c.allFireArmamnetsActive {
		return c.bikeChargeAttack(p)
	}
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Mult:               3.83,
		Element:            attributes.Physical,
		Durability:         25,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player().Pos(), geometry.Point{Y: 0.45}, 3.0, 5.0),
		chargeHitmark,
		chargeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionAttack],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) bikeChargeAttack(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = 1
	}
	switch {
	case hold < 1:
		hold = 1
	case hold > 8:
		hold = 8
	}

	var apCyclic combat.AttackPattern
	holdDur := 0
	for i := 0; i < hold; i++ {
		aiCyclic := c.getCyclicAi()
		if i == 0 {
			holdDur += bikeChargeHitmarkInitial
		} else {
			holdDur += bikeChargeHitmarkInterval
		}
		c.QueueCharTask(func() {
			if c.StatusIsActive(crucibleOfDeathAndLifeStatus) {
				aiCyclic.FlatDmg += 0.0144 * c.TotalAtk() * float64(c.consumedFightingSpirit)
			}
			aiCyclic.FlatDmg += c.c2FlatIncrease(attacks.AttackTagExtra)
			apCyclic = combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 4)
			c.Core.QueueAttack(aiCyclic, apCyclic, 0, 0)
		}, holdDur)
	}

	aiFinal := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Flamestrider Charged Attack Final DMG",
		AttackTag:      attacks.AttackTagExtra,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		Element:        attributes.Pyro,
		Durability:     25,
		FlatDmg:        c.TotalAtk() * 3.162,
		IgnoreInfusion: true,
	}
	apFinal := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 4.5)
	c.QueueCharTask(func() {
		if c.StatusIsActive(crucibleOfDeathAndLifeStatus) {
			aiFinal.FlatDmg += 0.0144 * c.TotalAtk() * float64(c.consumedFightingSpirit)
		}
		aiFinal.FlatDmg += c.c2FlatIncrease(attacks.AttackTagExtra)
		c.Core.QueueAttack(aiFinal, apFinal, 0, 0)
	}, holdDur+bikeChargeFinalHitmark)

	return action.Info{
		Frames: func(next action.Action) int {
			return holdDur + bikeChargeFinalHitmark + chargeBikeFrames[next]
		},
		AnimationLength: holdDur + bikeChargeFinalHitmark,
		CanQueueAfter:   holdDur + bikeChargeFinalHitmark + chargeBikeFrames[action.ActionAttack], // change to earliest
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) getCyclicAi() combat.AttackInfo {
	return combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Flamestrider Charged Attack Cyclic DMG",
		AttackTag:      attacks.AttackTagExtra,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagExtraAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		Element:        attributes.Pyro,
		Durability:     25,
		FlatDmg:        c.TotalAtk() * 2.17,
		IgnoreInfusion: true,
	}
}
