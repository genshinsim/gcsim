package travelergeo

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
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f-1, f-1)

	if c.NormalCounter == c.NormalHitNum-1 {
		//add 60% as geo dmg
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "A1",
			AttackTag:  core.AttackTagNormal,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Geo,
			Durability: 25,
			Mult:       0.6,
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f-1, f-1)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Starfell Sword",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f+10)

	count := 3
	if c.Core.Rand.Float64() < 0.33 {
		count = 4
	}
	c.QueueParticle(c.Name(), count, core.Geo, f+100)

	c.AddTask(func() {
		dur := 30 * 60
		if c.Base.Cons == 6 {
			dur += 600
		}
		con := &stone{
			src:    c.Core.F,
			expiry: c.Core.F + dur,
			char:   c,
		}
		c.Core.Constructs.New(con, false)
	}, "geomc-construct", f+10)

	c.SetCD(core.ActionSkill, 360)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	hits, ok := p["hits"]
	if !ok {
		hits = 2
	}

	maxConstructCount, ok := p["construct_limit"]
	if !ok {
		maxConstructCount = 1
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wake of Earth",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagTravelerWakeOfEarth,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	//The shockwave triggered by Wake of Earth regenerates 5 Energy for every opponent hit.
	// A maximum of 25 Energy can be regenerated in this manner at any one time.
	var cb core.AttackCBFunc
	src := c.Core.F

	if c.Base.Cons >= 4 {
		cb = func(a core.AttackCB) {
			// TODO: A bit of a cludge to deal with frame 0 casts. Will have to think about this behavior a bit more
			if a.Target.GetTag("traveler-c4-src") == src && src > 0 {
				return
			}
			a.Target.SetTag("traveler-c4-src", src)
			c.AddEnergy("geo-traveler-c4", 5)
		}
	}

	//1.1 sec duration, tick every .25
	for i := 0; i < hits; i++ {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), (i+1)*15, cb)
	}

	c.AddTask(func() {
		dur := 15 * 60
		if c.Base.Cons == 6 {
			dur += 300
		}
		con := &barrier{
			src:    c.Core.F,
			expiry: c.Core.F + dur,
			char:   c,
			count:  maxConstructCount,
		}
		c.Core.Constructs.NewNoLimitCons(con, true)
		if c.Base.Cons >= 1 {
			c.Tags["wall"] = 1
		}
	}, "geomc-wall", f)

	c.ConsumeEnergy(43)
	c.SetCDWithDelay(core.ActionBurst, 900, 43)
	return f, a
}
