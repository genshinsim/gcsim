package hutao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1BuffKey = "hutao-a1"
)

// When a Paramita Papilio state activated by Guide to Afterlife ends,
// all allies in the party (excluding Hu Tao herself) will have their CRIT Rate increased by 12% for 8s.
func (c *char) a1() {
	if c.Base.Ascension < 1 || !c.applyA1 {
		return
	}
	c.applyA1 = false

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

// When Hu Tao's HP is equal to or less than 50%, her Pyro DMG Bonus is increased by 33%.
//
// - TODO: in-game this is actually a check every 0.3s. if hp is < 50% then buff is active until the next time check takes places
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4buff = make([]float64, attributes.EndStatType)
	c.a4buff[attributes.PyroP] = 0.33
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("hutao-a4", -1),
		AffectedStat: attributes.PyroP,
		Amount: func() ([]float64, bool) {
			if c.CurrentHPRatio() <= 0.5 {
				return c.a4buff, true
			}
			return nil, false
		},
	})
}
