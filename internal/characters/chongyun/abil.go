package chongyun

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, f)

	if c.Base.Cons >= 1 && c.NormalCounter == 3 {
		ai := core.AttackInfo{
			Abil:       "Chongyun C1",
			ActorIndex: c.Index,
			AttackTag:  core.AttackTagNone,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Cryo,
			Durability: 25,
			Mult:       .5,
		}
		//3 blades
		for i := 0; i < 3; i++ {
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, f+i*5)
		}
	}
	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Blade: Chonghua's Layered Frost",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, f-1)

	//TODO: energy count; lib says 3:4?
	c.QueueParticle("Chongyun", 4, core.Cryo, 100)

	//a4 delayed damage + cryo resist shred
	c.AddTask(func() {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Spirit Blade: Chonghua's Layered Frost (Ar)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Cryo,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}
		cb := func(t core.Target, ae *core.AttackEvent) {

			t.AddResMod("Chongyun A4", core.ResistMod{
				Duration: 480, //10 seconds
				Ele:      core.Cryo,
				Value:    -0.10,
			})
		}

		//TODO: this needs to be fixed still for sac gs
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, 0, cb)

	}, "Chongyun-Skill", f+600)

	c.Core.Status.AddStatus("chongyunfield", 600)

	//TODO: delay between when frost field start ticking?
	for i := 60; i <= 600; i += 60 {
		c.AddTask(func() {
			active := c.Core.Chars[c.Core.ActiveChar]
			c.infuse(active)
		}, "chongyun-field", i)
	}

	c.SetCD(core.ActionSkill, 900)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Blade: Cloud-Parting Star",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	count := 3
	if c.Base.Cons == 6 {
		count = 4

	}

	for i := 0; i < count; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, f+10*i)
	}

	c.SetCD(core.ActionBurst, 720)
	c.Energy = 0
	return f, a
}
