package candace

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// A1 is not implemented:
// TODO: If Candace is hit by an attack in the Hold duration of Sacred Rite: Heron's Sanctum, that skill will finish charging instantly.

const a4Key = "candace-a4"

// Characters affected by the Prayer of the Crimson Crown caused by Sacred Rite: Wagtail's Tide will deal 0.5% increased DMG
// to opponents for every 1,000 points of Candace's Max HP when they deal Elemental DMG with their Normal Attacks.
func (c *char) a4(char *character.CharWrapper) {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a4Key, -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if !c.StatusIsActive(burstKey) {
				return nil, false
			}
			if atk.Info.AttackTag != combat.AttackTagNormal {
				return nil, false
			}
			if atk.Info.Element == attributes.Physical || atk.Info.Element == attributes.NoElement {
				return nil, false
			}
			m[attributes.DmgP] = 0.005 * c.MaxHP() / 1000
			return m, true
		},
	})
}
