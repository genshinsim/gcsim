package xiangling

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("xiangling", NewChar)
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
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = def.WeaponClassSpear
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	if c.Base.Cons >= 6 {
		c.c6()
	}
}

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionSkill:
		return 26
	case def.ActionBurst:
		return 140
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
		return 25
	default:
		c.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) c6() {
	m := make([]float64, def.EndStatType)
	m[def.PyroP] = 0.15

	for _, char := range c.Sim.Characters() {
		char.AddMod(def.CharStatMod{
			Key:    "xl-c6",
			Expiry: -1,
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, c.Sim.Status("xlc6") > 0
			},
		})
	}
}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(def.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSpear,
		def.Physical,
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
		d1.Element = def.Pyro
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
	f := c.ActionFrames(def.ActionCharge, p)
	d := c.Snapshot(
		"Charge",
		def.AttackTagExtra,
		def.ICDTagExtraAttack,
		def.ICDGroupPole,
		def.StrikeTypeSpear,
		def.Physical,
		25,
		nc[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	//return animation cd
	return f
}

func (c *char) Skill(p map[string]int) int {
	//check if on cd first

	f := c.ActionFrames(def.ActionSkill, p)
	d := c.Snapshot(
		"Guoba",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeSpear,
		def.Pyro,
		25,
		guoba[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	if c.Base.Cons >= 1 {
		d.OnHitCallback = func(t def.Target) {
			t.AddResMod("xiangling-c1", def.ResistMod{
				Ele:      def.Pyro,
				Value:    -0.15,
				Duration: 6 * 60,
			})
		}
	}

	delay := 120
	c.Sim.AddStatus("xianglingguoba", 500)

	//lasts 73 seconds, shoots every 1.6 seconds
	for i := 0; i < 4; i++ {
		x := d.Clone()
		c.QueueDmg(&x, delay+i*86)
		c.QueueParticle("xiangling", 1, def.Pyro, delay+i*95+90+60)
	}

	c.SetCD(def.ActionSkill, 12*60)
	//return animation cd
	return f
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(def.ActionBurst, p)
	lvl := c.TalentLvlBurst()
	d := c.Snapshot(
		"Pyronado",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeSpear,
		def.Pyro,
		25,
		0,
	)

	delay := []int{20, 50, 75}
	for i := 0; i < len(pyronadoInitial); i++ {
		x := d.Clone()
		x.Abil = fmt.Sprintf("Pyronado Hit %v", i+1)
		x.Mult = pyronadoInitial[i][lvl]
		c.QueueDmg(&x, delay[i])
	}

	//spin to win; snapshot on cast
	d = c.Snapshot(
		"Pyronado",
		def.AttackTagElementalBurst,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeSpear,
		def.Pyro,
		25,
		pyronadoSpin[lvl],
	)
	d.Targets = def.TargetAll

	//ok for now we assume it's 80 (or 70??) frames per cycle, that gives us roughly 10s uptime
	//max is either 10s or 14s
	max := 10 * 60
	if c.Base.Cons >= 4 {
		max = 14 * 60
	}

	c.Sim.AddStatus("xianglingburst", max)

	for delay := 70; delay <= max; delay += 70 {
		c.QueueDmg(&d, delay)
	}

	//add an effect starting at frame 70 to end of duration to increase pyro dmg by 15% if c6
	if c.Base.Cons >= 6 {
		//wait 70 frames, add effect
		c.AddTask(func() {
			c.Sim.AddStatus("xlc6", max)
		}, "xl activate c6", 70)

	}

	//add cooldown to sim
	c.SetCD(def.ActionBurst, 20*60)
	//use up energy
	c.Energy = 0

	//return animation cd
	return f
}
