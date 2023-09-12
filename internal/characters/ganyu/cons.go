package ganyu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key = "ganyu-c1"
	c4Key = "ganyu-c4"
	c4Dur = 180
	c6Key = "ganyu-c6"
)

// Ganyu C1: Taking DMG from a Charge Level 2 Frostflake Arrow or Frostflake Arrow Bloom decreases opponents' Cryo RES by 15% for 6s.
// A hit regenerates 2 Energy for Ganyu. This effect can only occur once per Charge Level 2 Frostflake Arrow, regardless if Frostflake Arrow itself or its Bloom hit the target.
func (c *char) c1() combat.AttackCBFunc {
	if c.Base.Cons < 1 {
		return nil
	}
	done := false

	return func(a combat.AttackCB) {
		e := a.Target.(*enemy.Enemy)
		if e.Type() != targets.TargettableEnemy {
			return
		}
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag(c1Key, 300),
			Ele:   attributes.Cryo,
			Value: -0.15,
		})
		if done {
			return
		}
		done = true
		c.AddEnergy(c1Key, 2)
	}
}

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c4Key, -1),
			Amount: func(_ *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				// reset stacks on expiry
				if !x.StatusIsActive(c4Key) {
					x.RemoveTag(c4Key)
					return nil, false
				}
				m[attributes.DmgP] = float64(x.GetTag(c4Key)) * 0.05
				return m, true
			},
		})
	}
}
