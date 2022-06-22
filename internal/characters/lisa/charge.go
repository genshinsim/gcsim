package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const conductiveTag = "lisa-conductive-stacks"

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	f, a := c.ActionFrames(action.ActionCharge, p)

	//TODO: assumes this applies every time per
	//[7:53 PM] Hold ₼KLEE like others hold GME: CHarge is pyro every charge
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	done := false
	cb := func(a combat.AttackCB) {
		if done {
			return
		}
		count := a.Target.GetTag(conductiveTag)
		if count < 3 {
			a.Target.SetTag(conductiveTag, count+1)
		}
		done = true
	}

	c.Core.Combat.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 0, f, cb)

	return f, a
}
