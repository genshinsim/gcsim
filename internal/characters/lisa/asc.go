package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a4(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	t.AddDefMod(combat.DefMod{
		Base:  modifier.NewBaseWithHitlag("lisa-a4", 600),
		Value: -0.15,
	})
}
