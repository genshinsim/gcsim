package yunjin

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c2Key = "yunjin-c2"
	c6Key = "yunjin-c6"
)

// After Cliffbreaker's Banner is unleashed, all nearby party members' Normal Attack DMG is increased by 15% for 12s.
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = .15
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(c2Key, 12*60),
			Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if ae.Info.AttackTag == attacks.AttackTagNormal {
					return m, true
				}
				return nil, false
			},
		})
	}
}

func (c *char) deleteC2() {
	if c.Base.Cons < 2 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.DeleteStatus(c2Key)
	}
}

// When Yun Jin triggers the Crystallize Reaction, her DEF is increased by 20% for 12s.
func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c.c4bonus = make([]float64, attributes.EndStatType)
	c.c4bonus[attributes.DEFP] = .2
	charModFunc := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}

		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("yunjin-c4", 12*60),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return c.c4bonus, true
			},
		})

		return false
	}
	c.Core.Events.Subscribe(event.OnCrystallizeCryo, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(event.OnCrystallizePyro, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(event.OnCrystallizeHydro, charModFunc, "yunjin-c4")
}

// Characters under the effects of the Flying Cloud Flag Formation have their Normal ATK SPD increased by 12%.
func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		this := char
		this.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c6Key, 12*60),
			AffectedStat: attributes.AtkSpd,
			Amount: func() ([]float64, bool) {
				//TODO: i assume this buff should go away if stacks are gone?
				if this.Tags[burstBuffKey] == 0 {
					return nil, false
				}
				return c.c6bonus, true
			},
		})
	}
}

func (c *char) deleteC6() {
	if c.Base.Cons < 6 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.DeleteStatus(c6Key)
	}
}
