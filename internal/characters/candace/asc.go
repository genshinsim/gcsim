package candace

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a4Key = "candace-a4"

func (c *char) a4(char *character.CharWrapper, duration int) {
	char.AddAttackMod(character.AttackMod{ // TODO: is this right implementation?
		Base: modifier.NewBaseWithHitlag(a4Key, duration),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagNormal {
				return nil, false
			}
			if atk.Info.Element == attributes.Physical || atk.Info.Element == attributes.NoElement {
				return nil, false
			}
			m := make([]float64, attributes.EndStatType)
			m[attributes.DmgP] = 0.5 * math.Floor(c.MaxHP()/1000)
			return m, true
		},
	})
}
