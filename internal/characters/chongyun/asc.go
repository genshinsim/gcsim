package chongyun

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When the field created by Spirit Blade: Chonghua's Layered Frost disappears,
// another spirit blade will be summoned to strike nearby opponents, dealing 100% of Chonghua's Layered Frost's Skill DMG as AoE Cryo DMG.
func (c *char) a4(delay, src int, useOldSnapshot bool) {
	if c.Base.Ascension < 4 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spirit Blade: Chonghua's Layered Frost (A4)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   100,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	// need to snap both snapshot and skill area into the task closure
	var snap combat.Snapshot
	if useOldSnapshot {
		snap = c.a4Snap
	} else {
		snap = c.Snapshot(&ai)
		c.a4Snap = snap
	}
	skillPattern := c.skillArea
	c.Core.Tasks.Add(func() {
		// if src changed then that means the field changed already
		if src != c.fieldSrc {
			return
		}
		enemy := c.Core.Combat.ClosestEnemyWithinArea(skillPattern, nil)
		var ap combat.AttackPattern
		if enemy != nil {
			ap = combat.NewCircleHitOnTarget(enemy, nil, 3.5)
		} else {
			ap = combat.NewCircleHitOnTarget(skillPattern.Shape.Pos(), nil, 3.5)
		}
		c.Core.QueueAttackWithSnap(ai, snap, ap, 0, c.a4CB, c.makeC4Callback())
	}, delay)
}

// Opponents hit by this blade will have their Cryo RES decreased by 10% for 8s.
func (c *char) a4CB(a combat.AttackCB) {
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("chongyun-a4", 480),
		Ele:   attributes.Cryo,
		Value: -0.10,
	})
}
