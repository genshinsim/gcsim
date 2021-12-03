package zhongli

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) newStele(dur int, max int) {
	//deal damage when created
	d := c.Snapshot(
		"Stone Stele (Initial)",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		50,
		skill[c.TalentLvlSkill()],
	)
	d.FlatDmg = 0.019 * c.HPMax
	d.Targets = core.TargetAll
	// Damage proc is near instant upon creation
	c.QueueDmg(&d, 0)

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
		d := c.Snapshot(
			"Stone Stele (Tick)",
			core.AttackTagElementalArt,
			core.ICDTagElementalArt,
			core.ICDGroupDefault,
			core.StrikeTypeBlunt,
			core.Geo,
			25,
			skillTick[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll
		d.FlatDmg = 0.019 * c.HPMax
		//check how many times to hit
		count := c.Core.Constructs.Count()
		if count > max {
			count = max
		}
		orb := false
		for i := 0; i < count; i++ {
			x := d.Clone()
			c.Core.Combat.ApplyDamage(&x)
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
