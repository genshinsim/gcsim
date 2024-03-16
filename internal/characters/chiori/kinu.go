package chiori

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

// Kinu will attack nearby opponents, dealing AoE Geo DMG equivalent to 170% of
// Tamoto's DMG. DMG dealt this way is considered Elemental Skill DMG.
//
// Kinu will leave the field after 1 attack or after lasting 3s.
func (c *char) createKinu() {
	// TODO:  is this right???
	t := newTicker(c.Core, 180)
	// add a on dmg sub based on this ticker
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if !t.alive {
			// in theory we should never be here
			return true
		}
		atk := args[1].(*combat.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}

		ai := combat.AttackInfo{
			Abil:       "Fluttering Hasode (Kinu)",
			ActorIndex: c.Index,
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagChioriSkill,
			ICDGroup:   attacks.ICDGroupChioriSkill,
			StrikeType: attacks.StrikeTypeBlunt,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       turretAtkScaling[c.TalentLvlSkill()] * 1.7,
		}
		snap := c.Snapshot(&ai)
		ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		ai.FlatDmg *= turretDefScaling[c.TalentLvlSkill()] * 1.7
		//TODO: hit box size
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2), 0)

		t.kill()
		c.cleanupKinu()

		return true
	}, fmt.Sprintf("chiori-kinu-%p", t))
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
