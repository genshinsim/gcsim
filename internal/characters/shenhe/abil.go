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
	// First hit comes out 20 frames before second
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-20, f-20)

	// Particles are emitted after the second hit lands
	c.QueueParticle("shenhe", 3, core.Cryo, f+100)

	c.skillPressBuff()
	c.SetCD(core.ActionSkill, 10*60)

	return f, a
}

func (c *char) skillHold(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	// No ICD to the 2 hits
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spring Spirit Summoning (Hold)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}
	// First hit comes out 20 frames before second
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-20, f-20)

	// Particles are emitted after the second hit lands
	c.QueueParticle("shenhe", 4, core.Cryo, f+100)

	c.skillHoldBuff()
	c.SetCD(core.ActionSkill, 15*60)

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
	var cb core.AttackCBFunc
	// Hit 1 comes out on frame 10
	// 2nd hit comes after lance drop animation finishes
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 1, false, core.TargettableEnemy), 10, 10, cb)

	//duration is 12 second (extended by c2 by 6s)
	dur := 12 * 60
	if c.Base.Cons >= 2 {
		dur += 6 * 60
	}

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
		// dot every 2 second after lance lands
		for i := 120; i < dur; i += 120 {
			c.Core.Combat.QueueAttack(ai, core.NewCircleHit(0, 0, 2, false, core.TargettableEnemy), 0, i+10, cb)
		}
	}, "shenhe-snapshot", f-10)

	c.Core.Status.AddStatus("shenheburst", dur)

	c.SetCD(core.ActionBurst, 20*60)
	c.ConsumeEnergy(12)

	return f, a
}

func (c *char) skillPressBuff() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15
	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue
		}
		char.AddMod(core.CharStatMod{
			Key:    "shenhe-a2",
			Expiry: c.Core.F + 10*60,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if a != core.AttackTagElementalBurst && a != core.AttackTagElementalArt && a != core.AttackTagElementalArtHold {
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
		if i == c.Index {
			continue
		}
		char.AddMod(core.CharStatMod{
			Key:    "shenhe-a2",
			Expiry: c.Core.F + 15*60,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if a != core.AttackTagNormal && a != core.AttackTagExtra && a != core.AttackTagPlunge {
					return nil, false
				}
				return val, true
			},
		})
	}
}
