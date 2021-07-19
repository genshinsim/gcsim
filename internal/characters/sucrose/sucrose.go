package sucrose

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("sucrose", NewChar)
}

type char struct {
	*character.Tmpl
	a4EM []float64
	// a4snap   def.Snapshot
	qInfused def.EleType
	//charges
	eCharge      int
	eChargeMax   int
	eNextRecover int
	eTickSrc     int
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
	c.Weapon.Class = def.WeaponClassCatalyst
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

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stat(def.AtkSpd)))
		return f
	case def.ActionCharge:
		return 53 //frames from keqing lib
	case def.ActionSkill:
		return 55 //ok
	case def.ActionBurst:
		return 46 //ok
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
		return 50
	default:
		return 0
	}
}

func (c *char) a2() {
	val := make([]float64, def.EndStatType)
	val[def.EM] = 50
	for _, char := range c.Sim.Characters() {
		if char.Ele() == def.Anemo || char.Ele() == def.Geo {
			continue //nothing for geo/anemo char
		}
		char.AddMod(def.CharStatMod{
			Key:    "sucrose-a4",
			Expiry: -1,
			Amount: func(a def.AttackTag) ([]float64, bool) {
				var f int
				var ok bool

				switch char.Ele() {
				case def.Pyro:
					f, ok = c.Tags["a2-pyro"]
				case def.Cryo:
					f, ok = c.Tags["a2-cryo"]
				case def.Hydro:
					f, ok = c.Tags["a2-hydro"]
				case def.Electro:
					f, ok = c.Tags["a2-electro"]
				default:
					return nil, false
				}
				return val, f > c.Sim.Frame() && ok
			},
		})
	}

	c.Sim.AddOnReaction(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex != c.Index {
			return
		}
		switch ds.ReactionType {
		case def.SwirlCryo:
			c.Tags["a2-cryo"] = c.Sim.Frame() + 480
			c.Log.Debugw("sucrose a2 triggered", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "reaction", ds.ReactionType, "expiry", c.Sim.Frame()+480)
		case def.SwirlElectro:
			c.Tags["a2-electro"] = c.Sim.Frame() + 480
			c.Log.Debugw("sucrose a2 triggered", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "reaction", ds.ReactionType, "expiry", c.Sim.Frame()+480)
		case def.SwirlHydro:
			c.Tags["a2-hydro"] = c.Sim.Frame() + 480
			c.Log.Debugw("sucrose a2 triggered", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "reaction", ds.ReactionType, "expiry", c.Sim.Frame()+480)
		case def.SwirlPyro:
			c.Tags["a2-pyro"] = c.Sim.Frame() + 480
			c.Log.Debugw("sucrose a2 triggered", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "reaction", ds.ReactionType, "expiry", c.Sim.Frame()+480)
		}
	}, "sucrose-a2-trigger")
}

func (c *char) a4() {
	c.a4EM = make([]float64, def.EndStatType)

	for i, char := range c.Sim.Characters() {
		if i == c.Index {
			continue //nothing for sucrose
		}
		char.AddMod(def.CharStatMod{
			Key:    "sucrose-a4",
			Expiry: -1,
			Amount: func(a def.AttackTag) ([]float64, bool) {
				if c.Sim.Status("sucrosea4") == 0 {
					return nil, false
				}
				return c.a4EM, true
			},
		})
	}
}

