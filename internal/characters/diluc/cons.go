package diluc

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	if c.Core.Combat.DamageMode {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.15
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("diluc-c1", -1),
			Amount: func(_ *info.AttackEvent, t info.Target) []float64 {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil
				}
				if x.HP()/x.MaxHP() > 0.5 {
					return m
				}
				return nil
			},
		})
	}
}

const (
	c2ICDKey  = "diluc-c2-icd"
	c2BuffKey = "diluc-c2"
)

func (c *char) c2() {
	c.c2buff = make([]float64, attributes.EndStatType)
	// we use OnPlayerHit here because he just has to get hit but triggers even if shielded
	c.Core.Events.Subscribe(event.OnPlayerHit, func(args ...any) bool {
		char := args[0].(int)
		// don't trigger if diluc was not hit
		if char != c.Index() {
			return false
		}
		if c.StatusIsActive(c2ICDKey) {
			return false
		}
		// if buff no longer active, reset stack back to 0
		if !c.StatModIsActive(c2BuffKey) {
			c.c2stack = 0
		}
		c.c2stack++
		if c.c2stack > 3 {
			c.c2stack = 3
		}
		c.c2buff[attributes.ATKP] = 0.1 * float64(c.c2stack)
		c.c2buff[attributes.AtkSpd] = 0.05 * float64(c.c2stack)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c2BuffKey, 600),
			AffectedStat: attributes.NoStat,
			Amount: func() []float64 {
				return c.c2buff
			},
		})
		return false
	}, "diluc-c2")
}

const c4BuffKey = "diluc-c4"

func (c *char) c4() {
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(c4BuffKey, 120),
		Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
			// should only affect skill dmg
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil
			}
			return c.c4buff
		},
	})
}
