package tighnari

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Tighnari's Charged Attack CRIT Rate is increased by 15%.
func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("tighnari-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			return m, true
		},
	})
}

// When there are opponents within the Vijnana-Khanda Field created by Vijnana-Phala Mine, Tighnari gains 20% Dendro DMG Bonus.
// The effect will last up to 6s if the field's duration ends or if it no longer has opponents within it.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DendroP] = .2
	for i := 0; i < 8*60; i += 30 {
		c.Core.Tasks.Add(func() {
			if !c.Core.Combat.Player().IsWithinArea(c.skillArea) {
				return
			}
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBase("tighnari-c2", 6*60),
				AffectedStat: attributes.DendroP,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}, i)
	}
}

// When Fashioner's Tanglevine Shaft is unleashed, all nearby party members gain 60 Elemental Mastery for 8s.
// TODO: If the Fashioner's Tanglevine Shaft triggers a Burning, Bloom, Quicken, or Spread reaction, their Elemental Mastery
// will be further increased by 60. This latter case will also refresh the buff state's duration.
func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 60
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("tighnari-c4", 8*60),
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}

		return false
	}, "tighnari-c4")

	f := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
			return false
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 120
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("tighnari-c4", 8*60),
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}

		return false
	}
	c.Core.Events.Subscribe(event.OnBurning, f, "tighnari-c4-burning")
	c.Core.Events.Subscribe(event.OnBloom, f, "tighnari-c4-bloom")
	c.Core.Events.Subscribe(event.OnQuicken, f, "tighnari-c4-quicken")
	c.Core.Events.Subscribe(event.OnSpread, f, "tighnari-c4-spread")
}
