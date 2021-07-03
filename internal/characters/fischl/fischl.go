package fischl

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("fischl", NewChar)
}

type char struct {
	*character.Tmpl
	//field use for calculating oz damage
	ozSnapshot def.Snapshot

	ozSource      int //keep tracks of source of oz aka resets
	ozActiveUntil int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.MaxEnergy = 60
	c.Weapon.Class = def.WeaponClassBow
	c.NormalHitNum = 5

	c.ozSource = -1

	//register A4
	c.a4()

	if p.Base.Cons == 6 {
		c.c6()
	}

	// f.turbo()

	return &c, nil
}

// func (c *char) ozAttack() {
// 	d := c.ozSnapshot.Clone()
// 	d.Durability = 0
// 	if c.ozAttackCounter%4 == 0 {
// 		//apply aura, add to timer
// 		d.Durability = 25
// 		c.ozICD = c.Sim.Frame() + 300 //add 300 second to skill ICD
// 	}
// 	//so oz is active and ready to shoot, we add damage
// 	c.S.AddTask(func(s *def.Sim) {
// 		s.ApplyDamage(d)
// 	}, "Fischl Oz (Damage)", 1)
// 	//put shoot on cd
// 	c.ozNextShootReady = c.Sim.Frame() + 50
// 	//increment hit counter
// 	c.ozAttackCounter++
// 	//assume fischl has 60% chance of generating orb every attack;
// 	if c.S.Rand.Float64() < .6 {
// 		c.S.AddEnergyParticles("Fischl", 1, def.Electro, 120)
// 	}
// }

func (c *char) Attack(p map[string]int) int {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f := c.ActionFrames(def.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypePierce,
		def.Physical,
		25,
		auto[c.NormalCounter][c.TalentLvlAttack()],
	)
	c.QueueDmg(&d, travel+f)
	c.AdvanceNormalIndex()

	//check for c1
	if c.Base.Cons >= 1 && c.ozActiveUntil < c.Sim.Frame() {
		d := c.Snapshot(
			"Fischl C1",
			def.AttackTagNormal,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypePierce,
			def.Physical,
			100,
			0.22,
		)
		c.QueueDmg(&d, travel+f+1)
	}

	return f
}

func (c *char) queueOz(src string) {

	dur := 600
	if c.Base.Cons == 6 {
		dur += 120
	}
	c.ozActiveUntil = c.Sim.Frame() + dur
	c.ozSource = c.Sim.Frame()
	c.ozSnapshot = c.Snapshot(
		fmt.Sprintf("Oz (%v)", src),
		def.AttackTagElementalArt,
		def.ICDTagElementalArt,
		def.ICDGroupFischl,
		def.StrikeTypePierce,
		def.Electro,
		25,
		birdAtk[c.TalentLvlSkill()],
	)
	c.AddTask(c.ozTick(c.Sim.Frame()), "oz", 60)
	c.Log.Debugw("Oz activated", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "source", src, "expected end", c.ozActiveUntil, "next expected tick", c.Sim.Frame()+60)

	c.Sim.AddStatus("fischloz", dur)

}

func (c *char) ozTick(src int) func() {
	return func() {
		c.Log.Debugw("Oz checking for tick", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "src", src)
		//if src != ozSource then this is no longer the same oz, do nothing
		if src != c.ozSource {
			return
		}
		c.Log.Debugw("Oz ticked", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "next expected tick", c.Sim.Frame()+60, "active", c.ozActiveUntil, "src", src)
		//trigger damage
		d := c.ozSnapshot.Clone()
		c.Sim.ApplyDamage(&d)
		//check for orb
		if c.Sim.Rand().Float64() < .67 {
			c.QueueParticle("fischl", 1, def.Electro, 120)
		}

		//queue up next hit only if next hit oz is still active
		if c.Sim.Frame()+60 < c.ozActiveUntil {
			c.AddTask(c.ozTick(src), "oz", 60)
		}
	}
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)
	//always trigger electro no ICD on initial summon
	d := c.Snapshot(
		"Oz (Summon)",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupFischl,
		def.StrikeTypePierce,
		def.Electro,
		25,
		birdSum[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	if c.Base.Cons >= 2 {
		d.Mult += 2
	}
	c.QueueDmg(&d, 1) //queue initial damage

	//set on field oz to be this one
	c.AddTask(func() {
		c.queueOz("Skill")
	}, "oz-skill", f-20)

	c.SetCD(def.ActionSkill, 25*60)
	//return animation cd
	return f
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(def.ActionBurst, p)

	//set on field oz to be this one
	c.AddTask(func() {
		c.queueOz("Burst")
	}, "oz-skill", f-10)

	//initial damage; part of the burst tag
	d := c.Snapshot(
		"Midnight Phantasmagoria",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupFischl,
		def.StrikeTypeBlunt,
		def.Electro,
		25,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll
	c.QueueDmg(&d, 1)

	//check for C4 damage
	if c.Base.Cons >= 4 {
		d := c.Snapshot(
			"Midnight Phantasmagoria",
			def.AttackTagElementalBurst,
			def.ICDTagElementalBurst,
			def.ICDGroupFischl,
			def.StrikeTypePierce,
			def.Electro,
			50,
			2.22,
		)
		c.QueueDmg(&d, 1)
		//heal at end of animation
		heal := c.MaxHP() * 0.2
		c.AddTask(func() {
			c.Sim.HealActive(heal)
		}, "c4heal", f-1)

	}

	c.Energy = 0
	c.SetCD(def.ActionBurst, 15*60)
	return f //this is if you cancel immediately
}
