package shenhe

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
		Mult:       charged[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f-1)

	//return animation cd
	return f, a
}

// Skill attack damage queue generator
// Includes optional argument "nobehind" for whether Rosaria appears behind her opponent or not (for her A1).
// Default behavior is to appear behind enemy - set "nobehind=1" to diasble A1 proc
func (c *char) Skill(p map[string]int) (int, int) {
	hold := p["hold"]
	if hold == 1 {
		return c.skillHold(p)
	}
	return c.skillPress(p)
}

func (c *char) skillPress(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	// No ICD to the 2 hits
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spring Spirit Summoning (Press)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	// c.AddTask(c.skillPressBuff, "shenhe (press) quill start", f-1)
	c.skillPressBuff()
	c.Core.Status.AddStatus(quillKey, 10*60)
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)

	// Particles are emitted after the second hit lands
	c.QueueParticle("shenhe", 3, core.Cryo, f+20)

	if c.eChargeMax == 1 {
		c.eNextRecover = 15 * 60
	}
	// Handle E charges
	if c.Tags["eCharge"] == 1 {
		c.SetCD(core.ActionSkill, c.eNextRecover)
	} else {
		c.eNextRecover = c.Core.F + 10*60
		c.Core.Log.Debugw("shenhe e (press) charge used, queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "recover at", c.eNextRecover)
		c.AddTask(c.recoverCharge(c.Core.F, 10*60), "charge", 10*60)
		c.eTickSrc = c.Core.F
	}
	c.Tags["eCharge"]--
	return f, a
}

func (c *char) skillHold(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	// No ICD to the 2 hits
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spring Spirit Summoning (Hold)",
		AttackTag:  core.AttackTagElementalArtHold,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	// First hit comes out 20 frames before second
	// c.AddTask(c.skillHoldBuff, "shenhe (hold) quill start", f-1)
	c.skillHoldBuff()
	c.Core.Status.AddStatus(quillKey, 15*60)
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), f, f)

	// Particles are emitted after the second hit lands
	c.QueueParticle("shenhe", 4, core.Cryo, f+40)

	// Handle E charges

	if c.eChargeMax == 1 {
		c.eNextRecover = 15 * 60
	}
	if c.Tags["eCharge"] == 1 {
		c.SetCD(core.ActionSkill, c.eNextRecover)
	} else {
		c.eNextRecover = c.Core.F + 15*60
		c.Core.Log.Debugw("shenhe e (hold) charge used, queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "recover at", c.eNextRecover)
		c.AddTask(c.recoverCharge(c.Core.F, 15*60), "charge", 15*60)
		c.eTickSrc = c.Core.F
	}
	c.Tags["eCharge"]--
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
		Abil:       "Divine Maiden's Deliverance (Hit 1)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	x, y := c.Core.Targets[0].Shape().Pos()

	//duration is 12 second (extended by c2 by 6s)
	dur := 12 * 60
	count := 5
	if c.Base.Cons >= 2 {
		dur += 6 * 60
		count = 6
	}
	// Hit 1 comes out on frame 10
	// 2nd hit comes after lance drop animation finishes
	cb := func(a core.AttackCB) {
		a.Target.AddResMod("Shenhe Burst Shred (Cryo)", core.ResistMod{
			Duration: dur,
			Ele:      core.Cryo,
			Value:    -burstrespp[c.TalentLvlBurst()],
		})
	}
	cb2 := func(a core.AttackCB) {
		a.Target.AddResMod("Shenhe Burst Shred (Phys)", core.ResistMod{
			Duration: dur,
			Ele:      core.Physical,
			Value:    -burstrespp[c.TalentLvlBurst()],
		})
	}
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 2, false, core.TargettableEnemy), 0, 15, cb, cb2)

	// Burst is snapshot when the lance lands (when the 2nd damage proc hits)
	ai = core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Divine Maiden's Deliverance (DoT)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burstdot[c.TalentLvlBurst()],
	}

	c.AddTask(func() {
		snap := c.Snapshot(&ai)
		//TODO: check this accuracy? Siri's sheet has 137 per
		// dot every 2 second, double tick shortly after another
		for i := 0; i < count; i++ {
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewCircleHit(0, 0, 5, false, core.TargettableEnemy), i*120+50)
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewCircleHit(0, 0, 5, false, core.TargettableEnemy), i*120+80)
		}
	}, "shenhe-snapshot", f)

	c.Core.Status.AddStatus("shenheburst", dur)

	c.SetCD(core.ActionBurst, 20*60)
	c.ConsumeEnergy(1)

	return f, a
}

