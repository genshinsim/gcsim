package sucrose

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
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
	eChargeMax          int
	c4Count             int
	eChargeLastRecovery int
	eLastUsed           int
	eChargeNextRecovery int
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
	c.Tags["eCharge"] = c.eChargeMax
	c.eChargeLastRecovery = 0
	c.eLastUsed = 0

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

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
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
		return f, f
	case core.ActionCharge:
		return 53, 53 //frames from keqing lib
	case core.ActionSkill:
		return 55, 55 //ok
	case core.ActionBurst:
		return 46, 46 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0, 0
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
	var val [core.EndStatType]float64
	val[core.EM] = 50
	for _, char := range c.Core.Chars {
		this := char
		if this.Ele() == core.Anemo || this.Ele() == core.Geo {
			continue //nothing for geo/anemo char
		}
		this.AddMod(core.CharStatMod{
			Key:    "sucrose-a2",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
				var f int
				var ok bool

				// c.Core.Log.Debugw("sucrose a2 check", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", this.CharIndex(), "ele", this.Ele())
				switch this.Ele() {
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
				// c.Core.Log.Debugw("sucrose a2 adding", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", this.CharIndex(), "ele", this.Ele(), "expiry", f, "ok", ok)
				return val, f > c.Core.F && ok
			},
		})
	}

	c.Core.Events.Subscribe(core.OnReactionOccured, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		switch ds.ReactionType {
		case core.SwirlCryo:
			c.Tags["a2-cryo"] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "reaction", ds.ReactionType, "expiry", c.Core.F+480)
		case core.SwirlElectro:
			c.Tags["a2-electro"] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "reaction", ds.ReactionType, "expiry", c.Core.F+480)
		case core.SwirlHydro:
			c.Tags["a2-hydro"] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "reaction", ds.ReactionType, "expiry", c.Core.F+480)
		case core.SwirlPyro:
			c.Tags["a2-pyro"] = c.Core.F + 480
			c.Core.Log.Debugw("sucrose a2 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "reaction", ds.ReactionType, "expiry", c.Core.F+480)
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
			Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
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
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			if c.Core.Status.Duration("sucrosec6") == 0 {
				return nil, false
			}
			p := core.EleToDmgP(c.qInfused)
			var val [core.EndStatType]float64
			val[p] = 0.2
			return val, true
		},
		Expiry: -1,
	})
}

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)

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

	c.c4()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)
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

	c.c4()

	return f, a
}

// Handles C4: Every 7 Normal and Charged Attacks, Sucrose will reduce the CD of Astable Anemohypostasis Creation-6308 by 1-7s
func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	c.c4Count++
	if c.c4Count < 7 {
		return
	}
	c.c4Count = 0

	// Change can be in float. See this Terrapin video for example
	// https://youtu.be/jB3x5BTYWIA?t=20
	cdReduction := 60 * int(c.Core.Rand.Float64()*6+1)

	c.eChargeNextRecovery -= cdReduction

	c.Core.Log.Debugw("sucrose c4 reducing E CD", "frame", c.Core.F, "event", core.LogCharacterEvent, "cd_reduction", cdReduction, "next_recovery", c.eChargeNextRecovery)
	if c.Core.F >= c.eChargeNextRecovery {
		c.recoverCharge()
		return
	}

	c.AddTask(c.recoverChargeWithCheck(c.Core.F), "sucrose-charge-recovery", c.eChargeNextRecovery-c.Core.F)
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
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
		return f, a
	}

	// Check
	switch c.Tags["eCharge"] {
	case 2:
		c.eChargeNextRecovery = c.Core.F + 900
	case 1:
		// When going from 1 to 0 charge, in game maintains the current CD of the skill
		// Need to add 1 to avoid same frame collision issues
		c.SetCD(core.ActionSkill, c.eChargeNextRecovery-c.Core.F+1)
	}
	c.eLastUsed = c.Core.F
	c.AddTask(c.recoverChargeWithCheck(c.Core.F), "sucrose-charge", c.eChargeNextRecovery-c.Core.F)
	c.Tags["eCharge"]--

	return f, a
}

func (c *char) recoverCharge() {
	c.Tags["eCharge"]++
	c.eChargeLastRecovery = c.Core.F
	c.eChargeNextRecovery = c.Core.F + 15*60
	c.ResetActionCooldown(core.ActionSkill)
	c.Core.Log.Debugw("sucrose recovered E charge", "frame", c.Core.F, "event", core.LogCharacterEvent, "charges", c.Tags["eCharge"], "nextRecovery", c.eChargeNextRecovery)
}

func (c *char) recoverChargeWithCheck(src int) func() {
	return func() {
		c.Core.Log.Debugw("sucrose E charge recovery check", "frame", c.Core.F, "event", core.LogCharacterEvent, "charges", c.Tags["eCharge"], "nextRecovery", c.eChargeNextRecovery, "src", src, "lastrecov", c.eChargeLastRecovery, "lastused", c.eLastUsed)
		if c.Core.F < c.eChargeNextRecovery {
			return
		}
		// Ensure that if C4 CD reduction procs, you don't recover multiple times
		if !((src >= c.eChargeLastRecovery) && (src >= c.eLastUsed)) {
			return
		}
		c.recoverCharge()
	}
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
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
	return f, a
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
