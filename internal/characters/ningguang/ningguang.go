package ningguang

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("ningguang", NewChar)
}

type char struct {
	*character.Tmpl
	c2reset     int
	lastScreen  int
	particleICD int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 40
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 1
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	c.a4()

	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 10 //TODO frames
		return int(float64(f) / (1 + c.Stats[core.AtkSpd]))
	case core.ActionCharge:
		return 50 //TODO frames
	case core.ActionSkill:
		return 60 //counted
	case core.ActionBurst:
		return 97 //counted, this is when you can swap but prob not when you can attack again
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
		if c.Tags["jade"] > 0 {
			return 0
		}
		return 50
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		"Normal",
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		attack[c.TalentLvlAttack()],
	)
	if c.Base.Cons > 0 {
		d.Targets = core.TargetAll
	}

	c.AddTask(func() {
		count := c.Tags["jade"]
		if count != 7 {
			count++
			if count > 3 {
				count = 3
			}
			c.Tags["jade"] = count
		}
		//refresh cooldown of seal every hit regardless if we got more stacks
		// c.CD["seal"] = s.F + 600
		x := d.Clone()
		c.Core.Combat.ApplyDamage(&d)
		c.Core.Combat.ApplyDamage(&x)
	}, "ningguang-attack", f+travel)

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {
	f := c.ActionFrames(core.ActionCharge, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		"Charge",
		core.AttackTagExtra,
		core.ICDTagExtraAttack,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		charge[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f+travel)

	d = c.Snapshot(
		"Charge (Gems)",
		core.AttackTagExtra,
		core.ICDTagExtraAttack,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		50,
		jade[c.TalentLvlAttack()],
	)
	j := c.Tags["jade"]
	for i := 0; i < j; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f+travel)
	}
	c.Tags["jade"] = 0

	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)

	d := c.Snapshot(
		"Jade Screen",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f)

	//put skill on cd first then check for construct/c2
	c.SetCD(core.ActionSkill, 720)

	//create a construct
	c.Core.Constructs.NewConstruct(c.newScreen(1800), true) //30 seconds

	c.lastScreen = c.Core.F

	//check if particles on icd

	if c.Core.F > c.particleICD {
		//3 balls, 33% chance of a fourth
		count := 3
		if c.Core.Rand.Float64() < .33 {
			count = 4
		}
		c.QueueParticle("ningguang", count, core.Geo, f+100)
		c.particleICD = c.Core.F + 360
	}

	return f
}

func (c *char) a4() {
	//activate a4 if screen is down and character uses dash
	c.Core.Events.Subscribe(core.OnDash, func(args ...interface{}) bool {
		if c.Core.Constructs.ConstructCountType(core.GeoConstructNingSkill) > 0 {
			val := make([]float64, core.EndStatType)
			val[core.GeoP] = 0.12
			char := c.Core.Chars[c.Core.ActiveChar]
			char.AddMod(core.CharStatMod{
				Key: "ning-screen",
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return val, true
				},
				Expiry: c.Core.F + 600,
			})
		}
		return false
	}, "ningguang-a4")
}

func (c *char) newScreen(dur int) core.Construct {
	return &construct{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
	}
}

type construct struct {
	src    int
	expiry int
	char   *char
}

func (c *construct) Key() int {
	return c.src
}

func (c *construct) Type() core.GeoConstructType {
	return core.GeoConstructNingSkill
}

func (c *construct) OnDestruct() {
	if c.char.Base.Cons >= 2 {
		//make sure last reset is more than 6 seconds ago
		if c.char.c2reset <= c.char.Core.F-360 {
			//reset cd
			c.char.ResetActionCooldown(core.ActionSkill)
			c.char.c2reset = c.char.Core.F
		}
	}
}
func (c *construct) Expiry() int {
	return c.expiry
}

func (c *construct) IsLimited() bool {
	return true
}

func (c *construct) Count() int {
	return 1
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)

	//fires 6 normally, + 6 if jade screen is active
	count := 6
	if c.Core.Constructs.Destroy(c.lastScreen) {
		c.Core.Log.Debugw("12 jade on burst", "event", core.LogCharacterEvent, "frame", c.Core.F, "char", c.Index)
		count += 6
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		"Starshatter",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		50,
		burst[c.TalentLvlBurst()],
	)

	//geo applied 1 4 7 10, +3 pattern; or 0 3 6 9
	for i := 0; i < count; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f+travel)
	}

	if c.Base.Cons == 6 {
		c.Tags["jade"] = 7
	}

	c.Energy = 0
	c.SetCD(core.ActionBurst, 720)
	return f
}
