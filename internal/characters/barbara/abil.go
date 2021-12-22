package barbara

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	done := false
	// Taken from Noelle code
	cb := func(a core.AttackCB) {
		if done { //why do we need this @srl
			return
		}
		//check for healing
		if c.Core.Status.Duration("barbskill") > 0 {
			//heal target
			heal := (prochpp[c.TalentLvlSkill()] + prochp[c.TalentLvlSkill()])
			c.Core.Health.HealAll(c.Index, heal)
			done = true
		}

	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f, cb)
	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Charge attack function - handles seal use
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       charge[c.NormalCounter][c.TalentLvlAttack()],
	}

	done := false
	// Taken from Noelle code
	cb := func(a core.AttackCB) {
		if done { //why do we need this @srl
			return
		}
		//check for healing
		if c.Core.Status.Duration("barbskill") > 0 {
			//heal target
			heal := (prochpp[c.TalentLvlSkill()] + prochp[c.TalentLvlSkill()])
			c.Core.Health.HealAll(c.Index, 4*heal)
			done = true
		}

	}
	var cbenergy func(a core.AttackCB) = nil
	energyCount := 0
	if c.Base.Cons >= 4 {
		cbenergy = func(a core.AttackCB) {
			//check for healing
			if c.Core.Status.Duration("barbskill") > 0 && energyCount < 5 {
				//regen energy
				c.AddEnergy(1)
				energyCount++
			}

		}
	}

	// TODO: Not sure of snapshot timing
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f, cb, cbenergy)

	return f, a
}

// barbara skill - copied from bennett burst

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	//add field effect timer
	//assumes a4
	c.Core.Status.AddStatus("barbskill", 15*60)
	//hook for buffs; active right away after cast

	c.stacks = 0
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Let the Show Begin♪",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25, //TODO: what is 1A GU?
		Mult:       skill[c.TalentLvlSkill()],
	}
	//TODO: review barbara AOE size?
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 5, 5)

	stats := c.SnapshotStats("Let the Show Begin♪ (Heal)", core.AttackTagNone)

	//apply right away
	c.applyBarbaraField(stats)()

	c.onSkillStackCount(stats)
	//add 1 tick each 5s
	//first tick starts at 0
	for i := 0; i <= 1200; i += 5 * 60 {
		c.AddTask(c.applyBarbaraField(stats), "barbara-field", i)
	}

	c.Energy = 0
	if c.Base.Cons >= 2 {
		c.SetCD(core.ActionSkill, 32*60*0.85)
	} else {
		c.SetCD(core.ActionSkill, 32*60)
	}
	return f, a //todo fix field cast time
}

func (c *char) applyBarbaraField(stats [core.EndStatType]float64) func() {
	hpplus := stats[core.Heal]
	heal := (skillhp[c.TalentLvlBurst()] + skillhpp[c.TalentLvlBurst()]*c.MaxHP()) * (1 + hpplus)
	val := make([]float64, core.EndStatType)
	val[core.HydroP] = 0.0
	if c.Base.Cons >= 2 {
		val[core.HydroP] += 0.2
	}
	return func() {
		c.Core.Log.Debugw("barbara field ticking", "frame", c.Core.F, "event", core.LogCharacterEvent)

		active := c.Core.Chars[c.Core.ActiveChar]
		c.Core.Health.HealActive(c.Index, heal)

		active.AddMod(core.CharStatMod{
			Key: "barbara-field",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return val, true // @srl what does this boolean mean?
			},
			Expiry: c.Core.F + 5*60, // this is for each application of the field.. is this correct @srliao
		})
		// Additional per-character status for config conditionals
		c.Core.Status.AddStatus(fmt.Sprintf("barbarabuff%v", active.Name()), 5*60)
		// missing wet self-reaction
	}
}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)
	//hook for buffs; active right away after cast

	stats := c.SnapshotStats("Shining Miracle♪ (Heal)", core.AttackTagNone)

	hpplus := stats[core.Heal]
	heal := (bursthp[c.TalentLvlBurst()] + bursthpp[c.TalentLvlBurst()]*c.MaxHP()) * (1 + hpplus)
	c.Core.Health.HealAll(c.Index, heal)

	c.Energy = 0
	c.SetCD(core.ActionBurst, 20*60)
	return f, a //todo fix field cast time
}

//inspired from raiden
func (c *char) onSkillStackCount(stats [core.EndStatType]float64) {
	particleStack := 0
	done := false // this is some copium implementation for adding a tick but will think about improving it later
	c.Core.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		if c.Core.Status.Duration("barbskill") > 0 && particleStack < 5 {
			particleStack++
			// add a second per particle
			c.Core.Status.ExtendStatus("barbskill", 60)
			return false
		}
		if particleStack == 5 && !done {
			c.AddTask(c.applyBarbaraField(stats), "barbara-field", c.Core.Status.Duration("barbskill")+c.Core.F)
			done = true
		}
		//this doesn't do anything yet
		return false
	}, "barbara-skill-extend")
}
