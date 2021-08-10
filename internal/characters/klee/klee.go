package klee

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("klee", NewChar)
}

type char struct {
	*character.Tmpl
	c1Chance     float64
	eCharge      int
	eChargeMax   int
	eNextRecover int
	eTickSrc     int
}

func NewChar(s core.Sim, log *zap.SugaredLogger, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 3
	c.eChargeMax = 2
	c.eCharge = 2

	c.a4()

	if c.Base.Cons >= 4 {
		c.c4()
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
			f = 11
		case 1:
			f = 33
		case 2:
			f = 60
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 84
	case core.ActionSkill:
		return 67
	case core.ActionBurst:
		return 101
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Tags["spark"] > 0 {
			return 0
		}
		return 50
	default:
		c.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) a4() {
	c.Sim.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.AttackTag != core.AttackTagExtra {
			return
		}
		if !crit {
			return
		}
		for _, x := range c.Sim.Characters() {
			x.AddEnergy(2)
		}

	}, "kleea2")
}

func (c *char) c1(delay int) {
	if c.Base.Cons < 1 {
		return
	}
	//0.1 base change, + 0.08 every failure
	if c.Sim.Rand().Float64() > c.c1Chance {
		//failed
		c.c1Chance += 0.08
		return
	}
	c.c1Chance = 0.1
	d := c.Snapshot(
		"Sparks'n'Splash c1",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Pyro,
		25,
		1.2*burst[c.TalentLvlBurst()],
	)
	//trigger dmg
	c.QueueDmg(&d, delay)
}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagKleeFireDamage,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Pyro,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.AddTask(func() {
		c.Sim.ApplyDamage(&d)
		c.addSpark()
	}, "klee normal", f+travel)

	c.c1(f + travel)

	c.AdvanceNormalIndex()

	if _, ok := p["walk"]; ok {
		//reset normal counter here
		c.ResetNormalCounter()
	}

	return f
}

func (c *char) addSpark() {
	if c.Sim.Rand().Float64() < 0.5 {
		c.Tags["spark"] = 1
	}
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
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Pyro,
		25,
		charge[c.TalentLvlAttack()],
	)

	//stam is calculated before this func is called so it's safe to
	//set spark to 0 here

	if c.Tags["spark"] == 1 {
		c.Tags["spark"] = 0
		d.Stats[core.DmgP] += 0.5
	}

	c.QueueDmg(&d, f+travel)

	c.c1(f + travel)

	return f
}

//p is the number of mines that hit, up to ??
func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)

	bounce, ok := p["bounce"]
	if !ok {
		bounce = 1
	}

	//mine lives for 5 seconds
	//3 bounces, roughly 30, 70, 110 hits
	d := c.Snapshot(
		"Jumpy Dumpty",
		core.AttackTagElementalArt,
		core.ICDTagKleeFireDamage,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Pyro,
		25,
		jumpy[c.TalentLvlSkill()],
	)

	for i := 0; i < bounce; i++ {
		x := d.Clone()
		c.AddTask(func() {
			c.Sim.ApplyDamage(&x)
			c.addSpark()
		}, "klee bomb", f+30+i*40)
	}

	if bounce > 0 {
		c.QueueParticle("klee", 4, core.Pyro, 130)
	}

	minehits, ok := p["mine"]
	if !ok {
		minehits = 2
	}

	//8 mines.. no idea how many normally hits
	d = c.Snapshot(
		"Jumpy Dumpty",
		core.AttackTagElementalArt,
		core.ICDTagKleeFireDamage,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Pyro,
		25,
		mine[c.TalentLvlSkill()],
	)
	if c.Base.Cons >= 2 {
		d.OnHitCallback = func(t core.Target) {
			t.AddDefMod("kleec2", -.233, 600)
		}
	}

	//roughly 160 frames after mines are laid
	for i := 0; i < minehits; i++ {
		x := d.Clone()
		c.AddTask(func() {
			c.Sim.ApplyDamage(&x)
			c.addSpark()
		}, "klee mine", f+160)

	}

	c.c1(f + 30)

	switch c.eCharge {
	case c.eChargeMax:
		c.Log.Debugw("klee at max charge, queuing next recovery", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "recover at", c.Sim.Frame()+721)
		c.eNextRecover = c.Sim.Frame() + 1201
		c.AddTask(c.recoverCharge(c.Sim.Frame()), "charge", 1200)
		c.eTickSrc = c.Sim.Frame()
	case 1:
		c.SetCD(core.ActionSkill, c.eNextRecover)
	}

	c.eCharge--

	// c.SetCD(def.ActionSkill, 20*60)
	return f
}

