package alhaitham

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1IcdKey = "alhaitham-a1-icd"

// When Alhaitham's Charged or Plunging Attacks hit opponents, they will generate 1 Chisel-Light Mirror.
// This effect can be triggered once every 12s.
func (c *char) makeA1CB() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		// ignore if projection on icd
		if c.Core.Status.Duration(a1IcdKey) > 0 {
			return
		}

		c.Core.Status.Add(a1IcdKey, 720) // 12s
		c.mirrorGain(1)
	}
}

// Each point of Alhaitham's Elemental Mastery will increase the DMG dealt by
// Projection Attacks and Particular Field: Fetters of Phenomena by 0.1%.
// The maximum DMG increase for both these abilities is 100%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("alhaitham-a4", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			// only trigger on projection attack and burst damage
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst &&
				atk.Info.ICDGroup != attacks.ICDGroupAlhaithamProjectionAttack {
				return nil, false
			}

			m[attributes.DmgP] = 0.001 * c.Stat(attributes.EM)
			if m[attributes.DmgP] > 1 {
				m[attributes.DmgP] = 1
			}
			return m, true
		},
	})
}
