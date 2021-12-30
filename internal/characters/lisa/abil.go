package lisa

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagLisaElectro,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, f-1)

	c.AdvanceNormalIndex()

	return f, a
}

const a4tag = "lisa-a4"

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)

	//TODO: assumes this applies every time per
	//[7:53 PM] Hold ₼KLEE like others hold GME: CHarge is pyro every charge
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	done := false
	cb := func(a core.AttackCB) {
		if done {
			return
		}
		count := a.Target.GetTag(a4tag)
		if count < 3 {
			a.Target.SetTag(a4tag, count+1)
		}
		done = true
	}

	count := 0
	var c1cb func(a core.AttackCB)
	if c.Base.Cons > 0 {
		c1cb = func(a core.AttackCB) {
			if count == 5 {
				return
			}
			count++
			c.AddEnergy(2)
		}
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f-1, cb, c1cb)

	return f, a
}

//p = 0 for no hold, p = 1 for hold
func (c *char) Skill(p map[string]int) (int, int) {
	hold := p["hold"]
	if hold == 1 {
		return c.skillHold(p)
	}
	return c.skillPress(p)
}

//TODO: how long do stacks last?
func (c *char) skillPress(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Violet Arc",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagLisaElectro,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	done := false
	cb := func(a core.AttackCB) {
		if done {
			return
		}
		count := a.Target.GetTag(a4tag)
		if count < 3 {
			a.Target.SetTag(a4tag, count+1)
		}
		done = true
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, f-1, cb)

	if c.Core.Rand.Float64() < 0.5 {
		c.QueueParticle("Lisa", 1, core.Electro, f+100)
	}

	c.SetCD(core.ActionSkill, 60)
	return f, a
}

func (c *char) skillHold(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	//no multiplier as that's target dependent
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Violet Arc (Hold)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, f)

	//c2 add defense? no interruptions either way
	if c.Base.Cons >= 2 {
		//increase def for the duration of this abil in however many frames
		val := make([]float64, core.EndStatType)
		val[core.DEFP] = 0.25
		c.AddMod(core.CharStatMod{
			Key:    "lisa-c2",
			Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
			Expiry: c.Core.F + 126,
		})
	}

	//[8:31 PM] ArchedNosi | Lisa Unleashed: yeah 4-5 50/50 with Hold
	//[9:13 PM] ArchedNosi | Lisa Unleashed: @gimmeabreak actually wait, xd i noticed i misread my sheet, Lisa Hold E always gens 5 orbs
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, f)

	// count := 4
	// if c.Core.Rand.Float64() < 0.5 {
	// 	count = 5
	// }
	c.QueueParticle("Lisa", 5, core.Electro, f+100)

	// c.CD[def.SkillCD] = c.Core.F + 960 //16seconds
	c.SetCD(core.ActionSkill, 960)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	//first zap has no icd
	targ := c.Core.RandomTargetIndex(core.TargettableEnemy)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Rose (Initial)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       0.1,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(targ, core.TargettableEnemy), f, f)

	//duration is 15 seconds, tick every .5 sec
	//30 zaps once every 30 frame, starting at f

	ai = core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Rose (Tick)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	for i := 30; i <= 900; i += 30 {

		var cb core.AttackCBFunc
		if c.Base.Cons >= 4 {
			//random 1 to 3 jumps
			count := c.Rand.Intn(3) + 1
			cb = func(a core.AttackCB) {
				if count == 0 {
					return
				}
				//generate additional attack, random target
				//if we get -1 for a target then that just means there's no target
				//to jump to so that's fine; chain will terminate
				count++
				//grab a list of enemies by range; we assume it'll just hit the closest?

			}
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f-1, f+i, cb)
	}

	//add a status for this just in case someone cares
	c.AddTask(func() {
		c.Core.Status.AddStatus("lisaburst", 900)
	}, "lisa burst status", f)

	//on lisa c4
	//[8:11 PM] gimmeabreak: er, what does lisa c4 do?
	//[8:11 PM] ArchedNosi | Lisa Unleashed: allows each pulse of the ult to be 2-4 arcs
	//[8:11 PM] ArchedNosi | Lisa Unleashed: if theres enemies given
	//[8:11 PM] gimmeabreak: oh so it jumps 2 to 4 times?
	//[8:11 PM] gimmeabreak: i guess single target it does nothing then?
	//[8:12 PM] ArchedNosi | Lisa Unleashed: yeah single does nothing

	c.ConsumeEnergy(64)
	// c.CD[def.BurstCD] = c.Core.F + 1200
	c.SetCD(core.ActionBurst, 1200)
	return f, a
}
