package yoimiya

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "yoimiyaa1"

// During Niwabi Fire-Dance, shots from Yoimiya's Normal Attack will increase
// her Pyro DMG Bonus by 2% on hit. This effect lasts for 3s and can have a
// maximum of 10 stacks.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	// TODO: change this to add mod on each hit instead
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("yoimiya-a1", -1),
		AffectedStat: attributes.PyroP,
		Amount: func() ([]float64, bool) {
			if c.StatusIsActive(a1Key) {
				c.a1bonus[attributes.PyroP] = float64(c.a1stack) * 0.02
				return c.a1bonus, true
			}
			c.a1stack = 0
			return nil, false
		},
	})

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if !c.StatusIsActive(skillKey) {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		// here we can add stacks up to 10
		if c.a1stack < 10 {
			c.a1stack++
		}
		c.AddStatus(a1Key, 180, true)
		return false
	}, "yoimiya-a1")
}

// Using Ryuukin Saxifrage causes nearby party members (not including Yoimiya)
// to gain a 10% ATK increase for 15s. Additionally, a further ATK Bonus will be
// added on based on the number of "Tricks of the Trouble-Maker" stacks Yoimiya
// possesses when using Ryuukin Saxifrage. Each stack increases this ATK Bonus
// by 1%.
func (c *char) a4() {
	c.a4Bonus[attributes.ATKP] = 0.1 + float64(c.a1stack)*0.01
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
