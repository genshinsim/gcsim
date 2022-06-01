package lisa

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = []int{26, 18, 17, 31}

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

	//todo: Does it really snapshot immediately?
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, hitmarks[c.NormalCounter])

	c.AdvanceNormalIndex()

	return f, a
}

const conductiveTag = "lisa-conductive-stacks"

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
		count := a.Target.GetTag(conductiveTag)
		if count < 3 {
			a.Target.SetTag(conductiveTag, count+1)
		}
		done = true
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f, cb)

	return f, a
}

var skillHitmarks = []int{22, 117}

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
		count := a.Target.GetTag(conductiveTag)
		if count < 3 {
			a.Target.SetTag(conductiveTag, count+1)
		}
		done = true
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, skillHitmarks[0], cb)

	if c.Core.Rand.Float64() < 0.5 {
		c.QueueParticle("Lisa", 1, core.Electro, f+100)
	}

	c.SetCDWithDelay(core.ActionSkill, 60, 17)
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
		Durability: 50,
	}

	//c2 add defense? no interruptions either way
	if c.Base.Cons >= 2 {
		//increase def for the duration of this abil in however many frames
		val := make([]float64, core.EndStatType)
		val[core.DEFP] = 0.25
		c.AddMod(core.CharStatMod{
			Key:    "lisa-c2",
			Amount: func() ([]float64, bool) { return val, true },
			Expiry: c.Core.F + 126,
		})
	}

	count := 0
	var c1cb func(a core.AttackCB)
	if c.Base.Cons > 0 {
		c1cb = func(a core.AttackCB) {
			if count == 5 {
				return
			}
			count++
			c.AddEnergy("lisa-c1", 2)
		}
	}

	//[8:31 PM] ArchedNosi | Lisa Unleashed: yeah 4-5 50/50 with Hold
	//[9:13 PM] ArchedNosi | Lisa Unleashed: @gimmeabreak actually wait, xd i noticed i misread my sheet, Lisa Hold E always gens 5 orbs
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, skillHitmarks[1], c1cb)

	// count := 4
	// if c.Core.Rand.Float64() < 0.5 {
	// 	count = 5
	// }
	c.QueueParticle("Lisa", 5, core.Electro, f+100)

	// c.CD[def.SkillCD] = c.Core.F + 960 //16seconds, starts after 114 frames
	c.SetCDWithDelay(core.ActionSkill, 960, 114)
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
		Durability: 0,
		Mult:       0.1,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(targ, core.TargettableEnemy), f, f, a4cb)

	//duration is 15 seconds, tick every .5 sec
	//30 zaps once every 30 frame, starting at 119

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

	for i := 119; i <= 119+900; i += 30 { //first tick at 119

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
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f-1, i, cb, a4cb)
	}

	//add a status for this just in case someone cares
	c.AddTask(func() {
		c.Core.Status.AddStatus("lisaburst", 119+900)
	}, "lisa burst status", f)

	//on lisa c4
	//[8:11 PM] gimmeabreak: er, what does lisa c4 do?
	//[8:11 PM] ArchedNosi | Lisa Unleashed: allows each pulse of the ult to be 2-4 arcs
	//[8:11 PM] ArchedNosi | Lisa Unleashed: if theres enemies given
	//[8:11 PM] gimmeabreak: oh so it jumps 2 to 4 times?
	//[8:11 PM] gimmeabreak: i guess single target it does nothing then?
	//[8:12 PM] ArchedNosi | Lisa Unleashed: yeah single does nothing

	//burst cd starts 53 frames after executed
	//energy usually consumed after 63 frames
	c.ConsumeEnergy(63)
	// c.CD[def.BurstCD] = c.Core.F + 1200
	c.SetCDWithDelay(core.ActionBurst, 1200, 53)
	return f, a
}

func a4cb(a core.AttackCB) {
	a.Target.AddDefMod("lisa-a4", -0.15, 600)
}
