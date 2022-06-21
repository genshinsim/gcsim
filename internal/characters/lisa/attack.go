package lisa

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var hitmarks = []int{26, 18, 17, 31}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	f, a := c.ActionFrames(action.ActionAttack, p)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagLisaElectro,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	//todo: Does it really snapshot immediately?
	c.Core.Combat.QueueAttack(ai, combat.NewDefSingleTarget(1, combat.TargettableEnemy), 0, hitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return f, a
}
