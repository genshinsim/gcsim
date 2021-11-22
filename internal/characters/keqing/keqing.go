package keqing

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("keqing", NewChar)
}

type char struct {
	*character.Tmpl
	eStartFrame int
	c2ICD       int
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
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

var delay = [][]int{{8}, {20}, {25}, {25, 35}, {34}}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 11
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 1:
			f = 33 - 11
		case 2:
			f = 60 - 33
		case 3:
			f = 97 - 60
		case 4:
			f = 133 - 97
		}
		return f, f
	case core.ActionCharge:
		return 52, 52
	case core.ActionSkill:
		if c.Tags["e"] == 1 {
			//2nd part
			return 84, 84
		}
		//first part
		return 34, 34
	case core.ActionBurst:
		return 125, 125
	}
	c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
	return 0, 0
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

func (c *char) c4() {
	c.Core.Events.Subscribe(core.OnTransReaction, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		switch ds.ReactionType {
		case core.Overload:
			fallthrough
		case core.ElectroCharged:
			fallthrough
		case core.Superconduct:
			fallthrough
		case core.SwirlElectro:
			fallthrough
		case core.CrystallizeElectro:
			val := make([]float64, core.EndStatType)
			val[core.ATK] = 0.25
			c.AddMod(core.CharStatMod{
				Key:    "c4",
				Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
				Expiry: c.Core.F + 600,
			})
		}
		return false
	}, "keqingc4")

}

func (c *char) Attack(p map[string]int) (int, int) {
	//apply attack speed
	f, a := c.ActionFrames(core.ActionAttack, p)

	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		0,
	)

	for i, mult := range attack[c.NormalCounter] {
		x := d.Clone()
		x.Mult = mult[c.TalentLvlAttack()]
		c.QueueDmg(&x, delay[c.NormalCounter][i])
	}

	if c.Base.Cons == 6 {
		c.activateC6("attack")
	}

	c.AdvanceNormalIndex()
	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	d := c.Snapshot(
		"Charge 1",
		core.AttackTagExtra,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		0,
	)
	d.Targets = core.TargetAll

	for i, mult := range charge {
		x := d.Clone()
		x.Mult = mult[c.TalentLvlAttack()]
		x.Abil = fmt.Sprintf("Charge %v", i)
		c.QueueDmg(&x, f-i*10-5)
	}

	if c.Tags["e"] == 1 {
		//2 hits
		for i := 0; i < 2; i++ {
			d := c.Snapshot(
				"Stellar Restoration (Thunderclap)",
				core.AttackTagElementalArt,
				core.ICDTagElementalArt,
				core.ICDGroupDefault,
				core.StrikeTypeSlash,
				core.Electro,
				50,
				skillCA[c.TalentLvlSkill()],
			)
			d.Targets = core.TargetAll
			c.QueueDmg(&d, f)
		}

		//place on cooldown
		c.Tags["e"] = 0
		// c.CD[def.SkillCD] = c.eStartFrame + 100
		c.SetCD(core.ActionSkill, c.eStartFrame+450-c.Core.F)
	}

	if c.Base.Cons == 6 {
		c.activateC6("charge")
	}

	return f, a
}

func (c *char) c2() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		if c.Core.F < c.c2ICD {
			return false
		}
		if c.Core.Rand.Float64() < 0.5 {
			c.c2ICD = c.Core.F + 300
			c.QueueParticle("keqing", 1, core.Electro, 100)
			c.Core.Log.Debugw("keqing c2 proc'd", "frame", c.Core.F, "event", core.LogCharacterEvent, "next ready", c.c2ICD)
		}
		return false
	}, "keqingc2")
}

func (c *char) Skill(p map[string]int) (int, int) {
	if c.Tags["e"] == 1 {
		return c.skillNext(p)
	}
	return c.skillFirst(p)
}

func (c *char) skillFirst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	d := c.Snapshot(
		"Stellar Restoration",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		skill[c.TalentLvlSkill()],
	)

	c.QueueDmg(&d, f)

	c.Tags["e"] = 1
	c.eStartFrame = c.Core.F

	//place on cd after certain frames if started is still true
	//looks like the E thing lasts 5 seconds
	c.AddTask(func() {
		if c.Tags["e"] == 1 {
			c.Tags["e"] = 0
			// c.CD[def.SkillCD] = c.eStartFrame + 100
			c.SetCD(core.ActionSkill, c.eStartFrame+450-c.Core.F) //TODO: cooldown if not triggered, 7.5s
		}
	}, "keqing-skill-cd", c.Core.F+300) //TODO: check this

	if c.Base.Cons == 6 {
		c.activateC6("skill")
	}

	return f, a
}

func (c *char) skillNext(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	d := c.Snapshot(
		"Stellar Restoration (Slashing)",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Electro,
		50,
		skillPress[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f)

	//add electro infusion

	c.Core.Status.AddStatus("keqinginfuse", 300)

	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "a2",
		Ele:    core.Electro,
		Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.F + 300,
	})

	if c.Base.Cons >= 1 {
		//2 tick dmg at start to end
		hits, ok := p["c2"]
		if !ok {
			hits = 1 //default 1 hit
		}
		d := c.Snapshot(
			"Stellar Restoration (Slashing)",
			core.AttackTagElementalArtHold,
			core.ICDTagElementalArt,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Electro,
			25,
			0.5,
		)
		for i := 0; i < hits; i++ {
			x := d.Clone()
			c.QueueDmg(&x, f)
		}
	}

	//place on cooldown
	c.Tags["e"] = 0
	c.SetCD(core.ActionSkill, c.eStartFrame+450-c.Core.F)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	//a4 increase crit + ER
	val := make([]float64, core.EndStatType)
	val[core.CR] = 0.15
	val[core.ER] = 0.15
	c.AddMod(core.CharStatMod{
		Key:    "a4",
		Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Core.F + 480,
	})

	//first hit 70 frame
	//first tick 74 frame
	//last tick 168
	//last hit 211

	//initial
	initial := c.Snapshot(
		"Starward Sword",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Electro,
		25,
		burstInitial[c.TalentLvlBurst()],
	)
	initial.Targets = core.TargetAll

	c.QueueDmg(&initial, 70)

	//8 hits
	dot := c.Snapshot(
		"Starward Sword (Tick)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Electro,
		25,
		burstDot[c.TalentLvlBurst()],
	)
	dot.Targets = core.TargetAll
	for i := 70; i < 170; i += 13 {
		c.QueueDmg(&dot, i)
	}

	//final
	final := c.Snapshot(
		"Starward Sword (Tick)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Electro,
		25,
		burstFinal[c.TalentLvlBurst()],
	)
	final.Targets = core.TargetAll

	c.QueueDmg(&final, 211)

	if c.Base.Cons == 6 {
		c.activateC6("burst")
	}

	c.Energy = 0
	// c.CD[def.BurstCD] = c.Core.F + 720 //12s
	c.SetCD(core.ActionBurst, 720)

	return f, a
}

func (c *char) activateC6(src string) {
	val := make([]float64, core.EndStatType)
	val[core.ElectroP] = 0.06
	c.AddMod(core.CharStatMod{
		Key:    src,
		Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Core.F + 480,
	})
}
