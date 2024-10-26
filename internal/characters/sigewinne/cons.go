package sigewinne

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1FlatConvalescenceIncrease = 100.0 / 1000
	c1FlatConvalescenceCap      = 3500

	c6CDmgHpRatio  = 0.022 / 1000
	c6CRateHpRatio = 0.004 / 1000
	c6CDmgCap      = 1.1
	c6CRateCap     = 0.2
)

func (c *char) c2CB(a combat.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	if a.Damage == 0 {
		return
	}

	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("sigewinne-c2", 8*60),
		Ele:   attributes.Hydro,
		Value: -0.35,
	})
}

func (c *char) c6CritMode() {
	if c.Base.Cons < 6 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("sigewinne-c6", 15*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			m[attributes.CD] = min(c6CDmgCap, c.MaxHP()*c6CDmgHpRatio)
			m[attributes.CR] = min(c6CRateCap, c.MaxHP()*c6CRateHpRatio)
			return m, true
		},
	})
}
