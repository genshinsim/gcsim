package alhaitham

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Alhaitham's Charged or Plunging Attacks hit opponents, they will generate 1 Chisel-Light Mirror.
// This effect can be triggered once every 12s.
func (c *char) a1CB(a combat.AttackCB) {
	ae := a.AttackEvent
	//ignore if projection on icd
	if c.a1ICD > c.Core.F {
		return
	}
	//ignore if it isn't CA/Plunge
	if ae.Info.AttackTag != combat.AttackTagExtra && ae.Info.AttackTag != combat.AttackTagPlunge {
		return
	}
	c.a1ICD = c.Core.F + 720 //12s
	c.mirrorGain()
}

// Each point of Alhaitham's Elemental Mastery will increase the DMG dealt by
// Projection Attacks and Particular Field: Fetters of Phenomena by 0.12%.
// The maximum DMG increase for both these abilities is 100%.
func (c *char) a4() {
	c.a4bonus = 0.0012 * c.Stat(attributes.EM)
	if c.a4bonus > 1 {
		c.a4bonus = 1
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("alhaitham-a4", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			// only trigger on projection attack and burst damage
			if atk.Info.AttackTag != combat.AttackTagElementalBurst &&
				atk.Info.ICDTag != combat.ICDTagAlhaithamProjectionAttack {
				return nil, false
			}

			m[attributes.DmgP] = 0.0009 * c.Stat(attributes.EM)
			if m[attributes.DmgP] > 0.9 {
				m[attributes.DmgP] = 0.9
			}
			return m, true
		},
	})
}
