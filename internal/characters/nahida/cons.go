package nahida

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When the Shrine of Maya is unleashed and the Elemental Types of the party
// members are being tabulated, the count will add 1 to the number of Pyro,
// Electro, and Hydro characters respectively.
func (c *char) c1() {
	c.pyroCount++
	c.hydroCount++
	c.electroCount++
}

// Opponents that are marked by Nahida's own Seed of Skandha will be affected by the following effects:
//   - Burning, Bloom, Hyperbloom, Burgeon Reaction DMG can score CRIT Hits.
//     CRIT Rate and CRIT DMG are fixed at 20% and 100% respectively.
//   - Within 8s of being affected by Quicken, Aggravate, Spread, DEF is decreased by 30%.
func (c *char) c2() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		ae := args[1].(*combat.AttackEvent)

		switch ae.Info.AttackTag {
		case combat.AttackTagBurningDamage:
		case combat.AttackTagBloom:
		case combat.AttackTagHyperbloom:
		case combat.AttackTagBurgeon:
		default:
			return false
		}

		if !t.StatusIsActive(skillMarkKey) {
			return false
		}

		//TODO: should this really be +=??
		ae.Snapshot.Stats[attributes.CR] += 0.2
		ae.Snapshot.Stats[attributes.CD] += 1

		c.Core.Log.NewEvent("nahida c2 buff", glog.LogCharacterEvent, ae.Info.ActorIndex).
			Write("final_crit", ae.Snapshot.Stats[attributes.CR])

		return false
	}, "nahida-c2-reaction-dmg-buff")

	cb := func(rx event.Event) event.EventHook {
		return func(args ...interface{}) bool {
			t, ok := args[0].(*enemy.Enemy)
			if !ok {
				return false
			}
			if !t.StatusIsActive(skillMarkKey) {
				return false
			}
			t.AddDefMod(enemy.DefMod{
				Base:  modifier.NewBaseWithHitlag("nahida-c2", 480),
				Value: -0.3,
			})
			return false
		}
	}

	c.Core.Events.Subscribe(event.OnQuicken, cb(event.OnQuicken), "nahida-c2-def-shred-quicken")
	c.Core.Events.Subscribe(event.OnAggravate, cb(event.OnAggravate), "nahida-c2-def-shred-aggravate")
	c.Core.Events.Subscribe(event.OnSpread, cb(event.OnSpread), "nahida-c2-def-shred-spread")
}

// When 1/2/3/(4 or more) nearby opponents are affected by All Schemes to Know's
// Seeds of Skandha, Nahida's Elemental Mastery will be increased by
// 100/120/140/160.
func (c *char) c4() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("nahida-c4", -1),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			x, y := c.Core.Combat.Player().Pos()
			count := 0
			for _, v := range c.Core.Combat.EnemiesWithinRadius(x, y, 30) {
				enemy, ok := c.Core.Combat.Enemy(v).(*enemy.Enemy)
				if !ok {
					continue
				}
				if enemy.StatusIsActive(skillMarkKey) {
					count++
				}
			}
			if count > 4 {
				count = 4
			}
			if count == 0 {
				return nil, false
			}
			c.c4Buff[attributes.EM] = float64(80 + count*20)
			return c.c4Buff, true
		},
	})
}

const (
	c6ICDKey    = "nahida-c6-icd"
	c6ActiveKey = "nahida-c6"
)

// When Nahida hits an opponent affected by All Schemes to Know's Seeds of
// Skandha with Normal or Charged Attacks after unleashing Illusory Heart, she
// will use Tri-Karma Purification: Karmic Oblivion on this opponent and all
// connected opponents, dealing Dendro DMG based on 200% of Nahida's ATK and 400%
// of her Elemental Mastery. DMG dealt by Tri-Karma Purification: Karmic Oblivion
// is considered Elemental Skill DMG and can be triggered once every 0.2s. This
// effect can last up to 10s and will be removed after Nahida has unleashed 6
// instances of Tri-Karma Purification: Karmic Oblivion.
func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if c.c6Count >= 6 {
			return false
		}
		if c.StatusIsActive(c6ICDKey) {
			return false
		}
		if !c.StatusIsActive(c6ActiveKey) {
			return false
		}
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}
		switch ae.Info.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagExtra:
		default:
			return false
		}
		if !e.StatusIsActive(skillMarkKey) {
			return false
		}
		c.AddStatus(c6ICDKey, 12, true) //TODO: hitlag?
		c.c6Count++
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tri-Karma Purification: Karmic Oblivion",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNahidaC6,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Dendro,
			Durability: 25,
			Mult:       2,
		}
		snap := c.Snapshot(&ai)
		ai.FlatDmg = snap.Stats[attributes.EM] * 4
		for _, v := range c.Core.Combat.Enemies() {
			e, ok := v.(*enemy.Enemy)
			if !ok {
				continue
			}
			if !e.StatusIsActive(skillMarkKey) {
				continue
			}
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewDefSingleTarget(e.Key()),
				1,
			)
		}

		return false
	}, "nahida-c6")

}