func (c *char) skillPressBuff() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15
	for i, char := range c.Core.Chars {
		c.quillcount[i] = 5
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "shenhe-a2-press",
			Expiry: c.Core.F + 10*60,
			Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
				if a.Info.AttackTag != core.AttackTagElementalBurst && a.Info.AttackTag != core.AttackTagElementalArt && a.Info.AttackTag != core.AttackTagElementalArtHold {
					return nil, false
				}
				return val, true
			},
		})
	}
}

func (c *char) skillHoldBuff() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15
	for i, char := range c.Core.Chars {
		c.quillcount[i] = 7
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "shenhe-a2-hold",
			Expiry: c.Core.F + 15*60,
			Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
				if a.Info.AttackTag != core.AttackTagNormal && a.Info.AttackTag != core.AttackTagExtra && a.Info.AttackTag != core.AttackTagPlunge {
					return nil, false
				}
				return val, true
			},
		})
	}
}

func (c *char) quillDamageMod() {

	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		consumeStack := true
		if atk.Info.Element != core.Cryo {
			return false
		}

		switch atk.Info.AttackTag {
		case core.AttackTagElementalBurst:
		case core.AttackTagElementalArt:
		case core.AttackTagElementalArtHold:
		case core.AttackTagNormal:
			consumeStack = c.Base.Cons < 6
		case core.AttackTagExtra:
			consumeStack = c.Base.Cons < 6
		case core.AttackTagPlunge:
		default:
			return false
		}

		if c.Core.Status.Duration(quillKey) == 0 {
			return false
		}

		if c.quillcount[atk.Info.ActorIndex] > 0 {
			stats := c.SnapshotStats("Quills", core.AttackTagNone)
			amt := skillpp[c.TalentLvlSkill()] * ((c.Base.Atk+c.Weapon.Atk)*(1+stats[core.ATKP]) + stats[core.ATK])
			if consumeStack { //c6
				c.quillcount[atk.Info.ActorIndex]--
			}
			c.Core.Log.Debugw(
				"Shenhe Quill proc dmg add",
				"frame", c.Core.F,
				"event", core.LogPreDamageMod,
				"char", atk.Info.ActorIndex,
				"before", atk.Info.FlatDmg,
				"addition", amt,
				"effect_ends_at", c.Core.Status.Duration(quillKey),
				"quills left", c.quillcount[atk.Info.ActorIndex],
			)
			atk.Info.FlatDmg += amt
			if c.Base.Cons >= 4 {
				if c.c4count < 50 {
					c.c4count++
				}
				c.c4expiry = c.Core.F + 60*60
			}
		}

		return false //not sure of correctness here
	}, "shenhe-quill")
}

// Helper function that queues up shenhe e charge recovery - similar to other charge recovery functions
func (c *char) recoverCharge(src int, cd int) func() {
	return func() {
		// Required stopper for recursion
		if c.eTickSrc != src {
			c.Core.Log.Debugw("shenhe e recovery function ignored, src diff", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "src", src, "new src", c.eTickSrc)
			return
		}
		c.Tags["eCharge"]++
		c.Core.Log.Debugw("shenhe e recovering a charge", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "skill last used at", src, "total charges", c.Tags["eCharge"])
		c.SetCD(core.ActionSkill, 0)
		if c.Tags["eCharge"] >= c.eChargeMax {
			return
		}

		c.eNextRecover = c.Core.F + cd
		c.Core.Log.Debugw("shenhe e charge queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "recover at", c.eNextRecover)
		c.AddTask(c.recoverCharge(src, cd), "charge", cd)
	}
}
