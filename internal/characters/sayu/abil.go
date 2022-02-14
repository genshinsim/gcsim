package sayu

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(r, false, core.TargettableEnemy), f, f, cb)

	c.AdvanceNormalIndex()
	return f, a
}

type noelleShield struct {
	*shield.Tmpl
	c *char
}

func (n *noelleShield) OnExpire() {
	if n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
}

func (n *noelleShield) OnDamage(dmg float64, ele core.EleType, bonus float64) (float64, bool) {
	taken, ok := n.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok && n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
	return taken, ok
}

func (c *char) newShield(base float64, t core.ShieldType, dur int) *noelleShield {
	n := &noelleShield{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.Src = c.Core.F
	n.Tmpl.ShieldType = t
	n.Tmpl.Name = "Noelle Skill"
	n.Tmpl.HP = base
	n.Tmpl.Expires = c.Core.F + dur
	n.c = c
	return n
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Initial)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skillpress[c.TalentLvlSkill()],
	}

	c.SetCD(core.ActionSkill, 24*60)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// TODO: Not sure of the exact frames
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(4.5, false, core.TargettableEnemy), f-10, f-10)

	c.SetCDWithDelay(core.ActionBurst, 900, 8)
	c.ConsumeEnergy(8)
	return f, a
}
