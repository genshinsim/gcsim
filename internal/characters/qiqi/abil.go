package qiqi

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack - nothing special
func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	for _, mult := range attack[c.NormalCounter] {
		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				fmt.Sprintf("Normal %v", c.NormalCounter),
				core.AttackTagNormal,
				core.ICDTagNormalAttack,
				core.ICDGroupDefault,
				core.StrikeTypeSlash,
				core.Physical,
				25,
				mult[c.TalentLvlAttack()],
			)
			return &d
		}, f)
	}

	c.AdvanceNormalIndex()

	return f, a
}

// Standard charge attack
func (c *char) Charge(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	for _, mult := range charge {
		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				"Charge",
				core.AttackTagExtra,
				core.ICDTagExtraAttack,
				core.ICDGroupDefault,
				core.StrikeTypeSlash,
				core.Physical,
				25,
				mult[c.TalentLvlAttack()],
			)
			return &d
		}, f)
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	// +1 to avoid end duration issues
	c.Core.Status.AddStatus("qiqiskill", 15*60+1)
	c.skillLastUsed = c.Core.F
	src := c.Core.F

	// Initial damage
	c.QueueDmgDynamic(func() *core.Snapshot {
		d := c.Snapshot(
			"Herald of Frost: Initial Damage",
			core.AttackTagElementalArt,
			core.ICDTagElementalArt,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Cryo,
			25,
			skillInitialDmg[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll

		// One healing proc happens immediately on cast
		c.Core.Health.HealActive(c.Index, c.heal(&d, skillHealContPer, skillHealContFlat, c.TalentLvlSkill()))

		return &d
	}, f)

	// Queue up continuous healing instances
	// No exact frame data on when the healing ticks happen. Just roughly guessing here
	// Healing ticks happen 3 additional times during the skill - assume ticks are roughly 4.5s apart
	// so in sec (0 = skill cast), 1, 5.5, 10, 14.5
	c.AddTask(c.skillHealTickTask(src), "qiqi-skill-heal-tick", f+4.5*60)

	// Queue up damage swipe instances.
	// No exact frame data on when the damage ticks happen. Just roughly guessing here
	// Occurs 9 times over the course of the skill
	// Once shortly after initial cast, then 8 additional procs over the rest of the duration
	// Each proc occurs in "pairs" of two swipes each spaced around 2.25s apart
	// The time between each swipe in a pair is about 1s
	// No exact frame data available plus the skill duration is affected by hitlag
	// Damage procs occur (in sec 0 = skill cast): 1.5, 3.75, 4.75, 7, 8, 10.25, 11.25, 13.5, 14.5
	c.AddTask(c.skillDmgTickTask(src, 60), "qiqi-skill-dmg-tick", f+30)

	c.SetCD(core.ActionSkill, 30*60)

	return f, a
}

// Handles skill damage swipe instances
// Also handles C1:
// When the Herald of Frost hits an opponent marked by a Fortune-Preserving Talisman, Qiqi regenerates 2 Energy.
func (c *char) skillDmgTickTask(src int, lastTickDuration int) func() {
	return func() {
		if c.Core.Status.Duration("qiqiskill") == 0 {
			return
		}

		// TODO: Not sure how this interacts with sac sword... Treat it as only one instance can be up at a time for now
		if c.skillLastUsed > src {
			return
		}

		d := c.Snapshot(
			"Herald of Frost: Skill Damage",
			core.AttackTagElementalArt,
			core.ICDTagElementalArt,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Cryo,
			25,
			skillDmgCont[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll

		if c.Base.Cons >= 1 {
			d.OnHitCallback = c.c1
		}

		c.Core.Combat.ApplyDamage(&d)

		nextTick := 60
		if lastTickDuration == 60 {
			nextTick = 135
		}
		c.AddTask(c.skillDmgTickTask(src, nextTick), "qiqi-skill-dmg-tick", nextTick)
	}
}

func (c *char) c1(t core.Target) {
	if c.talismanExpiry[t.Index()] < c.Core.F {
		return
	}
	c.AddEnergy(2)

	c.Core.Log.Debugw("Qiqi C1 Activation - Adding 2 energy", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "target", t.Index())
}

// Handles skill auto healing ticks
func (c *char) skillHealTickTask(src int) func() {
	return func() {
		if c.Core.Status.Duration("qiqiskill") == 0 {
			return
		}

		// TODO: Not sure how this interacts with sac sword... Treat it as only one instance can be up at a time for now
		if c.skillLastUsed > src {
			return
		}

		// Create a temporary snapshot to get the atk mods included
		d := c.Snapshot(
			"Herald of Frost: Continuous Heal Proc",
			core.AttackTagElementalArt,
			core.ICDTagElementalArt,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Cryo,
			25,
			0,
		)
		c.Core.Health.HealActive(c.Index, c.heal(&d, skillHealContPer, skillHealContFlat, c.TalentLvlSkill()))

		// Queue next instance
		c.AddTask(c.skillHealTickTask(src), "qiqi-skill-heal-tick", 4.5*60)
	}
}

// Implements burst and applies Talisman. Main Talisman functions are handled elsewhere
func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	c.QueueDmgDynamic(func() *core.Snapshot {
		d := c.Snapshot(
			"Fortune-Preserving Talisman",
			core.AttackTagElementalBurst,
			core.ICDTagElementalBurst,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Cryo,
			50,
			burstDmg[c.TalentLvlBurst()],
		)
		d.Targets = core.TargetAll
		return &d
	}, f)

	c.SetCD(core.ActionBurst, 20*60)
	c.Energy = 0

	return f, a
}

// Helper function to calculate healing amount from a snapshot instance, which has all mods applied
func (c *char) heal(d *core.Snapshot, healScalePer []float64, healScaleFlat []float64, talentLevel int) float64 {
	atk := d.BaseAtk*(1+d.Stats[core.ATKP]) + d.Stats[core.ATK]
	heal := (healScaleFlat[talentLevel] + atk*healScalePer[talentLevel]) * (1 + d.Stats[core.Heal])
	return heal
}
