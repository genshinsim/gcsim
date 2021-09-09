package jean

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("jean", NewChar)
}

type char struct {
	*character.Tmpl
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t

	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 14 //frames from keqing lib
		case 1:
			f = 37 - 14
		case 2:
			f = 66 - 37
		case 3:
			f = 124 - 66
		case 4:
			f = 159 - 124
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 73
	case core.ActionSkill:

		hold := p["hold"]
		//hold for p up to 5 seconds
		if hold > 300 {
			hold = 300
		}

		return 23 + hold
	case core.ActionBurst:
		return 88
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		auto[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	//check for healing
	if c.Core.Rand.Float64() < 0.5 {
		heal := 0.15 * (d.BaseAtk*(1+d.Stats[core.ATKP]) + d.Stats[core.ATK])
		c.AddTask(func() {
			c.Core.Health.HealAll(heal)
		}, "jean-heal", f-1)
	}

	c.AdvanceNormalIndex()

	return f
}

func (c *char) Skill(p map[string]int) int {

	f := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Gale Blade",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		50,
		skill[c.TalentLvlSkill()],
	)

	if c.Base.Cons >= 1 && p["hold"] >= 60 {
		//add 40% dmg
		d.Stats[core.DmgP] += .4
		c.Core.Log.Debugw("jean c1 adding 40% dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "final dmg%", d.Stats[core.DmgP])
	}

	c.QueueDmg(&d, f-1)

	count := 2
	if c.Core.Rand.Float64() < .67 {
		count++
	}
	c.QueueParticle("Jean", count, core.Anemo, f+100)

	c.SetCD(core.ActionSkill, 360)
	return f //TODO: frames, + p for holding
}

func (c *char) Burst(p map[string]int) int {
	//p is the number of times enemy enters or exits the field
	enter := p["enter"]
	delay, ok := p["delay"]
	if !ok {
		delay = 10
	}

	f := c.ActionFrames(core.ActionBurst, p)
	d := c.Snapshot(
		"Dandelion Breeze",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		50,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f-10) //TODO: initial damage frames

	for i := 0; i < enter; i++ {
		x := d.Clone()
		x.Abil = "Dandelion Breeze (In/Out)"
		x.Mult = burstEnter[c.TalentLvlBurst()]

		c.QueueDmg(&x, f+i*delay)
	}

	c.Core.Status.AddStatus("jeanq", 630)

	if c.Base.Cons >= 4 {
		//add debuff to all target for ??? duration
		for _, t := range c.Core.Targets {
			t.AddResMod("jeanc4", core.ResistMod{
				Duration: 600, //10 seconds
				Ele:      core.Anemo,
				Value:    -0.4,
			})
		}
	}

	//heal on cast
	hpplus := d.Stats[core.Heal]
	atk := d.BaseAtk*(1+d.Stats[core.ATKP]) + d.Stats[core.ATK]
	heal := (burstInitialHealFlat[c.TalentLvlBurst()] + atk*burstInitialHealPer[c.TalentLvlBurst()]) * (1 + hpplus)
	healDot := (burstDotHealFlat[c.TalentLvlBurst()] + atk*burstDotHealPer[c.TalentLvlBurst()]) * (1 + hpplus)

	c.AddTask(func() {
		c.Core.Health.HealAll(heal)
	}, "Jean Heal Initial", f)

	//duration is 10.5s
	for i := 60; i < 630; i++ {
		c.AddTask(func() {
			c.Core.Log.Debugw("jean q healing", "frame", c.Core.F, "event", core.LogCharacterEvent, "+heal", hpplus, "atk", atk, "heal amount", healDot)
			c.Core.Health.HealActive(heal)
		}, "Jean Tick", i)
	}

	c.SetCD(core.ActionBurst, 1200)
	//c.Energy = 16 //jean a4
	c.ConsumeEnergy(16, 0) //at 0,0 value acts the same as c.Energy = 0
	return f               //TODO: frames, + p for holding
}

func (c *char) c6() {
	//reduce dmg by 35% if q active, ignoring the lingering affect
	// c.Sim.AddDRFunc(func() float64 {
	// 	if c.S.StatusActive("jeanq") {
	// 		return 0.35
	// 	}
	// 	return 0
	// })
	c.Core.Log.Warnw("jean c6 not implemented", "frame", c.Core.F, "event", core.LogCharacterEvent)
}

func (c *char) ReceiveParticle(p core.Particle, isActive bool, partyCount int) {
	c.Tmpl.ReceiveParticle(p, isActive, partyCount)
	if c.Base.Cons >= 2 {
		//only pop this if jean is active
		if !isActive {
			return
		}
		for _, active := range c.Core.Chars {
			val := make([]float64, core.EndStatType)
			val[core.AtkSpd] = 0.15
			active.AddMod(core.CharStatMod{
				Key:    "jean-c2",
				Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
				Expiry: c.Core.F + 900,
			})
			c.Core.Log.Debugw("c2 - adding atk spd", "frame", c.Core.F, "event", core.LogCharacterEvent, "character", c.Name())
		}
	}
}
