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
	chargeHitmarks           = []int{41, 41}
	chargeHitlagHaltFrames   = []float64{0.0, 0.09}
	chargeCanBeDefenseHalted = []bool{false, true}
	chargePoiseDmg           = []float64{78, 42}

	azureDevourFrames             []int
	azureDevourHitmarks           = []int{40, 40, 40 + 20, 40 + 20}
	azureDevourHitlagHaltFrames   = []float64{0.0, 0.9, 0.0, 0.9}
	azureDevourCanBeDefenseHalted = []bool{false, true, false, true}
	azureDevourPoiseDmg           = []float64{39, 21, 39, 21}
)

func init() {
	chargeFrames = frames.InitAbilSlice(67)
	chargeFrames[action.ActionAttack] = 58
	chargeFrames[action.ActionBurst] = 50
	chargeFrames[action.ActionSkill] = 59
	chargeFrames[action.ActionBurst] = 58
	chargeFrames[action.ActionDash] = chargeHitmarks[1]
	chargeFrames[action.ActionJump] = chargeHitmarks[1]
	chargeFrames[action.ActionSwap] = 49
	chargeFrames[action.ActionWalk] = 49

	azureDevourFrames = frames.InitAbilSlice(75)
	azureDevourFrames[action.ActionAttack] = 64
	azureDevourFrames[action.ActionCharge] = 74
	azureDevourFrames[action.ActionSkill] = 66
	azureDevourFrames[action.ActionBurst] = 65
	azureDevourFrames[action.ActionDash] = azureDevourHitmarks[3]
	azureDevourFrames[action.ActionJump] = azureDevourHitmarks[3]
	azureDevourFrames[action.ActionSwap] = azureDevourHitmarks[3]
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillCharge()
	}

	windup := 14
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.ChargeAttackState:
		windup = 0
	}

	for i, hitmark := range chargeHitmarks {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               "Charge",
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           chargePoiseDmg[i],
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
				info.Point{Y: 1},
				2.5,
			),
			hitmark,
			hitmark,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFuncWithOffset(chargeFrames, windup),
		AnimationLength: chargeFrames[action.InvalidAction] + windup,
		CanQueueAfter:   chargeHitmarks[1] + windup,
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) skillCharge() (action.Info, error) {
	if c.fourWindsCharges() > 0 || c.c6FreeCA() {
		return c.skillAzureDevour(c.c6FreeCA())
	}

	windup := 14
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.ChargeAttackState:
		windup = 0
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
			PoiseDMG:           chargePoiseDmg[i],
			Element:            ele[i],
			Durability:         25,
			Mult:               charge[i][c.TalentLvlSkill()] * c.a1SkillMulti(),
			HitlagHaltFrames:   chargeHitlagHaltFrames[i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: chargeCanBeDefenseHalted[i],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: 1},
				2.5,
			),
			hitmark,
			hitmark,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFuncWithOffset(chargeFrames, windup),
		AnimationLength: chargeFrames[action.InvalidAction] + windup,
		CanQueueAfter:   chargeHitmarks[1] + windup,
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) skillAzureDevour(c6Free bool) (action.Info, error) {
	windup := 14
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.ChargeAttackState:
		windup = 0
	}

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
			PoiseDMG:           azureDevourPoiseDmg[i],
			Element:            ele[i],
			Durability:         25,
			Mult:               skillAzureDevour[i][c.TalentLvlSkill()] * c.a1SkillMulti() * c1Mult,
			HitlagHaltFrames:   azureDevourHitlagHaltFrames[i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: azureDevourCanBeDefenseHalted[i],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: 1},
				3,
			),
			hitmark,
			hitmark,
		)
	}
	if !c6Free {
		c.useFourWindsCharge()
		c.c6OnSkillCA()
	}

	c.c2OnSpecialSkill()

	return action.Info{
		Frames:          frames.NewAbilFuncWithOffset(azureDevourFrames, windup),
		AnimationLength: azureDevourFrames[action.InvalidAction] + windup,
		CanQueueAfter:   azureDevourHitmarks[3] + windup,
		State:           action.ChargeAttackState,
	}, nil
}
