package jahoda

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
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
	c.c6Buff[attributes.CR] = 0.05
	c.c6Buff[attributes.CD] = 0.4

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c6Key, 20*60),
			Amount: func(atk *info.AttackEvent, _ info.Target) ([]float64, bool) {
				if char.Moonsign < 1 {
					return nil, false
				}
				return c.c6Buff, true
			},
		})
	}

}
