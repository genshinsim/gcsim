package zhongli

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) newStele(dur int, max int) {
	//deal damage when created
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stone Stele (Initial)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
		FlatDmg:    0.019 * c.MaxHP(),
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 3}, 2), 0, 0)

	//create a construct
	con := &stoneStele{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		c:      c,
	}

	num := c.Core.Constructs.CountByType(construct.GeoConstructZhongliSkill)

	c.Core.Constructs.New(con, num == c.maxStele)

	c.steleCount = c.Core.Constructs.CountByType(construct.GeoConstructZhongliSkill)

	c.Core.Log.NewEvent(
		"Stele added",
		glog.LogCharacterEvent,
		c.Index,
	).
		Write("orig_count", num).
		Write("cur_count", c.steleCount).
		Write("max_hit", max).
		Write("next_tick", c.Core.F+120)

	// Snapshot buffs for resonance ticks
	aiSnap := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stone Stele (Tick)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       skillTick[c.TalentLvlSkill()],
		FlatDmg:    0.019 * c.MaxHP(),
	}
	snap := c.Snapshot(&aiSnap)
	// stele spawns with an offset of Y: 3 but the box has a Y: -4 offset, so it's Y: -1 relative to player as a result
	c.steleSnapshot = combat.AttackEvent{
		Info:        aiSnap,
		Snapshot:    snap,
		Pattern:     combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -1}, 8, 8),
		SourceFrame: c.Core.F,
	}

	c.Core.Tasks.Add(c.resonance(c.Core.F, max), 120)
}

func (c *char) resonance(src, max int) func() {
	return func() {
		c.Core.Log.NewEvent("Stele checking for tick", glog.LogCharacterEvent, c.Index).
			Write("src", src).
			Write("char", c.Index)
		if !c.Core.Constructs.Has(src) {
			return
		}
		c.Core.Log.NewEvent("Stele ticked", glog.LogCharacterEvent, c.Index).
			Write("next expected", c.Core.F+120).
			Write("src", src).
			Write("char", c.Index)

		// Use snapshot for damage
		ae := c.steleSnapshot

		//check how many times to hit
		count := c.Core.Constructs.Count() - c.Core.Constructs.CountByType(construct.GeoConstructZhongliSkill) + 1
		if count > max {
			count = max
		}
		orb := false
		for i := 0; i < count; i++ {
			c.Core.QueueAttackEvent(&ae, 0)
			if c.energyICD < c.Core.F && !orb && c.Core.Rand.Float64() < .5 {
				orb = true
			}
		}
		if orb {
			c.energyICD = c.Core.F + 90
			c.Core.QueueParticle("zhongli", 1, attributes.Geo, 20+c.ParticleDelay)
		}
		c.Core.Tasks.Add(c.resonance(src, max), 120)
	}
}
