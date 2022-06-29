package xiangling

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1(a combat.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(enemy.ResistMod{
		Base:  modifier.NewBaseWithHitlag("xiangling-c1", 6*60),
		Ele:   attributes.Pyro,
		Value: -0.15,
	})

}

func (c *char) c6(dur int) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.15

	c.Core.Status.Add("xlc6", dur)

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("xiangling-c6", dur),
			AffectedStat: attributes.PyroP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
