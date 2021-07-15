package diona

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"github.com/genshinsim/gsim/pkg/shield"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("diona", NewChar)
}

type char struct {
	*character.Tmpl
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.MaxEnergy = 60
	c.Weapon.Class = def.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.a2()

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionAim:
		if c.Base.Cons >= 4 && c.Sim.Status("dionaburst") > 0 {
			return 34 //reduced by 60%
		}
		return 84 //kqm
	case def.ActionBurst:
		return 21
	case def.ActionSkill:
		switch p["hold"] {
		case 1:
			return 24
		default:
			return 15
		}
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) a2() {
	c.Sim.AddStamMod(func(a def.ActionType) float64 {
		if c.Sim.GetShield(def.ShieldDionaSkill) != nil {
			return -0.1
		}
		return 0
	})
}

func (c *char) c6() {
	c.Sim.AddIncHealBonus(func() float64 {
		if c.Sim.Status("dionaburst") == 0 {
			return 0
		}
		char, _ := c.Sim.CharByPos(c.Sim.ActiveCharIndex())
		if char.HP()/char.MaxHP() <= 0.5 {
			c.Log.Debugw("diona c6 activated", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent)
			return 0.3
		}
		return 0
	})
}

func (c *char) Attack(p map[string]int) int {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f := c.ActionFrames(def.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypePierce,
		def.Physical,
		25,
		auto[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, travel+f)

	c.AdvanceNormalIndex()

	return f
}

func (c *char) Aimed(p map[string]int) int {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f := c.ActionFrames(def.ActionAim, p)

	d := c.Snapshot(
		"Aim (Charged)",
		def.AttackTagExtra,
		def.ICDTagExtraAttack,
		def.ICDGroupDefault,
		def.StrikeTypePierce,
		def.Cryo,
		25,
		aim[c.TalentLvlAttack()],
	)

	d.HitWeakPoint = true
	d.AnimationFrames = f

	c.QueueDmg(&d, travel+f)

	return f
}

func (c *char) Skill(p map[string]int) int {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}
	f := c.ActionFrames(def.ActionSkill, p)

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
		def.AttackTagElementalArt,
		def.ICDTagElementalArt,
		def.ICDGroupDefault,
		def.StrikeTypePierce,
		def.Cryo,
		25,
		paw[c.TalentLvlSkill()],
	)

	count := 0

	for i := 0; i < pawCount; i++ {
		x := d.Clone()
		if c.Base.Cons >= 2 {
			d.Stats[def.DmgP] += 0.15
		}
		c.QueueDmg(&x, travel+f-5+i)

		if c.Sim.Rand().Float64() < 0.8 {
			count++
		}
	}

	//particles
	c.QueueParticle("Diona", count, def.Cryo, 90) //90s travel time

	//add shield
	c.AddTask(func() {
		c.Sim.AddShield(&shield.Tmpl{
			Src:        c.Sim.Frame(),
			ShieldType: def.ShieldDionaSkill,
			HP:         shd,
			Ele:        def.Cryo,
			Expires:    c.Sim.Frame() + pawDur[c.TalentLvlSkill()], //15 sec
		})
	}, "Diona-Paw-Shield", f)

	c.SetCD(def.ActionSkill, cd)
	return f
}

func (c *char) Burst(p map[string]int) int {

	f := c.ActionFrames(def.ActionBurst, p)

	//initial hit
	d := c.Snapshot(
		"Signature Mix (Initial)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Cryo,
		25,
		burst[c.TalentLvlBurst()],
	)
	c.QueueDmg(&d, f-10)

	d = c.Snapshot(
		"Signature Mix (Tick)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Cryo,
		25,
		burstDot[c.TalentLvlBurst()],
	)
	hpplus := d.Stats[def.Heal]
	maxhp := c.MaxHP()
	heal := (burstHealPer[c.TalentLvlBurst()]*maxhp + burstHealFlat[c.TalentLvlBurst()]) * (1 + hpplus)

	//ticks every 2s, first tick at t=1s, then t=3,5,7,9,11, lasts for 12.5
	for i := 0; i < 6; i++ {
		c.AddTask(func() {
			x := d.Clone()
			c.Sim.ApplyDamage(&x)
			c.Log.Debugw("diona healing", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "+heal", hpplus, "max hp", maxhp, "heal amount", heal)
			c.Sim.HealActive(heal)
		}, "Diona Burst (DOT)", 60+i*120)
	}

	//apparently lasts for 12.5
	c.Sim.AddStatus("dionaburst", f+750) //TODO not sure when field starts, is it at animation end? prob when it lands...

	//c1
	if c.Base.Cons >= 1 {
		//15 energy after ends, flat not affected by ER
		c.AddTask(func() {
			c.Energy += 15
			if c.Energy > c.MaxEnergy {
				c.Energy = c.MaxEnergy
			}
			c.Log.Debugw("diona c1 regen 15 energy", "frame", c.Sim.Frame(), "event", def.LogEnergyEvent, "new energy", c.Energy)
		}, "Diona C1", f+750)
	}

	if c.Base.Cons == 6 {
		c.AddTask(func() {
			for _, char := range c.Sim.Characters() {
				val := make([]float64, def.EndStatType)
				val[def.EM] = 200
				char.AddMod(def.CharStatMod{
					Key:    "diona-c6",
					Expiry: 750,
					Amount: func(a def.AttackTag) ([]float64, bool) {
						return val, char.HP()/char.MaxHP() > 0.5
					},
				})
			}
		}, "c6-em-share", f)
	}

	c.SetCD(def.ActionBurst, 1200+f)
	c.Energy = 0
	return f
}
