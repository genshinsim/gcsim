package barbara

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f+travel)

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
	// TODO: Not sure of snapshot timing
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f)

	return f, a
}

// barbara skill - copied from bennett burst

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	//add field effect timer
	//assumes a4
	c.Core.Status.AddStatus("barbskill", 20)
	//hook for buffs; active right away after cast

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

	//add 1 tick each 5s
	//first tick starts at 0
	for i := 0; i <= 1200; i += 300 {
		c.AddTask(c.applyBarbaraField(stats), "barbara-field", i)
	}

	c.Energy = 0
	c.SetCD(core.ActionSkill, 32*60)
	return f, a //todo fix field cast time
}

func (c *char) applyBarbaraField(stats [core.EndStatType]float64) func() {
	hpplus := stats[core.Heal]
	heal := (skillhp[c.TalentLvlBurst()] + skillhpp[c.TalentLvlBurst()]*c.MaxHP()) * (1 + hpplus)
	var val [core.EndStatType]float64
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
			Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
				return val, true
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
