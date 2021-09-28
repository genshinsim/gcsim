package travelerelectro

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Lightning Blade",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		skill[c.TalentLvlSkill()],
	)
	//d.Targets = def.TargetAll

	hits, ok := p["hits"]
	if !ok {
		hits = 1
	}

	c.QueueParticle("travelerelectro", 1, core.Cryo, f+100)

	for i := 0; i < hits; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f)
	}

	c.SetCD(core.ActionSkill, 810+21) //13.5s, starts 21 frames in
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	d := c.Snapshot(
		"Bellowing Thunder",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		burst[c.TalentLvlBurst()],
	)

	c.QueueDmg(&d, f)

	//1573 start, 1610 cd starts, 1612 energy drained, 1633 first swapable
	c.ConsumeEnergy(39)
	c.SetCD(core.ActionBurst, 1200+37)
	return f, a
}
