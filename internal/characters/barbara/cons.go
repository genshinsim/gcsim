package barbara

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1(delay int) {
	c.Core.Tasks.Add(func() {
		c.AddEnergy("barbara-c1", 1)
		c.c1(0)
	}, delay+10*60)
}

func (c *char) c2() {
	for i, char := range c.Core.Player.Chars() {
		if i == c.Index {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("barbara-c2", skillDuration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return c.c2buff, true
			},
		})
	}
}

// inspired from hutao c6
// TODO: does this even work?
func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if di.Amount <= 0 {
			return false
		}
		if c.Core.Player.Active() != c.Index { // trigger only when not barbara
			c.checkc6()
		}
		return false
	}, "barbara-c6")
}

func (c *char) checkc6() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Core.F < c.c6icd && c.c6icd != 0 {
		return
	}
	// grab the active char
	char := c.Core.Player.ActiveChar()
	// if dead, revive back to 1 hp
	if char.CurrentHPRatio() <= 0 {
		char.SetHPByAmount(1)
	}

	c.c6icd = c.Core.F + 60*60*15
}
