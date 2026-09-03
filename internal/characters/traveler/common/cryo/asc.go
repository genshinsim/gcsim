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
