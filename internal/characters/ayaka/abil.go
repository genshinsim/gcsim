package ayaka

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			fmt.Sprintf("Normal %v", c.NormalCounter),
			core.AttackTagNormal,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			core.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, f-5+i)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	d := c.Snapshot(
		"Charge",
		core.AttackTagNormal,
		core.ICDTagExtraAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		ca[c.TalentLvlAttack()],
	)
	d.Targets = core.TargetAll

	for i := 0; i < 3; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f-3+i)
	}

	return f, a
}

func (c *char) Dash(p map[string]int) (int, int) {
	f, ok := p["f"]
	if !ok {
		f = 36
	}
	//no dmg attack at end of dash
	d := c.Snapshot(
		"Dash",
		core.AttackTagNone,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		25,
		0,
	)
	c.QueueDmg(&d, f)
	//since we always hit, just restore the stam and add bonus...
	c.AddTask(func() {
		c.Core.RestoreStam(10)
		var val [core.EndStatType]float64
		val[core.CryoP] = 0.18
		//a2 increase normal + ca dmg by 30% for 6s
		c.AddMod(core.CharStatMod{
			Key:    "ayaka-a4",
			Expiry: c.Core.F + 600,
			Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
				return val, true
			},
		})
	}, "ayaka-dash", f+1)
	//add cryo infuse
	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "ayaka-dash",
		Ele:    core.Cryo,
		Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.F + 300,
	})
	return f, f
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Hyouka",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		50,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	//2 or 3 1:1 ratio
	count := 4
	if c.Core.Rand.Float64() < 0.5 {
		count = 5
	}
	c.QueueParticle("ayaka", count, core.Cryo, f+100)

	var val [core.EndStatType]float64
	val[core.DmgP] = 0.3
	//a2 increase normal + ca dmg by 30% for 6s
	c.AddMod(core.CharStatMod{
		Key:    "ayaka-a2",
		Expiry: c.Core.F + 360,
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			return val, a == core.AttackTagNormal || a == core.AttackTagExtra
		},
	})

	c.QueueDmg(&d, f)

	c.SetCD(core.ActionSkill, 600)
	return f, a

}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	d := c.Snapshot(
		"Soumetsu",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		25,
		burstCut[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll
	db := c.Snapshot(
		"Soumetsu",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		25,
		burstBloom[c.TalentLvlBurst()],
	)
	db.Targets = core.TargetAll

	//5 second, 20 ticks, so once every 15 frames, bloom after 5 seconds
	c.QueueDmg(&db, f+300)
	for i := 0; i < 300; i += 15 {
		x := d.Clone()
		c.QueueDmg(&x, f+i)
	}

	c.SetCD(core.ActionBurst, 20*60)
	c.Energy = 0

	return f, a
}
