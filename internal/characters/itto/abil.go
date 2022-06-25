package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	if c.dasshuUsed {
		c.NormalCounter = c.dasshuCount
	}

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	// Check burst status
	r := 1.0
	if c.Core.Status.Duration("ittoq") > 0 {
		r = 2
		// Burst can expire during normals
		// If burst lasts to hitlag, extend burst
		if f < c.Core.Status.Duration("ittoq") {
			c.Core.Status.ExtendStatus("ittoq", 1)
		}
	}

	// Add superlative strength stacks on damage
	if c.Core.Status.Duration("ittoq") > 0 && c.NormalCounter < 3 {
		c.Tags["strStack"]++
	} else if c.NormalCounter == 1 {
		c.Tags["strStack"]++
	} else if c.NormalCounter == 3 {
		c.Tags["strStack"] += 2
	}
	if c.Tags["strStack"] > 5 {
		c.Tags["strStack"] = 5
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(r, false, core.TargettableEnemy), f, a)

	c.sCACount = 0
	c.dasshuUsed = false
	defer c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	// Check burst status
	r := 1.0
	if c.Core.Status.Duration("ittoq") > 0 {
		// Unsure of range, it's huge though
		r = 3
		// If burst will expire, extend for the CA
		if f > c.Core.Status.Duration("ittoq") {
			c.Core.Status.ExtendStatus("ittoq", f)
		} else {
			// Extend burst by hitlag value
			c.Core.Status.ExtendStatus("ittoq", 1)
		}
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Charged %v Stacks %v", c.sCACount, c.Tags["strStack"]),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       akCombo[c.TalentLvlAttack()],
		FlatDmg:    0.35*c.Base.Def*(1+c.Stats[core.DEFP]) + c.Stats[core.DEF],
	}
	if c.Tags["strStack"] == 0 {
		ai.Mult = saichiSlash[c.TalentLvlAttack()]
		ai.FlatDmg = 0
	} else if c.Tags["strStack"] == 1 {
		ai.Mult = akFinal[c.TalentLvlAttack()]
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(r, false, core.TargettableEnemy), f, f)

	c.Tags["strStack"]--
	c.sCACount++
	if c.Tags["strStack"] <= 0 {
		c.sCACount = 0
		c.Tags["strStack"] = 0
	}

	c.dasshuUsed = false

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	// Added "travel" parameter for future, since Ushi is thrown and takes 12 frames to hit the ground from a press E
	travel, ok := p["travel"]
	if !ok {
		travel = 3
	}

	//deal damage when created
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ushi Throw",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	// Ushi callback to create construct
	done := false
	cb := func(a core.AttackCB) {
		if done {
			return
		}
		c.Core.Constructs.New(c.newUshi(360), true) // 6 seconds from hit/land
		done = true
	}

	// Assume that Ushi always hits for a stack
	c.Tags["strStack"]++
	if c.Tags["strStack"] > 5 {
		c.Tags["strStack"] = 5
	}

	count := 3
	if c.Core.Rand.Float64() < 0.33 {
		count = 4
	}
	c.QueueParticle(c.Name(), count, core.Geo, f+100)

	// Check if used as Dasshu
	if c.NormalCounter > 0 {
		c.dasshuUsed = true
		c.dasshuCount = c.NormalCounter
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f+travel, cb)
	c.sCACount = 0

	c.SetCD(core.ActionSkill, 10*60)
	return f, a
}

// Adapted from Noelle
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// Add mod for def to attack burst conversion
	val := make([]float64, core.EndStatType)

	// Generate a "fake" snapshot in order to show a listing of the applied mods in the debug
	aiSnapshot := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Royal Descent: Behold, Itto the Evil! (Stat Snapshot)",
	}
	snapshot := c.Snapshot(&aiSnapshot)
	burstDefSnapshot := snapshot.BaseDef*(1+snapshot.Stats[core.DEFP]) + snapshot.Stats[core.DEF]
	// burstDefSnapshot := c.Base.Def*(1+c.Stats[core.DEFP]) + c.Stats[core.DEF]
	mult := defconv[c.TalentLvlBurst()]
	fa := mult * burstDefSnapshot
	val[core.ATK] = fa

	// Not sure if something else in the code can modify this - to be safe, copy this for the burst extension
	valCopy := make([]float64, core.EndStatType)
	copy(valCopy, val)

	// TODO: Confirm exact timing of buff - for now matched to status duration previously set, which is 900 + animation frames
	// Buff lasts 11.55s after anim, padded to cover basic combo
	c.AddMod(core.CharStatMod{
		Key:    "itto-burst",
		Expiry: c.Core.F + 960 + f,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})
	c.Core.Log.NewEvent("itto burst", core.LogSnapshotEvent, c.Index, "frame", c.Core.F, "total def", burstDefSnapshot, "atk added", fa, "mult", mult)

	c.Core.Status.AddStatus("ittoq", 960+f) // inflated from 11.55 seconds to cover basic combo

	if c.Base.Cons >= 4 {
		val := make([]float64, core.EndStatType)
		val[core.ATKP] = 0.2
		val[core.DEFP] = 0.2
		c.c4cb(960+f, val)
	}

	if c.Base.Cons >= 1 {
		c.Tags["strStack"] += 2
		if c.Tags["strStack"] > 5 {
			c.Tags["strStack"] = 5
		}
		for frame := 140; frame <= 200; frame += 30 {
			c.AddTask(func() {
				if c.Tags["strStack"] <= 4 {
					c.Tags["strStack"]++
				}
			}, "c1-itto", frame)
		}
	}

	if c.Base.Cons >= 2 {
		c.AddTask(func() {
			count := 0
			for _, char := range c.Core.Chars {
				if char.Ele() == core.Geo {
					count++
				}
			}
			if count > 3 {
				c.AddEnergy("itto-c2", 3*6)
				c.ReduceActionCooldown(core.ActionBurst, 3*1.5*60)
			} else {
				c.AddEnergy("itto-c2", float64(count)*6)
				c.ReduceActionCooldown(core.ActionBurst, count*(1.5*60))
			}
		}, "c2-itto", 9)

	}
	c.SetCDWithDelay(core.ActionBurst, 1080, 8)
	c.ConsumeEnergy(8)
	return f, a
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		c.Core.Status.DeleteStatus("ittoq")
		// Re-add mod with 0 time to remove
		c.AddMod(core.CharStatMod{
			Key:    "itto-burst",
			Expiry: c.Core.F,
			Amount: func() ([]float64, bool) {
				return make([]float64, core.EndStatType), true
			},
		})
		return false
	}, "itto-exit")
}

func (c *char) c4cb(delay int, buff []float64) func() {
	return func() {
		c.AddTask(func() {
			if c.Core.Status.Duration("ittoq") > 0 {
				c.c4cb(c.Core.Status.Duration("ittoq"), buff)
			} else {
				for _, char := range c.Core.Chars {
					char.AddMod(core.CharStatMod{
						Key:    "itto-c4",
						Expiry: c.Core.F + 10*60,
						Amount: func() ([]float64, bool) {
							return buff, true
						},
					})
				}
			}
		}, "ittoqend", delay)
	}
}
