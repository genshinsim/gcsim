package yoimiya

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/def"
)

func (c *char) Attack(p map[string]int) int {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f := c.ActionFrames(def.ActionAttack, p)

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			fmt.Sprintf("Normal %v", c.NormalCounter),
			def.AttackTagNormal,
			def.ICDTagNormalAttack,
			def.ICDGroupDefault,
			def.StrikeTypePierce,
			def.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, travel+f-5+i)
	}

	c.AdvanceNormalIndex()

	if c.Sim.Status("yoimiyaskill") > 0 {
		if c.lastPart < c.Sim.Frame() || c.lastPart == 0 {
			c.lastPart = c.Sim.Frame() + 300 //every 5 second
			count := 2
			if c.Sim.Rand().Float64() < 0.5 {
				count = 3
			}
			c.QueueParticle("yoimiya", count, def.Pyro, travel+f)
		}
	}

	return f
}

func (c *char) onExit() {
	c.Sim.AddEventHook(func(s def.Sim) bool {
		//do nothing if yoi becomes active
		if s.ActiveCharIndex() == c.Index {
			return false
		}
		//do nothing if skill not active
		if s.Status("yoimiyaskill") == 0 {
			return false
		}
		//so here we have active char != yoi and skill is still
		//active; so we need to deactivate
		s.DeleteStatus("yoimiyaskill")
		return false
	}, "yoimiya-off", def.PostSwapHook)
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)

	c.Sim.AddStatus("yoimiyaskill", 600) //activate for 10
	// log.Println(c.Sim.Status("yoimiyaskill"))

	if c.Sim.Status("yoimiyaa2") == 0 {
		c.a2stack = 0
	}

	c.SetCD(def.ActionSkill, 1080)
	return f
}

func (c *char) Burst(p map[string]int) int {

	f := c.ActionFrames(def.ActionBurst, p)

	//assume it does skill dmg at end of it's animation
	d := c.Snapshot(
		"Aurous Blaze",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Pyro,
		50,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, f)

	//marker an opponent after first hit
	//ignore the bouncing around for now (just assume it's always target 0)
	//icd of 2s, removed if down
	duration := 600
	if c.Base.Cons > 0 {
		duration = 840
	}
	c.AddTask(func() {
		c.Sim.AddStatus("aurous", duration)
		//attack buff if stacks
		if c.Sim.Status("yoimiyaa2") > 0 {
			val := make([]float64, def.EndStatType)
			val[def.ATKP] = 0.1 + float64(c.a2stack)*0.01
			for i, char := range c.Sim.Characters() {
				if i == c.Index {
					continue
				}
				char.AddMod(def.CharStatMod{
					Key:    "yoimiya-a4",
					Expiry: 900, //15s
					Amount: func(a def.AttackTag) ([]float64, bool) {
						return val, true
					},
				})
			}
		} else {
			c.a2stack = 0
		}
	}, "start-blaze", f)

	//add cooldown to sim
	c.SetCD(def.ActionBurst, 15*60)
	//use up energy
	c.Energy = 0

	return f
}

func (c *char) burstHook() {
	//check on attack landed for target 0
	//if aurous active then trigger dmg if not on cd
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if c.Sim.Status("aurous") == 0 {
			return
		}
		if ds.ActorIndex == c.Index {
			//ignore for self
			return
		}
		//ignore if on icd
		if c.Sim.Status("aurousicd") > 0 {
			return
		}
		//ignore if wrong tags
		switch ds.AttackTag {
		case def.AttackTagNormal:
		case def.AttackTagExtra:
		case def.AttackTagPlunge:
		case def.AttackTagElementalArt:
		case def.AttackTagElementalBurst:
		default:
			return
		}
		//do explosion, set icd
		d := c.Snapshot(
			"Aurous Blaze (Explode)",
			def.AttackTagElementalBurst,
			def.ICDTagElementalBurst,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Pyro,
			25,
			burstExplode[c.TalentLvlBurst()],
		)
		d.Targets = def.TargetAll
		c.QueueDmg(&d, 1)
		c.Sim.AddStatus("aurousicd", 120) //2 sec icd
	}, "yoimiya-burst-check")

	if c.Sim.Flags().HPMode {
		//add check for if yoimiya dies
		c.Sim.AddOnHurt(func(s def.Sim) {
			if c.HPCurrent <= 0 {
				c.Sim.DeleteStatus("aurous")
			}
		})
	}
}
