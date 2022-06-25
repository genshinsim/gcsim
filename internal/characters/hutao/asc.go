package hutao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.12
	for i, char := range c.Core.Player.Chars() {
		//does not affect hutao
		if c.Index == i {
			continue
		}
		char.AddStatMod("hutao-a1", 480, attributes.CR, func() ([]float64, bool) {
			return m, true
		})
	}
}

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.33
	c.AddStatMod("hutao-a4", -1, attributes.PyroP, func() ([]float64, bool) {
		if c.Core.Status.Duration("paramita") == 0 {
			return nil, false
		}
		if c.HPCurrent/c.MaxHP() <= 0.5 {
			return m, true
		}
		return nil, false
	})
}
