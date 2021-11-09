package xiao

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			fmt.Sprintf("Normal %v", c.NormalCounter),
			core.AttackTagNormal,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSpear,
			core.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, f-5+i)
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

	d := c.Snapshot(
		"Charge",
		core.AttackTagExtra,
		core.ICDTagExtraAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSpear,
		core.Physical,
		25,
		charge[c.TalentLvlAttack()],
	)
	// Kind of hits multiple, but radius not terribly big. Coded as single target for now

	c.QueueDmg(&d, f-1)

	//return animation cd
	return f, a
}

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) PlungeAttack(delay int) (int, int) {
	d := c.Snapshot(
		"Plunge (Normal)",
		core.AttackTagPlunge,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeSpear,
		core.Physical,
		25,
		plunge[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, delay)

	//return animation cd
	return delay, delay
}

// High Plunge attack damage queue generator
// Use the "plunge_hits" optional argument to determine how many plunge falling hits you do on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionHighPlunge, p)

	plunge_hits, ok := p["plunge_hits"]
	if !ok {
		plunge_hits = 0 // Number of normal plunge hits
	}

	for i := 0; i < plunge_hits; i++ {
		// Add plunge attack in each frame leading up to final hit for now - not sure we have clear mechanics on this
		// TODO: Perhaps amend later, but functionally in combat you usually get at most one of these anyway
		c.PlungeAttack(f - i - 1)
	}

	d := c.Snapshot(
		"High Plunge",
		core.AttackTagPlunge,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Physical,
		25,
		highplunge[c.TalentLvlAttack()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f-1)

	//return animation cd
	return f, a
}

// Low Plunge attack damage queue generator
// Use the "plunge_hits" optional argument to determine how many plunge falling hits you do on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionLowPlunge, p)

	plunge_hits, ok := p["plunge_hits"]
	if !ok {
		plunge_hits = 0 // Number of normal plunge hits
	}

	for i := 0; i < plunge_hits; i++ {
		// Add plunge attack in each frame leading up to final hit for now - not sure we have clear mechanics on this
		// TODO: Perhaps amend later, but functionally in combat you usually get at most one of these anyway
		c.PlungeAttack(f - i - 1)
	}

	d := c.Snapshot(
		"Low Plunge",
		core.AttackTagPlunge,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Physical,
		25,
		lowplunge[c.TalentLvlAttack()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f-1)

	//return animation cd
	return f, a
}

// Skill attack damage queue generator
// Additionally implements A4
// Using Lemniscatic Wind Cycling increases the DMG of subsequent uses of Lemniscatic Wind Cycling by 15%. This effect lasts for 7s and has a maximum of 3 stacks. Gaining a new stack refreshes the duration of this effect.
func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)
	// According to KQM library, double E is only 60 frames long whereas single cast is 36
	// No idea how this works, but add a special case to reduce the frames of the 2nd cast
	// TODO: No data listed on how 3 casts work - this might be too few frames compared to actual 3x usage
	if c.Core.LastAction.Target == "xiao" && c.Core.LastAction.Typ == core.ActionSkill {
		f = 60 - f
		a = 60 - a
	}

	d := c.Snapshot(
		"Lemniscatic Wind Cycling",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeSpear,
		core.Anemo,
		25,
		skill[c.TalentLvlSkill()],
	)
	// Pierces enemies, so I guess it targets all?
	d.Targets = core.TargetAll

	// Add damage based on A4
	if c.a4Expiry <= c.Core.F {
		c.Tags["a4"] = 0
	}
	stacks := c.Tags["a4"]
	d.Stats[core.DmgP] += float64(stacks) * 0.15

	// Text is not explicit, but assume that gaining a stack while at max still refreshes duration
	c.Tags["a4"]++
	c.a4Expiry = c.Core.F + 420
	if c.Tags["a4"] > 3 {
		c.Tags["a4"] = 3
	}
	c.Core.Log.Debugw("Xiao A4 adding damage", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "stacks", stacks, "expiry", c.a4Expiry)

	c.QueueDmg(&d, f)

	// Cannot create energy during burst uptime
	if c.Core.Status.Duration("xiaoburst") > 0 {
	} else {
		c.QueueParticle("xiao", 3, core.Anemo, f+100)
	}

	// C6 handling - can use skill ignoring CD and without draining charges
	// Can simply return early
	if c.Base.Cons == 6 && c.Core.Status.Duration("xiaoc6") > 0 {
		c.Core.Log.Debugw("xiao c6 active, Xiao E used, no charge used, no CD", "frame", c.Core.F, "event", core.LogCharacterEvent, "c6 remaining duration", c.Core.Status.Duration("xiaoc6"))
		return f, a
	}

	// Handle E charges
	if c.eCharge == 1 {
		c.SetCD(core.ActionSkill, c.eNextRecover)
	} else {
		c.eNextRecover = c.Core.F + 601
		c.Core.Log.Debugw("xiao e charge used, queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "recover at", c.eNextRecover)
		c.AddTask(c.recoverCharge(c.Core.F), "charge", 600)
		c.eTickSrc = c.Core.F
	}
	c.eCharge--

	return f, a
}

