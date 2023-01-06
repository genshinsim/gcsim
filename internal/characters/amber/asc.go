package amber

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = .1
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("amber-a1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			return m, atk.Info.AttackTag == combat.AttackTagElementalBurst
		},
	})
}

func (c *char) a4(a combat.AttackCB) {
	if !a.AttackEvent.Info.HitWeakPoint {
		return
	}
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("amber-a4", 600),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}
