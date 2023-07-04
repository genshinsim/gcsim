package mika

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Starfrost Swirl's Flowfrost Arrow first hits an opponent, or its Rimestar Flare hits an opponent,
// 1 Detector stack from Passive Talent "Suppressive Barrage" will be generated.
// You must have unlocked the Passive Talent "Suppressive Barrage" first.
func (c *char) c2() func(combat.AttackCB) {
	if c.Base.Cons < 2 || c.Base.Ascension < 1 {
		return nil
	}

	done := false
	return func(_ combat.AttackCB) {
		if done {
			return
		}

		c.addDetectorStack()
		done = true
	}
}

// The maximum number of Detector stacks that Starfrost Swirl's Soulwind can gain is increased by 1.
// You need to have unlocked the Passive Talent "Suppressive Barrage" first.
// Additionally, active characters affected by Soulwind will deal 60% more Physical CRIT DMG.
func (c *char) c6(char *character.CharWrapper) {
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("mika-c6", skillBuffDuration),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if c.Core.Player.Active() != char.Index {
				return nil, false
			}
			if atk.Info.Element != attributes.Physical {
				return nil, false
			}
			return c.c6buff, true
		},
	})
}
