package rosaria

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Normal",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f-5+i)
	}

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupPole,
		Element:    core.Physical,
		Durability: 25,
		Mult:       nc[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f-1)

	//return animation cd
	return f, a
}

// Skill attack damage queue generator
// Includes optional argument "nobehind" for whether Rosaria appears behind her opponent or not (for her A1).
// Default behavior is to appear behind enemy - set "nobehind=1" to diasble A1 proc
func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	// No ICD to the 2 hits
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ravaging Confession (Hit 1)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       skill[0][c.TalentLvlSkill()],
	}
	// First hit comes out 20 frames before second
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-20, f-20)

	// A1 activation
	// When Rosaria strikes an opponent from behind using Ravaging Confession, Rosaria's CRIT RATE increases by 12% for 5s.
	// We always assume that it procs on hit 1 to simplify
	if p["nobehind"] != 1 {
		val := make([]float64, core.EndStatType)
		val[core.CR] = 0.12
		c.AddMod(core.CharStatMod{
			Key:    "rosaria-a1",
			Expiry: c.Core.F + 300,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return val, true
			},
		})
		c.Core.Log.Debugw("Rosaria A1 activation", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "ends_on", c.Core.F+300)
	}

	// Rosaria E is dynamic, so requires a second snapshot
	//TODO: check snapshot timing here
	ai = core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ravaging Confession (Hit 2)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       skill[1][c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-1, f-1)

	// Particles are emitted after the second hit lands
	c.QueueParticle("rosaria", 3, core.Cryo, f+100)

	c.SetCD(core.ActionSkill, 360)

	return f, a
}

// Burst attack damage queue generator
// Rosaria swings her weapon to slash surrounding opponents, then she summons a frigid Ice Lance that strikes the ground. Both actions deal Cryo DMG.
// While active, the Ice Lance periodically releases a blast of cold air, dealing Cryo DMG to surrounding opponents.
// Also includes the following effects: A4, C6
func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	// Note - if a more advanced targeting system is added in the future
	// hit 1 is technically only on surrounding enemies, hits 2 and dot are on the lance
	// For now assume that everything hits all targets
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rites of Termination (Hit 1)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burst[0][c.TalentLvlBurst()],
	}
	x, y := c.Core.Targets[0].Shape().Pos()
	var cb core.AttackCBFunc
	if c.Base.Cons == 6 {
		cb = c6cb
	}
	// Hit 1 comes out on frame 10
	// 2nd hit comes after lance drop animation finishes
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 1, false, core.TargettableEnemy), 10, 10, cb)

	ai.Abil = "Rites of Termination (Hit 2)"
	ai.Mult = burst[1][c.TalentLvlBurst()]
	// Note old code set the hit 10 frames before the recorded one - not sure why
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(0, 0, 2, false, core.TargettableEnemy), f-10, f-10, cb)

	//duration is 8 second (extended by c2 by 4s), + 0.5
	dur := 510
	if c.Base.Cons >= 2 {
		dur += 240
	}

	// Burst is snapshot when the lance lands (when the 2nd damage proc hits)
	ai = core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rites of Termination (DoT)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burst[0][c.TalentLvlBurst()],
	}

	c.AddTask(func() {
		// dot every 2 second after lance lands
		for i := 120; i < dur; i += 120 {
			c.Core.Combat.QueueAttack(ai, core.NewCircleHit(0, 0, 2, false, core.TargettableEnemy), 0, i+10, cb)
		}
	}, "rosaria-snapshot", f-10)

	c.Core.Status.AddStatus("rosariaburst", dur)

	// Handle A4
	// Casting Rites of Termination increases CRIT RATE of all nearby party members, excluding Rosaria herself, by 15% of Rosaria's CRIT RATE for 10s. CRIT RATE bonus gained this way cannot exceed 15%.
	// Uses the snapshot generated by hit #1 to ensure all mods are accounted for.
	// Confirmed via testing that mods like Rosaria A1 are accounted for, and Blizzard Strayer modifications are not
	crit_share := 0.15 * c.Stats[core.CR]
	if crit_share > 0.15 {
		crit_share = 0.15
	}
	val := make([]float64, core.EndStatType)
	val[core.CR] = crit_share
	c.Core.Log.Debugw("Rosaria A4 activation", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "ends_on", c.Core.F+600, "crit_share", crit_share)

	for i, char := range c.Core.Chars {
		// skip Rosaria
		if i == c.Index {
			continue
		}
		char.AddMod(core.CharStatMod{
			Key:    "rosaria-a4",
			Expiry: c.Core.F + 600,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return val, true
			},
		})
	}

	c.SetCD(core.ActionBurst, 15*60)
	c.ConsumeEnergy(12)

	return f, a
}

// Applies C6 effect to enemies hit by it
// Rites of Termination's attack decreases opponent's Physical RES by 20% for 10s.
// Takes in a snapshot definition, and returns the same snapshot with an on hit callback added to apply the debuff
func c6cb(a core.AttackCB) {
	a.Target.AddResMod("rosaria-c6", core.ResistMod{
		Ele:      core.Physical,
		Value:    -0.2,
		Duration: 600,
	})
}

// func (c *char) applyC6(snap *core.Snapshot) {
// 	if c.Base.Cons == 6 {
// 		// Functions similarly to Guoba
// 		snap.OnHitCallback = func(t core.Target) {
// 			t.AddResMod("rosaria-c6", core.ResistMod{
// 				Ele:      core.Physical,
// 				Value:    -0.2,
// 				Duration: 600,
// 			})
// 		}
// 	}
// }
