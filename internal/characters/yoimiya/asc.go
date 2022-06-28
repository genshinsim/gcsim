package yoimiya

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// During Niwabi Fire-Dance, shots from Yoimiya's Normal Attack will increase her Pyro DMG Bonus by 2% on hit.
// This effect lasts for 3s and can have a maximum of 10 stacks.
func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{Base: modifier.NewBase("yoimiya-a1", -1), AffectedStat: attributes.PyroP, Amount: func() ([]float64, bool) {
		if c.Core.Status.Duration("yoimiyaa1") > 0 {
			m[attributes.PyroP] = float64(c.a1stack) * 0.02
			return m, true
		}
		c.a1stack = 0
		return nil, false
	}})

	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.Status.Duration("yoimiyaskill") == 0 {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		//here we can add stacks up to 10
		if c.a1stack < 10 {
			c.a1stack++
		}
		c.Core.Status.Add("yoimiyaa1", 180)
		// c.a1expiry = c.Core.F + 180 // 3 seconds
		return false
	}, "yoimiya-a1")
}

// Using Ryuukin Saxifrage causes nearby party members (not including Yoimiya) to gain a 10% ATK increase for 15s.
// Additionally, a further ATK Bonus will be added on based on the number of "Tricks of the Trouble-Maker" stacks Yoimiya possesses when using Ryuukin Saxifrage.
// Each stack increases this ATK Bonus by 1%.
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.1 + float64(c.a1stack)*0.01
	for _, x := range c.Core.Player.Chars() {
		if x.Index == c.Index {
			continue
		}
		x.AddStatMod(character.StatMod{Base: modifier.NewBase("yoimiya-a4", 900), AffectedStat: attributes.ATKP, Amount: func() ([]float64, bool) {
			return m, true
		}})
	}
}
