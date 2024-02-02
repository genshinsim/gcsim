package venti

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C1:
// Fires 2 additional arrows per Aimed Shot, each dealing 33% of the original arrow's DMG.
func (c *char) c1(ai combat.AttackInfo, travel int) {
	ai.Abil = "Aim (Charged) C1"
	ai.Mult /= 3.0
	for i := 0; i < 2; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			aimedHitmark,
			aimedHitmark+travel,
		)
	}
}

// C2:
// Skyward Sonnet decreases opponents' Anemo RES and Physical RES by 12% for 10s.
// Opponents launched by Skyward Sonnet suffer an additional 12% Anemo RES and Physical RES decrease while airborne.
// TODO: the airborne part isn't implemented
func (c *char) c2(a combat.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}

	e.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("venti-c2-anemo", 600),
		Ele:   attributes.Anemo,
		Value: -0.12,
	})
	e.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("venti-c2-phys", 600),
		Ele:   attributes.Physical,
		Value: -0.12,
	})
}

// C4:
// When Venti picks up an Elemental Orb or Particle, he receives a 25% Anemo DMG Bonus for 10s.
func (c *char) c4() {
	c.c4bonus = make([]float64, attributes.EndStatType)
	c.c4bonus[attributes.AnemoP] = 0.25
	c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		// only trigger if Venti catches the particle
		if c.Core.Player.Active() != c.Index {
			return false
		}
		// apply C4 to Venti
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("venti-c4", 600),
			AffectedStat: attributes.AnemoP,
			Amount: func() ([]float64, bool) {
				return c.c4bonus, true
			},
		})
		return false
	}, "venti-c4")
}

// C6:
// Targets who take DMG from Wind's Grand Ode have their Anemo RES decreased by 20%.
// If an Elemental Absorption occurred, then their RES towards the corresponding Element is also decreased by 20%.
func (c *char) c6(ele attributes.Element) func(a combat.AttackCB) {
	return func(a combat.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("venti-c6-"+ele.String(), 600),
			Ele:   ele,
			Value: -0.20,
		})
	}
}
