package xinyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.12

	icd := -1
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		// TODO: it works off field?
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if !crit {
			return false
		}
		if icd > c.Core.F {
			return false
		}

		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("xinyan-c1", 5*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				return m, true
			},
		})
		icd = c.Core.F + 5*60

		return false
	}, "xinyan-c1")
}

// Riff Revolution's Physical DMG has its CRIT Rate increased by 100%, and will form a shield at Shield Level 3: Rave when cast.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 1

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("xinyan-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalBurst {
				return nil, false
			}
			return m, true
		},
	})
}

// Sweeping Fervor's swing DMG decreases opponent's Physical RES by 15% for 12s.
func (c *char) c4(a combat.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}

	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(enemy.ResistMod{
		Base:  modifier.NewBase("xinyan-c4", 12*60),
		Ele:   attributes.Physical,
		Value: -0.15,
	})
}

// Decreases the Stamina Consumption of Xinyan's Charged Attacks by 30%. Additionally, Xinyan's Charged Attacks gain an ATK Bonus equal to 50% of her DEF.
// func (c *char) c6() {}
