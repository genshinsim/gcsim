package sucrose

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("sucrose", NewChar)
}

type char struct {
	*character.Tmpl
	a4EM []float64
	// a4snap   core.Snapshot
	qInfused core.EleType
	//charges
	eCharge      int
	eChargeMax   int
	eNextRecover int
	eTickSrc     int
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

	c.eChargeMax = 1
	if c.Base.Cons >= 1 {
		c.eChargeMax = 2
	}
	c.eCharge = c.eChargeMax

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	c.a2()
	c.a4()

	if c.Base.Cons == 6 {
		c.c6()
	}

}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 19 //frames from keqing lib
		case 1:
			f = 38 - 19
		case 2:
			f = 70 - 38
		case 3:
			f = 101 - 70
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f
	case core.ActionCharge:
		return 53 //frames from keqing lib
	case core.ActionSkill:
		return 55 //ok
	case core.ActionBurst:
		return 46 //ok
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
		return 0
	}
}

func (c *char) a2() {
	val := make([]float64, core.EndStatType)
	val[core.EM] = 50
	for _, char := range c.Core.Chars {
		if char.Ele() == core.Anemo || char.Ele() == core.Geo {
			continue //nothing for geo/anemo char
		}
		char.AddMod(core.CharStatMod{
			Key:    "sucrose-a4",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				var f int
				var ok bool

				switch char.Ele() {
				case core.Pyro:
					f, ok = c.Tags["a2-pyro"]
				case core.Cryo:
					f, ok = c.Tags["a2-cryo"]
				case core.Hydro:
					f, ok = c.Tags["a2-hydro"]
				case core.Electro:
					f, ok = c.Tags["a2-electro"]
				default:
					return nil, false
				}
				return val, f > c.Core.F && ok
			},
		})
	}

	c.Core.Events.Subscribe(core.OnReactionOccured, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		switch ds.ReactionType {
		case core.SwirlCryo:
			c.Tags["a2-cryo"] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "reaction", ds.ReactionType, "expiry", c.Core.F+480)
		case core.SwirlElectro:
			c.Tags["a2-electro"] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "reaction", ds.ReactionType, "expiry", c.Core.F+480)
		case core.SwirlHydro:
			c.Tags["a2-hydro"] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "reaction", ds.ReactionType, "expiry", c.Core.F+480)
		case core.SwirlPyro:
			c.Tags["a2-pyro"] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "reaction", ds.ReactionType, "expiry", c.Core.F+480)
		}
		return false
	}, "sucrose-a2-trigger")
}

func (c *char) a4() {
	c.a4EM = make([]float64, core.EndStatType)

	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue //nothing for sucrose
		}
		char.AddMod(core.CharStatMod{
			Key:    "sucrose-a4",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if c.Core.Status.Duration("sucrosea4") == 0 {
					return nil, false
				}
				return c.a4EM, true
			},
		})
	}
}

