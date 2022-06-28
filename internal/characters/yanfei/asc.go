package yanfei

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Yan Fei's Charged Attack consumes Scarlet Seals, each Scarlet Seal consumed will increase her Pyro DMG by 5% for 6 seconds. When this effect is repeatedly triggered it will overwrite the oldest bonus first.
// The Pyro DMG bonus from Proviso is applied before charged attack damage is calculated.
func (c *char) a1(stacks int) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = float64(stacks) * 0.05
	c.AddStatMod(character.StatMod{Base: modifier.NewBase("yanfei-a1", 360), AffectedStat: attributes.PyroP, Amount: func() ([]float64, bool) {
		return m, true
	}})
}

// When Yan Fei's Charged Attacks deal CRIT Hits, she will deal an additional instance of AoE Pyo DMG equal to 80% of her ATK. This DMG counts as Charged Attack DMG.
func (c *char) a4() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.Abil == "Blazing Eye (A4)" {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagExtra && !crit {
			return false
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Blazing Eye (A4)",
			AttackTag:  combat.AttackTagExtra,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       0.8,
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 1, 1)

		return false
	}, "yanfei-a4")
}
