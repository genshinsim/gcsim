package sayu

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C2:
// Yoohoo Art: Fuuin Dash gains the following effects:
// DMG of Fuufuu Whirlwind Kick in Press Mode increased by 3.3%.
// Every 0.5s in the Fuufuu Windwheel state will increase the DMG of this Fuufuu
// Whirlwind Kick by 3.3%. The maximum DMG increase possible through this method
// is 66%.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("sayu-c2", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.ActorIndex != c.Index {
				return nil, false
			}
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil, false
			}
			m[attributes.DmgP] = c.c2Bonus
			return m, true
		},
	})
}
