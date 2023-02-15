package xinyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1ICDKey = "xinyan-c1-icd"

// Upon scoring a CRIT Hit, increases ATK SPD of Xinyan's Normal and Charged Attacks by 12% for 5s.
// Can only occur once every 5s.
func (c *char) makeC1CB() combat.AttackCBFunc {
	if c.Base.Cons < 1 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if !a.IsCrit {
			return
		}
		if c.StatusIsActive(c1ICDKey) {
			return
		}
		c.AddStatus(c1ICDKey, 5*60, true)

		m := make([]float64, attributes.EndStatType)
		m[attributes.AtkSpd] = 0.12
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("xinyan-c1", 5*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				return m, true
			},
		})
	}
}

// Riff Revolution's Physical DMG has its CRIT Rate increased by 100%, and will form a shield at Shield Level 3: Rave when cast.
func (c *char) c2() {
	c.c2Buff = make([]float64, attributes.EndStatType)
	c.c2Buff[attributes.CR] = 1

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("xinyan-c2", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalBurst {
				return nil, false
			}
			return c.c2Buff, true
		},
	})
}

// Sweeping Fervor's swing DMG decreases opponent's Physical RES by 15% for 12s.
func (c *char) makeC4CB() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("xinyan-c4", 12*60),
			Ele:   attributes.Physical,
			Value: -0.15,
		})
	}
}

// Decreases the Stamina Consumption of Xinyan's Charged Attacks by 30%. Additionally, Xinyan's Charged Attacks gain an ATK Bonus equal to 50% of her DEF.
// func (c *char) c6() {}
