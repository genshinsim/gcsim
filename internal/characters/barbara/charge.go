package barbara

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// Charge attack function - handles seal use
func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {

	f, a := c.ActionFrames(action.ActionCharge, p)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	done := false
	// Taken from Noelle code
	cb := func(a combat.AttackCB) {
		if done { //why do we need this @srl
			return
		}
		//check for healing
		if c.Core.Status.Duration("barbskill") > 0 {
			//heal target
			c.Core.Health.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Melody Loop (Charged Attack)",
				Src:     4 * (prochpp[c.TalentLvlSkill()]*c.MaxHP() + prochp[c.TalentLvlSkill()]),
				Bonus:   c.Stat(attributes.Heal),
			})
			done = true
		}

	}
	var cbenergy func(a combat.AttackCB)
	energyCount := 0
	if c.Base.Cons >= 4 {
		cbenergy = func(a combat.AttackCB) {
			//check for healing
			if c.Core.Status.Duration("barbskill") > 0 && energyCount < 5 {
				//regen energy
				c.AddEnergy("barbara-c4", 1)
				energyCount++
			}

		}
	}

	// TODO: Not sure of snapshot timing
	c.Core.Combat.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, f, cb, cbenergy)

	return f, a
}
