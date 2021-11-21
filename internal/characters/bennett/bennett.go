package bennett

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("bennett", NewChar)
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
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5

	if c.Base.Cons >= 2 {
		c.c2()
	}

	//add effect for burst

	return &c, nil
}

func (c *char) c2() {
	val := make([]float64, core.EndStatType)
	val[core.ER] = .3

	c.AddMod(core.CharStatMod{
		Key: "bennett-c2",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, c.HPCurrent/c.HPMax < 0.7
		},
		Expiry: -1,
	})
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 12 //frames from keqing lib
		case 1:
			f = 20 + 5 //+5 i guess recovery from n1?
		case 2:
			f = 31
		case 3:
			f = 55
		case 4:
			f = 49
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		return 100, 100 //frames from keqing lib
	case core.ActionSkill:
		hold := p["hold"]
		switch hold {
		case 1:
			return 112, 112
		case 2:
			return 197, 197
		default:
			return 52, 52
		}
	case core.ActionBurst:
		return 51, 51 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0, 0
	}
}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	c.QueueDmgDynamic(func() *core.Snapshot {
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
		d.Targets = core.TargetAll

		return &d
	}, f-1)
	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	var cd int

	switch p["hold"] {
	case 1:
		c.skillHoldShort()
		cd = 450 - 90
	case 2:
		c.skillHoldLong()
		cd = 600 - 120
	default:
		c.skillPress()
		cd = 300 - 60
	}

	//A4
	if c.Core.Status.Duration("btburst") > 0 {
		cd = cd / 2
	}

	c.SetCD(core.ActionSkill, cd)

	return f, a

}

func (c *char) skillPress() {

	c.QueueDmgDynamic(func() *core.Snapshot {
		d := c.Snapshot(
			"Passion Overload (Press)",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			core.Pyro,
			50,
			skill[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll

		return &d
	}, 15)

	//25 % chance of 3 orbs
	count := 2
	if c.Core.Rand.Float64() < .25 {
		count++
	}
	c.QueueParticle("bennett", count, core.Pyro, 120)
}

func (c *char) skillHoldShort() {

	delay := []int{89, 115}

	for i, v := range skill1 {
		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				"Passion Overload (Hold)",
				core.AttackTagElementalArt,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeSlash,
				core.Pyro,
				25,
				v[c.TalentLvlSkill()],
			)
			d.Targets = core.TargetAll

			return &d
		}, delay[i])
	}

	//25 % chance of 3 orbs
	count := 2
	if c.Core.Rand.Float64() < .25 {
		count++
	}
	c.QueueParticle("bennett", count, core.Pyro, 215)
}

func (c *char) skillHoldLong() {
	//i think explode is guaranteed 3 orbs

	delay := []int{136, 154}

	for i, v := range skill2 {
		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				"Passion Overload (Hold)",
				core.AttackTagElementalArt,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeSlash,
				core.Pyro,
				25,
				v[c.TalentLvlSkill()],
			)
			d.Targets = core.TargetAll

			return &d
		}, delay[i])
	}

	c.QueueDmgDynamic(func() *core.Snapshot {
		d2 := c.Snapshot(
			"Passion Overload (Explode)",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Pyro,
			25,
			explosion[c.TalentLvlSkill()],
		)
		d2.Targets = core.TargetAll

		return &d2
	}, 198)

	//25 % chance of 3 orbs
	count := 2
	if c.Core.Rand.Float64() < .25 {
		count++
	}
	c.QueueParticle("bennett", count, core.Pyro, 298)

}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	//add field effect timer
	c.Core.Status.AddStatus("btburst", 720)
	//hook for buffs; active right away after cast

	c.AddTask(func() {
		d := c.Snapshot(
			"Fantastic Voyage",
			core.AttackTagElementalBurst,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Pyro,
			50,
			burst[c.TalentLvlBurst()],
		)
		d.Targets = core.TargetAll
		c.Core.Combat.ApplyDamage(&d)
	}, "bt-q", 33)

	d := c.Snapshot(
		"Fantastic Voyage (Heal)",
		core.AttackTagNone,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.NoElement,
		0,
		0,
	)

	//apply right away
	c.applyBennettField(d)()

	//add 12 ticks starting at t = 1 to t= 12
	//TODO confirm if starts at t=1 or after animation
	for i := 0; i <= 720; i += 60 {
		c.AddTask(c.applyBennettField(d), "bennett-field", i)
	}

	c.Energy = 0
	c.SetCD(core.ActionBurst, 900)
	return f, a //todo fix field cast time
}

func (c *char) applyBennettField(d core.Snapshot) func() {
	hpplus := d.Stats[core.Heal]
	heal := (bursthp[c.TalentLvlBurst()] + bursthpp[c.TalentLvlBurst()]*c.MaxHP()) * (1 + hpplus)
	pc := burstatk[c.TalentLvlBurst()]
	if c.Base.Cons >= 1 {
		pc += 0.2
	}
	atk := pc * float64(c.Base.Atk+c.Weapon.Atk)
	return func() {
		c.Core.Log.Debugw("bennett field ticking", "frame", c.Core.F, "event", core.LogCharacterEvent)

		active := c.Core.Chars[c.Core.ActiveChar]
		//heal if under 70%
		if active.HP()/active.MaxHP() < .7 {
			c.Core.Health.HealActive(c.Index, heal)
		}

		//add attack if over 70%
		threshold := .7
		if c.Base.Cons >= 1 {
			threshold = 1.1
		}
		if active.HP()/active.MaxHP() < threshold {
			//add 2.1s = 126 frames
			val := make([]float64, core.EndStatType)
			val[core.ATK] = atk
			active.AddMod(core.CharStatMod{
				Key: "bennett-field",
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return val, true
				},
				Expiry: c.Core.F + 126,
			})
			// Additional per-character status for config conditionals
			c.Core.Status.AddStatus(fmt.Sprintf("bennettbuff%v", active.Name()), 126)
			c.Core.Log.Debugw("bennett field - adding attack", "frame", c.Core.F, "event", core.LogCharacterEvent, "threshold", threshold)
			//if c6 add weapon infusion and 15% pyro
			if c.Base.Cons == 6 {
				switch active.WeaponClass() {
				case core.WeaponClassClaymore:
					fallthrough
				case core.WeaponClassSpear:
					fallthrough
				case core.WeaponClassSword:
					active.AddWeaponInfuse(core.WeaponInfusion{
						Key:    "bennett-fire-weapon",
						Ele:    core.Pyro,
						Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
						Expiry: c.Core.F + 126,
					})
				}

			}
		}
	}
}
