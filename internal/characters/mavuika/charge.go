package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int
var bikeChargeFrames []int
var bikeChargeHitmarks = []int{36, 78, 119, 165, 208, 252, 297, 341}
var bikeChargeFinalHitmark = 421

const chargeHitmark = 40

func init() {
	chargeFrames = frames.InitAbilSlice(48)
	chargeFrames[action.ActionBurst] = 50
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 50
	chargeFrames[action.ActionWalk] = 60

	bikeChargeFrames = frames.InitAbilSlice(430)
	bikeChargeFrames[action.ActionBurst] = 440
	bikeChargeFrames[action.ActionDash] = bikeChargeFinalHitmark
	bikeChargeFrames[action.ActionJump] = bikeChargeFinalHitmark
	bikeChargeFrames[action.ActionSwap] = 450
	bikeChargeFrames[action.ActionWalk] = 470
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		return c.bikeCharge(), nil
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   120.0,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: -1.8},
			2,
			4.5,
		),
		chargeHitmark,
		chargeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) bikeCharge() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Flamestride Charge",
		AttackTag:        attacks.AttackTagExtra,
		AdditionalTags:   []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:           attacks.ICDTagMavuikaFlamestrider,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         60.0,
		Element:          attributes.Pyro,
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.03 * 60,
		Durability:       25,
		Mult:             skillCharge[c.TalentLvlSkill()],
		IgnoreInfusion:   true,
	}

	for _, delay := range bikeChargeHitmarks {
		c.QueueCharTask(func() {
			ai.FlatDmg = c.burstBuffCA() + c.c2BikeCA()
			c.Core.QueueAttack(
				ai,
				combat.NewBoxHitOnTarget(
					c.Core.Combat.Player(),
					geometry.Point{Y: -1.8},
					2,
					4.5,
				),
				0,
				0,
			)
		}, delay)
	}

	c.QueueCharTask(func() {
		ai.Abil = "Flamestride Charge (Final)"
		ai.PoiseDMG = 120.0
		ai.HitlagHaltFrames = 0.04 * 60
		ai.Mult = skillChargeFinal[c.TalentLvlSkill()]
		ai.FlatDmg = c.burstBuffCA()

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: 1},
				4,
			),
			0,
			0,
		)
	}, bikeChargeFinalHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: bikeChargeFrames[action.InvalidAction],
		CanQueueAfter:   bikeChargeFinalHitmark,
		State:           action.ChargeAttackState,
	}
}
