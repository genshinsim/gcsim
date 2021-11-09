package diona

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
)

func init() {
	core.RegisterCharFunc("diona", NewChar)
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
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.a2()

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 16 //frames from keqing lib
		case 1:
			f = 37 - 16
		case 2:
			f = 67 - 37
		case 3:
			f = 101 - 67
		case 4:
			f = 152 - 101
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionAim:
		if c.Base.Cons >= 4 && c.Core.Status.Duration("dionaburst") > 0 {
			return 34, 34 //reduced by 60%
		}
		return 84, 84 //kqm
	case core.ActionBurst:
		return 21, 21
	case core.ActionSkill:
		switch p["hold"] {
		case 1:
			return 24, 24
		default:
			return 15, 15
		}
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0, 0
	}
}

func (c *char) a2() {
	c.Core.AddStamMod(func(a core.ActionType) (float64, bool) {
		if c.Core.Shields.Get(core.ShieldDionaSkill) != nil {
			return -0.1, false
		}
		return 0, false
	})
}

func (c *char) c6() {
	c.Core.Health.AddIncHealBonus(func() float64 {
		if c.Core.Status.Duration("dionaburst") == 0 {
			return 0
		}
		char := c.Core.Chars[c.Core.ActiveChar]
		if char.HP()/char.MaxHP() <= 0.5 {
			c.Core.Log.Debugw("diona c6 activated", "frame", c.Core.F, "event", core.LogCharacterEvent)
			return 0.3
		}
		return 0
	})
}

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

	return f, a
}

func (c *char) Aimed(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAim, p)

	d := c.Snapshot(
		"Aim (Charged)",
		core.AttackTagExtra,
		core.ICDTagExtraAttack,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Cryo,
		25,
		aim[c.TalentLvlAttack()],
	)

	d.HitWeakPoint = true
	d.AnimationFrames = f

	c.QueueDmg(&d, travel+f)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}
	f, a := c.ActionFrames(core.ActionSkill, p)

	// 2 paws
	var bonus float64 = 1
	cd := 360 + f
	pawCount := 2

	if p["hold"] == 1 {
		//5 paws, 75% absorption bonus
		bonus = 1.75
		cd = 900 + f
		pawCount = 5
	}

	shd := (pawShieldPer[c.TalentLvlSkill()]*c.MaxHP() + pawShieldFlat[c.TalentLvlSkill()]) * bonus
	if c.Base.Cons >= 2 {
		shd = shd * 1.15
	}

	d := c.Snapshot(
		"Icy Paw",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Cryo,
		25,
		paw[c.TalentLvlSkill()],
	)

	count := 0

	for i := 0; i < pawCount; i++ {
		x := d.Clone()
		if c.Base.Cons >= 2 {
			d.Stats[core.DmgP] += 0.15
		}
		c.QueueDmg(&x, travel+f-5+i)

		if c.Core.Rand.Float64() < 0.8 {
			count++
		}
	}

	//particles
	c.QueueParticle("Diona", count, core.Cryo, 90) //90s travel time

	//add shield
	c.AddTask(func() {
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldDionaSkill,
			HP:         shd,
			Ele:        core.Cryo,
			Expires:    c.Core.F + pawDur[c.TalentLvlSkill()], //15 sec
		})
	}, "Diona-Paw-Shield", f)

	c.SetCD(core.ActionSkill, cd)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	//initial hit
	d := c.Snapshot(
		"Signature Mix (Initial)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		25,
		burst[c.TalentLvlBurst()],
	)
	c.QueueDmg(&d, f-10)

	d = c.Snapshot(
		"Signature Mix (Tick)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		25,
		burstDot[c.TalentLvlBurst()],
	)
	hpplus := d.Stats[core.Heal]
	maxhp := c.MaxHP()
	heal := (burstHealPer[c.TalentLvlBurst()]*maxhp + burstHealFlat[c.TalentLvlBurst()]) * (1 + hpplus)

	//ticks every 2s, first tick at t=1s, then t=3,5,7,9,11, lasts for 12.5
	for i := 0; i < 6; i++ {
		c.AddTask(func() {
			x := d.Clone()
			c.Core.Combat.ApplyDamage(&x)
			c.Core.Log.Debugw("diona healing", "frame", c.Core.F, "event", core.LogCharacterEvent, "+heal", hpplus, "max hp", maxhp, "heal amount", heal)
			c.Core.Health.HealActive(heal)
		}, "Diona Burst (DOT)", 60+i*120)
	}

	//apparently lasts for 12.5
	c.Core.Status.AddStatus("dionaburst", f+750) //TODO not sure when field starts, is it at animation end? prob when it lands...

	//c1
	if c.Base.Cons >= 1 {
		//15 energy after ends, flat not affected by ER
		c.AddTask(func() {
			c.Energy += 15
			if c.Energy > c.EnergyMax {
				c.Energy = c.EnergyMax
			}
			c.Core.Log.Debugw("diona c1 regen 15 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "new energy", c.Energy)
		}, "Diona C1", f+750)
	}

	if c.Base.Cons == 6 {
		c.AddTask(func() {
			for _, char := range c.Core.Chars {
				this := char
				val := make([]float64, core.EndStatType)
				val[core.EM] = 200
				this.AddMod(core.CharStatMod{
					Key:    "diona-c6",
					Expiry: c.Core.F + 750,
					Amount: func(a core.AttackTag) ([]float64, bool) {
						return val, this.HP()/this.MaxHP() > 0.5
					},
				})
			}
		}, "c6-em-share", f)
	}

	c.SetCD(core.ActionBurst, 1200+f)
	c.Energy = 0
	return f, a
}
