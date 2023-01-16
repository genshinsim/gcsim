package hutao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1BuffKey = "hutao-a1"
)

func (c *char) a1() {
	for i, char := range c.Core.Player.Chars() {
		//does not affect hutao
		if c.Index == i {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a1BuffKey, 480),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return c.a1buff, true
			},
		})
	}
}

func (c *char) a4() {
	//TODO: in game this is actually a check every 0.3s. if hp is < 50% then buff is active until
	//the next time check takes places
	c.a4buff = make([]float64, attributes.EndStatType)
	c.a4buff[attributes.PyroP] = 0.33
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("hutao-a4", -1),
		AffectedStat: attributes.PyroP,
		Amount: func() ([]float64, bool) {
			if c.HPCurrent/c.MaxHP() <= 0.5 {
				return c.a4buff, true
			}
			return nil, false
		},
	})
}
