package diona

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = .15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("diona-c2", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag == attacks.AttackTagElementalArt {
				return m
			}
			return nil
		},
	})
}

func (c *char) c2BonusShield() float64 {
	if c.Base.Cons < 2 {
		return 0
	}
	return 0.15
}

// Characters within Signature Mix's radius will gain the following effects based on their HP levels:
// The Incoming Healing Bonus is increased by 30% when HP is at, or lower than, 50%.
// Elemental Mastery is increased by 200 when HP is above 50%.

// Additionally, Diona's Max HP is increased by 25%.

// Radiance: Stellar Glimmer: Characters within the Signature Mix field deal increased DMG for the following reactions:
// Radiance: Stellar-Conduct: Superconduct, Stellar-Conduct reaction DMG increased by 40%.
// Radiance: Stellar Swirl: Cryo Swirl, Stellar Swirl reaction DMG increased by 40%.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c.c6buff = make([]float64, attributes.EndStatType)
	c.c6buff[attributes.EM] = 200

	if !c.revelation {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = .25
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("diona-c6-revelation", -1),
		AffectedStat: attributes.HPP,
		Amount: func() []float64 {
			return m
		},
	})
}

func (c *char) c6() {
	// c6 should last for the duration of the burst
	// lasts 12.5 second, ticks every 0.5s; adds mod to active char for 2s
	for i := 30; i <= 750; i += 30 {
		c.Core.Tasks.Add(func() {
			if !c.Core.Combat.Player().IsWithinArea(c.burstBuffArea) {
				return
			}
			// add 200EM to active char
			active := c.Core.Player.ActiveChar()
			if active.CurrentHPRatio() > 0.5 {
				active.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("diona-c6", 120),
					AffectedStat: attributes.EM,
					Amount: func() []float64 {
						return c.c6buff
					},
				})
			} else {
				// add healing bonus if hp <= 0.5
				// bonus only lasts for 120 frames
				active.AddHealBonusMod(character.HealBonusMod{
					Base: modifier.NewBaseWithHitlag("diona-c6-healbonus", 120),
					Amount: func() float64 {
						return 0.3
					},
				})
			}

			switch c.getRadiance() {
			case radianceStellarConduct:
				active.AddReactBonusMod(character.ReactBonusMod{
					Base: modifier.NewBaseWithHitlag("diona-revelation-c6-ssc", 120),
					Amount: func(ai info.AttackInfo) float64 {
						switch ai.AttackTag {
						case
							attacks.AttackTagDirectStellarConduct,
							attacks.AttackTagSuperconductDamage:
							return 0.4
						}
						return 0
					},
				})
			case radianceStellarSwirl:
				active.AddReactBonusMod(character.ReactBonusMod{
					Base: modifier.NewBaseWithHitlag("diona-revelation-c6-ssw", 120),
					Amount: func(ai info.AttackInfo) float64 {
						switch ai.AttackTag {
						case
							attacks.AttackTagDirectStellarSwirl,
							attacks.AttackTagReactionStellarSwirl,
							attacks.AttackTagSwirlCryo:
							return 0.4
						}
						return 0
					},
				})
			}
		}, i)
	}
}
