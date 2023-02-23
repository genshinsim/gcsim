package ayaka

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// After using Kamisato Art: Hyouka, Kamisato Ayaka's Normal and Charged Attacks deal 30% increased DMG for 6s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.3
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("ayaka-a1", 360),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			return m, atk.Info.AttackTag == attacks.AttackTagNormal || atk.Info.AttackTag == attacks.AttackTagExtra
		},
	})
}

// When the Cryo application at the end of Kamisato Art: Senho hits an opponent, Kamisato Ayaka gains the following effects:
//
// - Restores 10 Stamina
//
// - Gains 18% Cryo DMG Bonus for 10s.
func (c *char) makeA4CB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		c.Core.Player.RestoreStam(10)

		m := make([]float64, attributes.EndStatType)
		m[attributes.CryoP] = 0.18
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("ayaka-a4", 600),
			AffectedStat: attributes.CryoP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
