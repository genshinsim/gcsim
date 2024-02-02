package rosaria

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Rosaria deals a CRIT Hit, her ATK Speed increase by 10% and her Normal Attack DMG increases by 10% for 4s (can trigger vs shielded enemies)
func (c *char) makeC1CB() combat.AttackCBFunc {
	if c.Base.Cons < 1 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if !a.IsCrit {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.1
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("rosaria-c1-dmg", 240), // 4s
			Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagNormal {
					return nil, false
				}
				return m, true
			},
		})

		mAtkSpd := make([]float64, attributes.EndStatType)
		mAtkSpd[attributes.AtkSpd] = 0.1
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("rosaria-c1-speed", 240), // 4s
			AffectedStat: attributes.AtkSpd,
			Amount: func() ([]float64, bool) {
				if c.Core.Player.CurrentState() != action.NormalAttackState {
					return nil, false
				}
				return mAtkSpd, true
			},
		})
	}
}

// Ravaging Confession's CRIT Hits regenerate 5 Energy for Rosaria. Can only be triggered once each time Ravaging Confession is cast.
func (c *char) makeC4CB() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if !a.IsCrit {
			return
		}
		if done {
			return
		}
		done = true
		c.AddEnergy("rosaria-c4", 5)
	}
}

// Rites of Termination's attack decreases opponent's Physical RES by 20% for 10s.
func (c *char) makeC6CB() combat.AttackCBFunc {
	if c.Base.Cons < 6 {
		return nil
	}
	return func(a combat.AttackCB) {
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
}
