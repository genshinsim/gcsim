package navia

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2IcdKey = "navia-c2-icd"

// Each stack of Crystal Shrapnel consumed when Navia uses Ceremonial Crystalshot will
// restore 3 Energy to her and decrease the CD of As the Sunlit Sky's Singing Salute by 1s.
// Up to 9 Energy can be gained this way, and the CD of "As the Sunlit Sky's Singing Salute"
// can be decreased by up to 3s.
func (c *char) c1(shrapnel int) {
	if c.Base.Cons < 1 {
		return
	}
	count := min(shrapnel, 3)
	c.ReduceActionCooldown(action.ActionBurst, count*60)
	c.AddEnergy("navia-c1-energy", float64(count*3))
}

// Each stack of Crystal Shrapnel consumed will increase the CRIT Rate of this
// Ceremonial Crystalshot instance by 12%. CRIT Rate can be increased by up to 36% in this way.
// In addition, when Ceremonial Crystalshot hits an opponent, one Cannon Fire Support shot from
// As the Sunlit Sky's Singing Salute will strike near the location of the hit.
// Up to one instance of Cannon Fire Support can be triggered each time Ceremonial Crystalshot is used,
// and DMG dealt by said Cannon Fire Support this way is considered Elemental Burst DMG.
func (c *char) c2() combat.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	return func(a combat.AttackCB) {
		e := a.Target.(*enemy.Enemy)
		if e.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(c2IcdKey) {
			return
		}
		c.AddStatus(c2IcdKey, 0.25*60, true)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "The President's Pursuit of Victory",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupNaviaBurst,
			StrikeType: attacks.StrikeTypeBlunt,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       burst[1][c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(geometry.CalcRandomPointFromCenter(e.Pos(), 0, 1.2, c.Core.Rand), nil, 3),
			0,
			30, // somewhere between 28-31
			c.burstCB(),
			c.c4(),
		)
	}
}

// When As the Sunlit Sky's Singing Salute hits an opponent,
// that opponent's Geo RES will be decreased by 20% for 8s.
func (c *char) c4() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		e := a.Target.(*enemy.Enemy)
		if e.Type() != targets.TargettableEnemy {
			return
		}
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("navia-c4-shred", 8*60),
			Ele:   attributes.Geo,
			Value: -0.2,
		})
	}
}
