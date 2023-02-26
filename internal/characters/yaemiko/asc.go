package yaemiko

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When casting Great Secret Art: Tenko Kenshin, each Sesshou Sakura destroyed
// resets the cooldown for 1 charge of Yakan Evocation: Sesshou Sakura.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.ResetActionCooldown(action.ActionSkill)
}

// Every point of Elemental Mastery Yae Miko possesses will increase Sesshou Sakura DMG by 0.15%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("yaemiko-a1", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			// only trigger on elemental art damage
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil, false
			}
			m[attributes.DmgP] = c.Stat(attributes.EM) * 0.0015
			return m, true
		},
	})
}
