package jahoda

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1Init() {
	eleCountMap := c.countElements()

	priority := []attributes.Element{
		attributes.Pyro,
		attributes.Hydro,
		attributes.Electro,
		attributes.Cryo,
	}

	highestEleCount := 0

	for _, ele := range priority {
		if eleCountMap[ele] > highestEleCount {
			highestEleCount = eleCountMap[ele]
			c.a1HighestEle = ele
		}
	}

	if highestEleCount == 0 {
		c.a1HighestEle = attributes.NoElement
	}

	if c.Base.Cons >= 2 && c.Core.Player.GetMoonsignLevel() > 2 {
		secondHighestEleCount := 0

		for _, ele := range priority {
			if ele == c.a1HighestEle {
				continue
			}

			if eleCountMap[ele] > secondHighestEleCount {
				secondHighestEleCount = eleCountMap[ele]
				c.c2NextHighestEle = ele
			}
		}

		if secondHighestEleCount == 0 {
			c.c2NextHighestEle = attributes.NoElement
		}
	}
}

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	c.applyA1Buff(c.a1HighestEle)

	if c.Base.Cons >= 2 && c.Core.Player.GetMoonsignLevel() > 2 {
		c.applyA1Buff(c.c2NextHighestEle)
	}
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
		c.robotHi.Src *= 1.2
	case attributes.Electro:
		c.robotCount += 1
	case attributes.Cryo:
		c.robotHitmarkInterval *= 0.9
	}
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Player.ActiveChar().AddStatMod(character.StatMod{
		Base:         modifier.NewBase("jahoda-a4", 6*60),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return c.a4Buff, true
		},
	})

	c.Core.Log.NewEvent("jahoda a4 triggered", glog.LogCharacterEvent, c.Index()).Write("em snapshot", c.a4Buff[attributes.EM]).Write("expiry", c.Core.F+6*60)
}
