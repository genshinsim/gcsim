package jahoda

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	eleCountMap := c.countElements()

	highestEleCount := 0

	for _, ele := range elePriority {
		if eleCountMap[ele] > highestEleCount {
			highestEleCount = eleCountMap[ele]
			c.a1HighestEle = ele
		}
	}

	if highestEleCount == 0 {
		c.a1HighestEle = attributes.NoElement
	}

	c.c2Init(eleCountMap)
}

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	c.applyA1Buff(c.a1HighestEle)
	c.c2()
}

func (c *char) countElements() map[attributes.Element]int {
	count := map[attributes.Element]int{
		attributes.Pyro:    0,
		attributes.Hydro:   0,
		attributes.Electro: 0,
		attributes.Cryo:    0,
	}

	for _, ch := range c.Core.Player.Chars() {
		if ch == nil {
			continue
		}

		switch ch.Base.Element {
		case attributes.Pyro,
			attributes.Hydro,
			attributes.Electro,
			attributes.Cryo:
			count[ch.Base.Element]++
		}
	}

	return count
}

func (c *char) applyA1Buff(ele attributes.Element) {
	switch ele {
	case attributes.Pyro:
		c.robotAi.FlatDmg *= 1.3
	case attributes.Hydro:
		c.robotHealCoeff = 1.2
	case attributes.Electro:
		c.robotCount += 1
	case attributes.Cryo:
		c.robotHitmarkInterval *= 0.9
	}
}

func (c *char) a4Init() {
	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4Buff[attributes.EM] = 100
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Player.ActiveChar().AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("jahoda-a4", 6*60),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return c.a4Buff
		},
	})
}
