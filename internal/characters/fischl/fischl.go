package fischl

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("fischl", NewChar)
}

type char struct {
	*character.Tmpl
	//field use for calculating oz damage
	ozSnapshot core.Snapshot

	ozSource      int //keep tracks of source of oz aka resets
	ozActiveUntil int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5

	c.ozSource = -1
	c.ozActiveUntil = -1

	//register A4
	c.a4()

	if p.Base.Cons == 6 {
		c.c6()
	}

	// f.turbo()

	return &c, nil
}

// func (c *char) ozAttack() {
// 	d := c.ozSnapshot.Clone()
// 	d.Durability = 0
// 	if c.ozAttackCounter%4 == 0 {
// 		//apply aura, add to timer
// 		d.Durability = 25
// 		c.ozICD = c.Core.F + 300 //add 300 second to skill ICD
// 	}
// 	//so oz is active and ready to shoot, we add damage
// 	c.S.AddTask(func(s *def.Sim) {
// 		s.ApplyDamage(d)
// 	}, "Fischl Oz (Damage)", 1)
// 	//put shoot on cd
// 	c.ozNextShootReady = c.Core.F + 50
// 	//increment hit counter
// 	c.ozAttackCounter++
// 	//assume fischl has 60% chance of generating orb every attack;
// 	if c.S.Rand.Float64() < .6 {
// 		c.S.AddEnergyParticles("Fischl", 1, def.Electro, 120)
// 	}
// }

func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Physical,
		25,
		auto[c.NormalCounter][c.TalentLvlAttack()],
	)
	c.QueueDmg(&d, travel+f)
	c.AdvanceNormalIndex()

	//check for c1
	if c.Base.Cons >= 1 && c.ozActiveUntil < c.Core.F {
		d := c.Snapshot(
			"Fischl C1",
			core.AttackTagNormal,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Physical,
			100,
			0.22,
		)
		c.QueueDmg(&d, travel+f+1)
	}

	return f, a
}

func (c *char) queueOz(src string) {

	dur := 600
	if c.Base.Cons == 6 {
		dur += 120
	}
	c.ozActiveUntil = c.Core.F + dur
	c.ozSource = c.Core.F
	c.ozSnapshot = c.Snapshot(
		fmt.Sprintf("Oz (%v)", src),
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupFischl,
		core.StrikeTypePierce,
		core.Electro,
		25,
		birdAtk[c.TalentLvlSkill()],
	)
	c.AddTask(c.ozTick(c.Core.F), "oz", 60)
	c.Core.Log.Debugw("Oz activated", "frame", c.Core.F, "event", core.LogCharacterEvent, "source", src, "expected end", c.ozActiveUntil, "next expected tick", c.Core.F+60)

	c.Core.Status.AddStatus("fischloz", dur)

}

func (c *char) ozTick(src int) func() {
	return func() {
		c.Core.Log.Debugw("Oz checking for tick", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src)
		//if src != ozSource then this is no longer the same oz, do nothing
		if src != c.ozSource {
			return
		}
		c.Core.Log.Debugw("Oz ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "next expected tick", c.Core.F+60, "active", c.ozActiveUntil, "src", src)
		//trigger damage
		d := c.ozSnapshot.Clone()
		c.Core.Combat.ApplyDamage(&d)
		//check for orb
		//Particle check is 67% for particle, from datamine
		if c.Core.Rand.Float64() < .67 {
			c.QueueParticle("fischl", 1, core.Electro, 120)
		}

		//queue up next hit only if next hit oz is still active
		if c.Core.F+60 <= c.ozActiveUntil {
			c.AddTask(c.ozTick(src), "oz", 60)
		}
	}
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	//always trigger electro no ICD on initial summon
	d := c.Snapshot(
		"Oz (Summon)",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupFischl,
		core.StrikeTypePierce,
		core.Electro,
		25,
		birdSum[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	if c.Base.Cons >= 2 {
		d.Mult += 2
	}
	c.QueueDmg(&d, 1) //queue initial damage

	//set on field oz to be this one
	c.AddTask(func() {
		c.queueOz("Skill")
	}, "oz-skill", f-20)

	c.SetCD(core.ActionSkill, 25*60)
	//return animation cd
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	//set on field oz to be this one
	c.AddTask(func() {
		c.queueOz("Burst")
	}, "oz-skill", f-10)

	//initial damage; part of the burst tag
	d := c.Snapshot(
		"Midnight Phantasmagoria",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupFischl,
		core.StrikeTypeBlunt,
		core.Electro,
		25,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, 1)

	//check for C4 damage
	if c.Base.Cons >= 4 {
		d := c.Snapshot(
			"Midnight Phantasmagoria",
			core.AttackTagElementalBurst,
			core.ICDTagElementalBurst,
			core.ICDGroupFischl,
			core.StrikeTypePierce,
			core.Electro,
			50,
			2.22,
		)
		c.QueueDmg(&d, 1)
		//heal at end of animation
		heal := c.MaxHP() * 0.2
		c.AddTask(func() {
			c.Core.Health.HealActive(c.Index, heal)
		}, "c4heal", f-1)

	}

	c.Energy = 0
	c.SetCD(core.ActionBurst, 15*60)
	return f, a
}
