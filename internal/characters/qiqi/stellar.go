package qiqi

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	stellarConductText = " (Stellar-Conduct)"
	radianceSwirlKey   = "radiance-stellar-swirl"
)

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

func (c *char) getRadiance() radianceState {
	if !c.revelation {
		return radianceNone
	}

	if c.StatusIsActive(reactable.PolestarFieldKey) {
		return radianceStellarConduct
	}

	if c.StatusIsActive(radianceSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}

// The cooldown of the Elemental Skill Adeptus Art: Herald of Frost is reduced to 15s.
// Radiance: Stellar Glimmer: While the Herald of Frost is on the field, Qiqi's own party members
// will have the DMG they deal via the corresponding reactions enhanced:
// Radiance: Stellar-Conduct: Superconduct and Stellar-Conduct reaction DMG is increased by 50%.
// Radiance: Stellar Swirl: Cryo Swirl and Stellar Swirl reaction DMG dealt is increased by 50%.
// Additionally, Qiqi will also enter the Radiance: Stellar-Conduct state when she is inside a
// Polestar Field, or the Radiance: Stellar Swirl state for 8s after a nearby party character
// triggers a Stellar Swirl reaction.
func (c *char) revelationInit() {
	if !c.revelation {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("qiqi-revelation-ssc", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !c.StatusIsActive(skillBuffKey) {
					return 0
				}
				if c.getRadiance() != radianceStellarConduct {
					return 0
				}
				switch ai.AttackTag {
				case attacks.AttackTagDirectStellarConduct,
					attacks.AttackTagSuperconductDamage:
					return 0.5
				}
				return 0
			},
		})

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("qiqi-revelation-ssw", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !c.StatusIsActive(skillBuffKey) {
					return 0
				}
				if c.getRadiance() != radianceStellarSwirl {
					return 0
				}
				switch ai.AttackTag {
				case attacks.AttackTagDirectStellarSwirl,
					attacks.AttackTagReactionStellarSwirl,
					attacks.AttackTagSwirlCryo:
					return 0.5
				}
				return 0
			},
		})
	}

	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		c.AddStatus(radianceSwirlKey, 8*60, false)
	}, "qiqi-ssw")
}

func (c *char) revelationSkillCDReduction() {
	if !c.revelation {
		return
	}

	// despite the description, this is implemented as a 15s CD reduction in game
	c.ReduceActionCooldown(action.ActionSkill, 15*60)
}
