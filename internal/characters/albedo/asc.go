package albedo

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Transient Blossoms generated by Abiogenesis: Solar
// Isotoma deal 25% more DMG to opponents whose HP is below 50%.
func (c *char) a1() {
	if !c.Core.Combat.DamageMode {
		return
	}
	if c.Base.Ascension < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.25
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("albedo-a1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil, false
			}
			// Can't be triggered by itself when refreshing
			if atk.Info.Abil == "Abiogenesis: Solar Isotoma" {
				return nil, false
			}
			if e, ok := t.(*enemy.Enemy); !(ok && e.HP()/e.MaxHP() < .5) {
				return nil, false
			}
			return m, true
		},
	})
}

// Using Rite of Progeniture: Tectonic Tide increases the Elemental Mastery of nearby party members by 125 for 10s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 125
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("albedo-a4", 600),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