func (c *char) c6() {
	c.AddMod(core.CharStatMod{
		Key: "sucrose-c6",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if c.Core.Status.Duration("sucrosec6") == 0 {
				return nil, false
			}
			p := core.EleToDmgP(c.qInfused)
			val := make([]float64, core.EndStatType)
			val[p] = 0.2
			return val, true
		},
		Expiry: -1,
	})
}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(core.ActionAttack, p)

	d := c.Snapshot(
		"Normal",
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-5)

	c.AdvanceNormalIndex()

	if c.Base.Cons >= 4 {
		count := c.Tags["c4"]
		count++
		if count == 7 {
			if c.Cooldown(core.ActionSkill) > 0 {
				n := c.Core.Rand.Intn(7) + 1
				c.ReduceActionCooldown(core.ActionSkill, n*60)
			}
			count = 0
		}
		c.Tags["c4"] = count
	}

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {
	f := c.ActionFrames(core.ActionCharge, p)
	d := c.Snapshot(
		"Charge Attack",
		core.AttackTagExtra,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		25,
		charge[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	if c.Base.Cons >= 4 {
		count := c.Tags["c4"]
		count++
		if count == 7 {
			if c.Cooldown(core.ActionSkill) > 0 {
				n := c.Core.Rand.Intn(7) + 1
				c.ReduceActionCooldown(core.ActionSkill, n*60)
			}
			count = 0
		}
		c.Tags["c4"] = count
	}

	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)
	//41 frame delay
	d := c.Snapshot(
		"Astable Anemohypostasis Creation-6308",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		25,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	c.AddTask(func() {
		c.Core.Status.AddStatus("sucrosea4", 480)
		c.a4EM[core.EM] = 0.2 * c.Stat(core.EM)
		c.Core.Log.Debugw("sucrose a4 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "em snapshot", c.a4EM, "expiry", c.Core.F+480)
		c.Core.Combat.ApplyDamage(&d)
	}, "Sucrose - Skill", 41)

	c.QueueParticle("sucrose", 4, core.Anemo, 150)

	if c.Base.Cons < 1 {
		c.SetCD(core.ActionSkill, 900)
		return f
	}

	switch c.eCharge {
	case c.eChargeMax:
		c.Core.Log.Debugw("sucrose e at max charge, queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "recover at", c.Core.F+900)
		c.eNextRecover = c.Core.F + 901
		c.AddTask(c.recoverCharge(c.Core.F), "charge", 900)
		c.eTickSrc = c.Core.F
	case 1:
		c.SetCD(core.ActionSkill, c.eNextRecover)
	}
	c.eCharge--

	return f
}

func (c *char) recoverCharge(src int) func() {
	return func() {
		if c.eTickSrc != src {
			c.Core.Log.Debugw("sucrose e recovery function ignored, src diff", "frame", c.Core.F, "char", c.Index, "event", core.LogCharacterEvent, "src", src, "new src", c.eTickSrc)
			return
		}
		c.eCharge++
		c.Core.Log.Debugw("sucrose e recovering a charge", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "src", src, "total charge", c.eCharge)
		c.SetCD(core.ActionSkill, 0)
		if c.eCharge >= c.eChargeMax {
			//fully charged
			return
		}
		//other wise restore another charge
		c.Core.Log.Debugw("sucrose e queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "src", src, "recover at", c.Core.F+720)
		c.eNextRecover = c.Core.F + 901
		c.AddTask(c.recoverCharge(src), "charge", 900)

	}
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)
	//tag a4
	//3 hits, 135, 249, 368; all 3 applied swirl; c2 i guess adds 2 second so one more hit
	//let's just assume 120, 240, 360, 480

	duration := 360
	if c.Base.Cons >= 2 {
		duration = 480
	}

	c.qInfused = core.NoElement

	// c.S.Status["sucroseburst"] = c.Core.F + count
	c.Core.Status.AddStatus("sucroseburst", duration)
	d := c.Snapshot(
		"Forbidden Creation-Isomer 75/Type II",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		25,
		burstDot[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	for i := 120; i <= duration; i += 120 {
		x := d.Clone()

		c.AddTask(func() {
			c.Core.Status.AddStatus("sucrosea4", 480)
			c.a4EM[core.EM] = 0.2 * c.Stat(core.EM)
			c.Core.Log.Debugw("sucrose a4 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "em snapshot", c.a4EM, "expiry", c.Core.F+480)

			c.Core.Combat.ApplyDamage(&x)

			if c.qInfused != core.NoElement {
				d := c.Snapshot(
					"Forbidden Creation-Isomer 75/Type II (Absorb)",
					core.AttackTagElementalBurst,
					core.ICDTagNone,
					core.ICDGroupDefault,
					core.StrikeTypeDefault,
					c.qInfused,
					25,
					burstAbsorb[c.TalentLvlBurst()],
				)
				d.Targets = core.TargetAll
				c.Core.Combat.ApplyDamage(&d)
			}
			//check if infused
		}, "sucrose-burst-em", i)
	}

	//
	c.AddTask(c.absorbCheck(c.Core.F, 0, int(duration/18)), "absorb-check", f)

	c.SetCD(core.ActionBurst, 1200)
	c.Energy = 0
	return f
}

func (c *char) absorbCheck(src int, count int, max int) func() {
	return func() {
		//max number of scans reached
		if count == max {
			return
		}

		fire := false
		water := false
		electric := false
		ice := false

		//scan through all targets, order is fire > water > electric > ice/frozen
		for _, t := range c.Core.Targets {
			switch t.AuraType() {
			case core.Pyro:
				fire = true
			case core.Hydro:
				water = true
			case core.Electro:
				electric = true
			case core.Cryo:
				ice = true
			case core.EC:
				water = true
			case core.Frozen:
				ice = true
			}
		}

		switch {
		case fire:
			c.qInfused = core.Pyro
		case water:
			c.qInfused = core.Hydro
		case electric:
			c.qInfused = core.Electro
		case ice:
			c.qInfused = core.Cryo
		default:
			//nothing found, queue next
			c.AddTask(c.absorbCheck(src, count+1, max), "absorb-detect", 18) //every 0.3 seconds
		}
	}
}

func (c *char) ResetActionCooldown(a core.ActionType) {
	//we're overriding this b/c of the c1 charges
	switch a {
	case core.ActionBurst:
		c.ActionCD[a] = 0
	case core.ActionSkill:
		if c.Base.Cons == 0 {
			c.ActionCD[a] = 0
			return
		}
		//if full charge do nothing; should never happen though since it takes 1 charge to proc it
		if c.eCharge == c.eChargeMax {
			c.ActionCD[a] = 0
			return
		}

		//otherwise reset tick src and add refresh if charges < max after ++
		c.eCharge++
		c.eTickSrc = c.Core.F
		c.eNextRecover = c.Core.F + 901
		c.ActionCD[a] = 0
		if c.eCharge < c.eChargeMax {
			c.AddTask(c.recoverCharge(c.Core.F), "charge", 900)
		}
	}
}
