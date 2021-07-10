package ningguang

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("ningguang", NewChar)
}

type char struct {
	*character.Tmpl
	c2reset     int
	lastScreen  int
	particleICD int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 40
	c.MaxEnergy = 40
	c.Weapon.Class = def.WeaponClassCatalyst
	c.NormalHitNum = 1
	c.BurstCon = 3
	c.SkillCon = 5

	c.a4()

	return &c, nil
}

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
		f := 10 //TODO frames
		return int(float64(f) / (1 + c.Stats[def.AtkSpd]))
	case def.ActionCharge:
		return 50 //TODO frames
	case def.ActionSkill:
		return 60 //counted
	case def.ActionBurst:
		return 97 //counted, this is when you can swap but prob not when you can attack again
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) ActionStam(a def.ActionType, p map[string]int) float64 {
	switch a {
	case def.ActionDash:
		return 18
	case def.ActionCharge:
		if c.Tags["jade"] > 0 {
			return 0
		}
		return 50
	default:
		c.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(def.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		"Normal",
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		attack[c.TalentLvlAttack()],
	)
	if c.Base.Cons > 0 {
		d.Targets = def.TargetAll
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
		c.Sim.ApplyDamage(&d)
		c.Sim.ApplyDamage(&x)
	}, "ningguang-attack", f+travel)

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {
	f := c.ActionFrames(def.ActionCharge, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		"Charge",
		def.AttackTagExtra,
		def.ICDTagExtraAttack,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		charge[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f+travel)

	d = c.Snapshot(
		"Charge (Gems)",
		def.AttackTagExtra,
		def.ICDTagExtraAttack,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
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
	f := c.ActionFrames(def.ActionSkill, p)

	d := c.Snapshot(
		"Jade Screen",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, f)

	//put skill on cd first then check for construct/c2
	c.SetCD(def.ActionSkill, 720)

	//create a construct
	c.Sim.NewConstruct(c.newScreen(1800), true) //30 seconds

	c.lastScreen = c.Sim.Frame()

	//check if particles on icd

	if c.Sim.Frame() > c.particleICD {
		//3 balls, 33% chance of a fourth
		count := 3
		if c.Sim.Rand().Float64() < .33 {
			count = 4
		}
		c.QueueParticle("ningguang", count, def.Geo, f+100)
		c.particleICD = c.Sim.Frame() + 360
	}

	return f
}

func (c *char) a4() {
	//activate a4 if screen is down and character uses dash
	c.Sim.AddEventHook(func(s def.Sim) bool {
		if c.Sim.ConstructCountType(def.GeoConstructNingSkill) > 0 {
			val := make([]float64, def.EndStatType)
			val[def.GeoP] = 0.12
			c.AddMod(def.CharStatMod{
				Key: "ning-screen",
				Amount: func(a def.AttackTag) ([]float64, bool) {
					return val, true
				},
				Expiry: c.Sim.Frame() + 600,
			})
		}
		return false
	}, "ningguang-a4", def.PostDashHook)
}

func (c *char) newScreen(dur int) def.Construct {
	return &construct{
		src:    c.Sim.Frame(),
		expiry: c.Sim.Frame() + dur,
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

func (c *construct) Type() def.GeoConstructType {
	return def.GeoConstructNingSkill
}

func (c *construct) OnDestruct() {
	if c.char.Base.Cons >= 2 {
		//make sure last reset is more than 6 seconds ago
		if c.char.c2reset <= c.char.Sim.Frame()-360 {
			//reset cd
			c.char.ResetActionCooldown(def.ActionSkill)
			c.char.c2reset = c.char.Sim.Frame()
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
	f := c.ActionFrames(def.ActionBurst, p)

	//fires 6 normally, + 6 if jade screen is active
	count := 6
	if c.Sim.Destroy(c.lastScreen) {
		c.Log.Debugw("12 jade on burst", "event", def.LogCharacterEvent, "frame", c.Sim.Frame(), "char", c.Index)
		count += 6
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		"Starshatter",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
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
	c.SetCD(def.ActionBurst, 720)
	return f
}
