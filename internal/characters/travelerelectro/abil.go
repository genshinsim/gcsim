package travelerelectro

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/def"
)

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(def.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	c.AdvanceNormalIndex()

	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)
	d := c.Snapshot(
		"Lightning Blade",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Electro,
		25,
		skill[c.TalentLvlSkill()],
	)
	//d.Targets = def.TargetAll

	hits, ok := p["hits"]
	if !ok {
		hits = 1
	}

	c.QueueParticle("travelerelectro", 1, def.Cryo, f+100)

	for i := 0; i < hits; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f)
	}

	c.SetCD(def.ActionSkill, 810+21) //13.5s, starts 21 frames in
	return f
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(def.ActionBurst, p)
	d := c.Snapshot(
		"Bellowing Thunder",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Electro,
		25,
		burst[c.TalentLvlBurst()],
	)

	c.QueueDmg(&d, f)

	//1573 start, 1610 cd starts, 1612 energy drained, 1633 first swapable
	c.ConsumeEnergy(39)
	c.SetCD(def.ActionBurst, 1200+37)
	return f
}
