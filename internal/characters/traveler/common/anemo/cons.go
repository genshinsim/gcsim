package anemo

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *Traveler) c2() {
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

func c6cb(ele attributes.Element) func(a info.AttackCB) {
	return func(a info.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("amc-c6-"+ele.String(), 600),
			Ele:   ele,
			Value: -0.20,
		})
	}
}
