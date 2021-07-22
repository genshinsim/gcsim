package ayaka

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/def"
)

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(def.ActionAttack, p)

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			fmt.Sprintf("Normal %v", c.NormalCounter),
			def.AttackTagNormal,
			def.ICDTagNormalAttack,
			def.ICDGroupDefault,
			def.StrikeTypeSlash,
			def.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, f-5+i)
	}

	c.AdvanceNormalIndex()

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {

	f := c.ActionFrames(def.ActionCharge, p)

	d := c.Snapshot(
		"Charge",
		def.AttackTagNormal,
		def.ICDTagExtraAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
		25,
		ca[c.TalentLvlAttack()],
	)
	d.Targets = def.TargetAll

	for i := 0; i < 3; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f-3+i)
	}

	return f
}

func (c *char) Dash(p map[string]int) int {
	f, ok := p["f"]
	if !ok {
		f = 36
	}
	//no dmg attack at end of dash
	d := c.Snapshot(
		"Dash",
		def.AttackTagNone,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Cryo,
		25,
		0,
	)
	c.QueueDmg(&d, f)
	//since we always hit, just restore the stam and add bonus...
	c.AddTask(func() {
		c.Sim.RestoreStam(10)
		val := make([]float64, def.EndStatType)
		val[def.CryoP] = 0.18
		//a2 increase normal + ca dmg by 30% for 6s
		c.AddMod(def.CharStatMod{
			Key:    "ayaka-a4",
			Expiry: c.Sim.Frame() + 600,
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return val, true
			},
		})
	}, "ayaka-dash", f+1)
	//add cryo infuse
	c.AddWeaponInfuse(def.WeaponInfusion{
		Key:    "ayaka-dash",
		Ele:    def.Cryo,
		Tags:   []def.AttackTag{def.AttackTagNormal, def.AttackTagExtra, def.AttackTagPlunge},
		Expiry: c.Sim.Frame() + 300,
	})
	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)
	d := c.Snapshot(
		"Hyouka",
		def.AttackTagElementalArt,
		def.ICDTagElementalArt,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Cryo,
		50,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	//2 or 3 1:1 ratio
	count := 4
	if c.Sim.Rand().Float64() < 0.5 {
		count = 5
	}
	c.QueueParticle("ayaka", count, def.Cryo, f+100)

	val := make([]float64, def.EndStatType)
	val[def.DmgP] = 0.3
	//a2 increase normal + ca dmg by 30% for 6s
	c.AddMod(def.CharStatMod{
		Key:    "ayaka-a2",
		Expiry: c.Sim.Frame() + 360,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return val, a == def.AttackTagNormal || a == def.AttackTagExtra
		},
	})

	c.QueueDmg(&d, f)

	c.SetCD(def.ActionSkill, 600)
	return f

}

func (c *char) Burst(p map[string]int) int {

	f := c.ActionFrames(def.ActionBurst, p)

	d := c.Snapshot(
		"Soumetsu",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Cryo,
		25,
		burstCut[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll
	db := c.Snapshot(
		"Soumetsu",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Cryo,
		25,
		burstBloom[c.TalentLvlBurst()],
	)
	db.Targets = def.TargetAll

	//5 second, 20 ticks, so once every 15 frames, bloom after 5 seconds
	c.QueueDmg(&db, f+300)
	for i := 0; i < 300; i += 15 {
		x := d.Clone()
		c.QueueDmg(&x, f+i)
	}

	c.SetCD(def.ActionBurst, 20*60)
	c.Energy = 0

	return f
}
