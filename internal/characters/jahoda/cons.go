package jahoda

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4Key = "jahoda-c4-flat-energy"
	c6Key = "jahoda-c6"
)

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	c.AddEnergy(c4Key, 4)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	c.c6Buff = make([]float64, attributes.EndStatType)
	c.c6Buff = make([]float64, attributes.EndStatType)

	c.c6Buff[attributes.CR] = 0.05
	c.c6Buff[attributes.CD] = 0.40

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c6Key, 20*60),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				if char.Moonsign < 1 {
					return nil, false
				}
				return c.c6Buff, true
			},
		})

		c.Core.Log.NewEvent("jahoda c6 triggered", glog.LogCharacterEvent, c.Index()).
			Write("cr", c.c6Buff[attributes.CR]).
			Write("cd", c.c6Buff[attributes.CD]).
			Write("expiry", c.Core.F+20*60)

	}
}
