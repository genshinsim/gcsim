package arlecchino

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 34

func init() {
	chargeFrames = frames.InitAbilSlice(59)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		// occurs before attack lands
		c.absorbDirectives()
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Charge",
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagExtraAttack,
			ICDGroup:           attacks.ICDGroupPoleExtraAttack,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagHaltFrames:   0.02,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
			Mult:               charge[c.TalentLvlAttack()],
		}

		if c.StatusIsActive(naBuffKey) {
			ai.Element = attributes.Pyro
			ai.IgnoreInfusion = true
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.8,
			),
			0,
			0,
		)
	}, chargeHitmark)
	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) absorbDirectives() {
	for _, e := range c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), nil) {
		if !e.StatusIsActive(directiveKey) {
			continue
		}

		level := e.GetTag(directiveKey)

		newDebt := a1Directive[level] * c.MaxHP()
		if c.StatusIsActive(directiveLimitKey) {
			newDebt = min(c.skillDebtMax-c.skillDebt, newDebt)
		}

		if newDebt > 0 {
			c.skillDebt += newDebt
			c.ModifyHPDebtByAmount(newDebt)
		}
		e.RemoveTag(directiveKey)
		e.RemoveTag(directiveSrcKey)
		e.DeleteStatus(directiveKey)

		c.c4OnAbsorb()
		if level >= 2 {
			c.c2OnAbsorbDue()
		}
	}
}
