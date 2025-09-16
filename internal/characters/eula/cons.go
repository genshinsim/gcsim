package eula

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c4() {
	if c.Core.Combat.DamageMode {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.25
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("eula-c4", -1),
			Amount: func(atk *info.AttackEvent, t info.Target) ([]float64, bool) {
				if atk.Info.Abil != "Glacial Illumination (Lightfall)" {
					return nil, false
				}
				if !c.Core.Combat.DamageMode {
					return nil, false
				}
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				if x.HP()/x.MaxHP() >= 0.5 {
					return nil, false
				}
				return m, true
			},
		})
	}
}
