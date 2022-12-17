package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"math"
)

func (c *char) c1() {
	// C1: Needs to be manually deleted when Windfavored state ends
	if c.Base.Cons >= 1 && c.StatusIsActive(skillKey) {
		m := make([]float64, attributes.EndStatType)
		m[attributes.AtkSpd] = 0.1
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("wanderer-c1-atkspd", 1200),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if !(atk.Info.AttackTag == combat.AttackTagNormal || atk.Info.AttackTag == combat.AttackTagExtra) {
					return nil, false
				}
				return m, true
			},
		})

	}
}

func (c *char) c2() {
	// C2: Buff stays active during entire burst animation
	if c.Base.Cons >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = math.Min((float64)(c.maxSkydwellerPoints-c.skydwellerPoints)*0.04, 2)
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("wanderer-c2-burstbonus", burstFramesE[action.InvalidAction]),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if !(atk.Info.AttackTag == combat.AttackTagElementalBurst) {
					return nil, false
				}
				return m, true
			},
		})

	}
}
