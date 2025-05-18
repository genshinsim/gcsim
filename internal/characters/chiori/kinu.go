package chiori

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	kinuDmgRatio       = 1.7
	kinuDuration       = int(3 * 60)
	kinuStartDelay     = int(0.6 * 60)
	kinuAttackInterval = int(0.5 * 60)
	kinuAttackDelay    = 5 // should be 0.08s
)

// Kinu will attack nearby opponents, dealing AoE Geo DMG equivalent to 170% of
// Tamoto's DMG. DMG dealt this way is considered Elemental Skill DMG.
//
// Kinu will leave the field after 1 attack or after lasting 3s.
func (c *char) createKinu(src int, centerOffset, minRandom, maxRandom float64) func() {
	return func() {
		// determine kinu pos
		player := c.Core.Combat.Player()
		center := geometry.CalcOffsetPoint(
			player.Pos(),
			geometry.Point{Y: centerOffset},
			player.Direction(),
		)
		kinuPos := geometry.CalcRandomPointFromCenter(center, minRandom, maxRandom, c.Core.Rand)

		c.Core.Log.NewEvent("kinu spawned", glog.LogCharacterEvent, c.Index).Write("src", src)

		// spawn kinu
		kinu := newTicker(c.Core, kinuDuration, nil)
		kinu.cb = c.kinuAttack(src, kinu, kinuPos)
		kinu.interval = kinuAttackInterval
		c.Core.Tasks.Add(kinu.tick, kinuStartDelay)
		c.kinus = append(c.kinus, kinu)
	}
}

func (c *char) kinuAttack(src int, kinu *ticker, pos geometry.Point) func() {
	return func() {
		c.Core.Tasks.Add(func() {
			ai := combat.AttackInfo{
				Abil:       "Fluttering Hasode (Kinu)",
				ActorIndex: c.Index,
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagChioriSkill,
				ICDGroup:   attacks.ICDGroupChioriSkill,
				StrikeType: attacks.StrikeTypeBlunt,
				PoiseDMG:   0,
				Element:    attributes.Geo,
				Durability: 25,
				Mult:       turretAtkScaling[c.TalentLvlSkill()] * kinuDmgRatio,
			}

			snap := c.Snapshot(&ai)
			ai.FlatDmg = snap.Stats.TotalDEF()
			ai.FlatDmg *= turretDefScaling[c.TalentLvlSkill()] * kinuDmgRatio

			// if the player has an attack target it will always choose this enemy
			// so just need to make sure that it is within the search AoE
			t := c.Core.Combat.PrimaryTarget()
			if !t.IsWithinArea(combat.NewCircleHitOnTarget(pos, nil, c.skillSearchAoE)) {
				return
			}

			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(t, nil, skillDollAoE), 0)

			c.Core.Log.NewEvent("kinu killed on attack", glog.LogCharacterEvent, c.Index).Write("src", src)

			kinu.kill()
			c.cleanupKinu()
		}, kinuAttackDelay)
	}
}

func (c *char) cleanupKinu() {
	n := 0
	for _, t := range c.kinus {
		if t.alive {
			c.kinus[n] = t
			n++
		}
	}
	c.kinus = c.kinus[:n]
}
