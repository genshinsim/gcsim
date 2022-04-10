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
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f)

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
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f+i*5, f+i*5)
		}
	}
	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	//if fieldSrc is < duration then this is prob a sac proc
	//we need to stop the old field from ticking (by changing fieldSrc)
	//and also trigger a4 delayed damage

	src := c.Core.F

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

	ai = core.AttackInfo{
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
	cb := func(a core.AttackCB) {
		a.Target.AddResMod("Chongyun A4", core.ResistMod{
			Duration: 480, //10 seconds
			Ele:      core.Cryo,
			Value:    -0.10,
		})
	}
	snap := c.Snapshot(&ai)

	//if field is overwriting last
	if src-c.fieldSrc < 600 {
		//we're overriding previous field so trigger a4 here
		atk := c.a4Snap
		c.Core.Combat.QueueAttackEvent(atk, 1)
	}
	c.fieldSrc = src
	//override previous snap
	c.a4Snap = &core.AttackEvent{
		Info:     ai,
		Snapshot: snap,
		Pattern:  core.NewDefCircHit(3, false, core.TargettableEnemy),
	}
	c.a4Snap.Callbacks = append(c.a4Snap.Callbacks, cb)

	//a4 delayed damage + cryo resist shred
	c.AddTask(func() {
		//if src changed then that means the field changed already
		if src != c.fieldSrc {
			return
		}
		//TODO: this needs to be fixed still for sac gs
		c.Core.Combat.QueueAttackEvent(c.a4Snap, 0)
	}, "Chongyun-Skill", f+600)

	c.Core.Status.AddStatus("chongyunfield", 600)

	//TODO: delay between when frost field start ticking?
	for i := f - 1; i <= 600; i += 60 {
		c.AddTask(func() {
			if src != c.fieldSrc {
				return
			}
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

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 50, 50)
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 57, 57)
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 65, 65)

	if c.Base.Cons == 6 {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 76, 76)
	}

	c.SetCDWithDelay(core.ActionBurst, 720, 10)
	c.ConsumeEnergy(10)
	return f, a
}
