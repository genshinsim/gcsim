package keqing

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("keqing", NewChar)
}

type char struct {
	*character.Tmpl
	eStartFrame int
	c2ICD       int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 40
	c.EnergyMax = 40
	c.Weapon.Class = def.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

var delay = [][]int{{8}, {20}, {25}, {25, 35}, {34}}

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			return 11
		case 1:
			return 33 - 11
		case 2:
			return 60 - 33
		case 3:
			return 97 - 60
		case 4:
			return 133 - 97
		}
	case def.ActionCharge:
		return 52
	case def.ActionSkill:
		if c.Tags["e"] == 1 {
			return 84 //2nd part
		}
		return 34 //first part
	case def.ActionBurst:
		return 125
	}
	c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
	return 0
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

func (c *char) c4() {
	c.Sim.AddOnTransReaction(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex != c.Index {
			return
		}
		switch ds.ReactionType {
		case def.Overload:
			fallthrough
		case def.ElectroCharged:
			fallthrough
		case def.Superconduct:
			fallthrough
		case def.SwirlElectro:
			fallthrough
		case def.CrystallizeElectro:
			val := make([]float64, def.EndStatType)
			val[def.ATK] = 0.25
			c.AddMod(def.CharStatMod{
				Key:    "c4",
				Amount: func(a def.AttackTag) ([]float64, bool) { return val, true },
				Expiry: c.Sim.Frame() + 600,
			})
		}

	}, "keqingc4")
}

func (c *char) Attack(p map[string]int) int {
	//apply attack speed
	f := c.ActionFrames(def.ActionAttack, p)

	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
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
	return f
}

func (c *char) ChargeAttack(p map[string]int) int {

	f := c.ActionFrames(def.ActionCharge, p)

	d := c.Snapshot(
		"Charge 1",
		def.AttackTagExtra,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
		25,
		0,
	)
	d.Targets = def.TargetAll

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
				def.AttackTagElementalArt,
				def.ICDTagElementalArt,
				def.ICDGroupDefault,
				def.StrikeTypeSlash,
				def.Electro,
				50,
				skillCA[c.TalentLvlSkill()],
			)
			d.Targets = def.TargetAll
			c.QueueDmg(&d, f)
		}

		//place on cooldown
		c.Tags["e"] = 0
		// c.CD[def.SkillCD] = c.eStartFrame + 100
		c.SetCD(def.ActionSkill, c.eStartFrame+450-c.Sim.Frame())
	}

	if c.Base.Cons == 6 {
		c.activateC6("charge")
	}

	return f
}

func (c *char) c2() {
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.Index {
			return
		}
		if c.Sim.Frame() < c.c2ICD {
			return
		}
		if c.Sim.Rand().Float64() < 0.5 {
			c.c2ICD = c.Sim.Frame() + 300
			c.QueueParticle("keqing", 1, def.Electro, 100)
			c.Log.Debugw("keqing c2 proc'd", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "next ready", c.c2ICD)
		}

	}, "keqingc2")
}

func (c *char) Skill(p map[string]int) int {
	if c.Tags["e"] == 1 {
		return c.skillNext(p)
	}
	return c.skillFirst(p)
}

func (c *char) skillFirst(p map[string]int) int {

	f := c.ActionFrames(def.ActionSkill, p)

	d := c.Snapshot(
		"Stellar Restoration",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Electro,
		25,
		skill[c.TalentLvlSkill()],
	)

	c.QueueDmg(&d, f)

	c.Tags["e"] = 1
	c.eStartFrame = c.Sim.Frame()

	//place on cd after certain frames if started is still true
	//looks like the E thing lasts 5 seconds
	c.AddTask(func() {
		if c.Tags["e"] == 1 {
			c.Tags["e"] = 0
			// c.CD[def.SkillCD] = c.eStartFrame + 100
			c.SetCD(def.ActionSkill, c.eStartFrame+450-c.Sim.Frame()) //TODO: cooldown if not triggered, 7.5s
		}
	}, "keqing-skill-cd", c.Sim.Frame()+300) //TODO: check this

	if c.Base.Cons == 6 {
		c.activateC6("skill")
	}

	return f
}

func (c *char) skillNext(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)

	d := c.Snapshot(
		"Stellar Restoration (Slashing)",
		def.AttackTagElementalArt,
		def.ICDTagElementalArt,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Electro,
		50,
		skillPress[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, f)

	//add electro infusion

	c.Sim.AddStatus("keqinginfuse", 300)

	c.AddWeaponInfuse(def.WeaponInfusion{
		Key:    "a2",
		Ele:    def.Electro,
		Tags:   []def.AttackTag{def.AttackTagNormal, def.AttackTagExtra, def.AttackTagPlunge},
		Expiry: c.Sim.Frame() + 300,
	})

	if c.Base.Cons >= 1 {
		//2 tick dmg at start to end
		hits, ok := p["c2"]
		if !ok {
			hits = 1 //default 1 hit
		}
		d := c.Snapshot(
			"Stellar Restoration (Slashing)",
			def.AttackTagElementalArtHold,
			def.ICDTagElementalArt,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Electro,
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
	c.SetCD(def.ActionSkill, c.eStartFrame+450-c.Sim.Frame())
	return f
}

func (c *char) Burst(p map[string]int) int {

	//a4 increase crit + ER
	val := make([]float64, def.EndStatType)
	val[def.CR] = 0.15
	val[def.ER] = 0.15
	c.AddMod(def.CharStatMod{
		Key:    "a4",
		Amount: func(a def.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Sim.Frame() + 480,
	})

	//first hit 70 frame
	//first tick 74 frame
	//last tick 168
	//last hit 211

	//initial
	initial := c.Snapshot(
		"Starward Sword",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Electro,
		25,
		burstInitial[c.TalentLvlBurst()],
	)
	initial.Targets = def.TargetAll

	c.QueueDmg(&initial, 70)

	//8 hits
	dot := c.Snapshot(
		"Starward Sword (Tick)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Electro,
		25,
		burstDot[c.TalentLvlBurst()],
	)
	dot.Targets = def.TargetAll
	for i := 70; i < 170; i += 13 {
		c.QueueDmg(&dot, i)
	}

	//final
	final := c.Snapshot(
		"Starward Sword (Tick)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Electro,
		25,
		burstFinal[c.TalentLvlBurst()],
	)
	final.Targets = def.TargetAll

	c.QueueDmg(&final, 211)

	if c.Base.Cons == 6 {
		c.activateC6("burst")
	}

	c.Energy = 0
	// c.CD[def.BurstCD] = c.Sim.Frame() + 720 //12s
	c.SetCD(def.ActionBurst, 720)
	return c.ActionFrames(def.ActionBurst, p)
}

func (c *char) activateC6(src string) {
	val := make([]float64, def.EndStatType)
	val[def.ElectroP] = 0.06
	c.AddMod(def.CharStatMod{
		Key:    src,
		Amount: func(a def.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Sim.Frame() + 480,
	})
}
