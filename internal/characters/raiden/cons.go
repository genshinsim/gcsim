package raiden

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When the Musou Isshin state applied by Secret Art: Musou Shinsetsu expires
// all nearby party members (excluding the Raiden Shogun) gain 30% bonus ATK for 10s.
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.3

	for i, char := range c.Core.Player.Chars() {
		if i == c.Index() {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("raiden-c4", 600),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

func (c *char) c6(ac info.AttackCB) {
	if c.Base.Cons < 6 {
		return
	}
	if c.Core.F < c.c6ICD {
		return
	}
	if c.c6Count == 5 {
		return
	}
	c.c6ICD = c.Core.F + 60
	c.c6Count++
	c.Core.Log.NewEvent("raiden c6 triggered", glog.LogCharacterEvent, c.Index()).
		Write("next_icd", c.c6ICD).
		Write("count", c.c6Count)
	for i, char := range c.Core.Player.Chars() {
		if i == c.Index() {
			continue
		}
		char.ReduceActionCooldown(action.ActionBurst, 60)
	}
}
