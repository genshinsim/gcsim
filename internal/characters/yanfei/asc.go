package yanfei

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Yanfei consumes Scarlet Seals by using a Charged Attack,
// each Scarlet Seal will increase Yanfei's Pyro DMG Bonus by 5%.
// This effects lasts for 6s. When a Charged Attack is used again
// during the effect's duration, it will dispel the previous effect.
func (c *char) a1(stacks int) {
	if c.Base.Ascension < 1 {
		return
	}
	c.a1Buff[attributes.PyroP] = float64(stacks) * 0.05
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("yanfei-a1", 360),
		AffectedStat: attributes.PyroP,
		Amount: func() ([]float64, bool) {
			return c.a1Buff, true
		},
	})
}

// When Yanfei's Charged Attack deals a CRIT Hit to opponents,
// she will deal an additional instance of AoE Pyro DMG equal to 80% of her ATK.
// This DMG counts as Charged Attack DMG.
func (c *char) makeA4CB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		trg := a.Target
		if trg.Type() != combat.TargettableEnemy {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if !a.IsCrit {
			return
		}
		if done {
			return
		}
		done = true

		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Blazing Eye (A4)",
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             combat.ICDTagNone,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeDefault,
			Element:            attributes.Pyro,
			Durability:         25,
			Mult:               0.8,
			HitlagFactor:       0.05,
			CanBeDefenseHalted: true,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3.5), 10, 10)
	}
}
