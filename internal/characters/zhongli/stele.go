package zhongli

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const particleICDKey = "zhongli-particle-icd"

func (c *char) newStele(dur int) {
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
	steleDir := c.Core.Combat.Player().Direction()
	stelePos := combat.CalcOffsetPoint(c.Core.Combat.Player().Pos(), combat.Point{Y: 3}, steleDir)
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(stelePos, nil, 2), 0, 0, c.steleEnergyCB())

	//create a construct
	con := &stoneStele{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		c:      c,
		dir:    steleDir,
		pos:    stelePos,
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
	c.steleSnapshot = combat.AttackEvent{
		Info:        aiSnap,
		Snapshot:    snap,
		SourceFrame: c.Core.F,
	}

	c.Core.Tasks.Add(c.resonance(c.Core.F), 120)
}

func (c *char) resonance(src int) func() {
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

		boxOffset := combat.Point{Y: -4}
		boxSize := 8.0
		boxSizeSquared := boxSize * boxSize

		ai := ae.Info
		snap := ae.Snapshot

		steles, others := c.Core.Constructs.ConstructsByType(construct.GeoConstructZhongliSkill)

		particleCB := c.steleEnergyCB()
		for _, s := range steles {
			// skip other stele
			if s.Key() != src {
				continue
			}
			steleDir := s.Direction()
			stelePos := s.Pos()

			// get all constructs except for the steles within radius 8 of each stele for resonance purposes
			var resonanceConstructs []construct.Construct
			for _, con := range others {
				if con.Pos().Sub(stelePos).MagnitudeSquared() > boxSizeSquared {
					continue
				}
				resonanceConstructs = append(resonanceConstructs, con)
			}

			// queue stele attack
			steleAttackPos := combat.CalcOffsetPoint(stelePos, boxOffset, steleDir)
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewBoxHitOnTarget(steleAttackPos, nil, boxSize, boxSize), 0, particleCB)

			// queue resonance attacks
			for _, con := range resonanceConstructs {
				resonanceAttackPos := combat.CalcOffsetPoint(con.Pos(), boxOffset, con.Direction())
				c.Core.QueueAttackWithSnap(ai, snap, combat.NewBoxHitOnTarget(resonanceAttackPos, nil, boxSize, boxSize), 0, particleCB)
			}
		}
		c.Core.Tasks.Add(c.resonance(src), 120)
	}
}

func (c *char) steleEnergyCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		// can only proc once per skill cast / stele tick
		if done {
			return
		}
		done = true
		if c.StatusIsActive(particleICDKey) {
			return
		}
		// 50% chance
		if c.Core.Rand.Float64() > 0.5 {
			return
		}
		c.AddStatus(particleICDKey, 90, true)
		c.Core.QueueParticle("zhongli", 1, attributes.Geo, 20+c.ParticleDelay)
	}
}