// Helper function that queues up Xiao e charge recovery - similar to other charge recovery functions
func (c *char) recoverCharge(src int) func() {
	return func() {
		// Required stopper for recursion
		if c.eTickSrc != src {
			c.Core.Log.Debugw("xiao e recovery function ignored, src diff", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "src", src, "new src", c.eTickSrc)
			return
		}
		c.eCharge++
		c.Core.Log.Debugw("xiao e recovering a charge", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "skill last used at", src, "total charges", c.eCharge)
		c.SetCD(core.ActionSkill, 0)
		if c.eCharge >= c.eChargeMax {
			return
		}

		c.eNextRecover = c.Core.F + 601
		c.Core.Log.Debugw("xiao e charge queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "recover at", c.eNextRecover)
		c.AddTask(c.recoverCharge(src), "charge", 600)
	}
}

// Sets Xiao's burst damage state
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// Per previous code, believe that the burst duration starts ticking down from after the animation is done
	// TODO: No indication of that in library though
	c.Core.Status.AddStatus("xiaoburst", 900+f)
	c.qStarted = c.Core.F

	// HP Drain - removes HP every 1 second tick after burst is activated
	// Per gameplay video, HP ticks start after animation is finished
	for i := f + 60; i < 900+f; i++ {
		c.AddTask(func() {
			if c.Core.Status.Duration("xiaoburst") > 0 {
				c.HPCurrent = c.HPCurrent * (1 - burstDrain[c.TalentLvlBurst()])
			}
		}, "xiaoburst-hp-drain", i)
	}

	// Checked gameplay - burst starts ticking down from activation. CD is 16.6 seconds after animation is done
	c.SetCD(core.ActionBurst, 18*60)
	c.Energy = 0

	return f, a
}

// Xiao specific Snapshot implementation for his burst bonuses. Similar to Hu Tao
// Implements burst anemo attack damage conversion and DMG bonus
// Also implements A1:
// While under the effects of Bane of All Evil, all DMG dealt by Xiao is increased by 5%. DMG is increased by an additional 5% for every 3s the ability persists. The maximum DMG Bonus is 25%
func (c *char) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	if c.Core.Status.Duration("xiaoburst") > 0 {
		// Calculate and add A1 damage bonus - applies to all damage
		// Fraction dropped in int conversion in go - acts like floor
		stacks := 1 + int((c.Core.F-c.qStarted)/180)
		if stacks > 5 {
			stacks = 5
		}
		ds.Stats[core.DmgP] += float64(stacks) * 0.05
		c.Core.Log.Debugw("a1 adding dmg %", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "stacks", stacks, "final", ds.Stats[core.DmgP], "time since burst start", c.Core.F-c.qStarted)

		// Anemo conversion and dmg bonus application to normal, charged, and plunge attacks
		// Also handle burst CA ICD change to share with Normal
		switch ds.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
			ds.ICDTag = core.ICDTagNormalAttack
		case core.AttackTagPlunge:
		default:
			return ds
		}
		ds.Element = core.Anemo
		bonus := burstBonus[c.TalentLvlBurst()]
		ds.Stats[core.DmgP] += bonus
		c.Core.Log.Debugw("xiao burst damage bonus", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "bonus", bonus, "final", ds.Stats[core.DmgP])
	}
	return ds
}
