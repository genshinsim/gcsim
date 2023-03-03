package yoimiya

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "yoimiya-a1"

// During Niwabi Fire-Dance, shots from Yoimiya's Normal Attack will increase
// her Pyro DMG Bonus by 2% on hit. This effect lasts for 3s and can have a
// maximum of 10 stacks.
func (c *char) makeA1CB() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}

		if !c.StatusIsActive(a1Key) {
			c.a1Stacks = 0
		}
		if c.a1Stacks < 10 {
			c.a1Stacks++
		}

		m := make([]float64, attributes.EndStatType)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(a1Key, 3*60),
			AffectedStat: attributes.PyroP,
			Amount: func() ([]float64, bool) {
				m[attributes.PyroP] = float64(c.a1Stacks) * 0.02
				return m, true
			},
		})
	}
}

// Using Ryuukin Saxifrage causes nearby party members (not including Yoimiya)
// to gain a 10% ATK increase for 15s. Additionally, a further ATK Bonus will be
// added on based on the number of "Tricks of the Trouble-Maker" stacks Yoimiya
// possesses when using Ryuukin Saxifrage. Each stack increases this ATK Bonus
// by 1%.
func (c *char) a4() {
	c.a4Bonus[attributes.ATKP] = 0.1 + float64(c.a1Stacks)*0.01
	for _, x := range c.Core.Player.Chars() {
		if x.Index == c.Index {
			continue
		}
		x.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("yoimiya-a4", 900),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return c.a4Bonus, true
			},
		})
	}
}
