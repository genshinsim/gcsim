package navia

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"math"
)

// Each charge of Crystal Shrapnel consumed when Navia uses Ceremonial Crystalshot will
// restore 2 Energy to her and decrease the CD of As the Sunlit Sky's Singing Salute by 1s.
// Up to 6 Energy can be gained this way, and the CD of Ceremonial Crystalshot can be
// decreased by up to 3s.
func (c *char) c1(shrapnel int) {
	count := math.Min(float64(shrapnel), 3)
	c.ReduceActionCooldown(action.ActionBurst, int(count*60))
	c.AddEnergy("navia-c1-energy", count*2)
	return
}

// The CRIT Rate of Ceremonial Crystalshot is increased by 8% for each charge of Crystal
// Shrapnel consumed. CRIT Rate can be increased by up to 24% in this way.
// In addition, when Ceremonial Crystalshot hits an opponent, one shot of Fire Support
// from As the Sunlit Sky's Singing Salute will strike near the location of the hit.
// Up to one instance of Fire Support can be triggered each time Ceremonial Crystalshot is used,
// and DMG dealt by Fire Support in this way is considered Elemental Burst DMG.
func (c *char) c2(a combat.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	if !c.naviaburst {
		return
	}

	// Function doesn't check for enemy type or limit as assumes that the CB will check or it.

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "The President's Pursuit of Victory",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       burst[1][c.TalentLvlSkill()],
	}
	c.Core.QueueAttackWithSnap(
		ai,
		c.artillerySnapshot.Snapshot,
		combat.NewCircleHitOnTarget(a.Target.Pos(), nil, 3),
		0,
		nil,
	)
}

// When As the Sunlit Sky's Singing Salute hits an opponent,
// that opponent's Geo RES will be decreased by 20% for 8s.
func (c *char) c4(a combat.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}
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
