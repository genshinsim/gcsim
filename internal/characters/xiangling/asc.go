package xiangling

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Increases the flame range of Guoba by 20%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.guobaFlameRange *= 1.2
}

// When Guoba Attack's effect ends, Guoba leaves a chili pepper on the spot where it disappeared. Picking up a chili pepper increases ATK by 10% for 10s.
func (c *char) a4(a4Delay int) {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.10
	// pick up chili pepper on active char after a user-specified delay relative to guoba expiry
	c.Core.Tasks.Add(func() {
		active := c.Core.Player.ActiveChar()
		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("xiangling-a4", 10*60),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		c.Core.Log.NewEvent(
			fmt.Sprintf("xiangling a4 chili pepper picked up by %v", active.Base.Key.String()),
			glog.LogCharacterEvent,
			c.Index,
		)
	}, a4Delay)
}
