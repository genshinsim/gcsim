package kokomi

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1(f, travel int) {
	if c.Base.Cons < 1 {
		return
	}
	if c.Core.Status.Duration("kokomiburst") == 0 {
		return
	}

	// TODO: Assume that these are 1A (not specified in library)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "At Water's Edge (C1)",
		AttackTag:  combat.AttackTagNone,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
	}
	ai.FlatDmg = 0.3 * c.MaxHP()

	// TODO: Is this snapshotted/dynamic?
	c.Core.QueueAttack(ai, combat.NewDefSingleTarget(1, combat.TargettableEnemy), f, f+travel)
}

// C2 handling
// Sangonomiya Kokomi gains the following Healing Bonuses with regard to characters with 50% or less HP via the following methods:
// Nereid's Ascension Normal and Charged Attacks: 0.6% of Kokomi's Max HP.
func (c *char) c2() {
	for _, char := range c.Core.Player.Chars() {
		if char.HPCurrent/char.MaxHP() > .5 {
			continue
		}
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  char.Index,
			Message: "The Clouds Like Waves Rippling",
			Src:     0.006 * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
	}
}

// C4 (Energy piece only) handling
// While donning the Ceremonial Garment created by Nereid's Ascension, Sangonomiya Kokomi's Normal Attack SPD is increased by 10%.
// and Normal Attacks that hit opponents will restore 0.8 Energy for her. This effect can occur once every 0.2s.
func (c *char) c4() {
	if c.Core.F < c.c4ICDExpiry {
		return
	}
	c.c4ICDExpiry = c.Core.F + 12
	c.AddEnergy("kokomi-c4", 0.8)
}

// C6 handling
// While donning the Ceremonial Garment created by Nereid's Ascension.
// Sangonomiya Kokomi gains a 40% Hydro DMG Bonus for 4s after her Normal and Charged Attacks heal a character with 80% or more HP.
func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.HydroP] = .4
	for _, char := range c.Core.Player.Chars() {
		if char.HPCurrent/char.MaxHP() < .8 {
			continue
		}
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("kokomi-c6", 480),
			AffectedStat: attributes.HydroP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		// No need to continue checking if we found one
		break
	}
}