func (c *char) recoverCharge(src int) func() {
	return func() {
		if c.eTickSrc != src {
			c.Log.Debugw("klee mine recovery function ignored, src diff", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "src", src, "new src", c.eTickSrc)
			return
		}
		c.eCharge++
		c.Log.Debugw("klee mine recovering a charge", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "src", src, "total charge", c.eCharge)
		c.SetCD(core.ActionSkill, 0)
		if c.eCharge >= c.eChargeMax {
			//fully charged
			return
		}
		//other wise restore another charge
		c.Log.Debugw("klee mine queuing next recovery", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "src", src, "recover at", c.Sim.Frame()+720)
		c.eNextRecover = c.Sim.Frame() + 1201
		c.AddTask(c.recoverCharge(src), "charge", 1200)

	}
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)

	d := c.Snapshot(
		"Sparks'n'Splash",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Pyro,
		25,
		burst[c.TalentLvlBurst()],
	)

	//lasts 10 seconds, starts after 2.2 seconds maybe?

	//every 1.8 second +on added shoots between 3 to 5, ignore the queue thing.. space it out .2 between each wave i guess

	for i := 132; i < 732; i += 108 {
		//wave 1 = 1
		x := d.Clone()
		c.AddTask(func() {
			//no more if klee is not on field
			if c.Sim.ActiveCharIndex() != c.Index {
				return
			}
			c.Sim.ApplyDamage(&x)
		}, "klee-burst", i)
		//wave 2 = 1 + 30% chance of 1
		x = d.Clone()
		c.AddTask(func() {
			//no more if klee is not on field
			if c.Sim.ActiveCharIndex() != c.Index {
				return
			}
			c.Sim.ApplyDamage(&x)
		}, "klee-burst", i+12)
		if c.Sim.Rand().Float64() < 0.3 {
			x = d.Clone()
			c.AddTask(func() {
				//no more if klee is not on field
				if c.Sim.ActiveCharIndex() != c.Index {
					return
				}
				c.Sim.ApplyDamage(&x)
			}, "klee-burst", i+12)
		}
		//wave 3 = 1 + 50% chance of 1
		x = d.Clone()
		c.AddTask(func() {
			//no more if klee is not on field
			if c.Sim.ActiveCharIndex() != c.Index {
				return
			}
			c.Sim.ApplyDamage(&x)
		}, "klee-burst", i+24)
		if c.Sim.Rand().Float64() < 0.5 {
			x = d.Clone()
			c.AddTask(func() {
				//no more if klee is not on field
				if c.Sim.ActiveCharIndex() != c.Index {
					return
				}
				c.Sim.ApplyDamage(&x)
			}, "klee-burst", i+24)
		}
	}

	c.AddTask(func() {
		c.Sim.AddStatus("kleeq", 600)
	}, "klee-burst-status", 132)

	//every 3 seconds add energy if c6
	if c.Base.Cons == 6 {
		for i := f + 180; i < f+600; i += 180 {
			c.AddTask(func() {
				//no more if klee is not on field
				if c.Sim.ActiveCharIndex() != c.Index {
					return
				}

				for i, x := range c.Sim.Characters() {
					if i == c.Index {
						continue
					}
					x.AddEnergy(3)
					c.Log.Debugw("klee c6 regen 3 energy", "frame", c.Sim.Frame(), "event", core.LogEnergyEvent, "char", x.CharIndex(), "new energy", x.CurrentEnergy())
				}

			}, "klee-c6", i)
		}

		//add 25% buff
		for _, x := range c.Sim.Characters() {
			val := make([]float64, core.EndStatType)
			val[core.PyroP] = .1
			x.AddMod(core.CharStatMod{
				Key:    "klee-c6",
				Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
				Expiry: c.Sim.Frame() + 1500,
			})
		}
	}

	c.c1(132)

	c.SetCD(core.ActionBurst, 15*60)
	c.Energy = 0
	return f
}

func (c *char) c4() {
	c.Sim.AddEventHook(func(s core.Sim) bool {
		//if burst is active and klee no longer active char
		if c.Sim.ActiveCharIndex() != c.Index && s.Status("kleeq") > 0 {
			s.DeleteStatus("kleeq")
			//blow up
			d := c.Snapshot(
				"Sparks'n'Splash c4",
				core.AttackTagNone,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Pyro,
				50,
				5.55,
			)
			s.ApplyDamage(&d)
		}
		return false
	}, "klee-c4", core.PostSwapHook)
}
