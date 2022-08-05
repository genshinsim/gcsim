package hutao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	if !c.applyA1 {
		return
	}
	c.applyA1 = false

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.12
	for i, char := range c.Core.Player.Chars() {
		//does not affect hutao
		if c.Index == i {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("hutao-a1", 480),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.33
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("hutao-a4", -1),
		AffectedStat: attributes.PyroP,
		Amount: func() ([]float64, bool) {
			if c.Core.Status.Duration("paramita") == 0 {
				return nil, false
			}
			if c.HPCurrent/c.MaxHP() <= 0.5 {
				return m, true
			}
			return nil, false
		},
	})
}
