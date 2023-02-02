package yanfei

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
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
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		trg := args[0].(combat.Target)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.Abil == "Blazing Eye (A4)" {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagExtra || !crit {
			return false
		}
		// make it so a4 only applies hitlag once per A4 proc and not everytime an enemy gets hit
		defhalt := !c.a4HitlagApplied
		c.a4HitlagApplied = true

		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Blazing Eye (A4)",
			AttackTag:          combat.AttackTagExtra,
			ICDTag:             combat.ICDTagNone,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeDefault,
			Element:            attributes.Pyro,
			Durability:         25,
			Mult:               0.8,
			HitlagFactor:       0.05,
			CanBeDefenseHalted: defhalt,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3.5), 10, 10)

		return false
	}, "yanfei-a4")
}
