package cryo

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a4Key = "travelercryo-a4"
)

// Radiance: Stellar-Conduct: When a Frostpierce Star is on the field, DMG from the Traveler's
// Normal Attacks, Charged Attacks, and Plunging Attacks are converted to Cryo DMG that cannot be
// overridden by another elemental infusion, and DMG dealt is also increased by 80% of the
// Traveler's ATK.
func (c *Traveler) a1Conversion(ai *info.AttackInfo) {
	if c.Base.Ascension < 1 {
		return
	}

	if c.getRadiance() != radianceStellarConduct {
		return
	}

	if !c.StatusIsActive(skillKey) {
		return
	}

	ai.Element = attributes.Cryo
	ai.IgnoreInfusion = true
	ai.Mult += 0.8
}

// The Traveler's Elemental Mastery is increased by 8% of their ATK. Up to 160 Elemental Mastery can
// be gained in this way.
func (c *Traveler) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(a4Key, -1),
		Extra:        true,
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
			m[attributes.EM] = min(stats.TotalATK()*0.08, 160)
			return m
		},
	})
}
