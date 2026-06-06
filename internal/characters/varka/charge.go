package varka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	chargeFrames             []int
	chargeHitmarks           = []int{30, 40}
	chargeHitlagHaltFrames   = []float64{0.0, 0.1}
	chargeCanBeDefenseHalted = []bool{false, true}

	skillChargeFrames             []int
	skillChargeHitmarks           = []int{30, 40}
	skillChargeHitlagHaltFrames   = []float64{0.0, 0.1}
	skillChargeCanBeDefenseHalted = []bool{false, true}

	azureDevourFrames             []int
	azureDevourHitmarks           = []int{28, 36, 50, 58}
	azureDevourHitlagHaltFrames   = []float64{0.0, 0.1, 0.0, 0.1}
	azureDevourCanBeDefenseHalted = []bool{false, true, false, true}
)

func init() {
	chargeFrames = frames.InitAbilSlice(48)
	chargeFrames[action.ActionBurst] = 50
	chargeFrames[action.ActionDash] = chargeHitmarks[1]
	chargeFrames[action.ActionJump] = chargeHitmarks[1]
	chargeFrames[action.ActionSwap] = 50
	chargeFrames[action.ActionWalk] = 60

	skillChargeFrames = frames.InitAbilSlice(80)
	skillChargeFrames[action.ActionBurst] = 60
	skillChargeFrames[action.ActionDash] = skillChargeHitmarks[1]
	skillChargeFrames[action.ActionJump] = skillChargeHitmarks[1]
	skillChargeFrames[action.ActionSwap] = 60
	skillChargeFrames[action.ActionWalk] = 60

	azureDevourFrames = frames.InitAbilSlice(80)
	azureDevourFrames[action.ActionBurst] = 60
	azureDevourFrames[action.ActionDash] = azureDevourHitmarks[3]
	azureDevourFrames[action.ActionJump] = azureDevourHitmarks[3]
	azureDevourFrames[action.ActionSwap] = 60
	azureDevourFrames[action.ActionWalk] = 60
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillCharge()
	}

	for i, hitmark := range chargeHitmarks {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               "Charge",
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           120.0,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               charge[i][c.TalentLvlAttack()],
			HitlagHaltFrames:   chargeHitlagHaltFrames[i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: chargeCanBeDefenseHalted[i],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: 0.3},
				3.3,
			),
			hitmark,
			hitmark,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[1],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) skillCharge() (action.Info, error) {
	if c.fourWindsCharges() > 0 {
		return c.skillAzureDevour()
	}
	ele := []attributes.Element{c.conversionElem, attributes.Anemo}
	for i, hitmark := range chargeHitmarks {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               "Sturm und Drang Charge",
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           120.0,
			Element:            ele[i],
			Durability:         25,
			Mult:               charge[i][c.TalentLvlAttack()] * c.a1SkillMulti(),
			HitlagHaltFrames:   skillChargeHitlagHaltFrames[i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: skillChargeCanBeDefenseHalted[i],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: 0.3},
				3.3,
			),
			hitmark,
			hitmark,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillChargeFrames),
		AnimationLength: skillChargeFrames[action.InvalidAction],
		CanQueueAfter:   skillChargeHitmarks[1],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) skillAzureDevour() (action.Info, error) {
	ele := []attributes.Element{c.conversionElem, attributes.Anemo, c.conversionElem, attributes.Anemo}

	c1Mult := c.c1OnSpecialSkill()
	for i, hitmark := range azureDevourHitmarks {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               "Azure Devour",
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           120.0,
			Element:            ele[i],
			Durability:         25,
			Mult:               skillAzureDevour[i][c.TalentLvlAttack()] * c.a1SkillMulti() * c1Mult,
			HitlagHaltFrames:   azureDevourHitlagHaltFrames[i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: azureDevourCanBeDefenseHalted[i],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: 0.3},
				3.3,
			),
			hitmark,
			hitmark,
		)
	}
	c.useFourWindsCharge()
	return action.Info{
		Frames:          frames.NewAbilFunc(azureDevourFrames),
		AnimationLength: azureDevourFrames[action.InvalidAction],
		CanQueueAfter:   azureDevourHitmarks[3],
		State:           action.ChargeAttackState,
	}, nil
}
