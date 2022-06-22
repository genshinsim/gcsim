package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) Burst(p map[string]int) action.ActionInfo {

	f, a := c.ActionFrames(action.ActionBurst, p)

	//first zap has no icd
	targ := c.Core.RandomTargetIndex(combat.TargettableEnemy)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Rose (Initial)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 0,
		Mult:       0.1,
	}
	c.Core.Combat.QueueAttack(ai, combat.NewDefSingleTarget(targ, combat.TargettableEnemy), f, f, a4cb)

	//duration is 15 seconds, tick every .5 sec
	//30 zaps once every 30 frame, starting at 119

	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Rose (Tick)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	for i := 119; i <= 119+900; i += 30 { //first tick at 119

		var cb core.AttackCBFunc
		if c.Base.Cons >= 4 {
			//random 1 to 3 jumps
			count := c.Rand.Intn(3) + 1
			cb = func(a combat.AttackCB) {
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
		c.Core.Combat.QueueAttack(ai, combat.NewDefSingleTarget(c.Core.RandomEnemyTarget(), combat.TargettableEnemy), f-1, i, cb, a4cb)
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
	c.SetCDWithDelay(action.ActionBurst, 1200, 53)
	return f, a
}

func a4cb(a combat.AttackCB) {
	a.Target.AddDefMod("lisa-a4", -0.15, 600)
}
