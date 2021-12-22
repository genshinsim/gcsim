package zhongli

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) newStele(dur int, max int) {
	//deal damage when created
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stone Stele (Initial)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
		FlatDmg:    0.019 * c.HPMax,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 0)

	//create a construct
	con := &stoneStele{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		c:      c,
	}

	num := c.Core.Constructs.CountByType(core.GeoConstructZhongliSkill)

	c.Core.Constructs.New(con, num == c.maxStele)

	c.steleCount = c.Core.Constructs.CountByType(core.GeoConstructZhongliSkill)

	c.Core.Log.Debugw(
		"Stele added",
		"frame", c.Core.F,
		"event", core.LogCharacterEvent,
		"char", c.Index,
		"orig_count", num,
		"cur_count", c.steleCount,
		"max_hit", max,
		"next_tick", c.Core.F+120,
	)

	c.AddTask(c.resonance(c.Core.F, max), "stele", 120)
}

func (c *char) resonance(src, max int) func() {
	return func() {
		c.Core.Log.Debugw("Stele checking for tick", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src, "char", c.Index)
		if !c.Core.Constructs.Has(src) {
			return
		}
		c.Core.Log.Debugw("Stele ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "next expected", c.Core.F+120, "src", src, "char", c.Index)
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Stone Stele (Tick)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Geo,
			Durability: 25,
			Mult:       skillTick[c.TalentLvlSkill()],
			FlatDmg:    0.019 * c.HPMax,
		}

		//check how many times to hit
		count := c.Core.Constructs.Count()
		if count > max {
			count = max
		}
		orb := false
		for i := 0; i < count; i++ {
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 0)
			if c.energyICD < c.Core.F && !orb && c.Core.Rand.Float64() < .5 {
				orb = true
			}
		}
		if orb {
			c.energyICD = c.Core.F + 90
			c.QueueParticle("zhongli", 1, core.Geo, 120)
		}
		c.AddTask(c.resonance(src, max), "stele", 120)
	}
}
