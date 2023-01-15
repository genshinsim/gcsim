package traveleranemo

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = .16

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("amc-c2", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func c6cb(ele attributes.Element) func(a combat.AttackCB) {
	return func(a combat.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("amc-c6-"+ele.String(), 600),
			Ele:   ele,
			Value: -0.20,
		})
	}
}
