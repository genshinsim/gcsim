package barbara

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) action.ActionInfo {
	f, a := c.ActionFrames(action.ActionAttack, p)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
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
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Melody Loop (Normal Attack)",
				Src:     prochpp[c.TalentLvlSkill()]*c.MaxHP() + prochp[c.TalentLvlSkill()],
				Bonus:   c.Stat(attributes.Heal),
			})
			done = true
		}

	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 0, f, cb)
	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}
