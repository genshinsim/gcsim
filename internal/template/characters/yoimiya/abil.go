package yoimiya

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{{15, 24}, {17}, {25}, {11, 26}, {17}}

func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	var totalMV float64

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter+1),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		totalMV += mult[c.TalentLvlAttack()]

		// TODO - double check snapshotDelay
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), hitmarks[c.NormalCounter][i], travel+hitmarks[c.NormalCounter][i])
	}

	if c.Base.Cons == 6 && c.Core.Status.Duration("yoimiyaskill") > 0 && c.Core.Rand.Float64() < 0.5 {
		//trigger attack
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Kindling (C6) - N%v", c.NormalCounter+1),
			AttackTag:  core.AttackTagNormal,
			ICDTag:     core.ICDTagNormalAttack,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Pyro,
			Durability: 25,
			Mult:       totalMV * 0.6,
		}
		//TODO: frames?
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, travel+f+5)

	}

	c.AdvanceNormalIndex()

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

	c.Core.Status.AddStatus("yoimiyaskill", 600+f) //activate for 10
	// log.Println(c.Core.Status.Duration("yoimiyaskill"))

	if c.Core.Status.Duration("yoimiyaa1") == 0 {
		c.a1stack = 0
	}

	c.SetCDWithDelay(core.ActionSkill, 1080, 11)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	//assume it does skill dmg at end of it's animation
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Aurous Blaze",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, f)

	//marker an opponent after first hit
	//ignore the bouncing around for now (just assume it's always target 0)
	//icd of 2s, removed if down
	duration := 600
	if c.Base.Cons > 0 {
		duration = 840
	}
	c.AddTask(func() {
		c.Core.Status.AddStatus("aurous", duration)
		val := make([]float64, core.EndStatType)
		//attack buff
		val[core.ATKP] = 0.1 + float64(c.a1stack)*0.01
		for i, char := range c.Core.Chars {
			if i == c.Index {
				continue
			}
			char.AddMod(core.CharStatMod{
				Key:    "yoimiya-a4",
				Expiry: c.Core.F + 900, //15s
				Amount: func() ([]float64, bool) {
					return val, true
				},
			})
		}
	}, "start-blaze", f)

	//add cooldown to sim
	c.SetCD(core.ActionBurst, 15*60)
	//use up energy
	c.ConsumeEnergy(5)

	return f, a
}

func (c *char) burstHook() {
	//check on attack landed for target 0
	//if aurous active then trigger dmg if not on cd
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		if c.Core.Status.Duration("aurous") == 0 {
			return false
		}
		if ae.Info.ActorIndex == c.Index {
			//ignore for self
			return false
		}
		//ignore if on icd
		if c.Core.Status.Duration("aurousicd") > 0 {
			return false
		}
		//ignore if wrong tags
		switch ae.Info.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
		case core.AttackTagPlunge:
		case core.AttackTagElementalArt:
		case core.AttackTagElementalBurst:
		default:
			return false
		}
		//do explosion, set icd
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Aurous Blaze (Explode)",
			AttackTag:  core.AttackTagElementalBurst,
			ICDTag:     core.ICDTagElementalBurst,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Pyro,
			Durability: 25,
			Mult:       burstExplode[c.TalentLvlBurst()],
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, 1)

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
			if c.HP() <= 0 {
				c.Core.Status.DeleteStatus("aurous")
			}
			return false
		}, "yoimiya-died")
	}
}

func (c *char) Dash(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionDash, p)
	return f, a
}

func (c *char) Jump(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionJump, p)
	return f, a
}
