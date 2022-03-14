package zhongli

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
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
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 0, 0)

	//create a construct
	con := &stoneStele{
		src:    c.Core.Frame,
		expiry: c.Core.Frame + dur,
		c:      c,
	}

	num := c.Core.Constructs.CountByType(core.GeoConstructZhongliSkill)

	c.Core.Constructs.New(con, num == c.maxStele)

	c.steleCount = c.Core.Constructs.CountByType(core.GeoConstructZhongliSkill)

	c.coretype.Log.NewEvent(
		"Stele added",
		coretype.LogCharacterEvent,
		c.Index,
		"orig_count", num,
		"cur_count", c.steleCount,
		"max_hit", max,
		"next_tick", c.Core.Frame+120,
	)
	// Snapshot buffs for resonance ticks
	aiSnap := core.AttackInfo{
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
	snap := c.Snapshot(&aiSnap)
	c.steleSnapshot = coretype.AttackEvent{
		Info:        aiSnap,
		Snapshot:    snap,
		Pattern:     core.NewDefCircHit(1, false, coretype.TargettableEnemy),
		SourceFrame: c.Core.Frame,
	}

	c.AddTask(c.resonance(c.Core.Frame, max), "stele", 120)
}

func (c *char) resonance(src, max int) func() {
	return func() {
		c.coretype.Log.NewEvent("Stele checking for tick", coretype.LogCharacterEvent, c.Index, "src", src, "char", c.Index)
		if !c.Core.Constructs.Has(src) {
			return
		}
		c.coretype.Log.NewEvent("Stele ticked", coretype.LogCharacterEvent, c.Index, "next expected", c.Core.Frame+120, "src", src, "char", c.Index)

		// Use snapshot for damage
		ae := c.steleSnapshot

		//check how many times to hit
		count := c.Core.Constructs.Count()
		if count > max {
			count = max
		}
		orb := false
		for i := 0; i < count; i++ {
			c.Core.Combat.QueueAttackEvent(&ae, 0)
			if c.energyICD < c.Core.Frame && !orb && c.Core.Rand.Float64() < .5 {
				orb = true
			}
		}
		if orb {
			c.energyICD = c.Core.Frame + 90
			c.QueueParticle("zhongli", 1, core.Geo, 120)
		}
		c.AddTask(c.resonance(src, max), "stele", 120)
	}
}
