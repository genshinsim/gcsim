package rosaria

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Adds event checker for C1: Unholy Revelation
// When Rosaria deals a CRIT Hit, her ATK Speed increase by 10% and her Normal Attack DMG increases by 10% for 4s (can trigger vs shielded enemies)
// TODO: Description is unclear whether attack speed affects NA + CA - assume that it only affects NA for now
func (c *char) c1() {
	c.c1bonus = make([]float64, attributes.EndStatType)
	c.c1bonus[attributes.AtkSpd] = 0.1
	c.c1bonus[attributes.DmgP] = 0.1
	// Add hook that monitors for crit hits. Mirrors existing favonius code
	// No log value saved as stat mod already shows up in debug view
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if !crit {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		// doesn't work off-field
		if c.Core.Player.Active() != c.Index {
			return false
		}

		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("rosaria-c1", 240), //4s
			Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != combat.AttackTagNormal {
					return nil, false
				}
				return c.c1bonus, true
			},
		})

		return false
	}, "rosaria-c1")
}

// Adds event checker for C4 Painful Grace
// Ravaging Confession's CRIT Hits regenerate 5 Energy for Rosaria. Can only be triggered once each time Ravaging Confession is cast.
// Only applies when a crit hit is resolved, so can't be handled within skill code directly
// TODO: Since this only is needed for her E, can change this so it spawns a subscription in her E code
// Then it can return true, which kills the callback
// However, would also need a timeout function as well since her E can not crit
// Requires additional work and references - will leave implementation for later
// TODO: conver this into a callback on first skill?
func (c *char) c4(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.c4completed {
		return
	}
	//check for crit
	if !a.IsCrit {
		return
	}
	c.AddEnergy("rosaria-c4", 5)
	c.c4completed = true
}

// Applies C6 effect to enemies hit by it
// Rites of Termination's attack decreases opponent's Physical RES by 20% for 10s.
// Takes in a snapshot definition, and returns the same snapshot with an on hit callback added to apply the debuff
func (c *char) c6(a combat.AttackCB) {
	if c.Base.Cons < 6 {
		return
	}
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("rosaria-c6", 600),
		Ele:   attributes.Physical,
		Value: -0.2,
	})
}
