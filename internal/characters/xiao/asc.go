package xiao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "xiao-a1"

// While under the effects of Bane of All Evil, all DMG dealt by Xiao increases
// by 5%. DMG increases by a further 5% for every 3s the ability persists. The
// maximum DMG Bonus is 25%.
func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a1Key, 900+burstStart),
		AffectedStat: attributes.DmgP,
		Amount: func() ([]float64, bool) {
			stacks := 1 + int((c.Core.F-c.qStarted)/180)
			if stacks > 5 {
				stacks = 5
			}
			m[attributes.DmgP] = float64(stacks) * 0.05
			return m, true
		},
	})
}
