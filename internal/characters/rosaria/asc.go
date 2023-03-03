package rosaria

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Rosaria strikes an opponent from behind using Ravaging Confession, Rosaria's CRIT Rate increases by 12% for 5s.
// TODO: does this need to change if we add player position?
func (c *char) makeA1CB() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		m := make([]float64, attributes.EndStatType)
		m[attributes.CR] = 0.12
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("rosaria-a1", 300),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		c.Core.Log.NewEvent("Rosaria A1 activation", glog.LogCharacterEvent, c.Index).
			Write("ends_on", c.Core.F+300)
	}
}

// Casting Rites of Termination increases CRIT Rate of all nearby party members (except Rosaria herself)
// by 15% of Rosaria's CRIT Rate for 10s.
// CRIT Rate Bonus gained this way cannot exceed 15%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	crit_share := 0.15 * c.Stat(attributes.CR)
	if crit_share > 0.15 {
		crit_share = 0.15
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = crit_share
	for i, char := range c.Core.Player.Chars() {
		// skip Rosaria
		if i == c.Index {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("rosaria-a4", 600),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	c.Core.Log.NewEvent("Rosaria A4 activation", glog.LogCharacterEvent, c.Index).
		Write("ends_on", c.Core.F+600).
		Write("crit_share", crit_share)
}
