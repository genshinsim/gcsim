package xiao

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{{4, 17}, {15}, {15}, {14, 31}, {16}, {39}}

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), hitmarks[c.NormalCounter][i], hitmarks[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
// Note that his CAs share an ICD with his NAs when he is under the effects of his burst
// TODO: No information available on whether regular CAs follow a similar pattern
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)

	//return animation cd
	return f, a
}

var collisionFrame = 38

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) PlungeAttack(delay int) (int, int) {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Plunge Collision",
		AttackTag:  core.AttackTagPlunge,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 0,
		Mult:       plunge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), delay, delay)

	//return animation cd
	return delay, delay
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionHighPlunge, p)

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xiao does a collision hit
	}

	if collision > 0 {
		c.PlungeAttack(collisionFrame)
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "High Plunge",
		AttackTag:  core.AttackTagPlunge,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       highplunge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f)

	//return animation cd
	return f, a
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionLowPlunge, p)

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xiao does a collision hit
	}

	if collision > 0 {
		c.PlungeAttack(collisionFrame)
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge",
		AttackTag:  core.AttackTagPlunge,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       lowplunge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f)

	//return animation cd
	return f, a
}

// Skill attack damage queue generator
// Additionally implements A4
// Using Lemniscatic Wind Cycling increases the DMG of subsequent uses of Lemniscatic Wind Cycling by 15%. This effect lasts for 7s and has a maximum of 3 stacks. Gaining a new stack refreshes the duration of this effect.
func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	// Add damage based on A4
	if c.a4Expiry <= c.Core.F {
		c.Tags["a4"] = 0
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lemniscatic Wind Cycling",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupXiaoDash,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), f)

	// Text is not explicit, but assume that gaining a stack while at max still refreshes duration
	c.Tags["a4"]++
	c.a4Expiry = c.Core.F + 420
	if c.Tags["a4"] > 3 {
		c.Tags["a4"] = 3
	}

	// Cannot create energy during burst uptime
	if c.Core.Status.Duration("xiaoburst") > 0 {
	} else {
		c.QueueParticle("xiao", 3, core.Anemo, f+100)
	}

	// C6 handling - can use skill ignoring CD and without draining charges
	// Can simply return early
	if c.Base.Cons == 6 && c.Core.Status.Duration("xiaoc6") > 0 {
		c.Core.Log.NewEvent("xiao c6 active, Xiao E used, no charge used, no CD", core.LogCharacterEvent, c.Index, "c6 remaining duration", c.Core.Status.Duration("xiaoc6"))
		return f, a
	}

	c.SetCD(core.ActionSkill, 600)

	return f, a
}

// Sets Xiao's burst damage state
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	var HPicd int
	HPicd = 0

	// Per previous code, believe that the burst duration starts ticking down from after the animation is done
	// TODO: No indication of that in library though
	c.Core.Status.AddStatus("xiaoburst", 900+f)
	c.qStarted = c.Core.F

	// HP Drain - removes HP every 1 second tick after burst is activated
	// Per gameplay video, HP ticks start after animation is finished
	for i := f + 60; i < 900+f; i++ {
		c.AddTask(func() {
			if c.Core.Status.Duration("xiaoburst") > 0 && c.Core.F >= HPicd {
				HPicd = c.Core.F + 60
				c.Core.Health.Drain(core.DrainInfo{
					ActorIndex: c.Index,
					Abil:       "Bane of All Evil",
					Amount:     burstDrain[c.TalentLvlBurst()] * c.HP(),
				})
			}
		}, "xiaoburst-hp-drain", i)
	}

	c.SetCDWithDelay(core.ActionBurst, 18*60, 29)
	c.ConsumeEnergy(36)

	return f, a
}

func (c *char) Dash(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionDash, p)
	return f, a
}

func (c *char) Jump(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionJump, p)
	return f, a
}
