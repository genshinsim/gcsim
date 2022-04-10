package qiqi

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack - nothing special
func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
	}
	for _, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f, f)
	}

	c.AdvanceNormalIndex()

	return f, a
}

// Standard charge attack
func (c *char) Charge(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
	}
	for _, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f, f)
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
	// Both healing and damage are snapshot
	c.AddTask(func() {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Herald of Frost: Initial Damage",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    core.Cryo,
			Durability: 25,
			Mult:       skillInitialDmg[c.TalentLvlSkill()],
		}
		snap := c.Snapshot(&ai)

		// One healing proc happens immediately on cast
		c.Core.Health.Heal(core.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.ActiveChar,
			Message: "Herald of Frost (Tick)",
			Src:     c.healSnapshot(&snap, skillHealContPer, skillHealContFlat, c.TalentLvlSkill()),
			Bonus:   snap.Stats[core.Heal],
		})

		// Healing and damage instances are snapshot
		// Separately cloned snapshots are fed into each function to ensure nothing interferes with each other

		// Queue up continuous healing instances
		// No exact frame data on when the healing ticks happen. Just roughly guessing here
		// Healing ticks happen 3 additional times during the skill - assume ticks are roughly 4.5s apart
		// so in sec (0 = skill cast), 1, 5.5, 10, 14.5
		c.skillHealSnapshot = snap
		c.AddTask(c.skillHealTickTask(src), "qiqi-skill-heal-tick", 4.5*60)

		// Queue up damage swipe instances.
		// No exact frame data on when the damage ticks happen. Just roughly guessing here
		// Occurs 9 times over the course of the skill
		// Once shortly after initial cast, then 8 additional procs over the rest of the duration
		// Each proc occurs in "pairs" of two swipes each spaced around 2.25s apart
		// The time between each swipe in a pair is about 1s
		// No exact frame data available plus the skill duration is affected by hitlag
		// Damage procs occur (in sec 0 = skill cast): 1.5, 3.75, 4.75, 7, 8, 10.25, 11.25, 13.5, 14.5

		aiTick := ai
		aiTick.Abil = "Herald of Frost: Skill Damage"
		aiTick.Mult = skillDmgCont[c.TalentLvlSkill()]
		snapTick := c.Snapshot(&aiTick)
		tickAE := &core.AttackEvent{
			Info:        aiTick,
			Snapshot:    snapTick,
			Pattern:     core.NewDefCircHit(2, false, core.TargettableEnemy),
			SourceFrame: c.Core.F,
		}

		c.AddTask(c.skillDmgTickTask(src, tickAE, 60), "qiqi-skill-dmg-tick", 30)

		// Apply damage needs to take place after above takes place to ensure stats are handled correctly
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), 0)
	}, "qiqi-skill-activation", f)

	c.SetCD(core.ActionSkill, 30*60)

	return f, a
}

// Handles skill damage swipe instances
// Also handles C1:
// When the Herald of Frost hits an opponent marked by a Fortune-Preserving Talisman, Qiqi regenerates 2 Energy.
func (c *char) skillDmgTickTask(src int, ae *core.AttackEvent, lastTickDuration int) func() {
	return func() {
		if c.Core.Status.Duration("qiqiskill") == 0 {
			return
		}

		// TODO: Not sure how this interacts with sac sword... Treat it as only one instance can be up at a time for now
		if c.skillLastUsed > src {
			return
		}

		// Clones initial snapshot
		tick := *ae //deference the pointer here

		if c.Base.Cons >= 1 {
			tick.Callbacks = append(tick.Callbacks, c.c1)
		}

		c.Core.Combat.QueueAttackEvent(&tick, 0)

		nextTick := 60
		if lastTickDuration == 60 {
			nextTick = 135
		}
		c.AddTask(c.skillDmgTickTask(src, ae, nextTick), "qiqi-skill-dmg-tick", nextTick)
	}
}

func (c *char) c1(a core.AttackCB) {
	if a.Target.GetTag(talismanKey) < c.Core.F {
		return
	}
	c.AddEnergy("qiqi-c1", 2)

	c.Core.Log.NewEvent("Qiqi C1 Activation - Adding 2 energy", core.LogCharacterEvent, c.Index, "target", a.Target.Index())
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

		c.Core.Health.Heal(core.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.ActiveChar,
			Message: "Herald of Frost (Tick)",
			Src:     c.healSnapshot(&c.skillHealSnapshot, skillHealContPer, skillHealContFlat, c.TalentLvlSkill()),
			Bonus:   c.skillHealSnapshot.Stats[core.Heal],
		})

		// Queue next instance
		c.AddTask(c.skillHealTickTask(src), "qiqi-skill-heal-tick", 4.5*60)
	}
}

// Only applies burst damage. Main Talisman functions are handled in qiqi.go
func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fortune-Preserving Talisman",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f, f)
	c.SetCDWithDelay(core.ActionBurst, 20*60, 10)
	c.ConsumeEnergy(10)

	return f, a
}

// Helper function to calculate healing amount dynamically using current character stats, which has all mods applied
func (c *char) healDynamic(healScalePer []float64, healScaleFlat []float64, talentLevel int) float64 {
	atk := c.Base.Atk + c.Weapon.Atk*(1+c.Stat(core.ATKP)) + c.Stat(core.ATK)
	heal := healScaleFlat[talentLevel] + atk*healScalePer[talentLevel]
	return heal
}

// Helper function to calculate healing amount from a snapshot instance
func (c *char) healSnapshot(d *core.Snapshot, healScalePer []float64, healScaleFlat []float64, talentLevel int) float64 {
	atk := d.BaseAtk*(1+d.Stats[core.ATKP]) + d.Stats[core.ATK]
	heal := healScaleFlat[talentLevel] + atk*healScalePer[talentLevel]
	return heal
}
