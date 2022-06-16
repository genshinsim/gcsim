package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) HighPlungeAttack(p map[string]int) action.ActionInfo {
	f, a := c.ActionFrames(action.ActionHighPlunge, p)
	ele := attributes.Physical
	if c.Core.LastAction.Target == core.Kazuha && c.Core.LastAction.Typ == action.ActionSkill {
		ele = attributes.Anemo
	}

	_, ok := p["collide"]
	if ok {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Plunge (Collide)",
			AttackTag:      combat.AttackTagPlunge,
			ICDTag:         combat.ICDTagNone,
			ICDGroup:       combat.ICDGroupDefault,
			Element:        ele,
			Durability:     0,
			Mult:           plunge[c.TalentLvlAttack()],
			IgnoreInfusion: true,
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.3, false, combat.TargettableEnemy), f, f)
	}

	//aoe dmg
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Plunge",
		AttackTag:      combat.AttackTagPlunge,
		ICDTag:         combat.ICDTagNone,
		ICDGroup:       combat.ICDGroupDefault,
		StrikeType:     combat.StrikeTypeBlunt,
		Element:        ele,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		IgnoreInfusion: true,
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), f, f)

	// a1 if applies
	if c.a1Ele != attributes.NoElement {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Kazuha A1",
			AttackTag:      combat.AttackTagPlunge,
			ICDTag:         combat.ICDTagNone,
			ICDGroup:       combat.ICDGroupDefault,
			StrikeType:     combat.StrikeTypeDefault,
			Element:        c.a1Ele,
			Durability:     25,
			Mult:           2,
			IgnoreInfusion: true,
		}

		c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), f-1, f-1)
		c.a1Ele = attributes.NoElement
	}

	return f, a
}
