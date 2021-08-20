package lisa

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("lisa", NewChar)
}

type char struct {
	*character.Tmpl
	c6icd int
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
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 4
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
			f = 25
		case 1:
			f = 46 - 25
		case 2:
			f = 70 - 46
		case 3:
			f = 114 - 70
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 95
	case core.ActionSkill:
		hold := p["hold"]
		if hold == 0 {
			return 21 //no hold
		}
		return 116 //yes hold
	case core.ActionBurst:
		return 30
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
		return 50
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.F < c.c6icd && c.c6icd != 0 {
			return false
		}
		if c.Core.ActiveChar == c.CharIndex() {
			//swapped to lisa
			c.Tags["stack"] = 3
			c.c6icd = c.Core.F + 300
		}
		return false
	}, "lisa-c6")
}

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagLisaElectro,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)
	c.AdvanceNormalIndex()

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {
	f := c.ActionFrames(core.ActionCharge, p)

	//TODO: assumes this applies every time per
	//[7:53 PM] Hold ₼KLEE like others hold GME: CHarge is pyro every charge
	d := c.Snapshot(
		"Charge Attack",
		core.AttackTagExtra,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		charge[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	//a2 add a stack
	if c.Tags["stack"] < 3 {
		c.Tags["stack"]++
	}

	hits := p["hits"]
	if hits > 5 {
		hits = 5
	}
	if hits == 0 {
		hits = 1
	}

	//c1 adds energy
	if c.Base.Cons > 0 {
		c.AddEnergy(2 * float64(hits))
	}

	return c.ActionFrames(core.ActionCharge, p)
}

//p = 0 for no hold, p = 1 for hold
func (c *char) Skill(p map[string]int) int {
	hold := p["hold"]
	if hold == 1 {
		return c.skillHold(p)
	}
	return c.skillPress(p)
}

//TODO: how long do stacks last?
func (c *char) skillPress(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Violet Arc",
		core.AttackTagElementalArt,
		core.ICDTagLisaElectro,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		skillPress[c.TalentLvlSkill()],
	)
	//add 1 stack if less than 3
	if c.Tags["stack"] < 3 {
		c.Tags["stack"]++
	}
	c.QueueDmg(&d, f-1)

	if c.Core.Rand.Float64() < 0.5 {
		c.QueueParticle("Lisa", 1, core.Electro, f+100)
	}

	c.SetCD(core.ActionSkill, 60)
	return f
}

func (c *char) skillHold(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Violet Arc (Hold)",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		skillHold[c.Tags["stack"]][c.TalentLvlSkill()],
	)
	c.Tags["stack"] = 0 //consume all stacks
	c.QueueDmg(&d, f-1)

	//c2 add defense? no interruptions either way
	if c.Base.Cons >= 2 {
		//increase def for the duration of this abil in however many frames
		val := make([]float64, core.EndStatType)
		val[core.DEFP] = 0.25
		c.AddMod(core.CharStatMod{
			Key:    "lisa-c2",
			Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
			Expiry: c.Core.F + 126,
		})
	}

	//[8:31 PM] ArchedNosi | Lisa Unleashed: yeah 4-5 50/50 with Hold
	//[9:13 PM] ArchedNosi | Lisa Unleashed: @gimmeabreak actually wait, xd i noticed i misread my sheet, Lisa Hold E always gens 5 orbs
	// count := 4
	// if c.Core.Rand.Float64() < 0.5 {
	// 	count = 5
	// }
	c.QueueParticle("Lisa", 5, core.Electro, f+100)

	// c.CD[def.SkillCD] = c.Core.F + 960 //16seconds
	c.SetCD(core.ActionSkill, 960)
	return f
}

func (c *char) Burst(p map[string]int) int {

	f := c.ActionFrames(core.ActionBurst, p)

	//first zap has no icd
	d := c.Snapshot(
		"Lightning Rose (Initial)",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		0,
		0.1,
	)
	//duration is 15 seconds, tick every .5 sec
	c.QueueDmg(&d, f)
	//30 zaps once every 30 frame, starting at f

	d = c.Snapshot(
		"Lightning Rose (Tick)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		burst[c.TalentLvlBurst()],
	)

	for i := 30; i <= 900; i += 30 {
		x := d.Clone()
		c.QueueDmg(&x, f+i)
	}

	//add a status for this just in case someone cares
	c.AddTask(func() {
		c.Core.Status.AddStatus("lisaburst", 900)
	}, "lisa burst status", f)

	//on lisa c4
	//[8:11 PM] gimmeabreak: er, what does lisa c4 do?
	//[8:11 PM] ArchedNosi | Lisa Unleashed: allows each pulse of the ult to be 2-4 arcs
	//[8:11 PM] ArchedNosi | Lisa Unleashed: if theres enemies given
	//[8:11 PM] gimmeabreak: oh so it jumps 2 to 4 times?
	//[8:11 PM] gimmeabreak: i guess single target it does nothing then?
	//[8:12 PM] ArchedNosi | Lisa Unleashed: yeah single does nothing

	c.Energy = 0
	// c.CD[def.BurstCD] = c.Core.F + 1200
	c.SetCD(core.ActionBurst, 1200)
	return f //TODO: frames
}
