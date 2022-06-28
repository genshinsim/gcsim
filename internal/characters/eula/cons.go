package eula

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.25
	c.AddAttackMod(character.AttackMod{Base: modifier.NewBase("eula-c4", -1), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.Abil != "Glacial Illumination (Lightfall)" {
			return nil, false
		}
		if !c.Core.Combat.DamageMode {
			return nil, false
		}
		if t.HP()/t.MaxHP() >= 0.5 {
			return nil, false
		}
		return m, true
	}})
}
