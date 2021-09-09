package xiangling

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("xiangling", NewChar)
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
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	if c.Base.Cons >= 6 {
		c.c6()
	}
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		case 0:
			f = 12
		case 1:
			f = 38 - 12
		case 2:
			f = 72 - 38
		case 3:
			f = 141 - 72
		case 4:
			f = 167 - 141
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionSkill:
		return 26
	case core.ActionBurst:
		return 99
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) c6() {
	m := make([]float64, core.EndStatType)
	m[core.PyroP] = 0.15

	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Key:    "xl-c6",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, c.Core.Status.Duration("xlc6") > 0
			},
		})
	}
}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSpear,
		core.Physical,
		25,
		0,
	)

	for i, mult := range attack[c.NormalCounter] {
		x := d.Clone()
		x.Mult = mult[c.TalentLvlAttack()]
		c.QueueDmg(&x, f-i)
	}

	//if n = 5, add explosion for c2
	if c.Base.Cons >= 2 && c.NormalCounter == 4 {
		d1 := d.Clone()
		d1.Element = core.Pyro
		d1.Mult = 0.75
		c.QueueDmg(&d1, 120)
	}
	//add a 75 frame attackcounter reset
	c.AdvanceNormalIndex()
	//return animation cd
	//this also depends on which hit in the chain this is
	return f
}

func (c *char) ChargeAttack(p map[string]int) int {
	f := c.ActionFrames(core.ActionCharge, p)
	d := c.Snapshot(
		"Charge",
		core.AttackTagExtra,
		core.ICDTagExtraAttack,
		core.ICDGroupPole,
		core.StrikeTypeSpear,
		core.Physical,
		25,
		nc[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	//return animation cd
	return f
}

func (c *char) Skill(p map[string]int) int {
	//check if on cd first

	f := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Guoba",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeSpear,
		core.Pyro,
		25,
		guoba[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	if c.Base.Cons >= 1 {
		d.OnHitCallback = func(t core.Target) {
			t.AddResMod("xiangling-c1", core.ResistMod{
				Ele:      core.Pyro,
				Value:    -0.15,
				Duration: 6 * 60,
			})
		}
	}

	delay := 120
	c.Core.Status.AddStatus("xianglingguoba", 500)

	//lasts 73 seconds, shoots every 1.6 seconds
	for i := 0; i < 4; i++ {
		x := d.Clone()
		c.QueueDmg(&x, delay+i*90)
		c.QueueParticle("xiangling", 1, core.Pyro, delay+i*95+90+60)
	}

	c.SetCD(core.ActionSkill, 12*60)
	//return animation cd
	return f
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)
	lvl := c.TalentLvlBurst()

	delay := []int{34, 50, 75}
	for i := 0; i < len(pyronadoInitial); i++ {
		j := i
		c.AddTask(func() {
			x := c.Snapshot(
				fmt.Sprintf("Pyronado Hit %v", j+1),
				core.AttackTagElementalBurst,
				core.ICDTagElementalBurst,
				core.ICDGroupDefault,
				core.StrikeTypeSpear,
				core.Pyro,
				25,
				pyronadoInitial[j][lvl],
			)
			c.Core.Combat.ApplyDamage(&x)
		}, "pyronado initial", delay[i])

		// c.QueueDmg(&x, delay[i])
	}

	//ok for now we assume it's 80 (or 70??) frames per cycle, that gives us roughly 10s uptime
	//max is either 10s or 14s
	max := 10 * 60
	if c.Base.Cons >= 4 {
		max = 14 * 60
	}

	var d core.Snapshot

	c.AddTask(func() {
		//spin to win; snapshot between 2nd and 3rd hit
		d = c.Snapshot(
			"Pyronado",
			core.AttackTagElementalBurst,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeSpear,
			core.Pyro,
			25,
			pyronadoSpin[lvl],
		)
		d.Targets = core.TargetAll
	}, "pyronado-snap", 69)

	c.Core.Status.AddStatus("xianglingburst", max)

	for delay := 70; delay <= max; delay += 70 {
		c.QueueDmg(&d, delay)
	}

	//add an effect starting at frame 70 to end of duration to increase pyro dmg by 15% if c6
	if c.Base.Cons >= 6 {
		//wait 70 frames, add effect
		c.AddTask(func() {
			c.Core.Status.AddStatus("xlc6", max)
		}, "xl activate c6", 70)

	}

	//add cooldown to sim
	c.SetCD(core.ActionBurst, 20*60)
	//use up energy
	//c.Energy = 0  forcing every character to comsume energy after burts in the energy.go to make my life easier
	c.ConsumeEnergy(0, 0) //at 0,0 value acts the same as c.Energy = 0

	//return animation cd
	return f
}
