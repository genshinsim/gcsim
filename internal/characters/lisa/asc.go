package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func (c *char) a4(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	t.AddDefMod("lisa-a4", 600, -0.15)
}
