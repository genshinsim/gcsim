package ganyu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key = "ganyu-c1"
	c4Key = "ganyu-c4"
	c6Key = "ganyu-c6"
	c1ICD = "ganyu-c1-energy-icd"
)

func (c *char) c1() combat.AttackCBFunc {
	return func(a combat.AttackCB) {
		e:= a.Target.(*enemy.Enemy)
		if e.Type() != combat.TargettableEnemy {
			return
		}
		e.AddResistMod(enemy.ResistMod{
			Base:  modifier.NewBaseWithHitlag("ganyu-c1", 300),
			Ele:   attributes.Cryo,
			Value: -0.15,
		})
		if !c.StatusIsActive(c1ICD) {
			c.AddEnergy(c1Key, 2)
			c.AddStatus(c1ICD, 24, false)
		}
		c.Core.Log.NewEvent("Rosaria A1 activation", glog.LogCharacterEvent, c.Index).
			Write("ends_on", c.Core.F+300)
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
				if c.Core.F > x.GetTag(c4Key) {
					c.c4Stacks = 0
				}
				m[attributes.DmgP] = float64(c.c4Stacks) * 0.05
				return m, c.c4Stacks > 0
			},
		})
	}
}
