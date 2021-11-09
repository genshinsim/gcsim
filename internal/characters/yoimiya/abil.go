package yoimiya

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	var totalMV float64

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			fmt.Sprintf("Normal %v", c.NormalCounter+1),
			// "Normal",
			core.AttackTagNormal,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, travel+f-5+i)
		totalMV += mult[c.TalentLvlAttack()]
	}

	c.AdvanceNormalIndex()

	if c.Base.Cons == 6 && c.Core.Rand.Float64() < 0.5 {
		//trigger attack
		d := c.Snapshot(
			fmt.Sprintf("Kindling (C6) - N%v", c.NormalCounter+1),
			// "Kindling (C6)",
			core.AttackTagNormal,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Pyro,
			25,
			totalMV,
		)
		c.QueueDmg(&d, travel+f+5)

	}

	if c.Core.Status.Duration("yoimiyaskill") > 0 {
		if c.lastPart < c.Core.F || c.lastPart == 0 {
			c.lastPart = c.Core.F + 300 //every 5 second
			count := 2
			if c.Core.Rand.Float64() < 0.5 {
				count = 3
			}
			c.QueueParticle("yoimiya", count, core.Pyro, travel+f)
		}
	}

	return f, a
}

func (c *char) onExit() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		next := args[1].(int)
		if prev == c.Index && next != c.Index {
			if c.Core.Status.Duration("yoimiyaskill") > 0 {
				c.Core.Status.DeleteStatus("yoimiyaskill")
			}
		}
		return false
	}, "yoimiya-exit")
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	c.Core.Status.AddStatus("yoimiyaskill", 600) //activate for 10
	// log.Println(c.Core.Status.Duration("yoimiyaskill"))

	if c.Core.Status.Duration("yoimiyaa2") == 0 {
		c.a2stack = 0
	}

	c.SetCD(core.ActionSkill, 1080)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	//assume it does skill dmg at end of it's animation
	d := c.Snapshot(
		"Aurous Blaze",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Pyro,
		50,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f)

	//marker an opponent after first hit
	//ignore the bouncing around for now (just assume it's always target 0)
	//icd of 2s, removed if down
	duration := 600
	if c.Base.Cons > 0 {
		duration = 840
	}
	c.AddTask(func() {
		c.Core.Status.AddStatus("aurous", duration)
		//attack buff if stacks
		if c.Core.Status.Duration("yoimiyaa2") > 0 {
			val := make([]float64, core.EndStatType)
			val[core.ATKP] = 0.1 + float64(c.a2stack)*0.01
			for i, char := range c.Core.Chars {
				if i == c.Index {
					continue
				}
				char.AddMod(core.CharStatMod{
					Key:    "yoimiya-a4",
					Expiry: c.Core.F + 900, //15s
					Amount: func(a core.AttackTag) ([]float64, bool) {
						return val, true
					},
				})
			}
		} else {
			c.a2stack = 0
		}
	}, "start-blaze", f)

	//add cooldown to sim
	c.SetCD(core.ActionBurst, 15*60)
	//use up energy
	c.Energy = 0

	return f, a
}

func (c *char) burstHook() {
	//check on attack landed for target 0
	//if aurous active then trigger dmg if not on cd
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if c.Core.Status.Duration("aurous") == 0 {
			return false
		}
		if ds.ActorIndex == c.Index {
			//ignore for self
			return false
		}
		//ignore if on icd
		if c.Core.Status.Duration("aurousicd") > 0 {
			return false
		}
		//ignore if wrong tags
		switch ds.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
		case core.AttackTagPlunge:
		case core.AttackTagElementalArt:
		case core.AttackTagElementalBurst:
		default:
			return false
		}
		//do explosion, set icd
		d := c.Snapshot(
			"Aurous Blaze (Explode)",
			core.AttackTagElementalBurst,
			core.ICDTagElementalBurst,
			// core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Pyro,
			25,
			burstExplode[c.TalentLvlBurst()],
		)
		d.Targets = core.TargetAll
		c.QueueDmg(&d, 1)
		c.Core.Status.AddStatus("aurousicd", 120) //2 sec icd

		//check for c4

		if c.Base.Cons >= 4 {
			c.ReduceActionCooldown(core.ActionSkill, 72)
		}

		return false

	}, "yoimiya-burst-check")

	if c.Core.Flags.DamageMode {
		//add check for if yoimiya dies
		c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
			if c.HPCurrent <= 0 {
				c.Core.Status.DeleteStatus("aurous")
			}
			return false
		}, "yoimiya-died")
	}
}