func (c *char) c6() {
	c.AddMod(def.CharStatMod{
		Key: "sucrose-c6",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if c.Sim.Status("sucrosec6") == 0 {
				return nil, false
			}
			p := def.EleToDmgP(c.qInfused)
			val := make([]float64, def.EndStatType)
			val[p] = 0.2
			return val, true
		},
		Expiry: -1,
	})
}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(def.ActionAttack, p)

	d := c.Snapshot(
		"Normal",
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Anemo,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-5)

	c.AdvanceNormalIndex()

	if c.Base.Cons >= 4 {
		count := c.Tags["c4"]
		count++
		if count == 7 {
			if c.Cooldown(def.ActionSkill) > 0 {
				n := c.Sim.Rand().Intn(7) + 1
				c.ReduceActionCooldown(def.ActionSkill, n*60)
			}
			count = 0
		}
		c.Tags["c4"] = count
	}

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {
	f := c.ActionFrames(def.ActionCharge, p)
	d := c.Snapshot(
		"Charge Attack",
		def.AttackTagExtra,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Anemo,
		25,
		charge[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	if c.Base.Cons >= 4 {
		count := c.Tags["c4"]
		count++
		if count == 7 {
			if c.Cooldown(def.ActionSkill) > 0 {
				n := c.Sim.Rand().Intn(7) + 1
				c.ReduceActionCooldown(def.ActionSkill, n*60)
			}
			count = 0
		}
		c.Tags["c4"] = count
	}

	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)
	//41 frame delay
	d := c.Snapshot(
		"Astable Anemohypostasis Creation-6308",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Anemo,
		25,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	c.AddTask(func() {
		c.Sim.AddStatus("sucrosea4", 480)
		c.a4EM[def.EM] = 0.2 * c.Stat(def.EM)
		c.Log.Debugw("sucrose a4 triggered", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "em snapshot", c.a4EM, "expiry", c.Sim.Frame()+480)
		c.Sim.ApplyDamage(&d)
	}, "Sucrose - Skill", 41)

	c.QueueParticle("sucrose", 4, def.Anemo, 150)

	if c.Base.Cons < 1 {
		c.SetCD(def.ActionSkill, 900)
		return f
	}

	switch c.eCharge {
	case c.eChargeMax:
		c.Log.Debugw("sucrose e at max charge, queuing next recovery", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "recover at", c.Sim.Frame()+900)
		c.eNextRecover = c.Sim.Frame() + 901
		c.AddTask(c.recoverCharge(c.Sim.Frame()), "charge", 900)
		c.eTickSrc = c.Sim.Frame()
	case 1:
		c.SetCD(def.ActionSkill, c.eNextRecover)
	}
	c.eCharge--

	return f
}

func (c *char) recoverCharge(src int) func() {
	return func() {
		if c.eTickSrc != src {
			c.Log.Debugw("sucrose e recovery function ignored, src diff", "frame", c.Sim.Frame(), "char", c.Index, "event", def.LogCharacterEvent, "src", src, "new src", c.eTickSrc)
			return
		}
		c.eCharge++
		c.Log.Debugw("sucrose e recovering a charge", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "char", c.Index, "src", src, "total charge", c.eCharge)
		c.SetCD(def.ActionSkill, 0)
		if c.eCharge >= c.eChargeMax {
			//fully charged
			return
		}
		//other wise restore another charge
		c.Log.Debugw("sucrose e queuing next recovery", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "char", c.Index, "src", src, "recover at", c.Sim.Frame()+720)
		c.eNextRecover = c.Sim.Frame() + 901
		c.AddTask(c.recoverCharge(src), "charge", 900)

	}
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(def.ActionBurst, p)
	//tag a4
	//3 hits, 135, 249, 368; all 3 applied swirl; c2 i guess adds 2 second so one more hit
	//let's just assume 120, 240, 360, 480

	duration := 360
	if c.Base.Cons >= 2 {
		duration = 480
	}

	c.qInfused = def.NoElement

	// c.S.Status["sucroseburst"] = c.Sim.Frame() + count
	c.Sim.AddStatus("sucroseburst", duration)
	d := c.Snapshot(
		"Forbidden Creation-Isomer 75/Type II",
		def.AttackTagElementalBurst,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Anemo,
		25,
		burstDot[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll

	for i := 120; i <= duration; i += 120 {
		x := d.Clone()

		c.AddTask(func() {
			c.Sim.AddStatus("sucrosea4", 480)
			c.a4EM[def.EM] = 0.2 * c.Stat(def.EM)
			c.Log.Debugw("sucrose a4 triggered", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "em snapshot", c.a4EM, "expiry", c.Sim.Frame()+480)

			c.Sim.ApplyDamage(&x)

			if c.qInfused != def.NoElement {
				d := c.Snapshot(
					"Forbidden Creation-Isomer 75/Type II (Absorb)",
					def.AttackTagElementalBurst,
					def.ICDTagNone,
					def.ICDGroupDefault,
					def.StrikeTypeDefault,
					c.qInfused,
					25,
					burstAbsorb[c.TalentLvlBurst()],
				)
				d.Targets = def.TargetAll
				c.Sim.ApplyDamage(&d)
			}
			//check if infused
		}, "sucrose-burst-em", i)
	}

	//
	c.AddTask(c.absorbCheck(c.Sim.Frame(), 0, int(duration/18)), "absorb-check", f)

	c.SetCD(def.ActionBurst, 1200)
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
		for _, t := range c.Sim.Targets() {
			switch t.AuraType() {
			case def.Pyro:
				fire = true
			case def.Hydro:
				water = true
			case def.Electro:
				electric = true
			case def.Cryo:
				ice = true
			case def.EC:
				water = true
			case def.Frozen:
				ice = true
			}
		}

		switch {
		case fire:
			c.qInfused = def.Pyro
		case water:
			c.qInfused = def.Hydro
		case electric:
			c.qInfused = def.Electro
		case ice:
			c.qInfused = def.Cryo
		default:
			//nothing found, queue next
			c.AddTask(c.absorbCheck(src, count+1, max), "absorb-detect", 18) //every 0.3 seconds
		}
	}
}

func (c *char) ResetActionCooldown(a def.ActionType) {
	//we're overriding this b/c of the c1 charges
	switch a {
	case def.ActionBurst:
		c.ActionCD[a] = 0
	case def.ActionSkill:
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
		c.eTickSrc = c.Sim.Frame()
		c.eNextRecover = c.Sim.Frame() + 901
		c.ActionCD[a] = 0
		if c.eCharge < c.eChargeMax {
			c.AddTask(c.recoverCharge(c.Sim.Frame()), "charge", 900)
		}
	}
}
